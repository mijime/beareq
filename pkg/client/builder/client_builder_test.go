package builder

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	"github.com/mijime/beareq/pkg/beareq"
	"golang.org/x/oauth2"
)

func Test_osGetEnv(t *testing.T) {
	type args struct {
		key          string
		defaultValue string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := osGetEnv(tt.args.key, tt.args.defaultValue); got != tt.want {
				t.Errorf("osGetEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewClientBuilder(t *testing.T) {
	tests := []struct {
		name string
		want *ClientBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewClientBuilder(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClientBuilder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClientBuilder_fetchConfigByProfile(t *testing.T) {
	type fields struct {
		Profile      string
		ProfilesPath string
		TokenDir     string
	}
	tests := []struct {
		name    string
		fields  fields
		want    profileConfig
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &ClientBuilder{
				Profile:      tt.fields.Profile,
				ProfilesPath: tt.fields.ProfilesPath,
				TokenDir:     tt.fields.TokenDir,
			}
			got, err := b.fetchConfigByProfile()
			if (err != nil) != tt.wantErr {
				t.Errorf("ClientBuilder.fetchConfigByProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ClientBuilder.fetchConfigByProfile() = %v, want %v", got, tt.want)
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
		t.Run(tt.name, func(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
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

func TestClientBuilder_saveToken(t *testing.T) {
	type fields struct {
		Profile      string
		ProfilesPath string
		TokenDir     string
	}
	type args struct {
		tokSrc oauth2.TokenSource
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
			b := &ClientBuilder{
				Profile:      tt.fields.Profile,
				ProfilesPath: tt.fields.ProfilesPath,
				TokenDir:     tt.fields.TokenDir,
			}
			if err := b.saveToken(tt.args.tokSrc); (err != nil) != tt.wantErr {
				t.Errorf("ClientBuilder.saveToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClientBuilder_fetchToken(t *testing.T) {
	type fields struct {
		Profile      string
		ProfilesPath string
		TokenDir     string
	}
	type args struct {
		config *oauth2.Config
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *oauth2.Token
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &ClientBuilder{
				Profile:      tt.fields.Profile,
				ProfilesPath: tt.fields.ProfilesPath,
				TokenDir:     tt.fields.TokenDir,
			}
			got, err := b.fetchToken(tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClientBuilder.fetchToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ClientBuilder.fetchToken() = %v, want %v", got, tt.want)
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
		t.Run(tt.name, func(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
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

func Test_fetchCode(t *testing.T) {
	type args struct {
		code   chan<- string
		config *oauth2.Config
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := fetchCode(tt.args.code, tt.args.config); (err != nil) != tt.wantErr {
				t.Errorf("fetchCode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_generateToken(t *testing.T) {
	type args struct {
		config *oauth2.Config
	}
	tests := []struct {
		name    string
		args    args
		want    *oauth2.Token
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateToken(tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("generateToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
