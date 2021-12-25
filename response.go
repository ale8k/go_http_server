package main

import (
	"errors"
	"strings"
)

// Generic wrapper for Http/1.1 response object
type HttpResponse struct {
	headers map[string]string
	body    []byte
}

// Adds a header to the response object
func (hr *HttpResponse) AddHeader(key string, value string) error {
	if strings.ToLower(key) == "content-length" {
		return errors.New("cannot set 'Content-Length' header manually")
	}
	hr.headers[key] = value
	return nil
}

// Sets the body of the response, overwriting any existing data in the buffer
func (hr *HttpResponse) SetBody(val []byte) {
	hr.body = val
}
