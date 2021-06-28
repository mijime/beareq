package main

import (
	"reflect"
	"testing"
)

func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}

func Test_fetchConfigByProfile(t *testing.T) {
	type args struct {
		b profileConfig
	}
	tests := []struct {
		name    string
		args    args
		want    openapiConfig
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fetchConfigByProfile(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("fetchConfigByProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fetchConfigByProfile() = %v, want %v", got, tt.want)
			}
		})
	}
}
