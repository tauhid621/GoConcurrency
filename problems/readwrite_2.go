package main

import (
	"fmt"
	"sync"
	"time"
)

type ReadWriteLocker interface {
	AcquireReadLock()
	ReleaseReadLock()
	AcquireWriteLock()
	ReleaseWriteLock()
}

type ReadWriteLock struct {
	numReaders int
	readMu     sync.Mutex
	writeMu    sync.Mutex
	syncMu     sync.Mutex
}

func (rwl *ReadWriteLock) AcquireReadLock() {
	rwl.syncMu.Lock()
	rwl.readMu.Lock()
	rwl.numReaders++

	if rwl.numReaders == 1 {
		rwl.writeMu.Lock()
	}

	rwl.readMu.Unlock()
	rwl.syncMu.Unlock()
}

func (rwl *ReadWriteLock) ReleaseReadLock() {
	rwl.readMu.Lock()
	rwl.numReaders--

	if rwl.numReaders == 0 {
		rwl.writeMu.Unlock()
	}

	rwl.readMu.Unlock()
}

func (rwl *ReadWriteLock) AcquireWriteLock() {
	rwl.syncMu.Lock()
	rwl.writeMu.Lock()
	rwl.syncMu.Unlock()
}

func (rwl *ReadWriteLock) ReleaseWriteLock() {
	rwl.writeMu.Unlock()
}

func reader(rId int, rwl *ReadWriteLock, wg *sync.WaitGroup) {
	defer wg.Done()

	rwl.AcquireReadLock()
	defer rwl.ReleaseReadLock()

	fmt.Printf("Reader %d using the lock\n", rId)
	time.Sleep(200 * time.Millisecond)
	fmt.Printf("Reader %d releasing the lock\n", rId)
}

func writer(wId int, rwl *ReadWriteLock, wg *sync.WaitGroup) {
	defer wg.Done()

	rwl.AcquireWriteLock()
	defer rwl.ReleaseWriteLock()

	fmt.Printf("Writer %d using the lock\n", wId)
	time.Sleep(1 * time.Second)
	fmt.Printf("Writer %d releasing the lock\n", wId)
}

func main() {
	fmt.Println("Test")

	rwl := &ReadWriteLock{}

	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		thId := i
		go writer(thId, rwl, &wg)
	}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		thId := i
		go reader(thId, rwl, &wg)
	}

	wg.Wait()
}
