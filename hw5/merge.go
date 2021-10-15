package main

import (
	"fmt"
	"sync"
	"time"
)

func fillChan(msg ...string) <-chan string {
	ch := make(chan string)

	go func() {
		for i := 0; i < len(msg); i++ {
			ch <- msg[i]
			time.Sleep(3 * time.Second)
		}
		close(ch)
	}()

	return ch
}

func merge(cs ...<-chan string) <-chan string {
	ch := make(chan string)
	wg := new(sync.WaitGroup)
	wg.Add(len(cs))
	for _, c := range cs {

		go func(lc <-chan string) {
			defer wg.Done()

			for in := range lc {
				ch <- in
			}
		}(c)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	return ch
}

func main() {
	mergedChannels := merge(fillChan("hello", "world", "golang"), fillChan("testing"))

	for val := range mergedChannels {
		fmt.Println(val)
	}

}
