package tree_test

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"testing"
	"time"

	tree "github.com/apoydence/go-binarytree"
	"github.com/apoydence/onpar"
	. "github.com/apoydence/onpar/expect"
	. "github.com/apoydence/onpar/matchers"
)

type TT struct {
	*testing.T
	bt *tree.BinaryTree
}

func TestTree(t *testing.T) {
	t.Parallel()
	o := onpar.New()
	defer o.Run(t)

	o.BeforeEach(func(t *testing.T) TT {
		bt := tree.New()

		return TT{
			T:  t,
			bt: bt,
		}
	})

	o.Spec("it maintains a required binary tree structure", func(t TT) {
		values := []int64{7, 0, 5, 1, 9, 2, 6, 12, 11, 3, 8, 13, 10, 4}

		for i, x := range values {
			t.bt.Insert(x, fmt.Sprintf("%d", i))
			if i%10 == 0 {
				t.bt.DropLeft()
			}

			var keys []int64
			tree.Traverse(t.bt.Root(), func(key int64, value interface{}) bool {
				keys = append(keys, key)
				return true
			})

			Expect(t, sort.IsSorted(ints(keys))).To(BeTrue())
		}

		var keys []int64
		tree.Traverse(t.bt.Root(), func(key int64, value interface{}) bool {
			keys = append(keys, key)
			return true
		})

		Expect(t, sort.IsSorted(ints(keys))).To(BeTrue())
	})

	o.Spec("drops the left node and keeps the right", func(t TT) {
		//        5
		//      /   \
		//     3     6
		//    / \
		//   2   4
		for _, i := range []int64{5, 3, 6, 4, 2} {
			t.bt.Insert(i, fmt.Sprintf("%d", i))
		}

		t.bt.DropLeft()

		Expect(t, t.bt.Stats()).To(Equal(tree.Stat{
			Added:   5,
			Dropped: 1,
			Size:    4,
		}))

		var keys []int64
		tree.Traverse(t.bt.Root(), func(key int64, value interface{}) bool {
			keys = append(keys, key)
			return true
		})

		Expect(t, keys).To(Equal([]int64{
			3, 4, 5, 6,
		}))
	})

	o.Spec("it balances for Left Left", func(t TT) {
		//  T1, T2, T3 and T4 are subtrees.
		//        z                                      y
		//       / \                                   /   \
		//      y   T4      Right Rotate (z)          x      z
		//     / \          - - - - - - - - ->      /  \    /  \
		//    x   T3                               T1  T2  T3  T4
		//   / \
		// T1   T2

		//        9
		//      /   \
		//     7     11
		//    / \
		//   4   8
		//  /
		// 1

		for _, i := range []int64{9, 7, 11, 4, 8, 1} {
			t.bt.Insert(i, fmt.Sprintf("%d", i))
		}

		Expect(t, t.bt.Stats().Size).To(Equal(6))
		Expect(t, tree.HeightFrom(1, t.bt.Root())).To(Equal(3))

		var keys []int64
		tree.Traverse(t.bt.Root(), func(key int64, value interface{}) bool {
			keys = append(keys, key)
			return true
		})

		Expect(t, keys).To(Equal([]int64{
			1, 4, 7, 8, 9, 11,
		}))
	})

	o.Spec("it balances for Left Right", func(t TT) {
		//      z                               z                           x
		//     / \                            /   \                        /  \
		//    y   T4  Left Rotate (y)        x    T4  Right Rotate(z)    y      z
		//   / \      - - - - - - - - ->    /  \      - - - - - - - ->  / \    / \
		// T1   x                          y    T3                    T1  T2 T3  T4
		//     / \                        / \
		//   T2   T3                    T1   T2

		//       11                            11                           8
		//      /   \                         /  \                         / \
		//     7     12                      8   12                       7   11
		//    / \                           /  \                         /   /  \
		//   4   8                         7   10                       4   10  12
		//        \                       /
		//         10                    4

		for _, i := range []int64{11, 7, 12, 4, 8, 10} {
			t.bt.Insert(i, fmt.Sprintf("%d", i))
		}

		Expect(t, t.bt.Stats().Size).To(Equal(6))
		Expect(t, tree.HeightFrom(12, t.bt.Root())).To(Equal(3))

		var keys []int64
		tree.Traverse(t.bt.Root(), func(key int64, value interface{}) bool {
			keys = append(keys, key)
			return true
		})

		Expect(t, keys).To(Equal([]int64{
			4, 7, 8, 10, 11, 12,
		}))
	})

	o.Spec("it balances for Right Right", func(t TT) {
		//   z                                y
		//  /  \                            /   \
		// T1   y     Left Rotate(z)       z      x
		//     /  \   - - - - - - - ->    / \    / \
		//    T2   x                     T1  T2 T3  T4
		//        / \
		//      T3  T4

		//    5
		//  /   \
		// 4    12
		//      /  \
		//    11   14
		//           \
		//            15

		for _, i := range []int64{5, 4, 12, 11, 14, 15} {
			t.bt.Insert(i, fmt.Sprintf("%d", i))
		}

		Expect(t, t.bt.Stats().Size).To(Equal(6))
		Expect(t, tree.HeightFrom(15, t.bt.Root())).To(Equal(3))

		var keys []int64
		tree.Traverse(t.bt.Root(), func(key int64, value interface{}) bool {
			keys = append(keys, key)
			return true
		})

		Expect(t, keys).To(Equal([]int64{
			4, 5, 11, 12, 14, 15,
		}))
	})

	o.Spec("it balances for Right Left", func(t TT) {
		//    z                            z                            x
		//   / \                          / \                          /  \
		// T1   y   Right Rotate (y)    T1   x      Left Rotate(z)   z      y
		//     / \  - - - - - - - - ->     /  \   - - - - - - - ->  / \    / \
		//    x   T4                      T2   y                  T1  T2  T3  T4
		//   / \                              /  \
		// T2   T3                           T3   T4

		//    5
		//  /   \
		// 4    14
		//      /  \
		//    11   15
		//      \
		//       12

		for _, i := range []int64{5, 4, 14, 11, 15, 12} {
			t.bt.Insert(i, fmt.Sprintf("%d", i))
		}

		Expect(t, t.bt.Stats().Size).To(Equal(6))
		Expect(t, tree.HeightFrom(15, t.bt.Root())).To(Equal(3))

		var keys []int64
		tree.Traverse(t.bt.Root(), func(key int64, value interface{}) bool {
			keys = append(keys, key)
			return true
		})

		Expect(t, keys).To(Equal([]int64{
			4, 5, 11, 12, 14, 15,
		}))
	})

	o.Spec("fuzz", func(t TT) {
		rand.Seed(time.Now().UnixNano())

		var j int64
		for i := int64(0); i < 1000; i++ {
			value := rand.Int63()
			t.bt.Insert(value, "")

			if i%10 == 0 {
				t.bt.DropLeft()
				j++
			}

			perfectHeight := int(math.Ceil(math.Log2(float64(i-j+2))) - 1)
			h := tree.HeightFrom(value, t.bt.Root())

			// We'll allow for some wiggle room.
			Expect(t, h-perfectHeight).To(BeBelow(5))

			var keys []int64
			tree.Traverse(t.bt.Root(), func(key int64, value interface{}) bool {
				keys = append(keys, key)
				return true
			})

			Expect(t, sort.IsSorted(ints(keys))).To(BeTrue())
		}
	})

	o.Spec("it survives the race detector", func(t TT) {
		go func() {
			for i := int64(0); i < 10000; i++ {
				t.bt.Insert(i, fmt.Sprintf("%d", i))

				if i%10 == 0 {
					t.bt.DropLeft()
				}
			}
		}()

		var result []int64
		Expect(t, func() []int64 {
			var keys []int64
			tree.Traverse(t.bt.Root(), func(key int64, value interface{}) bool {
				// Access and throw away
				t.bt.Stats()

				keys = append(keys, key)
				return true
			})
			return keys
		}).To(ViaPolling(Chain(HaveLen(9000), Fetch(&result))))

		// First 1000 have been pruned
		for i := int64(1000); i < 10000; i++ {
			Expect(t, result[int(i)-1000]).To(Equal(i))
		}
	})
}

type ints []int64

func (i ints) Len() int {
	return len(i)
}

func (i ints) Less(a, b int) bool {
	return i[a] < i[b]
}

func (i ints) Swap(a, b int) {
	t := i[a]
	i[a] = i[b]
	i[b] = t
}
