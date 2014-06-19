package eatest

import (
	"fmt"
	"sync"
	//	"syscall"
)

var cJobsdo = 0
var cantSee = 10

type EJob struct {
	Population int
	results    chan <- int
}
type RJob struct {
	Population int
	results    chan <- int
}

// PoolManager is the gorutine for control de workers. The island manager.
func TestParAlg(population []int,
	eCount int, rCount int,
	mSizeEval int, mSizeRep int,
	cEvals int, res chan <- int) {

	//	workers := eCount + rCount
	eJobs := make(chan EJob, 1)
	rJobs := make(chan RJob, 1)
	eResults := make(chan int, 1)
//	rResults := make(chan int, 1)

	control := make(chan struct{}, 1)
	p2Eval := make([]int, len(population))
	copy(p2Eval, population)

	// Siempre estarán ordenados: de mayor a menor.
	p2Rep := make([]int, 0)

	var mp2Eval sync.Mutex
	var mp2Rep sync.Mutex

	selPop2Eval := func() int {
		mp2Eval.Lock()
		res := 0
		if len(p2Eval) > 0 {
			res = p2Eval[0]
			p2Eval = p2Eval[1:]
			fmt.Println("Sel4 EJob", res)
		}
		mp2Eval.Unlock()
		return res
	}
	selPop2Rep := func() int {
		mp2Rep.Lock()
		res := 0
		if len(p2Rep) > 0 {
			res = p2Rep[0]
			p2Rep = p2Rep[1:]
			fmt.Println("Sel4 RJob", res)
		}
		mp2Rep.Unlock()
		return res
	}
	doEvalJobs := func() {
		for job := range eJobs {
			//		fmt.Println("Do job Eval")
			job.Do()
		}
	}

	doRepJobs := func() {
		for job := range rJobs {
			job.Do()
		}
	}
	_ = selPop2Rep
	_ = doRepJobs
	addJobsCEvals := func() {
		active := true
		for active {
			select {
			case <-control:
				active = false
			case eJobs <- EJob{selPop2Eval(), eResults}:
//			case rJobs <- RJob{selPop2Rep(), rResults}:
			}
		}
		close(eJobs)
//		close(rJobs)
	}
	waitAndProcessResults := func() {
		for ce := cEvals; ce > 0; {
			select { // Blocking
			case indEvals := <-eResults:
				mp2Rep.Lock()
				p2Rep = append(p2Rep, indEvals)
				fmt.Println("Evaluation arrived:", indEvals)
				mp2Rep.Unlock()
				ce--
//			case nInds := <-rResults:
//				mp2Eval.Lock()
//				fmt.Println("R rep:", nInds)
//				p2Eval = append(p2Eval, nInds)
//				mp2Eval.Unlock()
			}
		}
		control <- struct{}{}
		fmt.Println("En la lista p2Rep hay:", len(p2Rep), "elementos.")
		res <- p2Rep[len(p2Rep) - 1]
	}
	for i := 0; i < eCount; i++ {
		go doEvalJobs()
	}
//	for i := 0; i < rCount; i++ {
//		go doRepJobs()
//	}
	go addJobsCEvals()
	waitAndProcessResults()
}

func (job EJob) Do() {
	if job.Population != 0 {
		fmt.Println("EJob done:", job.Population*2)
		job.results <- job.Population * 2
	}else {
		job.results <- 0
	}
}

func (job RJob) Do() {
	if job.Population != 0 {
		fmt.Println("RJob done:", job.Population+1)
		job.results <- job.Population + 1
	}else {
		job.results <- 0
	}
}
