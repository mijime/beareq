package handler

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	"github.com/mijime/beareq/v2/pkg/response/attrs"
)

func TestNewResponseHandler(t *testing.T) {
	tests := []struct {
		name string
		want ResponseHandler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := NewResponseHandler(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewResponseHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResponseHandler_HandleResponse(t *testing.T) {
	type fields struct {
		JSONQuery attrs.JSONQuery
		Verbose   bool
		Fail      bool
	}
	type args struct {
		ctx  context.Context
		resp *http.Response
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
			h := ResponseHandler{
				JSONQuery: tt.fields.JSONQuery,
				Verbose:   tt.fields.Verbose,
				Fail:      tt.fields.Fail,
			}
			if err := h.HandleResponse(tt.args.ctx, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("ResponseHandler.HandleResponse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
