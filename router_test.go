package main

import (
	"testing"
)

func TestRouter_AddHandler(t *testing.T) {
	type fields struct {
		Handlers map[RequestPath]HandlerCallback
	}
	type args struct {
		method   string
		path     string
		callback HandlerCallback
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			"a handler is added to the handlers map",
			fields{
				Handlers: map[RequestPath]HandlerCallback{},
			},
			args{
				"GET", "/diglett", func(req *Request, res *Response) {},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Router{
				Handlers: tt.fields.Handlers,
			}
			r.AddHandler(tt.args.method, tt.args.path, tt.args.callback)
			if r.Handlers[RequestPath{tt.args.method, tt.args.path}] == nil {
				t.Errorf("Router.AddHandler() did not add handler")
			}
		})
	}
}
