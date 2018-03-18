package tree

import (
	"sync/atomic"
	"unsafe"
)

// BinaryTree is a self balancing AVL tree.
type BinaryTree struct {
	root unsafe.Pointer

	// Stat object
	stats unsafe.Pointer
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
	var s Stat

	return &BinaryTree{
		stats: unsafe.Pointer(&s),
	}
}

// Root returns the tree's root node. If the tree is empty, it will return
// nil.
func (t *BinaryTree) Root() *Node {
	return (*Node)(atomic.LoadPointer(&t.root))
}

// Stat is the result of calling the Stats method.
type Stat struct {
	Added   int
	Dropped int
	Size    int
}

// Stats returns the current stats of the tree.
func (t *BinaryTree) Stats() Stat {
	s := *(*Stat)(atomic.LoadPointer(&t.stats))
	s.Size = s.Added - s.Dropped
	return s
}

// Insert adds an entry to the BinaryTree. This can only be called by a single
// go-routine. However many go-routines can be reading while Insert is being
// called. Therefore it is a single producer, many consumer.
func (t *BinaryTree) Insert(key int64, value interface{}) {
	r := t.insert(key, value, (*Node)(t.root))
	atomic.StorePointer(&t.root, unsafe.Pointer(r))
}

func (t *BinaryTree) insert(key int64, value interface{}, n *Node) *Node {
	if n == nil {
		s := *(*Stat)(atomic.LoadPointer(&t.stats))
		s.Added++
		atomic.StorePointer(&t.stats, unsafe.Pointer(&s))
		return &Node{Key: key, Value: value, height: 1}
	}

	if key < n.Key {
		left := t.insert(key, value, n.left)
		n = &Node{
			Key:   n.Key,
			Value: n.Value,
			left:  left,
			right: n.right,
		}

		n.height = t.findHeight(n.left, n.right)

		return t.balance(n, key)
	}

	if key > n.Key {
		right := t.insert(key, value, n.right)

		n = &Node{
			Key:   n.Key,
			Value: n.Value,
			left:  n.left,
			right: right,
		}
		n.height = t.findHeight(n.left, n.right)

		return t.balance(n, key)
	}

	return &Node{
		Key:    key,
		Value:  value,
		left:   n.left,
		right:  n.right,
		height: n.height,
	}
}

func (t *BinaryTree) rightRotate(y *Node) *Node {
	x := y.left
	t2 := x.right

	y = &Node{
		Key:   y.Key,
		Value: y.Value,
		left:  t2,
		right: y.right,
	}

	x = &Node{
		Key:   x.Key,
		Value: x.Value,
		left:  x.left,
		right: y,
	}

	y.height = t.findNodeHeight(y)
	x.height = t.findNodeHeight(x)

	return x
}

func (t *BinaryTree) leftRotate(x *Node) *Node {
	y := x.right
	t2 := y.left

	x = &Node{
		Key:   x.Key,
		Value: x.Value,
		left:  x.left,
		right: t2,
	}

	y = &Node{
		Key:   y.Key,
		Value: y.Value,
		left:  x,
		right: y.right,
	}

	x.height = t.findNodeHeight(x)
	y.height = t.findNodeHeight(y)

	return y
}

func (t *BinaryTree) balance(n *Node, key int64) *Node {
	hl := t.findNodeHeight(n.left)
	hr := t.findNodeHeight(n.right)
	b := hl - hr

	// Left Left
	if b > 1 && key < n.left.Key {
		return t.rightRotate(n)
	}

	// Right Right
	if b < -1 && key > n.right.Key {
		return t.leftRotate(n)
	}

	// Left Right
	if b > 1 && key > n.left.Key {
		n.left = t.leftRotate(n.left)
		return t.rightRotate(n)
	}

	// Right Left
	if b < -1 && key < n.right.Key {

		// Check to see if we just dropped the left most node (without
		// balancing)
		if n.left == nil || n.right.left == nil {
			return n
		}

		n.right = t.rightRotate(n.right)
		return t.leftRotate(n)
	}

	return n
}

// DropLeft removes the left most node. If the tree is empty, then it is a
// nop. This can only be called on the same go-routine as the Insert
// go-routine. It can be called in parallel with consumers.
func (t *BinaryTree) DropLeft() {
	s := *(*Stat)(atomic.LoadPointer(&t.stats))
	s.Dropped++
	atomic.StorePointer(&t.stats, unsafe.Pointer(&s))

	r := t.dropLeft((*Node)(t.root))
	atomic.StorePointer(&t.root, unsafe.Pointer(r))
}

func (t *BinaryTree) dropLeft(n *Node) *Node {
	if n == nil {
		return nil
	}

	if n.left == nil {
		// Found left most node
		return n.right
	}

	n = &Node{
		Key:   n.Key,
		Value: n.Value,
		left:  t.dropLeft(n.left),
		right: n.right,
	}

	n.height = t.findNodeHeight(n)

	return n
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
