package eatest

import (
	"fmt"
	//	"sync"
	"time"
	//	"syscall"
)

var cJobsdo = 0
var cantSee = 50

type TPopulation []int
type EJob struct {
	Population TPopulation
	results    chan <- TPopulation
}
type RJob struct {
	Population TPopulation
	results    chan <- TPopulation
}
type IDo interface {
	Do()
}

// PoolManager is the gorutine for control de workers. The island manager.
func TestParAlg(population TPopulation,
	eCount int, rCount int,
	mSizeEval int, mSizeRep int,
	cEvals int, res chan <- int) {

	//	workers := eCount + rCount
	eJobs := make(chan IDo, 1)
	rJobs := make(chan IDo, 1)
	eResults := make(chan TPopulation, 1)
	rResults := make(chan TPopulation, 1)

	control := make(chan struct{}, 1)
	p2Eval := make(TPopulation, len(population))
	copy(p2Eval, population)

	p2Rep := make(TPopulation, 0)

	selPop2Eval := func() TPopulation {
		nSend2Eval := mSizeEval
		if len(p2Eval) < nSend2Eval {
			nSend2Eval = len(p2Eval)
		}
		res := TPopulation(nil)
		if len(p2Eval) > 0 {
			res = make(TPopulation, nSend2Eval)
			copy(res, p2Eval)
			//			res = append(TPopulation{}, p2Eval[:nSend2Eval]...)
			p2Eval = p2Eval[nSend2Eval:]
			fmt.Println("Sel4 E Job", res)
		}
		return res
	}
	selPop2Rep := func() TPopulation {
		nSend2Rep := mSizeRep
		if len(p2Rep) < nSend2Rep {
			nSend2Rep = len(p2Rep)
		}
		res := TPopulation(nil)
		if len(p2Rep) > 0 {
			res = make(TPopulation, nSend2Rep)
			copy(res, p2Rep)
			//			res = append(TPopulation{}, p2Rep[:nSend2Rep]...)
			p2Rep = p2Rep[nSend2Rep:]
			fmt.Println("Sel4 R Job", res)
		}
		return res
	}
	doJobs := func(jobs chan IDo) {
		for job := range jobs {
			job.Do()
		}
	}
	//	addJobsCEvals := func() {
	//		active := true
	//		for active {
	//			select {
	//			case <-control:
	//				active = false
	//			case eJobs <- EJob{selPop2Eval(), eResults}:
	//			case rJobs <- RJob{selPop2Rep(), rResults}:
	//			}
	//		}
	//		close(eJobs)
	//		close(rJobs)
	//	}
	bestSolution := -1
	waitAndProcessResults := func() {
		for ce := cEvals; ce > 0; {
			if cJobsdo < cantSee {
				cJobsdo++
				fmt.Println("Cant de p2Eval:", len(p2Eval), "Cant de p2Rep:", len(p2Rep))
			}
			select { // Blocking
			case eJobs <- EJob{selPop2Eval(), eResults}:
			case rJobs <- RJob{selPop2Rep(), rResults}:
			case indEvals := <-eResults:
				//				fmt.Println("evals")
				if indEvals != nil && len(indEvals) > 0 {
					if bestSolution < indEvals[0] {
						bestSolution = indEvals[0]
					}
					p2Rep = append(p2Rep, indEvals...)
					fmt.Println("Evaluation arrived:", indEvals)
					ce -= len(indEvals)
				}

			case nInds := <-rResults:
				//				fmt.Println("reps")
				if nInds != nil && len(nInds) > 0 {
					fmt.Println("R rep:", nInds)
					p2Eval = append(p2Eval, nInds...)
				}
			}
		}
		control <- struct{}{}

		fmt.Println("En la lista p2Rep hay:", len(p2Rep), "elementos.")
		res <- bestSolution
	}
	for i := 0; i < eCount; i++ {
		go doJobs(eJobs)
	}
	for i := 0; i < rCount; i++ {
		go doJobs(rJobs)
	}
	//	go addJobsCEvals()
	waitAndProcessResults()
}

func (job EJob) Do() {
	if job.Population != nil {
		fmt.Println("E Job done:", len(job.Population))
		res := make(TPopulation, len(job.Population))
		for i, v := range job.Population {
			res[i] = v*2
		}
		job.results <- res
	}else {
		finish := time.After(time.Duration(5))
		<-finish
		job.results <- nil
	}
}

func (job RJob) Do() {
	if job.Population != nil {
		fmt.Println("R Job done:", len(job.Population))
		res := make(TPopulation, len(job.Population))
		for i, v := range job.Population {
			res[i] = v+1
		}
		job.results <- res
	}else {
		finish := time.After(time.Duration(3))
		<-finish
		job.results <- nil
	}}
