package tools

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type HTTPTool struct {
	hostname string
	port     string
	timeout  time.Duration
	client   *http.Client
}

func (hT *HTTPTool) Init(hostname string, port string, timeout time.Duration) {
	hT.hostname = hostname
	hT.port = port
	hT.timeout = timeout

	hT.client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: hT.timeout,
	}
}

func (hT *HTTPTool) GET(url string) (statusCode int, respData []byte, err error) {
	var payload []byte
	return hT.request("GET", url, payload, nil)
}

func (hT *HTTPTool) POST(url string, payload []byte) (statusCode int, respData []byte, err error) {
	statusCode, respData, err = hT.request("POST", url, payload, nil)
	return
}

func (hT *HTTPTool) request(method string, url string, payload []byte, headers map[string]string) (statusCode int, respData []byte, err error) {
	fullPath := ""

	if hT.port == "" {
		fullPath = fmt.Sprintf("%v/%v", hT.hostname, url)
	} else {
		fullPath = fmt.Sprintf("%v:%v/%v", hT.hostname, hT.port, url)
	}

	var req *http.Request
	fmt.Println("path: ", fullPath)
	if len(payload) == 0 {
		req, err = http.NewRequest(method, fullPath, nil)
	} else {
		req, err = http.NewRequest(method, fullPath, bytes.NewBuffer(payload))

		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}

	if err != nil {
		log.Println("HTTPclient:", method, "NewRequest:", err)
		return
	}

	resp, err := hT.client.Do(req)
	if err != nil {
		log.Println("HTTPclient:", method, "Do:", err)
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("HTTPclient:", method, "Close:", err)
		}
	}(resp.Body)

	statusCode = resp.StatusCode
	respData, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Println("HTTPclient:", method, "ReadAll:", err)
	}
	return
}
