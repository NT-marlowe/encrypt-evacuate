package priority_queue

import (
	"log"
	"reflect"
	"testing"
)

func TestRestoreOrder(t *testing.T) {
	// values := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	values := []int{0, 9, 4, 5, 6, 3, 8, 1, 7, 2}

	// inputCh := make(chan Item)
	inputCh := make(chan Item, len(values))

	// reorder the values
	for _, value := range values {
		// go func() {
		// 	inputCh <- Item{value: value, index: value}
		// 	log.Printf("%d was written to inputCh", value)
		// }()
		tmpItem := Item{index: value, value: value}
		inputCh <- tmpItem
		log.Printf("%v was written to inputCh", tmpItem)
	}

	outputCh := RestoreOrder(inputCh)
	var restoredValues []int

	i := 0
	for item := range outputCh {
		t.Logf("Restored value: %v", item.value)
		restoredValues = append(restoredValues, item.value.(int))

		i++
		log.Printf("i: %d, value: %d", i, item.value.(int))
		if i == len(values) {
			close(inputCh)
		}
	}

	if !reflect.DeepEqual(restoredValues, values) {
		t.Errorf("Priority queue restore order is incorrect: Expected %v, got %v", values, restoredValues)
	}
}
