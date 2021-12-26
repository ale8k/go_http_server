package main

import (
	"reflect"
	"testing"
)

func TestResponse_parse(t *testing.T) {
	type fields struct {
		status  int
		headers map[string]string
		body    []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			"it parses the http msg correctly with all fields populated",
			fields{
				status:  200,
				headers: map[string]string{"bye": "bye"},
				body:    []byte("yolo"),
			},
			[]byte{72, 84, 84, 80, 47, 49, 46, 49, 32, 50, 48, 48, 32, 84, 79, 68, 79, 13, 10, 98, 121, 101, 58, 32, 98, 121, 101, 13, 10, 13, 10, 121, 111, 108, 111},
		},
		{
			"it parses the http msg correctly with no body fields populated",
			fields{
				status:  200,
				headers: map[string]string{"bye": "bye"},
				body:    []byte{},
			},
			[]byte{72, 84, 84, 80, 47, 49, 46, 49, 32, 50, 48, 48, 32, 84, 79, 68, 79, 13, 10, 98, 121, 101, 58, 32, 98, 121, 101, 13, 10, 13, 10},
		},
		{
			"it parses the http msg correctly with no body or header fields populated",
			fields{
				status:  200,
				headers: map[string]string{},
				body:    []byte{},
			},
			[]byte{72, 84, 84, 80, 47, 49, 46, 49, 32, 50, 48, 48, 32, 84, 79, 68, 79, 13, 10, 13, 10},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hr := &Response{
				status:  tt.fields.status,
				headers: tt.fields.headers,
				body:    tt.fields.body,
			}
			if got := hr.parse(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Response.parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResponse_AddHeader(t *testing.T) {
	type fields struct {
		status  int
		headers map[string]string
		body    []byte
	}
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"it adds a header as expected",
			fields{
				status:  200,
				headers: map[string]string{},
				body:    []byte{},
			},
			args{key: "diglett", value: "dugtrio"},
			false,
		},
		{
			"it rejects the header addition as it is content-length",
			fields{
				status:  200,
				headers: map[string]string{},
				body:    []byte{},
			},
			args{key: "content-length", value: "dugtrio"},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hr := &Response{
				status:  tt.fields.status,
				headers: tt.fields.headers,
				body:    tt.fields.body,
			}
			if err := hr.AddHeader(tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Response.AddHeader() error = %v, wantErr %v", err, tt.wantErr)
			} else if err == nil && hr.headers[tt.args.key] != tt.args.value {
				t.Errorf("Response.AddHeader() did not add header as expected")
			}
		})
	}
}

func TestResponse_SetBody(t *testing.T) {
	type fields struct {
		status  int
		headers map[string]string
		body    []byte
	}
	type args struct {
		val []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []byte
	}{
		{
			"it sets the body correctly",
			fields{
				status:  200,
				headers: map[string]string{},
				body:    []byte{69},
			},
			args{
				val: []byte{69},
			},
			[]byte{69},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hr := &Response{
				status:  tt.fields.status,
				headers: tt.fields.headers,
				body:    tt.fields.body,
			}
			hr.SetBody(tt.args.val)
			if !reflect.DeepEqual(hr.body, tt.want) {
				t.Errorf("Response.SetBody() did not set body correctly")
			}
		})
	}
}

func TestResponse_SetStatus(t *testing.T) {
	type fields struct {
		status  int
		headers map[string]string
		body    []byte
	}
	type args struct {
		status int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			"it sets the status-line correctly",
			fields{
				status:  200,
				headers: map[string]string{},
				body:    []byte{},
			},
			args{
				status: 200,
			},
			200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hr := &Response{
				status:  tt.fields.status,
				headers: tt.fields.headers,
				body:    tt.fields.body,
			}
			hr.SetStatus(tt.args.status)
			if hr.status != tt.want {
				t.Errorf("Response.SetStatus() did not set the status correctly")
			}
		})
	}
}
