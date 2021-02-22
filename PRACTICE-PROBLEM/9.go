package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("start main")
	ch1 := make(chan string)
	ch2 := make(chan string)
	go func() {
		for {
			ch1 <- "hello"

			time.Sleep(time.Second * 2)

		}

		close(ch1)
	}()

	go func() {
		for {
			ch2 <- "word"

			time.Sleep(time.Second)

		}

		close(ch2)
	}()

	// go func() {
	// 	for {
	// 		fmt.Println(<-ch1)
	// 		fmt.Println(<-ch2)
	// 	}
	// }()
	go func() {
		for {
			select {
			case msg1 := <-ch1:
				fmt.Println(msg1)
			case msg2 := <-ch2:
				fmt.Println(msg2)
			}
		}

	}()

	fmt.Scanln()
	fmt.Println("end of main")
}
