package priority_queue

import (
	"container/heap"
)

func RestoreOrder(reorderedChan <-chan Item) <-chan Item {
	sortedChan := make(chan Item)

	go func() {
		defer close(sortedChan)

		for {
			tmpItem, ok := <-reorderedChan
			if !ok {
				return
			}

			// do some sorting

			sortedChan <- tmpItem
		}
	}()

	return sortedChan
}

func minHeapSort(inputChan <-chan Item) <-chan Item {
	outputChan := make(chan Item)

	go func() {
		defer close(outputChan)
		currentMinIndex := 0

		pq := make(PriorityQueue, 0)
		heap.Init(&pq)

		for {
			tmpItem, ok := <-inputChan
			if !ok {
				return
			}

			if tmpItem.index == currentMinIndex {
				outputChan <- tmpItem
				currentMinIndex++
				continue
			}

			if pq.Len() == 0 {
				heap.Push(&pq, &tmpItem)
				continue
			}

			minItem := pq.Pop().(*Item)
			if minItem.index == currentMinIndex {
				outputChan <- *minItem
				currentMinIndex++
			} else {
				heap.Push(&pq, minItem)
			}

			heap.Push(&pq, &tmpItem)
		}

	}()

	return outputChan
}
