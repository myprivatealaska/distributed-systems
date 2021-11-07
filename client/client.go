package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"

	"github.com/myprivatealaska/distributed-systems/common"
	"github.com/myprivatealaska/distributed-systems/protocol"
)

func main() {

	tcpAddr, resolveErr := net.ResolveTCPAddr("tcp4", ":9000")
	common.CheckError(resolveErr)

	var src string
	var err error

	fmt.Println("Welcome to the Pile - your key/value store!\n To interact with the Pile, use the following commands:\n 'set key value' or 'get key'.")

	for {
		src = readStdin()

		action, key, val, parseErr := parseInput(src)

		if parseErr != nil {
			fmt.Printf("Error: %v \n", parseErr.Error())
			continue
		}

		conn, dialErr := net.DialTCP("tcp", nil, tcpAddr)
		common.CheckError(dialErr)

		payload := protocol.Encode(action, key, val)
		_, err = conn.Write(payload)
		common.CheckError(err)

		result, readErr := ioutil.ReadAll(conn)
		common.CheckError(readErr)
		fmt.Println(string(result))

		conn.Close()
	}
}

func readStdin() (buf string) {
	r := bufio.NewReader(os.Stdin)

	line, _ := r.ReadString('\n')
	return line
}

func parseInput(src string) (common.Action, string, string, error) {
	parts := strings.Split(strings.TrimSpace(src), " ")

	partsCount := len(parts)

	if partsCount < 2 || partsCount > 3 {
		return "", "", "", fmt.Errorf("invalid input. Should be of the form 'get key' or 'set key value'")
	}

	switch parts[0] {
	case string(common.Get):
		if partsCount > 2 {
			return "", "", "", fmt.Errorf("invalid input. Should be of the form 'get key'")
		}
		return common.Get, parts[1], "", nil
	case string(common.Set):
		if partsCount < 3 {
			return "", "", "", fmt.Errorf("invalid input. Should be of the form 'set key value'")
		}
		return common.Set, parts[1], parts[2], nil
	default:
		return "", "", "", fmt.Errorf("invalid action: %v. Should be get or set'", parts[0])
	}
}
