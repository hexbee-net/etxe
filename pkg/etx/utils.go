package etx

import (
	"fmt"
	"io"
)

const (
	indentationChar = "\t"
)

type Cloner[C any] interface {
	Clone() C
}

func cloneCollection[T Cloner[T]](src []T) []T {
	if src == nil {
		return nil
	}

	out := make([]T, 0, len(src))
	for _, item := range src {
		out = append(out, item.Clone())
	}

	return out
}

func cloneStrings(strings []string) []string {
	if strings == nil {
		return nil
	}
	out := make([]string, len(strings))
	copy(out, strings)

	return out
}

// indent inserts prefix at the beginning of each non-empty line of s.
func indent(s, prefix string) string {
	return string(indentBytes([]byte(s), []byte(prefix)))
}

// indentBytes inserts prefix at the beginning of each non-empty line of b.
func indentBytes(b, prefix []byte) []byte {
	var res []byte
	bol := true
	for _, c := range b {
		if bol && c != '\n' {
			res = append(res, prefix...)
		}
		res = append(res, c)
		bol = c == '\n'
	}

	return res
}

func mustFprintf(w io.Writer, format string, a ...any) {
	if _, err := fmt.Fprintf(w, format, a...); err != nil {
		panic(err)
	}
}
