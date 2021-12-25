package main

import (
	"fmt"
	"log"
	"syscall"
)

/*
Synchronous tcp/ip socket server with some http protocol implementation (RFC2616/1.1)

Purely just to experiment with syscalls in Go via libc wrapper and
to learn more about generic socket handling and prep for epoll/poll/select study.

see:

- https://tldp.org/LDP/tlk/net/net.html (great starter article, although I used it for binding)

- https://man7.org/linux/man-pages/man2/socketcall.2.html

- https://man7.org/linux/man-pages/man2/socket.2.html

- https://man7.org/linux/man-pages/man2/listen.2.html

- https://man7.org/linux/man-pages/man2/accept.2.html

- https://developer.ibm.com/articles/au-tcpsystemcalls/

- https://datatracker.ietf.org/doc/html/rfc2616

- https://www.w3.org/Protocols/HTTP/AsImplemented.html (more for general study of historic versions)
*/

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}

type HttpServer struct {
	// A request multiplexer capable of taking additional
	// paths to be registered for the API user to respond to
	Router Router
	// Internal use file descriptor for the server binding address
	serverSocketFd int
}

func (hs *HttpServer) createSocket() {
	var err error
	hs.serverSocketFd, err = syscall.Socket(
		syscall.AF_INET,
		syscall.SOCK_STREAM,
		0, // read more into SOCK_NONBLOCK & SOCK_CLOEXEC
	)
	handleErr(err)
}

func (hs *HttpServer) bindSocket(address []byte, port int) {
	err := syscall.Bind(hs.serverSocketFd, &syscall.SockaddrInet4{
		Port: port,
		Addr: [4]byte{
			address[0],
			address[1],
			address[2],
			address[3],
		},
	})
	handleErr(err)
}

func (hs *HttpServer) listenSocket(backlog int) {
	err := syscall.Listen(hs.serverSocketFd, backlog)
	handleErr(err)
	addr, err := syscall.Getsockname(hs.serverSocketFd)
	handleErr(err)
	addrInet4, ok := addr.(*syscall.SockaddrInet4)
	if ok {
		log.Printf(
			"server listening on address: %d.%d.%d.%d:%d",
			addrInet4.Addr[0],
			addrInet4.Addr[1],
			addrInet4.Addr[2],
			addrInet4.Addr[3],
			addrInet4.Port,
		)
	}
}

func (hs *HttpServer) Listen(address []byte, port int, backlog int) {
	hs.createSocket()
	hs.bindSocket(address, port)
	hs.listenSocket(backlog)
	hs.acceptIncomingConnections()
}

func (hs *HttpServer) acceptIncomingConnections() {
	for {
		incomingSocketFd, _, _ := syscall.Accept(hs.serverSocketFd)

		method, path, proto := getMethodPathProto(getStatusLine(incomingSocketFd))

		compliant := handleCompliance(incomingSocketFd, proto)
		if !compliant {
			break
		}

		headers, body, err := readIncomingPayload(incomingSocketFd)
		fmt.Println(headers, body, err)
		// Response
		syscall.Write(
			incomingSocketFd,
			[]byte("HTTP/1.1 200 OK\r\n"+"Content-Length: 8\r\n"+"Connection: close\r\n\r\nhi world"),
		)

		syscall.Close(incomingSocketFd)
	}
}

func main() {
	server := HttpServer{}
	server.Listen([]byte{127, 0, 0, 1}, 8000, 1)
}
