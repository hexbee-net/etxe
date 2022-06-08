package etx

import (
	"math/big"
	"testing"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNumber_Parsing(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr bool
		want    *ValueNumber
	}{
		{
			name:    "Int - Implicit Positive",
			input:   `1234`,
			wantErr: false,
			want:    &ValueNumber{Value: big.NewFloat(1234), Source: `1234`},
		},
		{
			name:    "Int - Implicit Positive with Underscores",
			input:   `12_34`,
			wantErr: false,
			want:    &ValueNumber{Value: big.NewFloat(1234), Source: `12_34`},
		},

		{
			name:    "Float - Implicit Positive",
			input:   `1234.56`,
			wantErr: false,
			want:    &ValueNumber{Value: big.NewFloat(1234.56), Source: `1234.56`},
		},
		{
			name:    "Float - Implicit Positive with Underscores",
			input:   `12_34.5_6`,
			wantErr: false,
			want:    &ValueNumber{Value: big.NewFloat(1234.56), Source: `12_34.5_6`},
		},

		{
			name:    "Float - Implicit Positive - Empty Integer",
			input:   `.56`,
			wantErr: false,
			want:    &ValueNumber{Value: big.NewFloat(.56), Source: `.56`},
		},

		{
			name:    "Float - Implicit Positive - Empty Fractional",
			input:   `1234.`,
			wantErr: true,
			want:    &ValueNumber{Value: nil, Source: `1234.`},
		},

		{
			name:    "Int - Exponent - Implicit Positive",
			input:   `1234e2`,
			wantErr: false,
			want:    &ValueNumber{Value: big.NewFloat(1234e2), Source: `1234e2`},
		},
		{
			name:    "Int - Exponent - Implicit Positive with Underscores",
			input:   `12_34e1_2`,
			wantErr: false,
			want:    &ValueNumber{Value: big.NewFloat(1234e12), Source: `12_34e1_2`},
		},

		{
			name:    "Float - Exponent - Implicit Positive",
			input:   `1234.56e2`,
			wantErr: false,
			want:    &ValueNumber{Value: big.NewFloat(1234.56e2), Source: `1234.56e2`},
		},

		{
			name:    "Float - Implicit Positive - Empty integer and fractional",
			input:   `.`,
			wantErr: true,
			want:    &ValueNumber{Value: nil, Source: `.`},
		},

		{
			name:    "Hex - Implicit Positive",
			input:   `0x1234`,
			wantErr: false,
			want:    &ValueNumber{Value: big.NewFloat(0x1234), Source: `0x1234`},
		},
		{
			name:    "Hex - Implicit Positive - Capital X",
			input:   `0X1234`,
			wantErr: false,
			want:    &ValueNumber{Value: big.NewFloat(0x1234), Source: `0X1234`},
		},
		{
			name:    "Hex - Implicit Positive with Underscores",
			input:   `0x12_34`,
			wantErr: false,
			want:    &ValueNumber{Value: big.NewFloat(0x1234), Source: `0x12_34`},
		},

		{
			name:    "Bin - Implicit Positive",
			input:   `0b1010`,
			wantErr: false,
			want:    &ValueNumber{Value: big.NewFloat(0b1010), Source: `0b1010`},
		},
		{
			name:    "Bin - Implicit Positive - Capital B",
			input:   `0B1010`,
			wantErr: false,
			want:    &ValueNumber{Value: big.NewFloat(0b1010), Source: `0B1010`},
		},
		{
			name:    "Bin - Implicit Positive with Underscores",
			input:   `0b10_10`,
			wantErr: false,
			want:    &ValueNumber{Value: big.NewFloat(0b1010), Source: `0b10_10`},
		},

		{
			name:    "Oct - Implicit Positive",
			input:   `0o1234`,
			wantErr: false,
			want:    &ValueNumber{Value: big.NewFloat(0o1234), Source: `0o1234`},
		},
		{
			name:    "Oct - Implicit Positive - Capital O",
			input:   `0O1234`,
			wantErr: false,
			want:    &ValueNumber{Value: big.NewFloat(0o1234), Source: `0O1234`},
		},
		{
			name:    "Oct - Implicit Positive with Underscores",
			input:   `0o12_34`,
			wantErr: false,
			want:    &ValueNumber{Value: big.NewFloat(0o1234), Source: `0o12_34`},
		},

		{
			name:    "Oct - 0-prefix - Implicit Positive",
			input:   `01234`,
			wantErr: false,
			want:    &ValueNumber{Value: big.NewFloat(0o1234), Source: `0o1234`},
		},
		{
			name:    "Oct - 0-prefix - Implicit Positive with Underscores",
			input:   `012_34`,
			wantErr: false,
			want:    &ValueNumber{Value: big.NewFloat(0o1234), Source: `0o12_34`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			type Num struct {
				Number *ValueNumber `parser:"@Number" json:"number,omitempty"`
			}

			parser := participle.MustBuild(&Num{}, participle.Lexer(lexer.MustStateful(lexRules())))

			res := &Num{}
			err := parser.ParseString("", tt.input, res)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				if tt.want != nil {
					require.NotNil(t, res.Number)
					require.NotNil(t, res.Number.Value)
					assert.Equal(t, tt.want.Source, res.Number.Source)

					if eq := tt.want.Value.Cmp(res.Number.Value); eq != 0 {
						if eq < 0 {
							assert.Failf(t, "want < res", "want: %s\nres: %s", tt.want.Value, res.Number.Value)
						} else {
							assert.Failf(t, "want(%s) > res(%s)", tt.want.Value.String(), res.Number.Value)
						}
					}
				} else {
					assert.Nil(t, res.Number)
				}
			}
		})
	}
}

