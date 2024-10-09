package main

// import (
// 	"container/heap"
// )

type PriorityQueue []*indexedDataBlock

func (pq PriorityQueue) Len() int {
	return len(pq)
}

// Min-heap.
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].index < pq[j].index
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

// Min-heap.
// |index| should not be updated because it is used to track the original position of the item.
func (pq *PriorityQueue) Push(x any) {
	// n := len(*pq)
	item := x.(*indexedDataBlock)
	*pq = append(*pq, item)
}

// Min-heap.
func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[0]
	old[0] = nil

	*pq = old[1:n]
	return item
}

// func (pq *PriorityQueue) update(item *indexedDataBlock, index int, value dataBlock) {
// 	item.index = index
// 	item.dataBlock = value
// 	heap.Fix(pq, item.index)
// }

// type Item struct {
// 	index int // The priority of the item in the queue. The less, the higher priority.
// 	value any
// }

// func MakeItem(index int, value any) Item {
// 	return Item{index: index, value: value}
// }

// // func (item *indexedDataBlock) GetValue() any {
// // 	return item.value
// // }

// func (item *indexedDataBlock) GetIndex() int {
// 	return item.index
// }
