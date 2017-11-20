package tree_test

import (
	"fmt"
	"math"
	"math/rand"
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
		t.Skip()
		for i := int64(1); i < 10; i++ {
			t.bt.Insert(i, fmt.Sprintf("%d", i))
		}

		// Insert 0 out of order
		t.bt.Insert(0, "0")

		// Replace key 5
		t.bt.Insert(5, "99")

		Expect(t, t.bt.Size()).To(Equal(10))

		var keys []int64
		var values []interface{}
		tree.Traverse(t.bt.Root(), func(key int64, value interface{}) bool {
			keys = append(keys, key)
			values = append(values, value)
			return true
		})

		Expect(t, keys).To(Equal([]int64{
			0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
		}))

		Expect(t, values).To(Equal([]interface{}{
			"0", "1", "2", "3", "4", "99", "6", "7", "8", "9",
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

		Expect(t, t.bt.Size()).To(Equal(6))
		Expect(t, tree.HeightFrom(11, t.bt.Root())).To(Equal(3))
	})

	o.Spec("it balances for Left Right", func(t TT) {
		//      z                               z                           x
		//     / \                            /   \                        /  \
		//    y   T4  Left Rotate (y)        x    T4  Right Rotate(z)    y      z
		//   / \      - - - - - - - - ->    /  \      - - - - - - - ->  / \    / \
		// T1   x                          y    T3                    T1  T2 T3  T4
		//     / \                        / \
		//   T2   T3                    T1   T2

		//       11
		//      /   \
		//     7     12
		//    / \
		//   4   9
		//        \
		//         10

		for _, i := range []int64{11, 7, 12, 4, 8, 10} {
			t.bt.Insert(i, fmt.Sprintf("%d", i))
		}

		Expect(t, t.bt.Size()).To(Equal(6))
		Expect(t, tree.HeightFrom(12, t.bt.Root())).To(Equal(3))
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

		Expect(t, t.bt.Size()).To(Equal(6))
		Expect(t, tree.HeightFrom(15, t.bt.Root())).To(Equal(3))
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

		Expect(t, t.bt.Size()).To(Equal(6))
		Expect(t, tree.HeightFrom(15, t.bt.Root())).To(Equal(3))
	})

	o.Spec("fuzz", func(t TT) {
		rand.Seed(time.Now().UnixNano())
		for i := int64(0); i < 100000; i++ {
			t.bt.Insert(i, "")
			h := tree.HeightFrom(i, t.bt.Root())
			x := int(math.Floor(math.Log2(float64(i + 1))))
			Expect(t, x-h).To(BeBelow(2))
		}
	})
}
