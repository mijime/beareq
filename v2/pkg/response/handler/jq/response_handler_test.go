package jq

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/itchyny/gojq"
)

func TestResponseHandler_HandleResponse(t *testing.T) {
	type fields struct {
		Writer io.Writer
		Query  *gojq.Query
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
			h := &ResponseHandler{
				Writer: tt.fields.Writer,
				Query:  tt.fields.Query,
			}
			if err := h.HandleResponse(tt.args.ctx, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("ResponseHandler.HandleResponse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
