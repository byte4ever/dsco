package hit

import (
	"reflect"
	"testing"
)

func TestNewIntNode(t *testing.T) {
	type args[T hash.Hash] struct {
		id           string
		value        int
		hashProvider HashProvider
	}
	type testCase[T hash.Hash] struct {
		name string
		args args[T]
		want *IntNode
	}
	tests := []testCase[ /* TODO: Insert concrete types here */ ]{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewIntNode(tt.args.id, tt.args.value, tt.args.hashProvider); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewIntNode() = %v, want %v", got, tt.want)
			}
		})
	}
}