package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("start main")
	ch := make(chan string, 2)
	ch <- "hello" //maximum size is 2 so > 2 give deadlock
	ch <- "hello"
	//ch <- "hello"
	time.Sleep(time.Second)
	fmt.Println(<-ch)
	time.Sleep(time.Second)
	fmt.Println(<-ch)
	ch <- "hello"
	fmt.Println(<-ch)
	// time.Sleep(time.Second)
	// fmt.Println(<-ch)
	fmt.Println("end of main")
}
