package attrs

import (
	"reflect"
	"testing"
)

func TestNewHTTPRequestMethod(t *testing.T) {
	tests := []struct {
		name string
		want HTTPRequestMethod
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := NewHTTPRequestMethod(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHTTPRequestMethod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHTTPRequestMethod_Exists(t *testing.T) {
	type fields struct {
		Method string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			a := &HTTPRequestMethod{
				Method: tt.fields.Method,
			}
			if got := a.Exists(); got != tt.want {
				t.Errorf("HTTPRequestMethod.Exists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHTTPRequestMethod_String(t *testing.T) {
	type fields struct {
		Method string
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
			a := &HTTPRequestMethod{
				Method: tt.fields.Method,
			}
			if got := a.String(); got != tt.want {
				t.Errorf("HTTPRequestMethod.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHTTPRequestMethod_Set(t *testing.T) {
	type fields struct {
		Method string
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			a := &HTTPRequestMethod{
				Method: tt.fields.Method,
			}
			if err := a.Set(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("HTTPRequestMethod.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
