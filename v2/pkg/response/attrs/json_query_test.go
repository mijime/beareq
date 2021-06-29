package attrs

import (
	"reflect"
	"testing"

	"github.com/itchyny/gojq"
)

func TestNewJSONQuery(t *testing.T) {
	tests := []struct {
		name string
		want JSONQuery
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := NewJSONQuery(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewJSONQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJSONQuery_Exists(t *testing.T) {
	type fields struct {
		Query *gojq.Query
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			q := &JSONQuery{
				Query: tt.fields.Query,
			}
			if got := q.Exists(); got != tt.want {
				t.Errorf("JSONQuery.Exists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJSONQuery_String(t *testing.T) {
	type fields struct {
		Query *gojq.Query
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
			q := &JSONQuery{
				Query: tt.fields.Query,
			}
			if got := q.String(); got != tt.want {
				t.Errorf("JSONQuery.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJSONQuery_Set(t *testing.T) {
	type fields struct {
		Query *gojq.Query
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
			q := &JSONQuery{
				Query: tt.fields.Query,
			}
			if err := q.Set(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("JSONQuery.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
