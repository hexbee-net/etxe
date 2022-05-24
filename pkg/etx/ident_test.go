package etx

import (
	"testing"
)

func TestIdent_Parsing(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr bool
		want    *Ident
	}{
		{
			name:    "one part",
			input:   "foo",
			wantErr: false,
			want: &Ident{
				Parts: []string{
					"foo",
				},
			},
		},
		{
			name:    "several parts",
			input:   "foo.bar.baz",
			wantErr: false,
			want: &Ident{
				Parts: []string{
					"foo",
					"bar",
					"baz",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testParser(t, tt.input, tt.want, tt.wantErr, true)
		})
	}

}

func TestIdent_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input Ident
		want  string
	}{
		{
			name: "One part",
			input: Ident{
				Parts: []string{
					"foo",
				},
			},
			want: "foo",
		},
		{
			name: "Several parts",
			input: Ident{
				Parts: []string{
					"foo",
					"bar",
					"baz",
				},
			},
			want: "foo.bar.baz",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStringer(t, false, tt.want, tt.input)
		})
	}
}
