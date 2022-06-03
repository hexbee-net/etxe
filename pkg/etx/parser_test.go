package etx

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		reader  io.Reader
		want    *AST
		wantErr bool
	}{
		{
			name:    "Read fail",
			reader:  iotest.ErrReader(errors.New("io error")),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Empty file",
			reader:  strings.NewReader(``),
			want:    &AST{},
			wantErr: false,
		},
		{
			name:   "Valid file",
			reader: strings.NewReader(`foo`),
			want: &AST{
				Items: []*RootItem{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Attribute: &Attribute{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Key:     "foo",
						},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := Parse(tt.reader)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.Equal(t, tt.want, res)
		})
	}
}

func TestParseString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *AST
		wantErr bool
	}{
		{
			name:  "Empty string",
			input: "",
			want:  &AST{},
		},
		{
			name:  "Valid string",
			input: "foo",
			want: &AST{
				Items: []*RootItem{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Attribute: &Attribute{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Key:     "foo",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:    "Invalid string",
			input:   "><",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := ParseString(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.Equal(t, tt.want, res)
		})
	}
}

func TestParseBytes(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    *AST
		wantErr bool
	}{
		{
			name:  "Empty string",
			input: []byte(""),
			want:  &AST{},
		},
		{
			name:  "Valid string",
			input: []byte("foo"),
			want: &AST{
				Items: []*RootItem{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Attribute: &Attribute{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Key:     "foo",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:    "Invalid string",
			input:   []byte("><"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := ParseBytes(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.Equal(t, tt.want, res)
		})
	}
}
