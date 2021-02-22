package main

import (
	"fmt"
	"sync"
	"time"
)

func count(thing string, ch chan string) {
	//defer close(ch)
	for i := 0; i < 5; i++ {
		ch <- thing
		time.Sleep(time.Second)
	}

}

func main() {
	var wg sync.WaitGroup
	ch := make(chan string)
	go func() {
		count("sheep", ch)
	}()

	for {
		msg, open := <-ch
		if !open {
			break
		}
		fmt.Println(msg)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		msg := <-ch
		fmt.Println(msg)
	}()
	ch <- "hello"
	wg.Wait()
	fmt.Println("end of main")
}
