package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
)

type ConcurrentRequest struct {
	Url         string
	MethodType  string
	NumRequests int
	Body        io.Reader
	Headers     map[string]string
	wg          *sync.WaitGroup
	response    chan any
}

func NewConcurrentRequest(url, methodType string, numRequests int, headers map[string]string, body io.Reader) ConcurrentRequestGenerator {
	return &ConcurrentRequest{
		Url:         url,
		MethodType:  methodType,
		NumRequests: numRequests,
		Body:        body,
		Headers:     headers,
		wg:          new(sync.WaitGroup),
		response:    make(chan any, numRequests),
	}
}

type ConcurrentRequestGenerator interface {
	Generate()
}

func (c *ConcurrentRequest) Generate() {

	defer close(c.response)

	c.wg.Add(1)
	go c.generate()

	c.wg.Add(1)
	go c.getResponse()

	c.wg.Wait()
}

func (c *ConcurrentRequest) generate() {
	for i := 0; i < c.NumRequests; i++ {
		c.wg.Add(1)
		go func() {
			defer c.wg.Done()
			fmt.Println("Sending request ->")
			req, err := http.NewRequest(c.MethodType, c.Url, c.Body)
			if err != nil {
				c.response <- err
				return
			}

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer <your_access_token>")
			for key, val := range c.Headers {
				req.Header.Set(key, val)
			}

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				c.response <- err
				return
			}
			c.response <- resp
		}()
	}
}

func (c *ConcurrentRequest) getResponse() {
	for {
		select {
		case response := <-c.response:
			fmt.Println("Reading response from the channel", response)
		default:
			break
		}
	}
}
