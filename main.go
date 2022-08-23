package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/rzaripov1990/simple-golang-proxy/middlewares"
	"github.com/rzaripov1990/simple-golang-proxy/requests"
	"github.com/rzaripov1990/simple-golang-proxy/responses"
	"github.com/rzaripov1990/simple-golang-proxy/utils"
	"github.com/rzaripov1990/simple-golang-proxy/utils/api"
	"github.com/rzaripov1990/simple-golang-proxy/utils/logger"
)

type bucket struct {
	In  interface{} `json:"in"`
	Out interface{} `json:"out"`
}

var storage sync.Map

func proxy(w http.ResponseWriter, r *http.Request) {
	var (
		err              error
		requestId        = middlewares.GetRequestId(r.Context())
		req              requests.ProxyRequests
		reqHeaders       http.Header
		reqBody          interface{}
		resp             *responses.ProxyResponse
		withResponseBody bool
	)

	if r.Method != http.MethodPost {
		err = fmt.Errorf("405 Method Not Allowed")
		logger.L(requestId, err)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	withResponseBody = r.URL.Query().Get("withResponseBody") == "true"

	err = utils.JsonDecode(r, &req, false)
	if err != nil {
		logger.L(requestId, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(req.Headers) > 0 {
		reqHeaders = make(http.Header, len(req.Headers))
		for k, v := range req.Headers {
			reqHeaders.Set(k, v.(string))
		}
	}

	if req.Body != nil {
		err = json.Unmarshal(req.Body, &reqBody)
		if err != nil {
			logger.L(requestId, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	storage.Store(requestId, bucket{In: req})

	resp, err = api.DoRequest(requestId, req.Method, req.Url, reqBody, reqHeaders)
	if err != nil || resp == nil {
		if err != nil {
			logger.L(requestId, err)
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !withResponseBody {
		resp.Response = ""
		resp.Length = 0
	}

	logger.L(requestId, resp)
	jbts, _ := json.Marshal(resp)

	bckt, ok := storage.Load(requestId)
	if ok {
		b := bckt.(bucket)
		b.Out = resp
		storage.Store(requestId, b)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jbts)
}

func print(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	m := map[string]bucket{}
	storage.Range(func(key, value interface{}) bool {
		m[key.(string)] = value.(bucket)
		return true
	})

	bts, err := json.Marshal(m)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(bts)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()

	mux.Handle("/", middlewares.LoggerMiddleware(http.HandlerFunc(proxy)))
	mux.HandleFunc("/print", http.HandlerFunc(print))

	logger.L("", "service is listening on a port: "+port)
	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		logger.L("", err)
		log.Fatal(err)
	}
}
