package priority_queue

import (
	"log"
	"reflect"
	"testing"
)

func TestRestoreOrder(t *testing.T) {
	values := []int{0, 9, 4, 5, 6, 3, 8, 1, 7, 2}
	reorderedChan := make(chan Item)

	// reorder the values
	for _, value := range values {
		go func() {
			reorderedChan <- Item{value: value, index: value}
		}()
	}

	restoredChan := RestoreOrder(reorderedChan)
	var restoredValues []int

	i := 0
	for item := range restoredChan {
		t.Logf("Restored value: %v", item.value)
		restoredValues = append(restoredValues, item.value.(int))

		log.Printf("i: %d, len(values): %d", i, len(values))
		if i == len(values) {
			close(reorderedChan)
		}
		i++
	}

	if !reflect.DeepEqual(restoredValues, values) {
		t.Errorf("Priority queue restore order is incorrect: Expected %v, got %v", values, restoredValues)
	}
}
