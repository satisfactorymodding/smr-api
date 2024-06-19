package logging

import (
	"encoding/json"
	"io"
	"strconv"
	"sync"
	"time"
)

var _ io.Writer = (*ClientJSON)(nil)

type promtailStream struct {
	Labels map[string]string `json:"stream"`
	Values [][]string        `json:"values"`
}

type promtailMsg struct {
	Streams []promtailStream `json:"streams"`
}

type ClientJSON struct {
	config    *ClientConfig
	quit      chan struct{}
	entries   chan []string
	waitGroup sync.WaitGroup
	client    httpClient
}

func NewClientProto(conf ClientConfig) (*ClientJSON, error) {
	client := ClientJSON{
		config:  &conf,
		quit:    make(chan struct{}),
		entries: make(chan []string, LogEntriesChanSize),
		client:  httpClient{},
	}

	client.waitGroup.Add(1)
	go client.run()

	return &client, nil
}

func (c *ClientJSON) Write(p []byte) (int, error) {
	c.Log(string(p))
	return len(p), nil
}

func (c *ClientJSON) Log(line string) {
	c.entries <- []string{
		strconv.FormatInt(time.Now().UnixNano(), 10),
		line,
	}
}

func (c *ClientJSON) Shutdown() {
	close(c.quit)
	c.waitGroup.Wait()
}

func (c *ClientJSON) run() {
	var batch [][]string
	batchSize := 0
	maxWait := time.NewTimer(c.config.BatchWait)

	defer func() {
		if batchSize > 0 {
			c.send(batch)
		}

		c.waitGroup.Done()
	}()

	for {
		select {
		case <-c.quit:
			return
		case entry := <-c.entries:
			batch = append(batch, entry)
			batchSize++
			if batchSize >= c.config.BatchEntriesNumber {
				c.send(batch)
				batch = make([][]string, 0)
				batchSize = 0
				maxWait.Reset(c.config.BatchWait)
			}
		case <-maxWait.C:
			if batchSize > 0 {
				c.send(batch)
				batch = make([][]string, 0)
				batchSize = 0
			}
			maxWait.Reset(c.config.BatchWait)
		}
	}
}

func (c *ClientJSON) send(entries [][]string) {
	var streams []promtailStream
	streams = append(streams, promtailStream{
		Labels: c.config.Labels,
		Values: entries,
	})

	msg := promtailMsg{Streams: streams}
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		return
	}

	resp, _, err := c.client.sendJSONReq("POST", c.config.PushURL, "application/json", jsonMsg)
	if err != nil {
		return
	}

	if resp.StatusCode != 204 {
		return
	}
}
