package attrs

import (
	"net/http"
	"reflect"
	"testing"
)

func TestNewHTTPHeader(t *testing.T) {
	tests := []struct {
		name string
		want HTTPHeader
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := NewHTTPHeader(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHTTPHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHTTPHeader_String(t *testing.T) {
	type fields struct {
		Values http.Header
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			h := &HTTPHeader{
				Values: tt.fields.Values,
			}
			if got := h.String(); got != tt.want {
				t.Errorf("HTTPHeader.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHTTPHeader_Set(t *testing.T) {
	type fields struct {
		Values http.Header
	}
	type args struct {
		values string
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			h := &HTTPHeader{
				Values: tt.fields.Values,
			}
			if err := h.Set(tt.args.values); (err != nil) != tt.wantErr {
				t.Errorf("HTTPHeader.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
