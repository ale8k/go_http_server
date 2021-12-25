package main

// Represents a function to 'handle' the incoming request
// provided it is multiplexed routed to
type HandlerCallback func(req *Request, res *Response)

// The request path
type RequestPath string

// Represents a multiplex router responsible for handling
// incoming requests based on their path (just in a map for now)
type Router struct {
	Handlers map[RequestPath]HandlerCallback
}

// Adds a handler to the router, if a path already exists for the provided handler
// it overwrites the existing callback
// TODO: validate path in regex before adding handler, and return err if no valid
func (r *Router) AddHandler(path RequestPath, callback HandlerCallback) {
	r.Handlers[path] = callback
}

// Find a handler for a given requestpath
func (r *Router) FindHandler(path RequestPath) HandlerCallback {
	return r.Handlers[path]
}
