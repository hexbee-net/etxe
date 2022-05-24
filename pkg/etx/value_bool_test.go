package etx

import (
	"testing"
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
			input: testValPtr[ValueBool](t, true),
			want:  testValPtr[ValueBool](t, true),
		},
		{
			name:  "False",
			input: testValPtr[ValueBool](t, false),
			want:  testValPtr[ValueBool](t, false),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCloner[*ValueBool](t, tt.want, tt.input)
		})
	}
}
