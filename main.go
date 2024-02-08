package main

import "io"

func main() {
	// Initialize all the parameters to see valid results
	var url, methodType string
	var numRequests int
	headers := make(map[string]string)
	var body io.Reader

	generator := NewConcurrentRequest(url, methodType, numRequests, headers, body)
	generator.Generate()
}
