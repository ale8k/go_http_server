package main

import (
	"errors"
	"strconv"
	"strings"
)

// Generic wrapper for Http/1.1 response object
type Response struct {
	status  int
	headers map[string]string
	body    []byte
}

// Parses the current response into a valid HTTP response msg
func (hr *Response) parse() []byte {
	msg := make([]byte, 0)
	// []byte("HTTP/1.1 200 OK\r\n"+"Content-Length: 8\r\n"+"Connection: close\r\n\r\nhi world")

	// Write status-line
	msg = append(msg, []byte("HTTP/1.1 "+strconv.Itoa(hr.status)+" "+"TODO"+"\r\n")...)
	// Write headers
	for k, v := range hr.headers {
		msg = append(msg, []byte(k+": "+v+"\r\n")...)
	}
	// Write header terminator
	msg = append(msg, []byte("\r\n")...)
	// Write body
	msg = append(msg, hr.body...)
	return msg
}

// Adds a header to the response object
func (hr *Response) AddHeader(key string, value string) error {
	if strings.ToLower(key) == "content-length" {
		return errors.New("cannot set 'Content-Length' header manually")
	}
	hr.headers[key] = value
	return nil
}

// Sets the body of the response, overwriting any existing data in the buffer
func (hr *Response) SetBody(val []byte) {
	hr.headers["Content-Length"] = strconv.Itoa(len(string(val)))
	hr.body = val
}

// Sets response status code
// TODO: handle valid status codes only
func (hr *Response) SetStatus(status int) {
	hr.status = status
}
