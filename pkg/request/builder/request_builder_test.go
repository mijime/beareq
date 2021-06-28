package builder

import (
	"context"
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/mijime/beareq/pkg/request/attrs"
)

func TestNewRequestBuilder(t *testing.T) {
	tests := []struct {
		name string
		want RequestBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRequestBuilder(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRequestBuilder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRequestBuilder_body(t *testing.T) {
	type fields struct {
		Method     attrs.HTTPRequestMethod
		Data       attrs.HTTPBody
		Header     attrs.HTTPHeader
		JSONObject attrs.JSONObject
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
		t.Run(tt.name, func(t *testing.T) {
			b := RequestBuilder{
				Method:     tt.fields.Method,
				Data:       tt.fields.Data,
				Header:     tt.fields.Header,
				JSONObject: tt.fields.JSONObject,
			}
			got, err := b.body()
			if (err != nil) != tt.wantErr {
				t.Errorf("RequestBuilder.body() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RequestBuilder.body() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRequestBuilder_BuildRequest(t *testing.T) {
	type fields struct {
		Method     attrs.HTTPRequestMethod
		Data       attrs.HTTPBody
		Header     attrs.HTTPHeader
		JSONObject attrs.JSONObject
	}
	type args struct {
		ctx context.Context
		url string
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
			b := RequestBuilder{
				Method:     tt.fields.Method,
				Data:       tt.fields.Data,
				Header:     tt.fields.Header,
				JSONObject: tt.fields.JSONObject,
			}
			got, err := b.BuildRequest(tt.args.ctx, tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("RequestBuilder.BuildRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RequestBuilder.BuildRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
