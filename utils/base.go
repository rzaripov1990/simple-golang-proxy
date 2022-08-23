package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	uuid "github.com/nu7hatch/gouuid"
)

func RequestID() string {
	reqID, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	return reqID.String()
}

func JsonDecode(r *http.Request, data interface{}, noClose bool) (err error) {
	if r.Header.Get("Content-Type") == "application/json" && r.Body != http.NoBody {
		rBody, err := io.ReadAll(r.Body)
		if err != nil {
			return err
		}

		err = json.Unmarshal(rBody, data)
		if err != nil {
			return err
		}

		if noClose {
			r.Body = io.NopCloser(bytes.NewBuffer(rBody))
		}
	} /*
		else {
			err = ...
		}
	*/

	return
}

func JsonEncode(v interface{}) (s string) {
	jsonb, err := json.Marshal(v)
	if err != nil {
		return err.Error()
	}
	return string(jsonb)
}
