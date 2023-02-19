package main

import (
	api "minitwit-backend/init/Api"
	simulator "minitwit-backend/init/Simulator"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		api.Start()
		defer wg.Done()
	}()
	go func() {
		simulator.Start()
		defer wg.Done()
	}()

	wg.Wait()
}
