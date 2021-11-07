package main

import (
	"fmt"
	"testing"
)

func TestEncode(t *testing.T) {
	tests := []struct {
		name     string
		action   Action
		key      string
		val      string
		expected []byte
	}{
		{
			name:     "Regular",
			action:   Set,
			key:      "key1",
			val:      "val1",
			expected: []byte{0x00},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := Encode(tt.action, tt.key, tt.val)
			fmt.Printf("%d\n", res)
		})
	}
}
