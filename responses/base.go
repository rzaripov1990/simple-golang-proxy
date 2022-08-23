package responses

import "net/http"

type ProxyResponse struct {
	RequestID  string      `json:"id"`
	StatusCode int         `json:"statusCode"`
	Status     string      `json:"status"`
	Headers    http.Header `json:"headers"`
	Length     int64       `json:"length"`
	Response   string      `json:"response,omitempty"`
}
