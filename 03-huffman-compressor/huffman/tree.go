package huffman

// HuffmanNode Represents a node of a Huffman tree
type HuffmanNode struct {
	Weight uint64
	Parent *HuffmanNode
	Left   *HuffmanNode
	Right  *HuffmanNode
	Byte   byte
	Code   *HuffmanCode
}

// IsLeaf Determines whether the current node is a leaf node
func (nd *HuffmanNode) IsLeaf() bool {
	return nd.Left == nil && nd.Right == nil
}

// IsLeft Determine whether the current node is a left subtree
func (nd *HuffmanNode) IsLeft() bool {
	if nd.Parent != nil {
		return nd.Parent.Left == nd
	}
	return false
}

// IsRight Determines whether the current node is a right subtree
func (nd *HuffmanNode) IsRight() bool {
	if nd.Parent != nil {
		return nd.Parent.Right == nd
	}
	return false
}

// setCode Set the Huffman encoding for this node
// Setup idea: Keep going up from the leaf node, add bits, and finally reverse all bits
func (nd *HuffmanNode) setCode() {
	cur := nd
	huffmanBits := &HuffmanCode{}
	for cur != nil {
		if cur.IsLeft() {
			huffmanBits.AppendZero()
		} else if cur.IsRight() {
			huffmanBits.AppendOne()
		}
		cur = cur.Parent
	}
	nd.Code = huffmanBits.ReverseNew()
}

// WeightLength Calculate the weighted path length
// This method is only available in nd. The code is set to get a valid value
func (nd *HuffmanNode) WeightLength() int {
	return int(nd.Weight) * nd.Code.BitLen()
}

// HuffmanTree Represents a Huffman tree
type HuffmanTree struct {
	Freq   Frequencies
	Root   *HuffmanNode
	Leaves []*HuffmanNode
}

// NewHuffmanTree Construct a new Huffman tree based on the specified frequency
func NewHuffmanTree(freq Frequencies) *HuffmanTree {
	root, leaves := ConstructHuffmanTree(freq)
	tree := &HuffmanTree{
		Freq:   freq,
		Root:   root,
		Leaves: leaves,
	}

	return tree
}

// ConstructHuffmanTree Create a Huffman tree based on frequency
// Returns the root node and all leaf nodes of the Huffman tree
func ConstructHuffmanTree(freq Frequencies) (*HuffmanNode, []*HuffmanNode) {
	// Cases where only one data is processed separately
	if len(freq) == 1 {
		var k byte
		var v uint64
		for k, v = range freq {
		}
		root := &HuffmanNode{}
		left := &HuffmanNode{Parent: root, Weight: v, Byte: k}
		root.Left = left

		left.Code = NewHuffmanCodeFromString("0")
		return root, []*HuffmanNode{left}
	}

	// 1. Build a priority queue
	pq := NewHuffmanPQ()

	// Inserts all leaf nodes
	var leaves []*HuffmanNode = make([]*HuffmanNode, 0, len(freq))
	for k, v := range freq {
		node := &HuffmanNode{Weight: v, Byte: k}
		leaves = append(leaves, node)
		pq.Push(node)
	}

	pq.UpdateOrder()

	// 2. Start building the Huffman tree
	for pq.Size() > 1 {
		// pop out the two nodes with the smallest weight
		nodeA := pq.Pop()
		nodeB := pq.Pop()
		if (nodeA.Weight == nodeB.Weight) && (nodeA.Byte > nodeB.Byte) {
			nodeA, nodeB = nodeB, nodeA
		}
		// Add a root node to merge this node
		nodeRoot := &HuffmanNode{
			Left:   nodeA,
			Right:  nodeB,
			Weight: nodeA.Weight + nodeB.Weight,
		}
		nodeA.Parent = nodeRoot
		nodeB.Parent = nodeRoot
		// Insert the new node back into the priority queue
		pq.Push(nodeRoot)
	}

	// Code the leaf nodes
	for _, leaf := range leaves {
		// fmt.Println(i)
		leaf.setCode()
	}

	return pq.Peek(), leaves
}
