package logger

import (
	"encoding/json"
	"fmt"
	"time"
)

type logStruct struct {
	T  time.Time   `json:"time"`
	ID string      `json:"requestId,omitempty"`
	M  interface{} `json:"message"`
}

// simple logger :)
func L(requestId string, msg interface{}) {
	ls := logStruct{
		T:  time.Now().Local(),
		ID: requestId,
	}

	switch m := msg.(type) {
	case error:
		ls.M = m.Error()
	case string:
		ls.M = m
	case []byte:
		ls.M = string(m)
	default:
		ls.M = m
	}

	bts, _ := json.Marshal(ls)
	fmt.Println(string(bts))
}
