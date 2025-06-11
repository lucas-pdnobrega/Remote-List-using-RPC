package main

import (
	"fmt"
	"remotelist/pkg"
	"sync"
	"time"
)

func main() {
	client := remotelist.NewClientStub()
	var wg sync.WaitGroup

	fmt.Println("Creating initial lists...")
	for i := 0; i < 3; i++ {
		client.CreateList()
	}

	for i := 0; i < 5; i++ {
		client.Append(0, i+1)
	}

	fmt.Println("Running concurrency edge tests...")

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id, value int) {
			defer wg.Done()
			client.Append(id, value)
		}(i%5, (i+1)*10) // some ids like 3,4 won't exist initially
	}

	for i := -1; i <= 6; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			client.Get(0, idx)
		}(i)
	}

	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			client.Remove(id)
		}(i) // id 3 may not exist
	}

	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			client.RemoveList(id)
		}(i)
	}

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			client.Size(id)
		}(i)
	}

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client.CreateList()
		}()
	}

	wg.Wait()
	time.Sleep(1 * time.Second)

	fmt.Println("\nAll concurrency edge case tests completed.")
}
