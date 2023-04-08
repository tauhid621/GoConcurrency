package main

import (
	"fmt"
	"sync"
	"time"
)

type UnisexBathroom interface {
	AcquireMaleLock()
	ReleaseMaleLock()
	AcquireFemaleLock()
	ReleaseFemaleLock()
}

type UnisexBathroomLock struct {
	numMale    int
	numFemale  int
	maleLock   sync.Mutex
	femaleLock sync.Mutex
	exLock     sync.Mutex
	gateLock   sync.Mutex
	space      chan struct{}
}

func (ubl *UnisexBathroomLock) AcquireMaleLock() {
	ubl.gateLock.Lock()
	ubl.maleLock.Lock()
	ubl.numMale++
	if ubl.numMale == 1 {
		ubl.exLock.Lock()
	}

	ubl.maleLock.Unlock()
	ubl.gateLock.Unlock()

	ubl.space <- struct{}{}
}

func (ubl *UnisexBathroomLock) ReleaseMaleLock() {
	ubl.maleLock.Lock()
	ubl.numMale--

	<-ubl.space

	if ubl.numMale == 0 {
		ubl.exLock.Unlock()
	}

	ubl.maleLock.Unlock()
}

func (ubl *UnisexBathroomLock) AcquireFemaleLock() {
	ubl.gateLock.Lock()
	ubl.femaleLock.Lock()
	ubl.numFemale++
	if ubl.numFemale == 1 {
		ubl.exLock.Lock()
	}

	ubl.femaleLock.Unlock()
	ubl.gateLock.Unlock()

	ubl.space <- struct{}{}
}

func (ubl *UnisexBathroomLock) ReleaseFemaleLock() {
	ubl.femaleLock.Lock()
	ubl.numFemale--

	<-ubl.space

	if ubl.numFemale == 0 {
		ubl.exLock.Unlock()
	}

	ubl.femaleLock.Unlock()
}

func maleUseBathroom(rId int, ubl *UnisexBathroomLock, wg *sync.WaitGroup) {
	defer wg.Done()

	ubl.AcquireMaleLock()
	defer ubl.ReleaseMaleLock()

	fmt.Printf("Male %d using the bathroom\n", rId)
	time.Sleep(200 * time.Millisecond)
	fmt.Printf("Male %d releasing the bathroom\n", rId)
}

func femaleUseBathroom(wId int, ubl *UnisexBathroomLock, wg *sync.WaitGroup) {
	defer wg.Done()

	ubl.AcquireFemaleLock()
	defer ubl.ReleaseFemaleLock()

	fmt.Printf("Female %d using the bathroom\n", wId)
	time.Sleep(200 * time.Millisecond)
	fmt.Printf("Female %d releasing the bathroom\n", wId)
}

func main() {
	fmt.Println("Yolo")

	ubl := &UnisexBathroomLock{
		space: make(chan struct{}, 3),
	}

	var wg sync.WaitGroup

	for i := 0; i < 20; i++ {
		wg.Add(1)
		thId := i
		go femaleUseBathroom(thId, ubl, &wg)
	}

	for i := 0; i < 20; i++ {
		wg.Add(1)
		thId := i
		go maleUseBathroom(thId, ubl, &wg)
	}

	wg.Wait()
}
