package openapi

import (
	"context"
	"flag"
	"net/http"
	"net/http/httputil"
	"reflect"
	"strings"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
)

func TestOperation_BuildRequest(t *testing.T) {
	type fields struct {
		Operation *openapi3.Operation
		PathItem  *openapi3.PathItem
		BaseURL   string
		Path      string
		Method    string
		args      map[string]map[string]flag.Value
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
		{
			fields: fields{
				Method:   http.MethodPost,
				PathItem: &openapi3.PathItem{},
				Operation: &openapi3.Operation{
					Parameters: openapi3.Parameters{
						{
							Value: &openapi3.Parameter{
								In: "query", Name: "page",
								Schema: &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "integer"}},
							},
						},
						{
							Value: &openapi3.Parameter{
								In: "formData", Name: "title",
								Schema: &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "string"}},
							},
						},
					},
				},
				args: map[string]map[string]flag.Value{
					"formData": {"title": &argString{V: "hello", IsSet: true}},
				},
			},
			args: args{
				ctx:     context.TODO(),
				baseURI: "http://localhost:3000",
			},
			want: func() *http.Request {
				req, _ := http.NewRequest(http.MethodPost, "http://localhost:3000", strings.NewReader("title=hello"))
				req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
				return req
			}(),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			op := &Operation{
				PathItem:  tt.fields.PathItem,
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
			gotDump, _ := httputil.DumpRequest(got, true)
			wantDump, _ := httputil.DumpRequest(tt.want, true)
			if !reflect.DeepEqual(gotDump, wantDump) {
				t.Errorf("Operation.BuildRequest() = %s, want %s", gotDump, wantDump)
			}
		})
	}
}

func TestNewOperation(t *testing.T) {
	type args struct {
		url    string
		path   string
		method string
		op     *openapi3.Operation
		pi     *openapi3.PathItem
	}
	tests := []struct {
		name string
		args args
		want *Operation
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := NewOperation(tt.args.url, tt.args.path, tt.args.method, tt.args.op, tt.args.pi); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewOperation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperation_Name(t *testing.T) {
	type fields struct {
		Operation *openapi3.Operation
		PathItem  *openapi3.PathItem
		BaseURL   string
		Path      string
		Method    string
		args      map[string]map[string]flag.Value
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
			op := &Operation{
				Operation: tt.fields.Operation,
				PathItem:  tt.fields.PathItem,
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

func TestOperation_FlagSet(t *testing.T) {
	type fields struct {
		Operation *openapi3.Operation
		PathItem  *openapi3.PathItem
		BaseURL   string
		Path      string
		Method    string
		args      map[string]map[string]flag.Value
	}
	type args struct {
		envPrefix string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *flag.FlagSet
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			op := &Operation{
				Operation: tt.fields.Operation,
				PathItem:  tt.fields.PathItem,
				BaseURL:   tt.fields.BaseURL,
				Path:      tt.fields.Path,
				Method:    tt.fields.Method,
				args:      tt.fields.args,
			}
			if got := op.FlagSet(tt.args.envPrefix); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Operation.FlagSet() = %v, want %v", got, tt.want)
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
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

func Test_argInteger_String(t *testing.T) {
	type fields struct {
		V     int
		IsSet bool
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
			a := &argInteger{
				V:     tt.fields.V,
				IsSet: tt.fields.IsSet,
			}
			if got := a.String(); got != tt.want {
				t.Errorf("argInteger.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_argInteger_Set(t *testing.T) {
	type fields struct {
		V     int
		IsSet bool
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
			a := &argInteger{
				V:     tt.fields.V,
				IsSet: tt.fields.IsSet,
			}
			if err := a.Set(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("argInteger.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_argBoolean_String(t *testing.T) {
	type fields struct {
		V     bool
		IsSet bool
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
			a := &argBoolean{
				V:     tt.fields.V,
				IsSet: tt.fields.IsSet,
			}
			if got := a.String(); got != tt.want {
				t.Errorf("argBoolean.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_argBoolean_Set(t *testing.T) {
	type fields struct {
		V     bool
		IsSet bool
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
			a := &argBoolean{
				V:     tt.fields.V,
				IsSet: tt.fields.IsSet,
			}
			if err := a.Set(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("argBoolean.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_argFloat_String(t *testing.T) {
	type fields struct {
		V     float64
		IsSet bool
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
			a := &argFloat{
				V:     tt.fields.V,
				IsSet: tt.fields.IsSet,
			}
			if got := a.String(); got != tt.want {
				t.Errorf("argFloat.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_argFloat_Set(t *testing.T) {
	type fields struct {
		V     float64
		IsSet bool
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
			a := &argFloat{
				V:     tt.fields.V,
				IsSet: tt.fields.IsSet,
			}
			if err := a.Set(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("argFloat.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_argString_String(t *testing.T) {
	type fields struct {
		V     string
		IsSet bool
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
			a := &argString{
				V:     tt.fields.V,
				IsSet: tt.fields.IsSet,
			}
			if got := a.String(); got != tt.want {
				t.Errorf("argString.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_argString_Set(t *testing.T) {
	type fields struct {
		V     string
		IsSet bool
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
			a := &argString{
				V:     tt.fields.V,
				IsSet: tt.fields.IsSet,
			}
			if err := a.Set(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("argString.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_argArray_String(t *testing.T) {
	type fields struct {
		V     []string
		IsSet bool
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
			a := &argArray{
				V:     tt.fields.V,
				IsSet: tt.fields.IsSet,
			}
			if got := a.String(); got != tt.want {
				t.Errorf("argArray.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_argArray_Set(t *testing.T) {
	type fields struct {
		V     []string
		IsSet bool
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
			a := &argArray{
				V:     tt.fields.V,
				IsSet: tt.fields.IsSet,
			}
			if err := a.Set(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("argArray.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_argObject_String(t *testing.T) {
	type fields struct {
		V map[string]flag.Value
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
			a := &argObject{
				V: tt.fields.V,
			}
			if got := a.String(); got != tt.want {
				t.Errorf("argObject.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_argObject_Set(t *testing.T) {
	type fields struct {
		V map[string]flag.Value
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
			a := &argObject{
				V: tt.fields.V,
			}
			if err := a.Set(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("argObject.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
