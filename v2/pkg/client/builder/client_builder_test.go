package builder

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	"github.com/mijime/beareq/v2/pkg/beareq"
	"golang.org/x/oauth2"
)

func TestNewClientBuilder(t *testing.T) {
	tests := []struct {
		name string
		want *ClientBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := NewClientBuilder(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClientBuilder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_headerOverride_RoundTrip(t *testing.T) {
	type args struct {
		req *http.Request
	}
	tests := []struct {
		name    string
		r       headerOverride
		args    args
		want    *http.Response
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := tt.r.RoundTrip(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("headerOverride.RoundTrip() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("headerOverride.RoundTrip() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClientBuilder_BuildClient(t *testing.T) {
	type fields struct {
		Profile      string
		ProfilesPath string
		TokenDir     string
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    beareq.Client
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			b := &ClientBuilder{
				Profile:      tt.fields.Profile,
				ProfilesPath: tt.fields.ProfilesPath,
				TokenDir:     tt.fields.TokenDir,
			}
			got, err := b.BuildClient(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClientBuilder.BuildClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ClientBuilder.BuildClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_profileClient_Do(t *testing.T) {
	type fields struct {
		client  *http.Client
		builder *ClientBuilder
		source  oauth2.TokenSource
	}
	type args struct {
		req *http.Request
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *http.Response
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := &profileClient{
				client:  tt.fields.client,
				builder: tt.fields.builder,
				source:  tt.fields.source,
			}
			got, err := c.Do(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("profileClient.Do() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("profileClient.Do() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_profileClient_Close(t *testing.T) {
	type fields struct {
		client  *http.Client
		builder *ClientBuilder
		source  oauth2.TokenSource
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := &profileClient{
				client:  tt.fields.client,
				builder: tt.fields.builder,
				source:  tt.fields.source,
			}
			if err := c.Close(); (err != nil) != tt.wantErr {
				t.Errorf("profileClient.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
