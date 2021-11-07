package main

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

func Encode(action Action, key string, val string) []byte {
	byteArr := []byte{}
	switch action {
	case Set:
		byteArr = append(byteArr, CommandSet)
	case Get:
		byteArr = append(byteArr, CommandGet)
	}

	byteArr = append(byteArr, byte(len(key)))
	byteArr = append(byteArr, byte(len(val)))

	byteArr = append(byteArr, []byte(key)...)
	byteArr = append(byteArr, []byte(val)...)

	return byteArr
}
