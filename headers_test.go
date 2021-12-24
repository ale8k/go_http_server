package main

import (
	"reflect"
	"testing"
)

func Test_getHeaderTermination(t *testing.T) {
	tests := []struct {
		name string
		args []byte
		want []byte
	}{
		{
			name: "it returns remaining body read",
			args: []byte{2, 13, 10, 13, 10, 30, 20, 10},
			want: []byte{30, 20, 10},
		},
		{
			name: "it returns empty slice when no over-read",
			args: []byte{2, 13, 10, 13, 10},
			want: []byte{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getHeaderTermination(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getHeaderTermination() = %v, want %v", got, tt.want)
			}
		})
	}
	t.Run("it panics when no \\r\\n can be found", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("getHeaderTermination() did not panic on index out of range")
			}
		}()
		getHeaderTermination([]byte{2, 10, 10, 20, 10})
	})
}

func Test_parseHeaders(t *testing.T) {
	type args struct {
		headers []byte
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseHeaders(tt.args.headers); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseHeaders() = %v, want %v", got, tt.want)
			}
		})
	}
}
