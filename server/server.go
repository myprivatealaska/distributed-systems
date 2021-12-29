package main

import (
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"sync"

	"github.com/myprivatealaska/distributed-systems/common"
	"github.com/myprivatealaska/distributed-systems/protocol"
)

const (
	LEADER_PORT         = 8081
	SYNC_FOLLOWER_PORT  = 8082
	ASYNC_FOLLOWER_PORT = 8083

	WAL_FILE_PATH = "wal"
)

type serverRole int

const (
	leader serverRole = iota
	asyncFollower
	syncFollower
)

var portMap = map[serverRole]int{
	leader:        LEADER_PORT,
	asyncFollower: ASYNC_FOLLOWER_PORT,
	syncFollower:  SYNC_FOLLOWER_PORT,
}

type server struct {
	role             serverRole
	lastUpdatedStamp int64
	mutex            sync.RWMutex
	memory           map[string]string
	dataFD           *os.File
	walFD            *os.File
}

func newServer(serverRole serverRole, storageFileName string) *server {
	// Create a file descriptor for writing to the file
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	serv := &server{
		role:   serverRole,
		mutex:  sync.RWMutex{},
		memory: map[string]string{},
	}

	serv.readFromDisk(storageFileName)

	dataFD, fDErr := os.OpenFile(fmt.Sprintf("%v/%v", currentDir, storageFileName), os.O_WRONLY|os.O_SYNC, 0644)
	if fDErr != nil {
		panic(fDErr)
	}

	walFD, walErr := os.OpenFile(fmt.Sprintf("%v/%v", currentDir, WAL_FILE_PATH), os.O_WRONLY|os.O_CREATE|os.O_APPEND|os.O_SYNC, 0644)
	if walErr != nil {
		panic(walErr)
	}

	serv.dataFD = dataFD
	serv.walFD = walFD
	return serv
}

func main() {
	args := os.Args[1:]
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run server.go <leader|async_follower|sync_follower> data.json")
		panic(errors.New("not enough arguments"))
	}

	var role serverRole

	switch args[0] {
	case "leader":
		role = leader
	case "async_follower":
		role = asyncFollower
	case "sync_follower":
		role = syncFollower
	default:
		panic(fmt.Errorf("role not supported: %v", role))
	}

	storageFileName := args[1]
	serv := newServer(role, storageFileName)
	defer serv.dataFD.Close()

	// Upon start, read the data into memory
	serv.readFromDisk(storageFileName)
	//serv.readFromWal()

	// Set up TCP connection
	port := portMap[serv.role]
	service := fmt.Sprintf(":%v", port)
	tcpAddr, resolveErr := net.ResolveTCPAddr("tcp4", service)
	common.CheckError(resolveErr)
	listener, listenErr := net.ListenTCP("tcp", tcpAddr)
	common.CheckError(listenErr)
	defer listener.Close()

	for {
		conn, acceptErr := listener.Accept()
		if acceptErr != nil {
			continue
		}
		go serv.handleClient(conn)
	}
}

func (s *server) handleClient(conn net.Conn) {
	request := make([]byte, 128) // set maximum request length to 128B to prevent flood based attacks

	read_len, err := conn.Read(request)

	if err != nil {
		common.CheckError(err)
	}

	if read_len == 0 {
		return // connection already closed by client
	} else {
		req := request[:read_len]
		action, key, val := protocol.Decode(req)
		log.Printf("Received request: %v %v %v\n", action, key, val)
		switch action {
		case common.Get:
			s.mutex.RLock()
			_, writerErr := conn.Write([]byte(s.memory[key]))
			common.CheckError(writerErr)
			log.Printf("Served response for key: %v\n", s.memory[key])
			s.mutex.RUnlock()
		case common.Set:
			s.mutex.Lock()
			if s.role == leader {
				s.writeToLog(key, val)
				s.replicateSync(key, val)
			}
			s.memory[key] = val
			s.writeToDisk()
			_, writerErr := conn.Write([]byte(s.memory[key]))
			s.mutex.Unlock()
			common.CheckError(writerErr)
			log.Printf("Stored new value %v. Served response: %v\n", val, key)
		}
		conn.Close()
	}
}

func (s *server) replicateSync(key string, val string) {
	log.Printf("Sync replication begin: %v", key)

	tcpAddr, resolveErr := net.ResolveTCPAddr("tcp4", fmt.Sprintf(":%d", SYNC_FOLLOWER_PORT))
	common.CheckError(resolveErr)
	conn, dialErr := net.DialTCP("tcp", nil, tcpAddr)
	common.CheckError(dialErr)

	payload := protocol.Encode(common.Set, key, val)
	_, err := conn.Write(payload)
	common.CheckError(err)

	_, readErr := ioutil.ReadAll(conn)
	common.CheckError(readErr)

	conn.Close()
	log.Printf("Sync replication finished: %v", key)
}

func (s *server) writeToDisk() {
	err := s.dataFD.Truncate(0)
	if err != nil {
		panic(fmt.Sprintf("Error truncating data file %e", err))
	}
	_, err = s.dataFD.Seek(0, 0)
	if err != nil {
		panic(fmt.Sprintf("Error seeking to the start of the data file %e", err))
	}

	encodeErr := gob.NewEncoder(s.dataFD).Encode(&s.memory)

	if encodeErr != nil {
		panic(fmt.Sprintf("Error encoding storage %e", encodeErr))
	}
}

func (s *server) readFromDisk(storageFileName string) {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	file, fileErr := os.OpenFile(fmt.Sprintf("%v/%v", currentDir, storageFileName), os.O_RDONLY|os.O_CREATE, 0777)
	if fileErr != nil {
		panic(fmt.Sprintf("Error reading data from disk: %e", fileErr))
	}
	defer file.Close()

	decodeErr := gob.NewDecoder(file).Decode(&s.memory)
	if decodeErr != nil && decodeErr != io.EOF {
		panic(fmt.Sprintf("Error decoding storage %e", decodeErr))
	}
}

func (s *server) readFromWal() {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	file, fileErr := os.OpenFile(fmt.Sprintf("%v/%v", currentDir, WAL_FILE_PATH), os.O_RDONLY, 0777)
	if fileErr != nil {
		panic(fmt.Sprintf("Error reading data from disk: %e", fileErr))
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)

	for err = nil; err != io.EOF; {
		var record logRecord
		err = decoder.Decode(&record)
		if err != io.EOF && err != nil {
			continue
		}
		s.memory[record.Key] = record.Val
	}
}
