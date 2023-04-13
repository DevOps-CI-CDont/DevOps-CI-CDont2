package main

import (
	_ "minitwit-backend/init/tzinit"

	api "minitwit-backend/init/Api"

	simulator "minitwit-backend/init/Simulator"
	"os"
	"sync"
)

func main() {
	var mode string
	if len(os.Args) > 1 {
		mode = os.Args[1]
	} else {
		mode = "prod"
	}
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		api.Start(mode)
		defer wg.Done()
	}()
	go func() {
		simulator.Start()
		defer wg.Done()
	}()

	wg.Wait()
}
