package protocol

import (
	"github.com/myprivatealaska/distributed-systems/common"
)

// +--------------+-------------+---------------+----------------+-----------------+
// |              |             |               |                |                 |
// | Command +    |  Key Length |  Value Length |      Key       |      Value      |
// | Collection   |			   	|			    |                |                 |
// |              |  		   	|               |                |                 |
// | (3 + 5 bits) |  (1 byte)   |   (1 byte)    |  (Variable,    |  (Variable,     |
// |              |            	|               |  30 bytes max) |   30 bytes max) |
// +--------------+-------------+---------------+----------------+-----------------+

type Command int

const (
	CommandSet = iota
	CommandGet
)

func Encode(action common.Action, key string, val string) []byte {
	byteArr := []byte{}
	switch action {
	case common.Set:
		byteArr = append(byteArr, CommandSet)
	case common.Get:
		byteArr = append(byteArr, CommandGet)
	}

	byteArr = append(byteArr, byte(len(key)))

	valLen := len(val)
	if valLen > 0 {
		byteArr = append(byteArr, byte(len(val)))
	}

	byteArr = append(byteArr, []byte(key)...)
	if valLen > 0 {
		byteArr = append(byteArr, []byte(val)...)
	}

	return byteArr
}

func Decode(message []byte) (action common.Action, key string, val string) {
	actionByte := int(message[0])

	switch actionByte {
	case CommandSet:
		action = common.Set
	case CommandGet:
		action = common.Get
	}

	keyLengthByte := int(message[1])

	if action == common.Set {
		valLengthByte := int(message[2])
		keyBytes := message[3:(3 + keyLengthByte)]
		key = string(keyBytes)

		valBytes := message[3+keyLengthByte : (3 + keyLengthByte + valLengthByte)]
		val = string(valBytes)
	}

	if action == common.Get {
		keyBytes := message[2:(2 + keyLengthByte)]
		key = string(keyBytes)
	}

	return
}
