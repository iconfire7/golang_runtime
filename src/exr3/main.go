package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Передайте аргумент K")
		fmt.Println("Пример: go run lru.go 2")
		return
	}

	kArg, err := strconv.ParseUint(os.Args[1], 10, 32)
	if err != nil {
		fmt.Printf("Невалидный K: %v\n", err)
		return
	}
	k := uint(kArg)

	startChan := make(chan os.Signal, 1)
	signal.Notify(startChan, os.Interrupt)

	stopChan := make(chan struct{})

	go func() {
		tkNum := 1
		start := time.Now()

		for {
			select {
			case <-stopChan:
				return
			default:
				time.Sleep(time.Duration(k) * time.Second)

				end := time.Since(start).Seconds()
				fmt.Printf("Tick %d since %.0f\n", tkNum, end)
				tkNum++
			}
		}
	}()

	<-startChan

	close(stopChan)

	fmt.Println("Все завершено")
}
