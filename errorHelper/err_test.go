package errorHelper

import (
	"errors"
	"fmt"
	"testing"
)

func TestErrorWithTrace_Error(t *testing.T) {
	type fields struct {
		error error
		trace string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ErrorWithTrace{
				error: tt.fields.error,
				trace: tt.fields.trace,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrorWithTrace_UnWarp(t *testing.T) {
	type fields struct {
		error error
		trace string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ErrorWithTrace{
				error: tt.fields.error,
				trace: tt.fields.trace,
			}
			if err := e.Unwrap(); (err != nil) != tt.wantErr {
				t.Errorf("UnWarp() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWithTrace(t *testing.T) {
	err := errors.New("test")
	err = WithTrace(err)

	err = fmt.Errorf("2 err:%w", err)

	target := &ErrorWithTrace{}
	a := errors.As(err, &target)

	trace := Trace(err)
	println(err.Error(), a)
	println(trace)
}
