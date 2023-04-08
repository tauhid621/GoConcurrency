package main

import (
	"fmt"
	"sync"
	"time"
)

type Semaphore struct {
	// maxPermits int
	// permits    int
	channel chan struct{}
}

func NewSemaphore(permits int) *Semaphore {
	return &Semaphore{
		// maxPermits: permits,
		channel: make(chan struct{}, permits),
	}
}

func (s *Semaphore) Acquire() {
	s.channel <- struct{}{}
}

func (s *Semaphore) Release() {
	<-s.channel
}

func acqureSemaphoreAndPrint(thId int, s *Semaphore) {
	s.Acquire()
	defer s.Release()

	fmt.Printf("Thread %d is using semaphore\n", thId)
	time.Sleep(500 * time.Millisecond)
	fmt.Printf("Thread %d is releasing semaphore\n", thId)

}

func main() {
	fmt.Println("Hello")

	s := NewSemaphore(5)
	// done := make(chan struct{})
	var wg sync.WaitGroup

	for i := 0; i < 20; i++ {
		wg.Add(1)
		thId := i
		go func() {
			defer wg.Done()
			acqureSemaphoreAndPrint(thId, s)

			// if thId == 19 {
			// 	done <- struct{}{}
			// }
		}()
	}

	// time.Sleep(5 * time.Second)

	// x := sync.NewSemaphore

	// <-done
	wg.Wait()
}
