package main

import (
	"container/heap"
	"reflect"
	"sort"
	"testing"
)

func TestPriorityQueuePopOrder(t *testing.T) {
	indices := []int{0, 9, 4, 5, 6, 3, 8, 1, 7, 2}

	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	for idx := range indices {
		heap.Push(&pq, &indexedDataBlock{index: idx, dataBlock: dataBlock{}})
	}

	var poppedValues []int
	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*indexedDataBlock)
		t.Logf("Popped idx: %v", item.dataBlock)
		poppedValues = append(poppedValues, item.index)
	}

	sortedValues := make([]int, len(indices))
	copy(sortedValues, indices)
	sort.Ints(sortedValues)

	if !reflect.DeepEqual(poppedValues, sortedValues) {
		t.Errorf("Priority queue pop order is incorrect: Expected %v, got %v", sortedValues, poppedValues)
	}
}

// func TestPriorityQueueUpdate(t *testing.T) {
// 	indices := []int{0, 9, 4, 5, 6, 3, 8, 1, 7, 2}

// 	pq := make(PriorityQueue, 0)
// 	heap.Init(&pq)

// 	for idx := range indices {
// 		heap.Push(&pq, &indexedDataBlock{index: idx, dataBlock: dataBlock{}})
// 	}

// 	// Update the idx of the item with index 5
// 	targetIndex := 5
// 	newDataBuf := [4096]byte{1, 2, 3, 4}
// 	newDataBlock := makeIndexedDataBlock(targetIndex, newDataBuf, 4)
// 	pq.update(pq[targetIndex], targetIndex, newDataBlock.index)

// 	var poppedValues []int
// 	for pq.Len() > 0 {
// 		item := heap.Pop(&pq).(*indexedDataBlock)
// 		t.Logf("Popped idx: %v", item.index)
// 		poppedValues = append(poppedValues, item.index)
// 	}

// 	// The updated idx should be the first one to be popped
// 	if poppedValues[targetIndex] != makeIndexedDataBlock(targetIndex, newValue).index {
// 		t.Errorf("Priority queue update is incorrect: Expected %v, got %v", newValue, poppedValues[0])
// 	}
// }

func TestEmptyPriorityQueue(t *testing.T) {
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	if pq.Len() != 0 {
		t.Errorf("Priority queue length is incorrect: Expected 0, got %v", pq.Len())
	}

	// this test throws a panic: runtime error: index out of range [0] with length 0
	// item := heap.Pop(&pq)
	// if item != nil {
	// 	t.Errorf("Popped item is not nil: Expected nil, got %v", item)
	// }
}
