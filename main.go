package main

import (
	"fmt"
	"sync"
)

func executeParallel(ch chan<- int, functions ...func() int) {
	var wg sync.WaitGroup

	//Create a buffered chanel to store the results
	resultsCh := make(chan int, len(functions))

	//Launch goroutine for each function
	for _, f := range functions {
		wg.Add(1)
		go func(f func() int) {
			defer wg.Done()
			result := f()
			resultsCh <- result
		}(f)
	}
	// Close resultsCh after all goroutine are done
	go func() {
		wg.Wait()
		close(resultsCh)
	}()
	// Read results from the buffered channel and send them to the main channel
	for result := range resultsCh {
		ch <- result
	}
	close(ch)
}

func exampleFunction(counter int) int {
	sum := 0
	for i := 0; i < counter; i++ {
		sum += 1
	}
	return sum
}

func main() {
	expensiveFuntion := func() int { return exampleFunction(200000000) }
	cheapFunction := func() int { return exampleFunction(100000000) }

	ch := make(chan int)

	go executeParallel(ch, expensiveFuntion, cheapFunction)

	for result := range ch {
		fmt.Printf("Result: %d\n", result)
	}
}
