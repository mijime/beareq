package attrs

import (
	"io"
	"reflect"
	"testing"
)

func TestNewHTTPBody(t *testing.T) {
	tests := []struct {
		name string
		want HTTPBody
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHTTPBody(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHTTPBody() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHTTPBody_Exists(t *testing.T) {
	type fields struct {
		Reader io.ReadCloser
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &HTTPBody{
				Reader: tt.fields.Reader,
			}
			if got := b.Exists(); got != tt.want {
				t.Errorf("HTTPBody.Exists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHTTPBody_String(t *testing.T) {
	type fields struct {
		Reader io.ReadCloser
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
			b := &HTTPBody{
				Reader: tt.fields.Reader,
			}
			if got := b.String(); got != tt.want {
				t.Errorf("HTTPBody.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHTTPBody_Set(t *testing.T) {
	type fields struct {
		Reader io.ReadCloser
	}
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &HTTPBody{
				Reader: tt.fields.Reader,
			}
			if err := b.Set(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("HTTPBody.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
