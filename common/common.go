package common

import (
	"fmt"
	"strings"
)

type Action string

const Get Action = "get"
const Set Action = "set"

func checkErr(err error) {
	if err != nil {
		panic(err)
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
