package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/rzaripov1990/simple-golang-proxy/responses"
	"github.com/rzaripov1990/simple-golang-proxy/utils/logger"
)

func DoRequest(requestId, method, url string, request interface{}, request_header http.Header) (response *responses.ProxyResponse, err error) {
	var (
		//bodystring string
		body  *bytes.Buffer
		inreq *http.Request
	)

	if request != nil {
		body = new(bytes.Buffer)

		switch reqtype := request.(type) {
		case string:
			body.WriteString(reqtype)
		case []byte:
			body.Write(reqtype)
		case json.RawMessage:
			body.Write(reqtype)
		default:
			json.NewEncoder(body).Encode(reqtype)
		}
	}

	if body != nil {
		//bodystring = body.String()
		inreq, err = http.NewRequestWithContext(context.TODO(), method, url, body)
	} else {
		inreq, err = http.NewRequestWithContext(context.TODO(), method, url, nil)
	}
	if err != nil {
		logger.L(requestId, err)
		return
	}

	/*
		var logReq = requests.ProxyRequests{
			Url:     url,
			Method:  method,
			Body:    json.RawMessage(bodystring),
			Headers: make(map[string]interface{}),
		}
	*/

	for k, v := range request_header {
		inreq.Header.Set(k, v[0])
		//logReq.Headers[k] = string(v[0])
	}
	//logger.L(requestId, logReq)

	c := &http.Client{}
	resp, err := c.Do(inreq)
	if err != nil {
		logger.L(requestId, err)
		return
	}
	defer resp.Body.Close()

	var bts []byte
	bts, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.L(requestId, err)
		return
	}

	length := resp.ContentLength
	if length < 0 {
		length = int64(len(bts))
	}

	response = &responses.ProxyResponse{
		Headers:    resp.Header,
		Length:     length,
		RequestID:  requestId,
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
		Response:   string(bts),
	}

	return
}
