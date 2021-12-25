package main

type Request struct {
	Headers map[string]string
	Body    []byte
}
