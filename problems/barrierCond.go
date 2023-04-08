package main

import (
	"fmt"
	"sync"
	"time"
)

type BarrierConder interface {
	PassBarrierCond()
}

type goBarrierCond struct {
	threshold     int
	activeThreads int
	waitChannel   chan struct{} // sync.Cond also be used here
	lock          sync.Mutex
	cond          sync.Cond
}

func (gb *goBarrierCond) PassBarrierCond() {
	needToWait := true
	gb.cond.L.Lock()

	gb.activeThreads++

	if gb.activeThreads >= gb.threshold {
		needToWait = false
		gb.activeThreads = 0

		gb.cond.Broadcast()
	}

	if needToWait {
		gb.cond.Wait()
	}
	gb.cond.L.Unlock()
}

func crossBarrierCond(thId int, gb BarrierConder, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("Thread %d tring to pass barrier\n", thId)

	gb.PassBarrierCond()

	fmt.Printf("Thread %d passed barrier\n", thId)
}

func main() {
	fmt.Println("Hello")

	gb := &goBarrierCond{
		threshold:   5,
		waitChannel: make(chan struct{}),
	}

	gb.cond = *sync.NewCond(&gb.lock)

	var wg sync.WaitGroup

	for i := 0; i < 15; i++ {
		wg.Add(1)
		go crossBarrierCond(i, gb, &wg)
		time.Sleep(200 * time.Millisecond)
	}

	wg.Wait()
}
