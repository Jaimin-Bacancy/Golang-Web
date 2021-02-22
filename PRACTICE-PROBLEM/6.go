package main

import (
	"fmt"
	"time"
)

func count(thing string, ch chan string) {
	defer close(ch)
	for i := 0; i < 5; i++ {
		ch <- thing
		time.Sleep(time.Second)
	}

}

func main() {
	ch := make(chan string)
	go func() {
		count("sheep", ch)
	}()

	//before this forloop use channel must be close
	//otherwise deadlock occured
	//at receving end
	for msg := range ch {
		fmt.Println(msg)
	}
	//on close channel can not send data
	//ch <- "hello"
	fmt.Println("end of main")
}
