package common

type Action string

const Get Action = "get"
const Set Action = "set"

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
