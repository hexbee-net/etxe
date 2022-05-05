package etx

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func TestNumber_Parsing(t *testing.T) {
	type args struct {
		Input string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    *big.Float
	}{
		{
			name: "Int - Implicit Positive",
			args: args{
				Input: `1234`,
			},
			wantErr: false,
			want:    big.NewFloat(1234),
		},
		{
			name: "Int - Implicit Positive with Underscores",
			args: args{
				Input: `12_34`,
			},
			wantErr: false,
			want:    big.NewFloat(1234),
		},
		{
			name: "Int - Explicit Positive",
			args: args{
				Input: `+1234`,
			},
			wantErr: false,
			want:    big.NewFloat(1234),
		},
		{
			name: "Int - Negative",
			args: args{
				Input: `-1234`,
			},
			wantErr: false,
			want:    big.NewFloat(-1234),
		},

		{
			name: "Float - Implicit Positive",
			args: args{
				Input: `1234.56`,
			},
			wantErr: false,
			want:    big.NewFloat(1234.56),
		},
		{
			name: "Float - Implicit Positive with Underscores",
			args: args{
				Input: `12_34.5_6`,
			},
			wantErr: false,
			want:    big.NewFloat(1234.56),
		},
		{
			name: "Float - Explicit Positive",
			args: args{
				Input: `+1234.56`,
			},
			wantErr: false,
			want:    big.NewFloat(1234.56),
		},
		{
			name: "Float - Negative",
			args: args{
				Input: `-1234.56`,
			},
			wantErr: false,
			want:    big.NewFloat(-1234.56),
		},

		{
			name: "Float - Implicit Positive - Empty Integer",
			args: args{
				Input: `.56`,
			},
			wantErr: false,
			want:    big.NewFloat(.56),
		},
		{
			name: "Float - Explicit Positive - Empty Integer",
			args: args{
				Input: `+.56`,
			},
			wantErr: false,
			want:    big.NewFloat(.56),
		},
		{
			name: "Float - Negative - Empty Integer",
			args: args{
				Input: `-.56`,
			},
			wantErr: false,
			want:    big.NewFloat(-.56),
		},

		{
			name: "Float - Implicit Positive - Empty Fractional",
			args: args{
				Input: `1234.`,
			},
			wantErr: true,
			want:    nil,
		},
		{
			name: "Float - Explicit Positive - Empty Fractional",
			args: args{
				Input: `+1234.`,
			},
			wantErr: true,
			want:    nil,
		},
		{
			name: "Float - Negative - Empty Fractional",
			args: args{
				Input: `-1234.`,
			},
			wantErr: true,
			want:    nil,
		},

		{
			name: "Int - Exponent - Implicit Positive",
			args: args{
				Input: `1234e2`,
			},
			wantErr: false,
			want:    big.NewFloat(1234e2),
		},
		{
			name: "Int - Exponent - Implicit Positive with Underscores",
			args: args{
				Input: `12_34e1_2`,
			},
			wantErr: false,
			want:    big.NewFloat(1234e12),
		},
		{
			name: "Int - Exponent - Explicit Positive",
			args: args{
				Input: `1234e+2`,
			},
			wantErr: false,
			want:    big.NewFloat(1234e2),
		},
		{
			name: "Int - Exponent - Negative",
			args: args{
				Input: `1234e-2`,
			},
			wantErr: false,
			want:    big.NewFloat(1234e-2),
		},

		{
			name: "Float - Exponent - Implicit Positive",
			args: args{
				Input: `1234.56e2`,
			},
			wantErr: false,
			want:    big.NewFloat(1234.56e2),
		},
		{
			name: "Float - Exponent - Explicit Positive",
			args: args{
				Input: `1234.56e+2`,
			},
			wantErr: false,
			want:    big.NewFloat(1234.56e2),
		},
		{
			name: "Float - Exponent - Negative",
			args: args{
				Input: `1234.56e-2`,
			},
			wantErr: false,
			want:    big.NewFloat(1234.56e-2),
		},

		{
			name: "Float - Implicit Positive - Empty integer and fractional",
			args: args{
				Input: `.`,
			},
			wantErr: true,
			want:    nil,
		},
		{
			name: "Float - Explicit Positive - Empty integer and fractional",
			args: args{
				Input: `+.`,
			},
			wantErr: true,
			want:    nil,
		},
		{
			name: "Float - Negative - Empty integer and fractional",
			args: args{
				Input: `-.`,
			},
			wantErr: true,
			want:    nil,
		},

		{
			name: "Hex - Implicit Positive",
			args: args{
				Input: `0x1234`,
			},
			wantErr: false,
			want:    big.NewFloat(0x1234),
		},
		{
			name: "Hex - Implicit Positive - Capital X",
			args: args{
				Input: `0X1234`,
			},
			wantErr: false,
			want:    big.NewFloat(0x1234),
		},
		{
			name: "Hex - Implicit Positive with Underscores",
			args: args{
				Input: `0x12_34`,
			},
			wantErr: false,
			want:    big.NewFloat(0x1234),
		},
		{
			name: "Hex - Explicit Positive",
			args: args{
				Input: `+0x1234`,
			},
			wantErr: false,
			want:    big.NewFloat(0x1234),
		},
		{
			name: "Hex -  Negative",
			args: args{
				Input: `-0x1234`,
			},
			wantErr: false,
			want:    big.NewFloat(-0x1234),
		},

		{
			name: "Bin - Implicit Positive",
			args: args{
				Input: `0b1010`,
			},
			wantErr: false,
			want:    big.NewFloat(0b1010),
		},
		{
			name: "Bin - Implicit Positive - Capital B",
			args: args{
				Input: `0B1010`,
			},
			wantErr: false,
			want:    big.NewFloat(0b1010),
		},
		{
			name: "Bin - Implicit Positive with Underscores",
			args: args{
				Input: `0b10_10`,
			},
			wantErr: false,
			want:    big.NewFloat(0b1010),
		},
		{
			name: "Bin - Explicit Positive",
			args: args{
				Input: `+0b1010`,
			},
			wantErr: false,
			want:    big.NewFloat(0b1010),
		},
		{
			name: "Bin - Negative",
			args: args{
				Input: `-0b1010`,
			},
			wantErr: false,
			want:    big.NewFloat(-0b1010),
		},

		{
			name: "Oct - Implicit Positive",
			args: args{
				Input: `0o1234`,
			},
			wantErr: false,
			want:    big.NewFloat(0o1234),
		},
		{
			name: "Oct - Implicit Positive - Capital O",
			args: args{
				Input: `0O1234`,
			},
			wantErr: false,
			want:    big.NewFloat(0o1234),
		},
		{
			name: "Oct - Implicit Positive with Underscores",
			args: args{
				Input: `0o12_34`,
			},
			wantErr: false,
			want:    big.NewFloat(0o1234),
		},
		{
			name: "Oct - Explicit Positive",
			args: args{
				Input: `+0o1234`,
			},
			wantErr: false,
			want:    big.NewFloat(0o1234),
		},
		{
			name: "Oct - Negative",
			args: args{
				Input: `-0o1234`,
			},
			wantErr: false,
			want:    big.NewFloat(-0o1234),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			type Num struct {
				Number *Number `parser:"@Number" json:"number,omitempty"`
			}
			parser := participle.MustBuild(&Num{}, participle.Lexer(lexer.MustStateful(lexRules())))

			res := &Num{}
			err := parser.ParseString("", tt.args.Input, res)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			if tt.want != nil {
				require.NotNil(t, res.Number)
				require.NotNil(t, res.Number.Float)
				assert.Zero(t, tt.want.Cmp(res.Number.Float))
			} else {
				assert.Nil(t, res.Number)
			}
		})
	}
}
