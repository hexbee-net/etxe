package etx

import (
	"fmt"
	"math/big"
	"regexp"
)

var (
	needsOctalPrefix = regexp.MustCompile(`^0\d+$`)
)

// Number of arbitrary precision.
type Number struct{ *big.Float }

func (n *Number) GoString() string {
	return n.String()
}

// Capture override because big.Float doesn't directly support 0-prefix octal parsing.
func (n *Number) Capture(values []string) error {
	value := values[0]
	if needsOctalPrefix.MatchString(value) {
		value = "0o" + value[1:]
	}
	n.Float = big.NewFloat(0)
	if _, _, err := n.Float.Parse(value, 0); err != nil {
		return fmt.Errorf("failed to parse number value '%s': %w", value, err)
	}

	return nil
}
