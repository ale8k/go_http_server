package main

// Represents a function to 'handle' the incoming request
// provided it is multiplexed routed to
type HandlerCallback func(req *Request, res *Response)

// Holds the method and path for a given request
type RequestPath struct {
	Method string
	Path   string
}

// Represents a multiplex router responsible for handling
// incoming requests based on their path (just in a map for now)
type Router struct {
	Handlers map[RequestPath]HandlerCallback
}

// Adds a handler to the router, if a path already exists for the provided handler
// it overwrites the existing callback
// TODO: validate path in regex before adding handler, and return err if no valid
func (r *Router) AddHandler(method string, path string, callback HandlerCallback) {
	r.Handlers[RequestPath{method, path}] = callback
}

// Find a handler for a given requestpath
func (r *Router) FindHandler(method string, path string) HandlerCallback {
	return r.Handlers[RequestPath{method, path}]
}
