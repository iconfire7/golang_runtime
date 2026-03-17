package main

import (
	"fmt"
	"os"
	"strconv"
)

func generator(k, n int) <-chan int {
	ch := make(chan int)

	go func() {
		defer close(ch)

		for i := k; i <= n; i++ {
			ch <- i
		}
	}()

	return ch
}

func square(inp <-chan int) <-chan int {
	ch := make(chan int)

	go func() {
		defer close(ch)

		for n := range inp {
			ch <- n * n
		}
	}()

	return ch
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Передайте аргументы K и N при вызове программы.")
		fmt.Println("Пример: go run main.go 2 5")
		return
	}

	k, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("Аргументы K и N должны быть целыми числами")
		return
	}
	n, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("Аргументы K и N должны быть целыми числами")
		return
	}

	if k > n {
		fmt.Println("K должен быть больше N")
		return
	}

	if k < 0 || n < 0 {
		fmt.Println("K и N натуральные числа")
		return
	}

	ch1 := generator(k, n)

	ch2 := square(ch1)

	for v := range ch2 {
		fmt.Println(v)
	}
}
