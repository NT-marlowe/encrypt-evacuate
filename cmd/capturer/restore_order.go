package main

import (
	"container/heap"
	"fmt"
	"time"
)

// slice, key: Item.index, value: time.TIme
var enqueueTime = make(map[int]time.Time)

func measureTime(index int) {
	elapsed := time.Since(enqueueTime[index])
	fmt.Printf("minHeapSort: %v\n", elapsed)
	delete(enqueueTime, index)
}

func restoreOrder(reorderedChan <-chan indexedDataBlock) <-chan indexedDataBlock {
	return minHeapSort(reorderedChan)
}

func minHeapSort(inputChan <-chan indexedDataBlock) <-chan indexedDataBlock {
	outputChan := make(chan indexedDataBlock)

	go func() {
		defer close(outputChan)
		currentMinIndex := 0

		pq := make(PriorityQueue, 0)
		heap.Init(&pq)

		for {
			select {
			case tmpItem, ok := <-inputChan:
				enqueueTime[tmpItem.index] = time.Now()

				if !ok {
					return
				}

				if tmpItem.index == currentMinIndex {
					outputChan <- tmpItem
					currentMinIndex++

					measureTime(tmpItem.index)

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

					measureTime(minItem.index)

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

					measureTime(minItem.index)
				} else {
					heap.Push(&pq, minItem)
				}
			}
		}

	}()

	return outputChan
}
