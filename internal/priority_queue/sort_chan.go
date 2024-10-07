package priority_queue

func RestoreOrder(reorderedChan <-chan Item) <-chan Item {
	sortedChan := make(chan Item)

	go func() {
		defer close(sortedChan)

		for {
			item, ok := <-reorderedChan
			if !ok {
				return
			}

			// do some sorting

			sortedChan <- item
		}
	}()

	return sortedChan
}
