package main

import (
	"fmt"
	"regexp"
	"strconv"
	"syscall"
)

// Protocol related utility functions, can have direct side-effects
// and contain syscalls

// Reads an incoming http/1.1 payload for the given fd, parses it, and returns the segregated pieces
// namely:
//	- headers in a map
//	- body as a raw []byte
//	- any errors that have occured
func readIncomingPayload(incomingSocketFd int) (map[string]string, []byte, error) {
	var headers map[string]string
	var body []byte
	var err error
	headerBuf := make([]byte, 0, 4096)
	for {
		buf := make([]byte, 20)
		_, err := syscall.Read(incomingSocketFd, buf)
		headerBuf = append(headerBuf, buf...)
		if err != nil {
			fmt.Println("error occured in read")
			break
		} else if body = getHeaderTermination(headerBuf); body != nil {
			headers = parseHeaders(headerBuf)
			if headers["content-length"] != "" {
				bodyLength, _ := strconv.Atoi(headers["content-length"])
				b := make([]byte, bodyLength)
				fmt.Println("attempting to read body")
				_, err := syscall.Read(incomingSocketFd, b)
				handleErr(err)
				body = append(body, b...)
			}
			break
		}
	}
	return headers, body, err
}

// Handles version compliance, currently server is 1.1 ONLY
// responds and kill the socket
func handleCompliance(fd int, statusLine string) bool {
	compliant := regexp.MustCompile(`HTTP/1.1`).MatchString(statusLine)
	if !compliant {
		syscall.Write(
			fd,
			[]byte("HTTP/1.1 505 HTTP Version Not Supported\r\n"+"Content-Length: 0\r\n\r\n"),
		)
		syscall.Close(fd)
	}
	return compliant
}
