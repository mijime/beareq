package raw

import (
	"context"
	"io"
	"net/http"
	"testing"
)

func TestResponseHandler_HandleResponse(t *testing.T) {
	type fields struct {
		Writer io.Writer
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
		t.Run(tt.name, func(t *testing.T) {
			h := &ResponseHandler{
				Writer: tt.fields.Writer,
			}
			if err := h.HandleResponse(tt.args.ctx, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("ResponseHandler.HandleResponse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
