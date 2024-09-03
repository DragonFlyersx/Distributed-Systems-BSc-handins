package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {

	forks := make([]Fork, 5)
	for index, f := range forks {
		f.id = index+1;
	}

	philosophers := make([]Philosopher, 5)
	for index, p := range philosophers {
		p.id = index+1;
		p.Rightfork = forks[index+1]
		p.Leftfork = forks[index+2]
	}

	// var mu1,mu2,mu3,mu4,mu5 sync.Mutex
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

// eat lock forks
func Eat(philo Philosopher) {
	// receive forks from channel
	// wait for forks are ready
	// lock forks
	philo.eatCount++
	fmt.Println("Philosopher", philo.id, "is eating", "he has eaten", philo.eatCount, "times")
	// unlock forks
	// send them in channel
}

// think Fjerne lock fork
func think(philo Philosopher) { // think delay 1-2

	// wait for forks are ready
	// receive from channel whether fork is being used by check whether it is locked?
	// recieve fork here from channel lock funch EAT

	fmt.Println("Philosopher", philo.id, "is thinking")
}

type Philosopher() struct {
	id        int
	eatCount  int
	Leftfork  chan bool
	Rightfork chan bool
}

type Fork struct {
	id   int
	lock sync.Mutex
}
