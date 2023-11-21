package huffman

import (
	"container/heap"
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHuffmanPQ_Primary(t *testing.T) {
	nodes := huffmanPQ{
		&HuffmanNode{Weight: 10},
		&HuffmanNode{Weight: 5},
		&HuffmanNode{Weight: 2},
		&HuffmanNode{Weight: 3},
		&HuffmanNode{Weight: 1},
		&HuffmanNode{Weight: 6},
		&HuffmanNode{Weight: 7},
		&HuffmanNode{Weight: 9},
	}

	required := []int{1, 2, 3, 5, 6, 7, 9, 10}

	heap.Init(&nodes)

	size := nodes.Len()
	for i := 0; i < size; i++ {
		l := heap.Pop(&nodes).(*HuffmanNode).Weight
		fmt.Printf("%d ", l)
		require.EqualValues(t, required[i], l)
	}
	require.Zero(t, nodes.Len())
	fmt.Println()

	push := []uint64{9, 8, 7, 6, 5, 4, 3, 2, 4, 6}
	required2 := []int{2, 3, 4, 4, 5, 6, 6, 7, 8, 9}
	for i := 0; i < len(push); i++ {
		heap.Push(&nodes, &HuffmanNode{Weight: push[i]})
	}
	require.EqualValues(t, len(push), nodes.Len())

	size = nodes.Len()
	for i := 0; i < size; i++ {
		l := heap.Pop(&nodes).(*HuffmanNode).Weight
		fmt.Printf("%d ", l)
		require.EqualValues(t, required2[i], l)
	}
	fmt.Println()
}

func TestHuffmanPQ_Secondary(t *testing.T) {
	nodes := huffmanPQ{
		&HuffmanNode{Weight: 10},
		&HuffmanNode{Weight: 5},
		&HuffmanNode{Weight: 2},
		&HuffmanNode{Weight: 3},
		&HuffmanNode{Weight: 1},
		&HuffmanNode{Weight: 6},
		&HuffmanNode{Weight: 7},
		&HuffmanNode{Weight: 9},
	}
	heap.Init(&nodes)

	// 1 2 3 5 6 7 9 10
	for i := 0; i < 4; i++ {
		heap.Pop(&nodes)
	}

	for i := 0; i < nodes.Len(); i++ {
		fmt.Printf("%d ", nodes[i].Weight)
	}
	fmt.Println()

	push := []uint64{10, 15, 2, 4, 70, 0}
	for i := 0; i < len(push); i++ {
		heap.Push(&nodes, &HuffmanNode{Weight: push[i]})
	}

	required := []uint64{0, 2, 4, 6, 7, 9, 10, 10, 15, 70}
	size := nodes.Len()
	require.EqualValues(t, size, len(required))
	for i := 0; i < size; i++ {
		l := heap.Pop(&nodes).(*HuffmanNode).Weight
		fmt.Printf("%d ", l)
		require.EqualValues(t, required[i], l)
	}
	fmt.Println()
}

func TestHuffmanPQ(t *testing.T) {

	testCases := []struct {
		Ints []uint64
	}{
		{Ints: []uint64{10, 4, 2, 3, 7, 8, 1, 0}},
		{Ints: []uint64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
		{Ints: []uint64{9, 8, 7, 6, 5, 4, 3, 2, 1}},
		{Ints: []uint64{2, 4, 0, 6, 8, 2, 0, 10}},
		{Ints: []uint64{1, 3, 5, 9, 7}},
	}

	for _, tc := range testCases {
		nodes := NewHuffmanPQFromUInts(tc.Ints)
		size := nodes.Size()
		require.EqualValues(t, size, len(tc.Ints))

		sort.Slice(tc.Ints, func(i, j int) bool { return tc.Ints[i] < tc.Ints[j] })

		for i := 0; i < size; i++ {
			pl := nodes.Peek().Weight
			l := nodes.Pop().Weight
			require.EqualValues(t, pl, l)
			require.EqualValues(t, tc.Ints[i], l)
		}
	}

}

func TestHuffmanPQ_PushPop(t *testing.T) {
	testCases2 := []struct {
		OriginalInts []uint64
		NumPop       int
		AfterPopInts []uint64
	}{
		{
			OriginalInts: []uint64{1, 3, 9, 2, 4},
			NumPop:       1,
			AfterPopInts: []uint64{5, 7, 9, 0},
		},
		{
			OriginalInts: []uint64{5, 9, 3, 52, 0, 71, 62, 36, 0},
			NumPop:       6,
			AfterPopInts: []uint64{8, 0, 41, 963, 41},
		},
		{
			OriginalInts: []uint64{9, 3, 6, 45, 46, 47, 20},
			NumPop:       7,
			AfterPopInts: []uint64{1, 0},
		},
		{
			OriginalInts: []uint64{11, 53, 91, 27, 44},
			NumPop:       0,
			AfterPopInts: []uint64{9, 6, 8, 55, 32},
		},
		{
			OriginalInts: []uint64{9, 7, 8, 6, 2},
			NumPop:       4,
			AfterPopInts: []uint64{63},
		},
	}

	for _, tc := range testCases2 {
		// Create a priority queue
		nodes := NewHuffmanPQFromUInts(tc.OriginalInts)
		size := nodes.Size()
		require.EqualValues(t, size, len(tc.OriginalInts))

		// Then pop NumPop times
		for i := 0; i < tc.NumPop; i++ {
			nodes.Pop()
		}

		// Then insert AfterPopInts
		for _, ii := range tc.AfterPopInts {
			nodes.Push(&HuffmanNode{Weight: ii})
		}

		sort.Slice(tc.OriginalInts, func(i, j int) bool {
			return tc.OriginalInts[i] < tc.OriginalInts[j]
		})

		tc.OriginalInts = tc.OriginalInts[tc.NumPop:]
		required := append(tc.OriginalInts, tc.AfterPopInts...)
		sort.Slice(required, func(i, j int) bool {
			return required[i] < required[j]
		})

		size = nodes.Size()
		require.EqualValues(t, size, len(required))

		for i := 0; i < size; i++ {
			pl := nodes.Peek().Weight
			l := nodes.Pop().Weight
			require.EqualValues(t, pl, l)
			require.EqualValues(t, required[i], l)
		}

	}
}

func TestHuffmanPQ_Abnormal(t *testing.T) {
	require.Panics(t, func() {
		NewHuffmanPQ().Pop()
	})

	require.Panics(t, func() {
		NewHuffmanPQ().Peek()
	})

}
