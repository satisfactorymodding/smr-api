package logging

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

const LogEntriesChanSize = 5000

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

// http.Client wrapper for adding new methods, particularly sendJSONReq
type httpClient struct {
	parent http.Client
}

// A bit more convenient method for sending requests to the HTTP server
func (client *httpClient) sendJSONReq(method, url string, ctype string, reqBody []byte) (*http.Response, []byte, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, nil, fmt.Errorf("unable to create http request: %w", err)
	}

	req.Header.Set("Content-Type", ctype)

	resp, err := client.parent.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to send http request: %w", err)
	}
	defer resp.Body.Close()

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to read response body: %w", err)
	}

	return resp, resBody, nil
}
