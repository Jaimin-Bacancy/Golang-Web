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
	go count("sheep")
	go count("hello")
	//fmt.Scanln()
	time.Sleep(time.Second)
}
