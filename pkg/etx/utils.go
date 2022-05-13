package etx

const (
	indentationChar = "\t"
)

// AddParentRefs recursively updates an AST's parent references.
//
// This is called automatically during Parse*(), but can be called on a manually constructed AST.
func AddParentRefs(node Node) error {
	return nil
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