func TestNumber_Clone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *ValueNumber
		want  *ValueNumber
	}{
		{
			name:  "Nil",
			input: nil,
			want:  nil,
		},
		{
			name:  "Empty",
			input: &ValueNumber{},
			want:  &ValueNumber{},
		},
		{
			name: "Value",
			input: &ValueNumber{
				Value: big.NewFloat(1),
			},
			want: &ValueNumber{
				Value: big.NewFloat(1),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCloner[*ValueNumber](t, tt.want, tt.input)
		})
	}
}

func TestNumber_Children(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *ValueNumber
		want  []Node
	}{
		{
			name:  "Nil",
			input: nil,
			want:  nil,
		},
		{
			name: "Value",
			input: &ValueNumber{
				Value: big.NewFloat(1),
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.input.Children())
		})
	}
}

func TestNumber_FormattedString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     *ValueNumber
		wantPanic bool
		want      string
	}{
		{
			name:      "Nil",
			input:     nil,
			wantPanic: true,
		},
		{
			name:      "Empty",
			input:     &ValueNumber{},
			wantPanic: true,
		},
		{
			name:  "Value and string match",
			input: &ValueNumber{Value: big.NewFloat(0), Source: "0"},
			want:  "0",
		},
		{
			name:  "Value and string mismatch",
			input: &ValueNumber{Value: big.NewFloat(0), Source: "1"},
			want:  "1",
		},
		{
			name:  "Empty string ",
			input: &ValueNumber{Value: big.NewFloat(0), Source: ""},
			want:  "0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStringer(t, tt.wantPanic, tt.want, tt.input)
		})
	}
}

func TestNumber_Capture(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     []string
		wantErr   bool
		wantPanic bool
		want      *ValueNumber
	}{
		{
			name:      "No values",
			input:     []string{},
			wantErr:   false,
			wantPanic: true,
		},
		{
			name:      "Regular number",
			input:     []string{"1"},
			wantErr:   false,
			wantPanic: false,
			want: &ValueNumber{
				Value:  big.NewFloat(1),
				Source: "1",
			},
		},
		{
			name:      "Non-prefixed octal number",
			input:     []string{"0123"},
			wantErr:   false,
			wantPanic: false,
			want: &ValueNumber{
				Value:  big.NewFloat(0o123),
				Source: "0o123",
			},
		},
		{
			name:      "Invalid number",
			input:     []string{"foo"},
			wantErr:   true,
			wantPanic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := &ValueNumber{}

			if tt.wantPanic {
				assert.Panics(t, func() {
					_ = res.Capture(tt.input)
				})
				return
			}

			err := res.Capture(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.Equal(t, tt.want, res)
		})
	}

}
