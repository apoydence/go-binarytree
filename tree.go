package tree

import (
	"sync/atomic"
	"unsafe"
)

// BinaryTree is a self balancing AVL tree.
type BinaryTree struct {
	root unsafe.Pointer

	added int64
}

// Node stores a tree's vertice.
type Node struct {
	Key   int64
	Value interface{}

	// These should be treated as final. They should never be altered once
	// set.
	left, right *Node

	height int
}

// Left returns the left node.
func (n *Node) Left() *Node {
	return n.left
}

// Right returns the right node.
func (n *Node) Right() *Node {
	return n.right
}

// New returns a new BinaryTree.
func New() *BinaryTree {
	return &BinaryTree{}
}

// Root returns the tree's root node. If the tree is empty, it will return
// nil.
func (t *BinaryTree) Root() *Node {
	return (*Node)(atomic.LoadPointer(&t.root))
}

// Insert adds an entry to the BinaryTree. This can only be called by a single
// go-routine. However many go-routines can be reading while Insert is being
// called. Therefore it is a single producer, many consumer.
func (t *BinaryTree) Insert(key int64, value interface{}) {
	r, _, _ := t.insert(key, value, (*Node)(t.root))
	atomic.StorePointer(&t.root, unsafe.Pointer(r))
}

func (t *BinaryTree) insert(key int64, value interface{}, n *Node) (*Node, *Node, bool) {
	if n == nil {
		atomic.AddInt64(&t.added, 1)
		return &Node{Key: key, Value: value, height: 1}, nil, true
	}

	if key < n.Key {
		left, x, ok := t.insert(key, value, n.left)
		n = &Node{
			Key:    n.Key,
			Value:  n.Value,
			height: t.findHeight(left, n.right),
			left:   left,
			right:  n.right,
		}

		n = t.balance(x, left, n)
		return n, left, ok
	}

	if key > n.Key {
		right, x, ok := t.insert(key, value, n.right)

		n = &Node{
			Key:    n.Key,
			Value:  n.Value,
			height: t.findHeight(n.left, right),
			left:   n.left,
			right:  right,
		}

		n = t.balance(x, right, n)
		return n, right, ok
	}

	return &Node{
		Key:    key,
		Value:  value,
		left:   n.left,
		right:  n.right,
		height: n.height,
	}, nil, false
}

func (t *BinaryTree) balance(x, y, z *Node) *Node {
	hl := t.findNodeHeight(z.left)
	hr := t.findNodeHeight(z.right)

	if hr-hl < -1 {
		// Left
		if x == y.left {
			// Left Left
			// Right rotate (z)

			z = &Node{
				Key:    z.Key,
				Value:  z.Value,
				height: z.height,
				left:   y.right,
				right:  z.right,
			}

			y = &Node{
				Key:    y.Key,
				Value:  y.Value,
				height: y.height,
				left:   y.left,
				right:  z,
			}

			// z.left = y.right
			// y.right = z

			z.height = t.findNodeHeight(z)
			y.height = t.findNodeHeight(y)
			return y
		}

		// Left Right
		// Left Rotate (y)

		y = &Node{
			Key:    y.Key,
			Value:  y.Value,
			height: y.height,
			left:   y.left,
			right:  x.left,
		}

		x = &Node{
			Key:    x.Key,
			Value:  x.Value,
			height: x.height,
			left:   y,
			right:  x.right,
		}

		// y.right = x.left
		// x.left = y

		// Right Rotate (z)

		z = &Node{
			Key:    z.Key,
			Value:  z.Value,
			height: z.height,
			left:   x.right,
			right:  z.right,
		}

		// x is already a copy
		x.right = z

		// z.left = x.right

		y.height = t.findNodeHeight(y)
		z.height = t.findNodeHeight(z)
		x.height = t.findNodeHeight(x)
		return x
	}

	if hr-hl > 1 {
		// Right
		if x == y.right {
			// Right Right
			// Left Rotate (z)

			z = &Node{
				Key:    z.Key,
				Value:  z.Value,
				height: z.height,
				left:   z.left,
				right:  y.left,
			}

			y = &Node{
				Key:    y.Key,
				Value:  y.Value,
				height: y.height,
				left:   z,
				right:  y.right,
			}

			// z.right = y.left
			// y.left = z

			z.height = t.findNodeHeight(z)
			y.height = t.findNodeHeight(y)
			return y
		}

		// Right Left
		// Right Rotate (y)

		y = &Node{
			Key:    y.Key,
			Value:  y.Value,
			height: y.height,
			left:   x.right,
			right:  y.right,
		}

		x = &Node{
			Key:    x.Key,
			Value:  x.Value,
			height: x.height,
			left:   x.left,
			right:  y,
		}

		// y.left = x.right
		// x.right = y

		// Left Rotate (z)

		z = &Node{
			Key:    z.Key,
			Value:  z.Value,
			height: z.height,
			left:   z.left,
			right:  x.left,
		}

		// z.right = x.left

		// x is already a copy
		x.left = z

		y.height = t.findNodeHeight(y)
		z.height = t.findNodeHeight(z)
		x.height = t.findNodeHeight(x)
		return x
	}

	// Balanced
	return z
}

// Size returns the number of entries the tree is currently storing.
func (t *BinaryTree) Size() int {
	return int(atomic.LoadInt64(&t.added))
}

func (t *BinaryTree) findNodeHeight(n *Node) int {
	if n == nil {
		return 0
	}

	return t.findHeight(n.left, n.right)
}

func (t *BinaryTree) findHeight(l, r *Node) int {
	var hl, hr int

	if l != nil {
		hl = l.height
	}

	if r != nil {
		hr = r.height
	}

	if hl > hr {
		return hl + 1
	}

	return hr + 1
}

// Traverse is used to traverse a tree starting at the given node.
func Traverse(n *Node, f func(key int64, value interface{}) (keepGoing bool)) bool {
	if n == nil {
		return true
	}

	if !Traverse(n.Left(), f) {
		return false
	}
	if !f(n.Key, n.Value) {
		return false
	}
	if !Traverse(n.Right(), f) {
		return false
	}

	return true
}

// HeightFrom measures the height from the given key via traversing.
func HeightFrom(key int64, n *Node) int {
	return heightFrom(key, 0, n)
}

func heightFrom(key int64, count int, n *Node) int {
	if n == nil {
		return 0
	}

	if key < n.Key {
		return heightFrom(key, count+1, (*Node)(n.left))
	}

	if key > n.Key {
		return heightFrom(key, count+1, (*Node)(n.right))
	}

	return count + 1
}
