package priority_queue

import (
	"container/heap"
	"reflect"
	"sort"
	"testing"
)

func TestPriorityQueuePopOrder(t *testing.T) {
	values := []int{0, 9, 4, 5, 6, 3, 8, 1, 7, 2}

	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	for _, value := range values {
		t.Logf("Pushed value: %v", value)
		heap.Push(&pq, &Item{value: value, index: value})
	}

	var poppedValues []int
	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*Item)
		t.Logf("Popped value: %v", item.value)
		poppedValues = append(poppedValues, item.value.(int))
	}

	sortedValues := make([]int, len(values))
	copy(sortedValues, values)
	sort.Ints(sortedValues)

	if !reflect.DeepEqual(poppedValues, sortedValues) {
		t.Errorf("Priority queue pop order is incorrect: Expected %v, got %v", sortedValues, poppedValues)
	}
}

func TestPriorityQueueUpdate(t *testing.T) {
	values := []int{0, 9, 4, 5, 6, 3, 8, 1, 7, 2}

	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	for _, value := range values {
		heap.Push(&pq, &Item{value: value, index: value})
	}

	// Update the value of the item with index 5
	newValue := 100
	targetIndex := 5
	pq.update(pq[targetIndex], newValue, pq[targetIndex].index)

	var poppedValues []int
	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*Item)
		t.Logf("Popped value: %v", item.value)
		poppedValues = append(poppedValues, item.value.(int))
	}

	// The updated value should be the first one to be popped
	if poppedValues[targetIndex] != newValue {
		t.Errorf("Priority queue update is incorrect: Expected %v, got %v", newValue, poppedValues[0])
	}
}
