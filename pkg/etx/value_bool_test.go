package etx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBool_Clone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *ValueBool
		want  *ValueBool
	}{
		{
			name:  "Nil",
			input: nil,
			want:  nil,
		},
		{
			name:  "True",
			input: &ValueBool{Value: true},
			want:  &ValueBool{Value: true},
		},
		{
			name:  "False",
			input: &ValueBool{Value: false},
			want:  &ValueBool{Value: false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCloner[*ValueBool](t, tt.want, tt.input)
		})
	}
}

func TestBool_Children(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *ValueBool
		want  []Node
	}{
		{
			name:  "Nil",
			input: nil,
			want:  nil,
		},
		{
			name:  "True",
			input: &ValueBool{Value: true},
			want:  nil,
		},
		{
			name:  "False",
			input: &ValueBool{Value: false},
			want:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.input.Children())
		})
	}
}
