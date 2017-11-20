package tree

// BinaryTree is a self balancing AVL tree.
type BinaryTree struct {
	root *Node

	added int
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
	return t.root
}

// Insert adds an entry to the BinaryTree.
func (t *BinaryTree) Insert(key int64, value interface{}) {
	t.root, _, _ = t.insert(key, value, t.root)
}

func (t *BinaryTree) insert(key int64, value interface{}, n *Node) (*Node, *Node, bool) {
	if n == nil {
		t.added++
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

		nn := t.balance(x, left, n)
		return nn, left, ok
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

		nn := t.balance(x, right, n)
		return nn, right, ok
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

			z.left = y.right
			y.right = z

			z.height = t.findHeight(z.left, z.right)
			y.height = t.findHeight(y.left, y.right)
			return y
		}

		// Left Right
		// Left Rotate (y)
		y.right = x.left
		x.left = y

		// Right Rotate (z)
		z.left = x.right
		x.right = z

		y.height = t.findHeight(y.left, y.right)
		z.height = t.findHeight(z.left, z.right)
		x.height = t.findHeight(x.left, x.right)
		return x
	}

	if hr-hl > 1 {
		// Right
		if x == y.right {
			// Right Right
			// Left Rotate (z)

			z.right = y.left
			y.left = z

			z.height = t.findHeight(z.left, z.right)
			y.height = t.findHeight(y.left, y.right)
			return y
		}

		// Right Left
		// Right Rotate (y)
		y.left = x.right
		x.right = y

		// Left Rotate (z)
		z.right = x.left
		x.left = z

		y.height = t.findHeight(y.left, y.right)
		z.height = t.findHeight(z.left, z.right)
		x.height = t.findHeight(x.left, x.right)
		return x
	}

	// Balanced
	return z
}

// Size returns the number of entries the tree is currently storing.
func (t *BinaryTree) Size() int {
	return t.added
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

// HeightFrom measures the height from the given key.
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
