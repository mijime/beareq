package beareq

import (
	"context"
	"testing"
)

func TestRun(t *testing.T) {
	type args struct {
		ctx  context.Context
		cb   ClientBuilder
		rb   RequestBuilder
		rh   ResponseHandler
		urls []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := Run(tt.args.ctx, tt.args.cb, tt.args.rb, tt.args.rh, tt.args.urls...); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
