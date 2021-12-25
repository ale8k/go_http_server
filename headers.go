package main

import (
	"strings"
	"syscall"
)

// Utilities for handling headers are kept here

// Reads up to status-line and returns string representation
func getStatusLine(fd int) string {
	var err error
	statusLineBuf := make([]byte, 0)
	for {
		// we read 1 at a time such that we don't over-read
		buf := make([]byte, 1)
		handleErr(err)
		_, err = syscall.Read(fd, buf)
		statusLineBuf = append(statusLineBuf, buf[0])
		if buf[0] == 10 {
			break
		}
	}
	return string(statusLineBuf)
}

// Gets method, path and protocol version from a provided status-line as per RFC2616/1.1
func getMethodPathProto(statusLine string) (string, string, string) {
	statusLineSplit := strings.Split(statusLine, " ")
	return statusLineSplit[0], statusLineSplit[1], statusLineSplit[2]
}

// Searches for header termination, when found, returns the trailing part of the buffer
// (the body) as a slice to begin reading into from your initial header designated buffer
func getHeaderTermination(buffer []byte) []byte {
	for i := range buffer {
		crlf := buffer[i : i+4]
		if crlf[0] == 13 && crlf[1] == 10 && crlf[2] == 13 && crlf[3] == 10 {
			return buffer[i+4:]
		}
	}
	return nil
}

// Parses headers into a map[string]string
// trims all whitespace and sets the keys to lowercase for consistency
func parseHeaders(headers []byte) map[string]string {
	headerMap := make(map[string]string)
	offset := 0
	for i, v := range headers {
		if v == 13 && headers[i+1] == 10 {
			kv := strings.Split(string(headers[offset:i+1]), ":")
			if len(kv) == 2 {
				key := strings.ToLower(strings.TrimSpace(kv[0]))
				value := strings.TrimSpace(kv[1])
				headerMap[key] = value
			}
			offset = i + 2
		}
	}
	return headerMap
}
