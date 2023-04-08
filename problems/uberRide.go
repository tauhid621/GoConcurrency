package main

import (
	"fmt"
	"sync"
	"time"
)

/*
	Uber ride problem

	NOTE: Approach is not exactly correct.

	4 seats
	4 R
	4 D
	2R 2D

	cRec
	cDemo

	toal -> ?

	lock on the count for both

	If you do not find a seat leave the lock and return

	If seat found then we invoke seated

	If count is 4 then we should signal/ invoke drive()

	What about blocking the thread? Should we just allow the threads to continue or block them untill all seats are filled?

*/

type party string

const (
	DEMOCRAT    party = "DEMOCRAT"
	REPUBLICIAN party = "REPUBLICIAN"
)

type Car struct {
	number         int
	numReplubcians int
	numDemocrats   int
	available      bool
	lock           sync.Mutex
}

func (c *Car) seated(par party) {
	fmt.Println("New member seated: ", par)
}

func (c *Car) drive() {
	fmt.Println("Starting drive: ")
	// c.lock.Lock()
	// defer c.lock.Unlock()

	fmt.Printf("Num rep: %d, num demo: %d\n", c.numReplubcians, c.numDemocrats)

	c.available = false
}

func tryToSit(c *Car, par party) bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	ableToSit := false

	if c.available {
		// Check if we can sit

		if par == DEMOCRAT && c.numReplubcians <= 2 {
			c.numDemocrats++
			c.seated(par)
			ableToSit = true
		}

		if par == REPUBLICIAN && c.numDemocrats <= 2 {
			c.numReplubcians++
			c.seated(par)
			ableToSit = true
		}

		if ableToSit && c.numDemocrats+c.numReplubcians == 4 {
			c.drive()
		}
	}

	return ableToSit
}

func main() {
	fmt.Println("Hello")

	car := &Car{
		number:    0,
		available: true,
	}

	for i := 0; i < 4; i++ {
		polNo := i
		go func() {
			result := tryToSit(car, REPUBLICIAN)

			fmt.Printf("Republician %d result: %v\n", polNo, result)
		}()
	}

	for i := 0; i < 4; i++ {
		polNo := i
		go func() {
			result := tryToSit(car, DEMOCRAT)

			fmt.Printf("Democrat %d result: %v\n", polNo, result)
		}()
	}

	time.Sleep(1 * time.Second)
}
