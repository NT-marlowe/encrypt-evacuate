package priority_queue

import (
	"container/heap"
)

type Item struct {
	index int // The priority of the item in the queue. The less, the higher priority.
	value any
}

func MakeItem(index int, value any) Item {
	return Item{index: index, value: value}
}

func (item *Item) GetValue() any {
	return item.value
}

func (item *Item) GetIndex() int {
	return item.index
}

type PriorityQueue []*Item

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
	item := x.(*Item)
	// item.index = n
	*pq = append(*pq, item)
}

// Min-heap.
func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil

	*pq = old[0 : n-1]
	return item
}

func (pq *PriorityQueue) update(item *Item, value any, index int) {
	item.value = value
	item.index = index
	heap.Fix(pq, item.index)
}
