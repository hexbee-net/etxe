package dag

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSet_Add1(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		source Set[string]
		want   Set[string]
	}{
		{
			name:   "Empty set",
			input:  "1",
			source: Set[string]{},
			want:   Set[string]{"1": "1"},
		},
		{
			name:   "Not Existing",
			input:  "2",
			source: Set[string]{"1": "1"},
			want: Set[string]{
				"1": "1",
				"2": "2",
			},
		},
		{
			name:   "Existing",
			input:  "1",
			source: Set[string]{"1": "1"},
			want:   Set[string]{"1": "1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.source.Add(tt.input)
			assert.Equal(t, tt.want, tt.source)
		})
	}
}

func TestSet_Delete(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		source Set[string]
		want   Set[string]
	}{
		{
			name:   "Empty set",
			input:  "1",
			source: Set[string]{},
			want:   Set[string]{},
		},
		{
			name:   "Not Existing",
			input:  "2",
			source: Set[string]{"1": "1"},
			want: Set[string]{
				"1": "1",
			},
		},
		{
			name:   "Existing",
			input:  "1",
			source: Set[string]{"1": "1"},
			want:   Set[string]{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.source.Delete(tt.input)
			assert.Equal(t, tt.want, tt.source)
		})
	}
}

func TestSet_Includes(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		source Set[string]
		want   bool
	}{
		{
			name:   "Empty set",
			input:  "1",
			source: Set[string]{},
			want:   false,
		},
		{
			name:   "Not Existing",
			input:  "2",
			source: Set[string]{"1": "1"},
			want:   false,
		},
		{
			name:   "Existing",
			input:  "1",
			source: Set[string]{"1": "1"},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := tt.source.Includes(tt.input)
			assert.Equal(t, tt.want, res)
		})
	}
}

func TestSet_Intersection(t *testing.T) {
	tests := []struct {
		name   string
		source Set[string]
		other  Set[string]
		want   Set[string]
	}{
		{
			name:   "Nil Source",
			source: nil,
			other: Set[string]{
				"1": "1",
				"3": "3",
			},
			want: Set[string]{},
		},
		{
			name:   "Empty Source",
			source: Set[string]{},
			other: Set[string]{
				"1": "1",
				"3": "3",
			},
			want: Set[string]{},
		},
		{
			name: "Nil Other",
			source: Set[string]{
				"1": "1",
				"2": "2",
				"3": "3",
				"4": "4",
			},
			other: nil,
			want:  Set[string]{},
		},
		{
			name: "Empty other",
			source: Set[string]{
				"1": "1",
				"2": "2",
				"3": "3",
				"4": "4",
			},
			other: Set[string]{},
			want:  Set[string]{},
		},
		{
			name: "Match",
			other: Set[string]{
				"1": "1",
				"3": "3",
			},
			source: Set[string]{
				"1": "1",
				"2": "2",
				"3": "3",
				"4": "4",
			},
			want: Set[string]{
				"1": "1",
				"3": "3",
			},
		},
		{
			name: "No Match",
			other: Set[string]{
				"5": "5",
				"6": "6",
			},
			source: Set[string]{
				"1": "1",
				"2": "2",
				"3": "3",
				"4": "4",
			},
			want: Set[string]{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := tt.source.Intersection(tt.other)
			assert.Equal(t, tt.want, res)
		})
	}
}

