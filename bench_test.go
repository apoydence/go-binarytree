package tree_test

import (
	"math/rand"
	"testing"
	"time"

	tree "github.com/apoydence/go-binarytree"
)

func BenchmarkTree(b *testing.B) {
	b.ReportAllocs()
	rand.Seed(time.Now().UnixNano())
	t := tree.New()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		value := rand.Int63()
		t.Insert(value, struct{}{})
	}
}

func BenchmarkTreeParallel(b *testing.B) {
	b.ReportAllocs()
	rand.Seed(time.Now().UnixNano())
	t := tree.New()

	b.ResetTimer()

	go func() {
		for {
			value := rand.Int63()
			t.Insert(value, struct{}{})
		}
	}()

	b.RunParallel(func(b *testing.PB) {
		for b.Next() {
			findRight(t.Root())
		}
	})
}

func findRight(n *tree.Node) {
	if n == nil {
		return
	}

	findRight(n.Right())
}
