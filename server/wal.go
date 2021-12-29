package main

import (
	"encoding/gob"
	"time"

	"github.com/myprivatealaska/distributed-systems/common"
)

type logRecord struct {
	Stamp int64
	Key   string
	Val   string
}

func (s *server) writeToLog(key string, val string) {
	r := logRecord{
		Stamp: time.Time{}.UnixNano(),
		Key:   key,
		Val:   val,
	}

	err := gob.NewEncoder(s.walFD).Encode(r)
	common.CheckError(err)
	s.lastUpdatedStamp = r.Stamp
}
