package openapi

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
)

func TestNewOperation(t *testing.T) {
	type args struct {
		url    string
		path   string
		method string
		op     *openapi3.Operation
	}
	tests := []struct {
		name string
		args args
		want *Operation
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewOperation(tt.args.url, tt.args.path, tt.args.method, tt.args.op); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewOperation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperation_Name(t *testing.T) {
	type fields struct {
		Operation *openapi3.Operation
		BaseURL   string
		Path      string
		Method    string
		args      map[string]map[string]*string
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
			op := &Operation{
				Operation: tt.fields.Operation,
				BaseURL:   tt.fields.BaseURL,
				Path:      tt.fields.Path,
				Method:    tt.fields.Method,
				args:      tt.fields.args,
			}
			if got := op.Name(); got != tt.want {
				t.Errorf("Operation.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperation_Parse(t *testing.T) {
	type fields struct {
		Operation *openapi3.Operation
		BaseURL   string
		Path      string
		Method    string
		args      map[string]map[string]*string
	}
	type args struct {
		envPrefix string
		args      []string
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
			op := &Operation{
				Operation: tt.fields.Operation,
				BaseURL:   tt.fields.BaseURL,
				Path:      tt.fields.Path,
				Method:    tt.fields.Method,
				args:      tt.fields.args,
			}
			if err := op.Parse(tt.args.envPrefix, tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("Operation.Parse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOperation_BuildRequest(t *testing.T) {
	type fields struct {
		Operation *openapi3.Operation
		BaseURL   string
		Path      string
		Method    string
		args      map[string]map[string]*string
	}
	type args struct {
		ctx     context.Context
		baseURI string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *http.Request
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op := &Operation{
				Operation: tt.fields.Operation,
				BaseURL:   tt.fields.BaseURL,
				Path:      tt.fields.Path,
				Method:    tt.fields.Method,
				args:      tt.fields.args,
			}
			got, err := op.BuildRequest(tt.args.ctx, tt.args.baseURI)
			if (err != nil) != tt.wantErr {
				t.Errorf("Operation.BuildRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Operation.BuildRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateOperationFromPath(t *testing.T) {
	type args struct {
		baseURL string
		path    string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]*Operation
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateOperationFromPath(tt.args.baseURL, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateOperationFromPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateOperationFromPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateOperationFromData(t *testing.T) {
	type args struct {
		baseURL string
		data    []byte
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]*Operation
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateOperationFromData(tt.args.baseURL, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateOperationFromData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateOperationFromData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_generateOperation(t *testing.T) {
	type args struct {
		baseURL string
		doc3    *openapi3.T
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]*Operation
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateOperation(tt.args.baseURL, tt.args.doc3)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateOperation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("generateOperation() = %v, want %v", got, tt.want)
			}
		})
	}
}
