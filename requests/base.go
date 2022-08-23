package requests

import "encoding/json"

type ProxyRequests struct {
	Method  string                 `json:"method"`
	Url     string                 `json:"url"`
	Headers map[string]interface{} `json:"headers,omitempty"`
	Body    json.RawMessage        `json:"body,omitempty"`
}
