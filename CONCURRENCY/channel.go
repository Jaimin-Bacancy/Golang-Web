package main

import (
	"log"
	"time"
)

var (
	ch = make(chan int, 10)
)

func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}

func apiCall(i int) {
	log.Println("API call for", i, "started")
	<-ch
}

func main() {
	numArray := makeRange(0, 1000)

	start := time.Now()

	for i, _ := range numArray {
		ch <- 1
		go apiCall(i)
	}

	for i := 1; i <= 10; i++ {
		ch <- i
	}
	elapsed := time.Since(start)
	log.Printf("Time taken %s", elapsed)
}
