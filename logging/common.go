package logging

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

const LOG_ENTRIES_CHAN_SIZE = 5000

type ClientConfig struct {
	PushURL            string
	Labels             map[string]string
	BatchWait          time.Duration
	BatchEntriesNumber int
}

func NewClientConfig(labels map[string]string) ClientConfig {
	return ClientConfig{
		PushURL:            viper.GetString("loki_endpoint"),
		BatchWait:          5 * time.Second,
		BatchEntriesNumber: 10000,
		Labels:             labels,
	}
}

// http.Client wrapper for adding new methods, particularly sendJsonReq
type httpClient struct {
	parent http.Client
}

// A bit more convenient method for sending requests to the HTTP server
func (client *httpClient) sendJsonReq(method, url string, ctype string, reqBody []byte) (resp *http.Response, resBody []byte, err error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Content-Type", ctype)

	resp, err = client.parent.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	resBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	return resp, resBody, nil
}
