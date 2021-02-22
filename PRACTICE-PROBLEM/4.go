package main

import (
	"fmt"
	"time"
)

var (
	ch = make(chan bool)
)

func count(thing string) {
	for i := 0; i < 5; i++ {
		fmt.Println(thing)
		time.Sleep(time.Second)
	}

}

func main() {
	go func() {
		count("sheep")
		ch <- true
	}()
	<-ch
	fmt.Println("end of main")
}
