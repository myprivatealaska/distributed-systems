package protocol

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/myprivatealaska/distributed-systems/common"
)

func TestEncode(t *testing.T) {
	tests := []struct {
		name     string
		action   common.Action
		key      string
		val      string
		expected []byte
	}{
		{
			name:     "Regular Set",
			action:   common.Set,
			key:      "key1",
			val:      "val1",
			expected: []byte{0, 4, 4, 107, 101, 121, 49, 118, 97, 108, 49},
		},
		{
			name:     "Regular Get",
			action:   common.Get,
			key:      "key1",
			expected: []byte{1, 4, 107, 101, 121, 49},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := Encode(tt.action, tt.key, tt.val)
			if !reflect.DeepEqual(tt.expected, res) {
				fmt.Printf("Failed %v\n", tt.name)
				t.Fail()
			}
		})
	}
}

func TestDecode(t *testing.T) {
	tests := []struct {
		name    string
		action  common.Action
		key     string
		val     string
		message []byte
	}{
		{
			name:    "Regular Set",
			action:  common.Set,
			key:     "key1",
			val:     "val1",
			message: []byte{0, 4, 4, 107, 101, 121, 49, 118, 97, 108, 49},
		},
		{
			name:    "Regular Get",
			action:  common.Get,
			key:     "key1",
			message: []byte{1, 4, 107, 101, 121, 49},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decodedAction, decodedKey, decodedVal := Decode(tt.message)
			if decodedAction != tt.action {
				fmt.Printf("Failed to decode action %v\n", tt.name)
				t.Fail()
			}
			if decodedKey != tt.key {
				fmt.Printf("Failed to decode key %v\n", tt.name)
				t.Fail()
			}
			if decodedVal != tt.val {
				fmt.Printf("Failed to decode val %v\n", tt.name)
				t.Fail()
			}
		})
	}
}
