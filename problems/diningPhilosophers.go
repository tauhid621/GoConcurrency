package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type table struct {
	forkLock    []sync.Mutex
	forkOnTable []bool
}

func contemplate(philId int) {
	fmt.Printf("Philosopher %d contemplating.\n", philId)

	rand.Seed(time.Now().Unix())

	dur := rand.Intn(500)
	time.Sleep(time.Duration(dur) * time.Millisecond)
}

func eat(philId int) {
	fmt.Printf("Philosopher %d eating.\n", philId)

	rand.Seed(time.Now().Unix())

	dur := rand.Intn(1000)
	time.Sleep(time.Duration(dur) * time.Millisecond)
}

func tryToEat(table *table, philId, numPhilosophers int) {
	leftIndex := (numPhilosophers + philId - 1) % numPhilosophers
	rightIndex := (numPhilosophers + philId + 1) % numPhilosophers

	if philId%2 == 0 {
		table.forkLock[leftIndex].Lock()
		table.forkLock[rightIndex].Lock()
	} else {
		table.forkLock[rightIndex].Lock()
		table.forkLock[leftIndex].Lock()
	}

	// if table.forkOnTable[leftIndex] && table.forkOnTable[rightIndex] {
	// 	// we can use the forks
	// 	table.forkOnTable[leftIndex] = false
	// 	table.forkOnTable[rightIndex] = false

	// 	table.forkLock[leftIndex].Unlock()
	// 	table.forkLock[rightIndex].Unlock()

	// 	eat(philId)

	// 	table.forkLock[leftIndex].Lock()
	// 	table.forkLock[rightIndex].Lock()

	// 	table.forkOnTable[leftIndex] = true
	// 	table.forkOnTable[rightIndex] = true

	// }

	eat(philId)

	if philId%2 == 0 {
		table.forkLock[rightIndex].Unlock()
		table.forkLock[leftIndex].Unlock()
	} else {
		table.forkLock[leftIndex].Unlock()
		table.forkLock[rightIndex].Unlock()
	}

}

func philosopherStart(table *table, philId, numPhilosophers int) {

	fmt.Printf("Philosopher %d starting\n", philId)

	for {
		contemplate(philId)
		tryToEat(table, philId, numPhilosophers)
	}
}

func main() {
	fmt.Println("Test")

	var philosophers int = 5

	table := &table{
		forkLock:    make([]sync.Mutex, philosophers),
		forkOnTable: make([]bool, philosophers),
	}

	for i := 0; i < philosophers; i++ {
		table.forkOnTable[i] = true
	}

	for i := 0; i < philosophers; i++ {
		go philosopherStart(table, i, philosophers)
	}

	time.Sleep(10 * time.Second)
}
