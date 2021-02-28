package main

import (
	"log"
	"time"
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
	time.Sleep(100 * time.Millisecond)
}

func main() {
	numArray := makeRange(0, 50)

	start := time.Now()

	for i, _ := range numArray {
		apiCall(i)
	}

	elapsed := time.Since(start)
	log.Printf("Time taken %s", elapsed)
}
