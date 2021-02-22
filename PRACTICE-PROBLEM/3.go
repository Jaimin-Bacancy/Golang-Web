package main

import (
	"fmt"
	"sync"
	"time"
)

func count(thing string) {
	for i := 0; i < 5; i++ {
		fmt.Println(thing)
		time.Sleep(time.Second)
	}
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		count("sheep")
		wg.Done()
	}()
	wg.Wait()
	fmt.Println("end of main")
}
