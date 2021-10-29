package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

var (
	mutex    sync.RWMutex
	memory   = map[string]string{}
	dataFile *os.File
)

type Action string

const Get Action = "get"
const Set Action = "set"

func main() {
	var src string
	var err error

	readFromDisk()

	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	dataFile, err = os.OpenFile(fmt.Sprintf("%v/data.json", currentDir), os.O_WRONLY, 0777)
	if err != nil {
		panic(err)
	}
	defer dataFile.Close()

	fmt.Println("PILE: what's up?")

	for true {
		src = readStdin()

		action, key, val, err := parseInput(src)
		if err == nil {
			println("---------------------------------")
			switch action {
			case Get:
				mutex.RLock()
				fmt.Println(memory[key])
				mutex.RUnlock()
			case Set:
				mutex.Lock()
				memory[key] = val
				fmt.Println("Set")
				writeToDisk()
				mutex.Unlock()
			}
		} else {
			fmt.Println("== Error ========")
			fmt.Println(err)
		}
	}
}

func parseInput(src string) (Action, string, string, error) {
	parts := strings.Split(strings.TrimSpace(src), " ")

	partsCount := len(parts)

	if partsCount < 2 || partsCount > 3 {
		return "", "", "", fmt.Errorf("invalid input. Should be of the form 'get key' or 'set key value'")
	}

	switch parts[0] {
	case string(Get):
		if partsCount > 2 {
			return "", "", "", fmt.Errorf("invalid input. Should be of the form 'get key'")
		}
		return Get, parts[1], "", nil
	case string(Set):
		if partsCount < 3 {
			return "", "", "", fmt.Errorf("invalid input. Should be of the form 'set key value'")
		}
		return Set, parts[1], parts[2], nil
	default:
		return "", "", "", fmt.Errorf("invalid action: %v. Should be get or set'", parts[0])
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

func readFromDisk() {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	file, fileErr := os.OpenFile(fmt.Sprintf("%v/data.json", currentDir), os.O_RDONLY|os.O_CREATE, 0777)
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

func readStdin() (buf string) {
	r := bufio.NewReader(os.Stdin)

	line, _ := r.ReadString('\n')
	return line
}
