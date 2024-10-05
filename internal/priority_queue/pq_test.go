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

	for value := range values {
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
