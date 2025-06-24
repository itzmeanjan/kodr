package gf256_test

import (
	"testing"

	"github.com/itzmeanjan/kodr"
	"github.com/itzmeanjan/kodr/kodr_internals/gf256"
)

// TestGf256Operations tests the properties of GF(2^8) field operations
func TestGf256Operations(t *testing.T) {
	const numTestIterations = 100_000

	for i := 0; i < numTestIterations; i++ {
		// Generate random Gf256 elements
		a := gf256.Random()
		b := gf256.Random()

		// Test Addition, Subtraction, Negation
		sum := a.Add(b)
		diff := sum.Sub(b)
		if !diff.Equal(a) {
			t.Errorf("Addition/Subtraction property failed: %v - %v != %v", sum, b, a)
		}

		// Test Multiplication, Division, Inversion
		mul := a.Mul(b)
		div, err := mul.Div(b)

		if b == gf256.Zero() {
			if err != kodr.ErrCannotInvertGf256AdditiveIndentity {
				t.Errorf("Division by zero should return nil")
			}
		} else {
			if !div.Equal(a) {
				t.Errorf("Multiplication/Division property failed: %v / %v != %v", mul, b, a)
			}
		}
	}
}
