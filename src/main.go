package main

import (
	"os"
	"strconv"
	"fmt"
	"time"
	"sync"
)

// Generates a sequence of numberts from 1 to 10.
func generate(ints chan int, name string) {
	for i := 1; i <= 10; i++ {
		fmt.Printf("%v sending %v...\n", name, i)
		ints <- i
	}
	close(ints)
}

// Handles the work for each given number.
// Simulates work by sleeping.
func work(ints chan int, name string) {
	for {
		fmt.Printf("%v waiting for data... \n", name)
		time.Sleep(1000 * time.Millisecond)
		foo, more := <-ints
		if !more {
			fmt.Println("Done!")
			return
		}
		fmt.Printf("%v recieved %v\n", name, foo)
		fmt.Printf("%v processing %v...\n", name, foo)
		time.Sleep(1000 * time.Millisecond)
	}
}

func main() {
	var generatorsAmount int
	var workersAmount int
	var err error

	if len(os.Args) > 1 {
		generatorsAmount, err = strconv.Atoi(os.Args[1])
		if err != nil {
			generatorsAmount = 0
		}
	} else {
		generatorsAmount = 0
	}

	if len(os.Args) > 2 {
		workersAmount, err = strconv.Atoi(os.Args[2])
		if err != nil {
			workersAmount = 0
		}
	} else {
		workersAmount = 0
	}

	var waitGroup sync.WaitGroup
	var ints chan int

	if len(os.Args) > 3 {
		bufferSize, err := strconv.Atoi(os.Args[3])
		if err != nil {
			ints = make(chan int)
		} else {
			ints = make(chan int, bufferSize)
		}
	} else {
		ints = make(chan int)
	}

	var name = 'A'
	for i := 1; i <= generatorsAmount; i, name = i + 1, name + 1 {
		go generate(ints, string(name))
	}

	time.Sleep(1000 * time.Millisecond)

	for i := 1; i <= workersAmount; i, name = i + 1, name + 1 {
		waitGroup.Add(1)
		go func(workerName rune) {
			work(ints, string(workerName))
			defer waitGroup.Done()
		}(name)
	}

	waitGroup.Wait()
}
