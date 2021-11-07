package client

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"os"

	"github.com/myprivatealaska/distributed-systems/common"
)

func main() {

	tcpAddr, resolveErr := net.ResolveTCPAddr("tcp4", ":9000")
	checkErr(resolveErr)

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
		checkErr(dialErr)

		_, err = conn.Write([]byte(fmt.Sprintf("%v %v %v", action, key, val)))
		checkErr(err)
		result, readErr := ioutil.ReadAll(conn)
		checkErr(readErr)
		fmt.Println(string(result))

		conn.Close()
	}
}

func readStdin() (buf string) {
	r := bufio.NewReader(os.Stdin)

	line, _ := r.ReadString('\n')
	return line
}
