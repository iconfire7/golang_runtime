package main

import (
	"flag"
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"
)

type Result struct {
	GoroutineID int
	Tm          int
}

func main() {
	n := flag.Int("n", 0, "количество горутин (N)")
	m := flag.Int("m", 0, "максимальное время сна в миллисекундах (M)")
	flag.Parse()

	if *n <= 0 || *m < 0 {
		fmt.Println("Ошибка: N должно быть > 0, M должно быть >= 0")
		fmt.Println("Пример запуска: go run lru.go -n=5 -m=1000")
		return
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	results := make([]Result, *n)
	var wg sync.WaitGroup
	wg.Add(*n)

	for i := 0; i < *n; i++ {
		id := i

		sleepMs := 0
		if *m > 0 {
			sleepMs = r.Intn(*m + 1)
		}

		results[id] = Result{
			GoroutineID: id,
			Tm:          sleepMs,
		}

		go func(idx, ms int) {
			defer wg.Done()
			time.Sleep(time.Duration(ms) * time.Millisecond)
		}(id, sleepMs)
	}

	wg.Wait()

	sort.Slice(results, func(i, j int) bool {
		return results[i].Tm > results[j].Tm
	})

	for _, res := range results {
		fmt.Printf("<%d, %d>\n", res.GoroutineID, res.Tm)
	}
}
