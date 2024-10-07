package priority_queue

import (
	"container/heap"
	"log"
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
				if !ok {
					return
				}
				log.Printf("tmpItem: %v", tmpItem)

				if tmpItem.index == currentMinIndex {
					outputChan <- tmpItem
					currentMinIndex++
					log.Printf("%v was written to outputChan, minIndex = %d\n", tmpItem, currentMinIndex)
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
					log.Printf("%v was written to outputChan, minIndex = %d\n", tmpItem, currentMinIndex)
				} else {
					heap.Push(&pq, minItem)
				}

				heap.Push(&pq, &tmpItem)

			default:
				// log.Printf("minIndex: %d, pq.Len(): %d", currentMinIndex, pq.Len())
				if pq.Len() == 0 {
					continue
				}

				minItem := pq.Pop().(*Item)
				if minItem.index == currentMinIndex {
					log.Printf("%v was written to outputChan, minIndex = %d\n", minItem, currentMinIndex)
					outputChan <- *minItem
					currentMinIndex++
				} else {
					heap.Push(&pq, minItem)
				}
			}
		}

	}()

	return outputChan
}
