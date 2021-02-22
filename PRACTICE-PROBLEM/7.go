package main

import (
	"fmt"
)

func main() {
	fmt.Println("start main")
	ch := make(chan string)

	ch <- "hello" //block beacuse wait for receving
	<-ch
	// never execute
	fmt.Println("end of main")
}
