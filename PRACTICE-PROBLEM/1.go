package main

import (
	"fmt"
	"time"
)

func count(thing string) {
	for {
		fmt.Println(thing)
		time.Sleep(time.Second)
	}
}

func main() {
	count("sheep")
	count("hello")
}
