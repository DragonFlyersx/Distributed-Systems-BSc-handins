package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	philoCount := 1
	forks := make([]Fork, 2)
	for index := range forks {
		forks[index].occupationCh = make(chan int, 2)
	}

	philosophers := make([]Philosopher, 1)
	for index := range philosophers {
		philosophers[index].id = index
		philosophers[index].eatCount = 0
		//p.Rightfork = forks[index]
		//p.Leftfork = forks[index+1]
		fmt.Print("Eat is about to be called")

		// Add philosopher to waitgroup, so that it knows to wait for philo to finish
		wg.Add(1)
		go philosophers[index].Eat(&forks[index], &forks[(index+1)%philoCount], &wg)
		// Wait here
		fmt.Print("\nP has eaten this much: ", philosophers[index].eatCount)
	}
	wg.Wait()
	fmt.Print("\nFinish")

}

// 1. Each philosopher must eat 3 times
// 2. A philosopher can only eat if he has both forks
// 3. A philosopher can only take one fork at a time
// 4. A philosopher can only take a fork if it is not being used
// 5. A philosopher can only eat if he has both forks
// 6. A philosopher can only put down a fork if he has it
// Mutex Asynchronous and Synchronous
// must communicate through channels

// Possible goroutine is that they lock themselves when communicated frm philosopher
// Possible ways to avoid deadlock could be randomizing intervals or saying that all but one philosopher attempt to pick up right fork and one picks up left

type Philosopher struct {
	id       int
	eatCount int
	// LeftFork  chan Fork // Use a channel to request the left fork
	// RightFork chan Fork // Use a channel to request the right fork
}

func (philo *Philosopher) think() { // think delay 1-2
	// wait for forks are ready
	// receive from channel whether fork is being used by check whether it is locked?
	// recieve fork here from channel lock funch EAT

	fmt.Println("Philosopher", philo.id, "is thinking")
	time.Sleep(5 * time.Second)

}

func (philo *Philosopher) Eat(Fork1 *Fork, Fork2 *Fork, wg *sync.WaitGroup) {
	defer wg.Done() // Mark as done for the waitgroup when finished
	// receive forks from channel
	// lock forks
	fmt.Print("\nAbout to send value to channels")
	Fork1.occupationCh <- 1
	Fork2.occupationCh <- 1

	fmt.Print("\nBefore fork go routines")
	go Fork1.UseFork()
	go Fork2.UseFork()
	fmt.Print("\nAfter fork go routines")

	forkOcuppationOne := <-Fork1.occupationCh
	forkOcuppationTwo := <-Fork2.occupationCh

	if forkOcuppationOne == 0 && forkOcuppationTwo == 0 {
		philo.eatCount++
		fmt.Println("Philosopher", philo.id, "is eating", "he has eaten", philo.eatCount, "times")
	}
	// unlock forks
	// send them in channel

	// He has eaten
	Fork1.occupationCh <- 0
	Fork2.occupationCh <- 0

	philo.think()
}

type Fork struct {
	id           int
	isOcuppied   int
	occupationCh chan int // 0 = unoccupied, 1 == occupied
	//lock sync.Mutex
}

func (fork *Fork) UseFork() {
	var order int
	// State stores 0 or 1 depending on if the fork is currently occupied
	fmt.Print("UseFork() was called")

	order = <-fork.occupationCh // philosopher has sent out either a 0 or 1 - 0 = placed fork down, 1 = want to pick up

	switch order {
	case 0:
		// if 0 philosopher drops the fork, and state of fork is changed to 0
		fork.isOcuppied = 0
	case 1:
		// if fork is occupied, send value 1 into the channel to let the philosopher know
		if fork.isOcuppied == 1 {
			fork.occupationCh <- 1 // Send a value from channel
		} else {
			// if fork isn't occupied, send value 0 into the channel to let the philosopher know
			fork.occupationCh <- 0 // Send a value from channel
			fork.isOcuppied = 1
		}

	}
}