func TestSet_Difference(t *testing.T) {
	tests := []struct {
		name   string
		source Set[string]
		other  Set[string]
		want   Set[string]
	}{
		{
			name:   "Nil Source",
			source: nil,
			other: Set[string]{
				"1": "1",
				"3": "3",
			},
			want: Set[string]{},
		},
		{
			name:   "Empty Source",
			source: Set[string]{},
			other: Set[string]{
				"1": "1",
				"3": "3",
			},
			want: Set[string]{},
		},
		{
			name: "Nil other",
			source: Set[string]{
				"1": "1",
				"2": "2",
				"3": "3",
				"4": "4",
			},
			other: nil,
			want: Set[string]{
				"1": "1",
				"2": "2",
				"3": "3",
				"4": "4",
			},
		},
		{
			name: "Empty other",
			source: Set[string]{
				"1": "1",
				"2": "2",
				"3": "3",
				"4": "4",
			},
			other: Set[string]{},
			want: Set[string]{
				"1": "1",
				"2": "2",
				"3": "3",
				"4": "4",
			},
		},
		{
			name: "Match",
			other: Set[string]{
				"1": "1",
				"3": "3",
			},
			source: Set[string]{
				"1": "1",
				"2": "2",
				"3": "3",
				"4": "4",
			},
			want: Set[string]{
				"2": "2",
				"4": "4",
			},
		},
		{
			name: "No Match",
			other: Set[string]{
				"5": "5",
				"6": "6",
			},
			source: Set[string]{
				"1": "1",
				"2": "2",
				"3": "3",
				"4": "4",
			},
			want: Set[string]{
				"1": "1",
				"2": "2",
				"3": "3",
				"4": "4",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := tt.source.Difference(tt.other)
			assert.Equal(t, tt.want, res)
		})
	}
}

func TestSet_Filter(t *testing.T) {
	tests := []struct {
		name   string
		f      func(v int) bool
		source Set[int]
		want   Set[int]
	}{
		{
			name: "Empty Source",
			f: func(v int) bool {
				return true
			},
			source: Set[int]{},
			want:   Set[int]{},
		},
		{
			name: "Always true",
			f: func(v int) bool {
				return true
			},
			source: Set[int]{
				1: 1,
				2: 2,
				3: 3,
				4: 4,
			},
			want: Set[int]{
				1: 1,
				2: 2,
				3: 3,
				4: 4,
			},
		},
		{
			name: "Always false",
			f: func(v int) bool {
				return false
			},
			source: Set[int]{
				1: 1,
				2: 2,
				3: 3,
				4: 4,
			},
			want: Set[int]{},
		},
		{
			name: "Even numbers",
			f: func(v int) bool {
				return v%2 == 0
			},
			source: Set[int]{
				1: 1,
				2: 2,
				3: 3,
				4: 4,
			},
			want: Set[int]{
				2: 2,
				4: 4,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := tt.source.Filter(tt.f)
			assert.Equal(t, tt.want, res)
		})
	}
}

func TestSet_List(t *testing.T) {
	tests := []struct {
		name   string
		source Set[string]
		want   []string
	}{
		{
			name:   "Nil  set",
			source: nil,
			want:   nil,
		},
		{
			name:   "Empty set",
			source: Set[string]{},
			want:   []string{},
		},
		{
			name: "Values",
			source: Set[string]{
				"1": "1",
				"2": "2",
			},
			want: []string{"1", "2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := tt.source.List()
			assert.Equal(t, tt.want, res)
		})
	}
}

func TestSet_Copy(t *testing.T) {
	tests := []struct {
		name   string
		source Set[string]
		want   Set[string]
	}{
		{
			name:   "Empty set",
			source: Set[string]{},
			want:   Set[string]{},
		},
		{
			name: "Values",
			source: Set[string]{
				"1": "1",
				"2": "2",
			},
			want: Set[string]{
				"1": "1",
				"2": "2",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := tt.source.Copy()
			assert.Equal(t, tt.want, res)
			assert.NotSame(t, tt.want, res)
		})
	}
}

func BenchmarkSet_SetIntersection_100_100000(b *testing.B) {
	small := makeBenchmarkSet(b, 100)
	large := makeBenchmarkSet(b, 100000)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		small.Intersection(large)
	}
}

func BenchmarkSet_SetIntersection_100000_100(b *testing.B) {
	small := makeBenchmarkSet(b, 100)
	large := makeBenchmarkSet(b, 100000)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		large.Intersection(small)
	}
}

func makeBenchmarkSet(b *testing.B, n int) Set[string] {
	b.Helper()

	ret := make(Set[string], n)
	for i := 0; i < n; i++ {
		ret.Add(strconv.Itoa(i))
	}
	return ret
}
