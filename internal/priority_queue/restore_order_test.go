package priority_queue

import (
	// "log"
	"reflect"
	"sort"
	"testing"
)

func TestRestoreOrder(t *testing.T) {
	values := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	inputCh := make(chan Item)

	// reorder the values
	for _, value := range values {
		go func() {
			inputCh <- Item{value: value, index: value}
		}()
	}

	outputCh := RestoreOrder(inputCh)
	var restoredValues []int

	i := 0
	for item := range outputCh {
		t.Logf("Restored value: %v", item.value)
		restoredValues = append(restoredValues, item.value.(int))

		i++
		if i == len(values) {
			close(inputCh)
		}
	}

	sortedValues := assendingSortIntSlice(values)
	if !reflect.DeepEqual(restoredValues, sortedValues) {
		t.Errorf("Priority queue restore order is incorrect: Expected %v, got %v", values, restoredValues)
	}
}

func assendingSortIntSlice(input []int) []int {
	sortedValues := make([]int, len(input))
	copy(sortedValues, input)
	sort.Ints(sortedValues)
	return sortedValues
}
