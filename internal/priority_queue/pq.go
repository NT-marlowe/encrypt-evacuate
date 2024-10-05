package priority_queue

import (
	"container/heap"
)

type Item struct {
	index int // The priority of the item in the queue. The less, the higher priority.
	value any
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].index < pq[j].index
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

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
