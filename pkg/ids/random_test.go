package ids

import (
	"bytes"
	"testing"
)

// This test asserts that there are no duplicate values in a large pool of
// random numbers. It has a very small but non-zero chance of failing randomly.
// A more consistent pattern of failure indicates that it is not generating
// random numbers.
func TestRandomBytes(t *testing.T) {
	iterations := 10000

	results := make([][]byte, iterations)
	for i := 0; i < iterations; i++ {
		results[i] = RandomBytes(16)
		for j := 0; j < i; j++ {
			if bytes.Equal(results[i], results[j]) {
				t.Fatal("Random number collision!")
			}
		}
	}
}

func BenchmarkRandomBytes(b *testing.B) {
	for n := 0; n < b.N; n++ {
		RandomBytes(16)
	}
}
