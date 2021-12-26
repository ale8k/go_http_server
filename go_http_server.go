package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"syscall"
	"time"
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
		// SOCK_NONBLOCK & SOCK_CLOEXEC ensure net polling and close on exec
		// we should be able to bitwise OR but uhh, no luck...
		syscall.SOCK_STREAM,
		0,
	)
	// apparently this is it in go's wrapper
	syscall.CloseOnExec(hs.serverSocketFd)
	syscall.SetNonblock(hs.serverSocketFd, true)

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

// Break into own func for better profiling
func (hs *HttpServer) respondToRequest(incomingSocketFd int) {
	method, path, proto := getMethodPathProto(getStatusLine(incomingSocketFd))
	compliant := handleCompliance(incomingSocketFd, proto)
	if !compliant {
		return
	}

	headers, body, err := readIncomingPayload(incomingSocketFd)
	handleErr(err)
	req := &Request{Headers: headers, Body: body}
	res := &Response{headers: make(map[string]string)}
	cb := hs.Router.FindHandler(method, path)
	if cb != nil {
		cb(req, res)
		_, err := syscall.Write(incomingSocketFd, res.parse())
		handleErr(err)
	} else {
		res.SetStatus(404)
		res.SetBody([]byte("No route matching " + method + ":" + path))
		_, err := syscall.Write(incomingSocketFd, res.parse())
		handleErr(err)
	}
	syscall.Close(incomingSocketFd)
}

func (hs *HttpServer) acceptIncomingConnections() {
	reqCount := &struct{ count int }{count: 0}
	// Can we have many accepts ...?
	for {
		// TODO: syscall.Select() open N descriptors to be selected from
		// for accepting the call, for now creating them on the fly is ok
		// https://man7.org/linux/man-pages/man2/select.2.html
		// socket returns a network file descriptor that is ready for asynchronous I/O using the network poller.
		incomingSocketFd, _, err := syscall.Accept(hs.serverSocketFd)

		reqCount.count++
		fmt.Printf("incoming request on: %v, err: %v, reqcount: %v\n", incomingSocketFd, err, reqCount.count)
		go hs.respondToRequest(incomingSocketFd)
	}
}

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	server := HttpServer{Router: Router{Handlers: make(map[RequestPath]HandlerCallback)}}

	server.Router.AddHandler("GET", "/1", func(req *Request, res *Response) {
		res.SetBody([]byte("Hello World"))
	})

	server.Router.AddHandler("GET", "/2", func(req *Request, res *Response) {
		res.AddHeader("Connection", "close")
		res.SetBody([]byte("test"))
		time.Sleep(time.Second * 10)
	})

	server.Router.AddHandler("GET", "/1/2", func(req *Request, res *Response) {
		res.AddHeader("Connection", "close")
		res.SetBody([]byte("test"))
		time.Sleep(time.Second * 10)
	})

	server.Listen([]byte{127, 0, 0, 1}, 8000, 100000)
}
