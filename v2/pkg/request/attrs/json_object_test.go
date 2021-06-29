package attrs

import (
	"io"
	"reflect"
	"testing"
)

func TestNewJSONObject(t *testing.T) {
	tests := []struct {
		name string
		want JSONObject
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := NewJSONObject(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewJSONObject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJSONObject_String(t *testing.T) {
	type fields struct {
		Args []string
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
			x := &JSONObject{
				Args: tt.fields.Args,
			}
			if got := x.String(); got != tt.want {
				t.Errorf("JSONObject.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJSONObject_Set(t *testing.T) {
	type fields struct {
		Args []string
	}
	type args struct {
		v string
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
			x := &JSONObject{
				Args: tt.fields.Args,
			}
			if err := x.Set(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("JSONObject.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestJSONObject_Exists(t *testing.T) {
	type fields struct {
		Args []string
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
			x := JSONObject{
				Args: tt.fields.Args,
			}
			if got := x.Exists(); got != tt.want {
				t.Errorf("JSONObject.Exists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJSONObject_Parse(t *testing.T) {
	type fields struct {
		Args []string
	}
	tests := []struct {
		name    string
		fields  fields
		want    io.Reader
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			x := JSONObject{
				Args: tt.fields.Args,
			}
			got, err := x.Parse()
			if (err != nil) != tt.wantErr {
				t.Errorf("JSONObject.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JSONObject.Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
