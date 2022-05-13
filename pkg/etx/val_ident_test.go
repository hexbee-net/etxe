package etx

import (
	"testing"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIdent_Parsing(t *testing.T) {
	type args struct {
		Input string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    *Ident
	}{
		{
			name: "one part",
			args: args{
				Input: "foo",
			},
			wantErr: false,
			want: &Ident{
				Parts: []string{
					"foo",
				},
			},
		},
		{
			name: "several parts",
			args: args{
				Input: "foo.bar.baz",
			},
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
			parser := participle.MustBuild(&Ident{}, participle.Lexer(lexer.MustStateful(lexRules())))

			res := &Ident{}
			err := parser.ParseString("", tt.args.Input, res)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.want, res)
		})
	}

}

func TestIdent_String(t *testing.T) {
	type args struct {
		Input Ident
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "One part",
			args: args{
				Input: Ident{
					Parts: []string{
						"foo",
					},
				},
			},
			want: "foo",
		},
		{
			name: "Several parts",
			args: args{
				Input: Ident{
					Parts: []string{
						"foo",
						"bar",
						"baz",
					},
				},
			},
			want: "foo.bar.baz",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.args.Input.String())
		})
	}
}
