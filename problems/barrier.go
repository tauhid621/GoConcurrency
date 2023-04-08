package main

import (
	"fmt"
	"sync"
	"time"
)

type Barrier interface {
	PassBarrier()
}

type goBarrier struct {
	threshold     int
	activeThreads int
	waitChannel   chan struct{} // sync.Cond also be used here
	lock          sync.Mutex
}

func (gb *goBarrier) PassBarrier() {
	needToWait := true
	gb.lock.Lock()

	gb.activeThreads++
	// fmt.Println(gb.activeThreads)

	if gb.activeThreads >= gb.threshold {
		needToWait = false
		gb.activeThreads = 0

		for i := 0; i < gb.threshold-1; i++ {
			gb.waitChannel <- struct{}{}
			// Time out can be added here using select statementt
			// done context as well
		}
	}

	gb.lock.Unlock()

	if needToWait {
		<-gb.waitChannel
	}
}

func crossBarrier(thId int, gb Barrier, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("Thread %d tring to pass barrier\n", thId)

	gb.PassBarrier()

	fmt.Printf("Thread %d passed barrier\n", thId)
}

func main() {
	fmt.Println("Hello")

	gb := &goBarrier{
		threshold:   5,
		waitChannel: make(chan struct{}),
	}

	var wg sync.WaitGroup

	for i := 0; i < 15; i++ {
		wg.Add(1)
		go crossBarrier(i, gb, &wg)
		time.Sleep(200 * time.Millisecond)
	}

	wg.Wait()
}
