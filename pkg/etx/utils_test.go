package etx

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRefFields[T any](t *testing.T, testFunc func(assert.TestingT, interface{}, interface{}, ...interface{}) bool, expected, actual T) {
	t.Helper()

	et := reflect.ValueOf(expected)
	at := reflect.ValueOf(actual)

	if et.Kind() == reflect.Ptr {
		et = et.Elem()
	}

	if at.Kind() == reflect.Ptr {
		at = at.Elem()
	}

	if et.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < et.NumField(); i++ {
		field := et.Field(i)
		typeField := et.Type().Field(i)

		if field.Kind() == reflect.Ptr && !field.IsNil() {
			testFunc(t, field.Interface(), at.Field(i).Interface(), fmt.Sprintf("Field name: %v", typeField.Name))
		}
	}
}

// func testParser[T any](t *testing.T, input string, want *T, wantErr bool) {

func testParser[T any](t *testing.T, input string, want *T, wantErr, compareNodes bool) {
	t.Helper()

	var res T
	parser := participle.MustBuild(&res, participle.Lexer(lexer.MustStateful(lexRules())))

	err := parser.ParseString("", input, &res)

	if wantErr {
		require.Error(t, err)
	} else {
		require.NoError(t, err)
	}

	if !wantErr || (wantErr && want != nil) {
		numberComparer := cmp.Comparer(func(x, y *ValueNumber) bool {
			if x == nil && y == nil {
				return true
			}

			if x == nil || y == nil {
				return false
			}

			return x.Value.Cmp(y.Value) == 0
		})

		posComparer := cmp.Comparer(func(x, y ASTNode) bool {
			if !compareNodes {
				return true
			}

			return reflect.DeepEqual(x, y)
		})

		if !cmp.Equal(want, &res, numberComparer, posComparer) {
			assert.Fail(t, "Not equal -want +res", cmp.Diff(want, &res, numberComparer, posComparer))
		}
	}
}

func testStringer(t *testing.T, wantPanic bool, want string, input fmt.Stringer, msgAndArgs ...interface{}) {
	t.Helper()

	if wantPanic {
		if want != "" {
			assert.PanicsWithValue(t, want, func() {
				_ = input.String()
			}, msgAndArgs)
		} else {
			assert.Panics(t, func() {
				_ = input.String()
			}, msgAndArgs)
		}
	} else {
		assert.Equal(t, want, input.String())
	}
}

func testCloner[T any](t *testing.T, want, input Cloner[T]) {
	t.Helper()

	clone := input.Clone()

	if reflect.ValueOf(want).IsNil() {
		assert.Nil(t, clone)
		return
	}

	assert.Equal(t, want, clone)
	assert.NotSame(t, want, clone)

	testRefFields(t, assert.NotSame, want, input)
}

// /////////////////////////////////////

func TestIndent(t *testing.T) {
	type args struct {
		s      string
		prefix string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "One char indent",
			args: args{
				s: `
Not tricks, Michael, illusions.
I care deeply for nature.
Well, what do you expect, mother? Get me a vodka rocks.
And a piece of toast.
Marry me.`[1:],
				prefix: "X",
			},
			want: `
XNot tricks, Michael, illusions.
XI care deeply for nature.
XWell, what do you expect, mother? Get me a vodka rocks.
XAnd a piece of toast.
XMarry me.`[1:],
		},
		{
			name: "Multiple chars indent",
			args: args{
				s: `
Not tricks, Michael, illusions.
I care deeply for nature.
Well, what do you expect, mother? Get me a vodka rocks.
And a piece of toast.
Marry me.`[1:],
				prefix: "XXX",
			},
			want: `
XXXNot tricks, Michael, illusions.
XXXI care deeply for nature.
XXXWell, what do you expect, mother? Get me a vodka rocks.
XXXAnd a piece of toast.
XXXMarry me.`[1:],
		},
		{
			name: "Empty lines",
			args: args{
				s: `
Not tricks, Michael, illusions.
I care deeply for nature.

Well, what do you expect, mother? Get me a vodka rocks.

And a piece of toast.
Marry me.`[1:],
				prefix: "XXX",
			},
			want: `
XXXNot tricks, Michael, illusions.
XXXI care deeply for nature.

XXXWell, what do you expect, mother? Get me a vodka rocks.

XXXAnd a piece of toast.
XXXMarry me.`[1:],
		},
		{
			name: "Trailing empty lines",
			args: args{
				s: `
Not tricks, Michael, illusions.
I care deeply for nature.
Well, what do you expect, mother? Get me a vodka rocks.
And a piece of toast.
Marry me.


`[1:],
				prefix: "XXX",
			},
			want: `
XXXNot tricks, Michael, illusions.
XXXI care deeply for nature.
XXXWell, what do you expect, mother? Get me a vodka rocks.
XXXAnd a piece of toast.
XXXMarry me.


`[1:],
		},
		{
			name: "Preceding empty lines",
			args: args{
				s: `


Not tricks, Michael, illusions.
I care deeply for nature.
Well, what do you expect, mother? Get me a vodka rocks.
And a piece of toast.
Marry me.`[1:],
				prefix: "XXX",
			},
			want: `


XXXNot tricks, Michael, illusions.
XXXI care deeply for nature.
XXXWell, what do you expect, mother? Get me a vodka rocks.
XXXAnd a piece of toast.
XXXMarry me.`[1:],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, indent(tt.args.s, tt.args.prefix), "indent(%v, %v)", tt.args.s, tt.args.prefix)
		})
	}
}
