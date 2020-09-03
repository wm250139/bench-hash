package hasher_test

import (
	"bench-hash/hasher"
	"testing"
)

func BenchmarkString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = hasher.String("Hello")
	}
}
