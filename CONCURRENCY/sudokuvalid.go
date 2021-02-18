package main

import (
	"fmt"
	"sync"
)

func isValidSlice(slice []int, results chan<- bool, wg *sync.WaitGroup) {
	defer wg.Done()
	seen := make(map[int]bool)
	for _, val := range slice {
		if seen[val] {
			results <- false
		} else {
			seen[val] = true
		}
	}

	results <- true
}

func isValidSudoku(board [][]int, wg *sync.WaitGroup) bool {
	c := make(chan bool, 55)

	// Check every row
	for x := 0; x < 9; x++ {
		wg.Add(1)
		go isValidSlice(append([]int{}, board[x]...), c, wg)
	}

	//column check
	for x := 0; x < 9; x++ {
		slice := make([]int, 0)
		for y := 0; y < 9; y++ {
			slice = append(slice, board[y][x])
		}
		wg.Add(1)
		go isValidSlice(slice, c, wg)
	}

	// Check every 3x3 block
	for x := 0; x <= 6; x += 3 {
		for y := 0; y <= 6; y += 3 {
			blockdigits := append([]int{}, board[x][y:y+3]...)
			blockdigits = append(blockdigits, board[x+1][y:y+3]...)
			blockdigits = append(blockdigits, board[x+2][y:y+3]...)
			wg.Add(1)

			go isValidSlice(blockdigits, c, wg)
		}
	}
	go func() {
		wg.Wait()
		close(c)
	}()
	for result := range c {
		if !result {
			return false
		}
	}

	return true
}

func main() {
	var wg1 sync.WaitGroup
	board := [][]int{
		{5, 3, 4, 6, 7, 8, 9, 1, 2},
		{6, 7, 2, 1, 9, 5, 3, 4, 8},
		{1, 9, 8, 3, 4, 2, 5, 6, 7},
		{8, 5, 9, 7, 6, 1, 4, 2, 3},
		{4, 2, 6, 8, 5, 3, 7, 9, 1},
		{7, 1, 3, 9, 2, 4, 8, 5, 6},
		{9, 6, 1, 5, 3, 7, 2, 8, 4},
		{2, 8, 7, 4, 1, 9, 6, 3, 5},
		{3, 4, 5, 2, 8, 6, 1, 7, 9},
	}
	if isValidSudoku(board, &wg1) {
		fmt.Println("valid sudoku")
	} else {
		fmt.Println("Invalid sudoku")
	}
	var wg2 sync.WaitGroup
	board = [][]int{
		{3, 2, 1, 4, 5, 6, 7, 8, 9},
		{2, 3, 4, 5, 6, 7, 8, 9, 1},
		{3, 4, 5, 6, 7, 8, 9, 1, 2},
		{4, 5, 6, 7, 8, 9, 1, 2, 3},
		{5, 6, 7, 8, 9, 1, 2, 3, 4},
		{6, 7, 8, 9, 1, 2, 3, 4, 5},
		{7, 8, 9, 1, 2, 3, 4, 5, 6},
		{8, 9, 1, 2, 3, 4, 5, 6, 7},
		{9, 1, 2, 3, 4, 5, 6, 7, 8},
	}
	if isValidSudoku(board, &wg2) {
		fmt.Println("valid sudoku")
	} else {
		fmt.Println("Invalid sudoku")
	}
}
