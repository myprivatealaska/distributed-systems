package main

import (
	"encoding/gob"
	"fmt"
	"io"
	"net"
	"os"
	"sync"

	"github.com/myprivatealaska/distributed-systems/common"
	"github.com/myprivatealaska/distributed-systems/protocol"
)

var (
	mutex    sync.RWMutex
	memory   = map[string]string{}
	dataFile *os.File
)

func main() {

	args := os.Args[1:]
	port := args[0]
	storageFileName := args[1]

	// Upon start, read the data into memory
	readFromDisk(storageFileName)

	// Create a file descriptor for writing to the file
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	dataFile, err = os.OpenFile(fmt.Sprintf("%v/%v", currentDir, storageFileName), os.O_WRONLY, 0777)
	if err != nil {
		panic(err)
	}
	defer dataFile.Close()

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
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
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
		switch action {
		case common.Get:
			mutex.RLock()
			_, writerErr := conn.Write([]byte(memory[key]))
			common.CheckError(writerErr)
			mutex.RUnlock()
		case common.Set:
			mutex.Lock()
			memory[key] = val
			fmt.Println("Set")
			writeToDisk()
			mutex.Unlock()
			_, writerErr := conn.Write([]byte(memory[key]))
			common.CheckError(writerErr)
		}
		conn.Close()
	}
}

func writeToDisk() {
	err := dataFile.Truncate(0)
	if err != nil {
		panic(fmt.Sprintf("Error truncating data file %e", err))
	}
	_, err = dataFile.Seek(0, 0)
	if err != nil {
		panic(fmt.Sprintf("Error seeking to the start of the data file %e", err))
	}

	encoder := gob.NewEncoder(dataFile)
	encodeErr := encoder.Encode(&memory)

	if encodeErr != nil {
		panic(fmt.Sprintf("Error encoding storage %e", encodeErr))
	}
}

func readFromDisk(storageFileName string) {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	file, fileErr := os.OpenFile(fmt.Sprintf("%v/%v", currentDir, storageFileName), os.O_RDONLY|os.O_CREATE, 0777)
	if fileErr != nil {
		panic(fmt.Sprintf("Error reading data from disk: %e", fileErr))
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	decodeErr := decoder.Decode(&memory)
	if decodeErr != nil && decodeErr != io.EOF {
		panic(fmt.Sprintf("Error decoding storage %e", decodeErr))
	}
}
