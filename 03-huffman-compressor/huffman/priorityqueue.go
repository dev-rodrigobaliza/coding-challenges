package huffman

import "container/heap"

// Define a priority queue
type huffmanPQ []*HuffmanNode

// Implement HEAP. Interface
func (pq huffmanPQ) Len() int {
	return len(pq)
}

func (pq huffmanPQ) Less(i, j int) bool {
	// Minimal heap
	return pq[i].Weight < pq[j].Weight
}

func (pq huffmanPQ) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *huffmanPQ) Push(x interface{}) {
	node := x.(*HuffmanNode)
	*pq = append(*pq, node)
}

func (pq *huffmanPQ) Pop() interface{} {
	old := *pq
	n := len(old)
	node := old[n-1]
	old[n-1] = nil // avoid memory leak
	*pq = old[0 : n-1]
	return node
}

// HuffmanPQ It is an encapsulated huffmanPQ for ease of use
type HuffmanPQ struct {
	data huffmanPQ
}

// NewHuffmanPQ Create and return a priority queue
func NewHuffmanPQ() *HuffmanPQ {
	pq := &HuffmanPQ{}
	pq.data = make(huffmanPQ, 0)
	heap.Init(&pq.data)

	return pq
}

// NewHuffmanPQFromUInts Create and return a priority queue from the uint64 slice
func NewHuffmanPQFromUInts(ints []uint64) *HuffmanPQ {
	pq := &HuffmanPQ{}
	pq.data = make(huffmanPQ, 0, len(ints))

	for i := 0; i < len(ints); i++ {
		pq.data = append(pq.data, &HuffmanNode{Weight: ints[i]})
	}

	heap.Init(&pq.data)

	return pq
}

// ?? Does this method have an impact after initialization? ??
func (pq *HuffmanPQ) UpdateOrder() {
	heap.Init(&pq.data)
}

func (pq *HuffmanPQ) Push(node *HuffmanNode) {
	heap.Push(&pq.data, node)
}

func (pq *HuffmanPQ) Pop() *HuffmanNode {
	return heap.Pop(&pq.data).(*HuffmanNode)
}

func (pq *HuffmanPQ) Size() int {
	return pq.data.Len()
}

func (pq *HuffmanPQ) Peek() *HuffmanNode {
	return pq.data[0]
}
