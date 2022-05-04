package etx

// Bool represents a parsed boolean value.
type Bool bool

func (b *Bool) Capture(values []string) error {
	*b = values[0] == "true"

	return nil
}
