package main

import (
	"container/heap"
	"ebpf-ssl/internal/priority_queue"
	// "log"
)

func restoreOrder(reorderedChan <-chan indexedDataBlock) <-chan indexedDataBlock {
	return minHeapSort(reorderedChan)
}

func minHeapSort(inputChan <-chan indexedDataBlock) <-chan indexedDataBlock {
	outputChan := make(chan indexedDataBlock)

	go func() {
		defer close(outputChan)
		currentMinIndex := 0

		pq := make(priority_queue.PriorityQueue, 0)
		heap.Init(&pq)

		for {
			select {
			case tmpItem, ok := <-inputChan:
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

				minItem := pq.Pop().(*indexedDataBlock)
				if minItem.index == currentMinIndex {
					outputChan <- *minItem
					currentMinIndex++
				} else {
					heap.Push(&pq, minItem)
				}

				heap.Push(&pq, &tmpItem)

			default:
				if pq.Len() == 0 {
					continue
				}

				minItem := pq.Pop().(*indexedDataBlock)
				if minItem.index == currentMinIndex {
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
