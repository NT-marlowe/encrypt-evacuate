package priority_queue

import (
	"container/heap"
	// "log"
)

func RestoreOrder(reorderedChan <-chan Item) <-chan Item {
	return minHeapSort(reorderedChan)
}

func minHeapSort(inputChan <-chan Item) <-chan Item {
	outputChan := make(chan Item)

	go func() {
		defer close(outputChan)
		currentMinIndex := 0

		pq := make(PriorityQueue, 0)
		heap.Init(&pq)

		for {
			select {
			case tmpItem, ok := <-inputChan:
				copyItem := tmpItem
				if !ok {
					return
				}
				// log.Printf("copyItem: %v", copyItem)

				if copyItem.index == currentMinIndex {
					outputChan <- copyItem
					// log.Printf("chan, copy: %v was written to outputChan, minIndex = %d\n", copyItem, currentMinIndex)
					currentMinIndex++
					continue
				}

				if pq.Len() == 0 {
					heap.Push(&pq, &copyItem)
					continue
				}

				minItem := pq.Pop().(*Item)
				// log.Printf("minItem.index: %d, currentMinIndex: %d", minItem.index, currentMinIndex)
				if minItem.index == currentMinIndex {
					outputChan <- *minItem
					// log.Printf("chan, pop: %v was written to outputChan, minIndex = %d\n", copyItem, currentMinIndex)
					currentMinIndex++
				} else {
					heap.Push(&pq, minItem)
				}

				heap.Push(&pq, &copyItem)

			default:
				// log.Printf("minIndex: %d, pq.Len(): %d", currentMinIndex, pq.Len())
				if pq.Len() == 0 {
					continue
				}

				minItemDefault := pq.Pop().(*Item)
				// log.Printf("minItemDefault.index: %d, currentMinIndex: %d", minItemDefault.index, currentMinIndex)
				if minItemDefault.index == currentMinIndex {
					outputChan <- *minItemDefault
					// log.Printf("default: %v was written to outputChan, minIndex = %d\n", minItemDefault, currentMinIndex)
					currentMinIndex++
				} else {
					heap.Push(&pq, minItemDefault)
				}
			}
		}

	}()

	return outputChan
}
