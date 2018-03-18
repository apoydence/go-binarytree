package tree_test

import (
	"math/rand"
	"sync"
	"sync/atomic"
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

func BenchmarkTreeHeavyRead(b *testing.B) {
	b.ReportAllocs()
	rand.Seed(time.Now().UnixNano())
	t := tree.New()

	b.ResetTimer()

	var done int64
	var wg sync.WaitGroup
	defer wg.Wait()
	wg.Add(1)

	defer func() {
		atomic.AddInt64(&done, 1)
	}()

	go func() {
		defer wg.Done()
		for atomic.LoadInt64(&done) == 0 {
			findRight(t.Root())
		}
	}()

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

	var done int64
	var wg sync.WaitGroup
	defer wg.Wait()
	wg.Add(1)

	defer func() {
		atomic.AddInt64(&done, 1)
	}()

	go func() {
		defer wg.Done()
		var i int
		for atomic.LoadInt64(&done) == 0 {
			i++
			value := rand.Int63()
			t.Insert(value, struct{}{})

			if i > 100000 {
				t.DropLeft()
			}
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
