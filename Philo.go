package main

import (
	"fmt"
	"sync"
)

type Philosopher struct {
	id        int
	eatCount  int
	LeftFork  chan bool // Use a channel to request the left fork
	RightFork chan bool // Use a channel to request the right fork
}

func main() {
	var wg sync.WaitGroup
	philoCount := 5
	eatLimit := 3
	forks := make([]chan bool, philoCount) // fork creation
	for i := range forks {
		forks[i] = make(chan bool, 1)
		go fork(i, forks[i])
	}

	philosophers := make([]Philosopher, philoCount)
	for index := range philosophers {
		philosophers[index] = Philosopher{
			id:        index,
			eatCount:  0,
			LeftFork:  forks[index],
			RightFork: forks[(index+1)%philoCount],
		}
		// Add philosopher to waitgroup, so that it knows to wait for philosopher to finish
		wg.Add(1)
		go philosophers[index].Eat(&wg, eatLimit)
	}
	wg.Wait()
	for _, phil := range philosophers {
		if phil.eatCount >= eatLimit {
			fmt.Println("Philosopher ", phil.id, "has finished eating")
		}
	}
}

func (philo *Philosopher) Eat(wg *sync.WaitGroup, eatLimit int) {
	defer wg.Done() // Mark as done for the waitgroup when finished

	for philo.eatCount < eatLimit {

		// request to pickup forks
		philo.LeftFork <- true
		philo.RightFork <- true

		// Simulate eating
		fmt.Println("Philosopher", philo.id, "is eating")

		// Release forks:
		<-philo.LeftFork
		<-philo.RightFork
		philo.eatCount++
		fmt.Println("Philosopher", philo.id, " has eaten", philo.eatCount, "times")

		philo.think()
	}
}

func (philo *Philosopher) think() {
	fmt.Println("Philosopher", philo.id, "is thinking")
	//If you want to wait time wise
	//time.Sleep(time.Microsecond * 100)
}

func fork(id int, forkChan chan bool) { // fork go rutine
	for {
		forkChan <- true
		<-forkChan
	}
}

/*
Why we will never reach a deadlock:
The reason why a deadlock never happens is because we have implemented buffered channels with a size of 1.
Meaning that we allow a philosopher to attempt to pick up a fork but it does not immediately block the fork, in case that it is occupied.

This also means that the channel will block the philosopher's request until the fork is released,
but since the buffer is released as soon as the fork is freed this is not an issue (think of it as an temporary wait or a queue, not an indefinite wait)

*/
