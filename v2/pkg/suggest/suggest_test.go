package suggest

import (
	"reflect"
	"testing"
)

func TestSuggest(t *testing.T) {
	type args struct {
		words  []string
		target string
		size   int
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := Suggest(tt.args.words, tt.args.target, tt.args.size); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Suggest() = %v, want %v", got, tt.want)
			}
		})
	}
}
