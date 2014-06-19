package eatest

import (
	"fmt"
	"sync"
	//	"syscall"
)

var cJobsdo = 0
var cantSee = 10

type TPopulation []int
type EJob struct {
	Population TPopulation
	results    chan <- IRes
}
type RJob struct {
	Population TPopulation
	results    chan <- IRes
}
type IRes interface {
	Pop() TPopulation
	Kind() int
}
type REval struct {
	Population TPopulation
}
type RRep struct {
	Population TPopulation
}

func (r REval) Pop() TPopulation {
	return r.Population
}
func (r RRep) Pop() TPopulation {
	return r.Population
}
func (r REval) Kind() int {
	return 1
}
func (r RRep) Kind() int {
	return 2
}

// PoolManager is the gorutine for control de workers. The island manager.
func TestParAlg(population TPopulation,
	eCount int, rCount int,
	mSizeEval int, mSizeRep int,
	cEvals int, res chan <- int) {

	//	workers := eCount + rCount
	eJobs := make(chan EJob, 1)
	rJobs := make(chan RJob, 1)
	results := make(chan IRes, 1)

	control := make(chan struct{}, 1)
	p2Eval := make(TPopulation, len(population))
	copy(p2Eval, population)

	p2Rep := make(TPopulation, 0)
	var mp2Eval sync.Mutex
	var mp2Rep sync.Mutex

	selPop2Eval := func() TPopulation {
		nSend2Eval := mSizeEval
		mp2Eval.Lock()
		if len(p2Eval) < nSend2Eval {
			nSend2Eval = len(p2Eval)
		}
		res := TPopulation(nil)
		if len(p2Eval) > 0 {
			res = append(TPopulation{}, p2Eval[:nSend2Eval]...)
			p2Eval = p2Eval[nSend2Eval:]
			fmt.Println("Sel4 E Job", res)
		}
		mp2Eval.Unlock()
		return res
	}
	selPop2Rep := func() TPopulation {
		nSend2Rep := mSizeRep
		mp2Rep.Lock()
		if len(p2Rep) < nSend2Rep {
			nSend2Rep = len(p2Rep)
		}
		res := TPopulation(nil)
		if len(p2Rep) > 0 {
			res = append(TPopulation{}, p2Rep[:nSend2Rep]...)
			p2Rep = p2Rep[nSend2Rep:]
			fmt.Println("Sel4 R Job", res)
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
	//	_ = selPop2Rep
	//	_ = doRepJobs
	addJobsCEvals := func() {
		active := true
		for active {
			select {
			case <-control:
				active = false
			case eJobs <- EJob{selPop2Eval(), results}:
			case rJobs <- RJob{selPop2Rep(), results}:
			}
		}
		close(eJobs)
		close(rJobs)
	}
	bestSolution := -1
	waitAndProcessResults := func() {
		for ce := cEvals; ce > 0; {
			select { // Blocking
			case res := <-results:
				if res != nil {
					if res.Kind() == 1 {
						indEvals := res.Pop()
						mp2Rep.Lock()
						if bestSolution < indEvals[0] {
							bestSolution = indEvals[0]
						}
						p2Rep = append(p2Rep, indEvals...)
						fmt.Println("Evaluation arrived:", indEvals)
						ce -= len(indEvals)
						mp2Rep.Unlock()
					}else {
						nInds := res.Pop()
						mp2Eval.Lock()
						fmt.Println("R rep:", nInds)
						p2Eval = append(p2Eval, nInds...)
						mp2Eval.Unlock()
					}
				}
			}
		}

		control <- struct{}{}

		fmt.Println("En la lista p2Rep hay:", len(p2Rep), "elementos.")
		res <- bestSolution
	}
	for i := 0; i < eCount; i++ {
		go doEvalJobs()
	}
	for i := 0; i < rCount; i++ {
		go doRepJobs()
	}
	go addJobsCEvals()
	waitAndProcessResults()
}

func (job EJob) Do() {
	if job.Population != nil {
		fmt.Println("E Job done:", len(job.Population))
		vals := make(TPopulation, len(job.Population))
		for i, v := range job.Population {
			vals[i] = v*2
		}
		res := REval{vals}
		job.results <- res
	}else {
		job.results <- nil
	}
}

func (job RJob) Do() {
	if job.Population != nil {
		fmt.Println("R Job done:", len(job.Population))
		vals := make(TPopulation, len(job.Population))
		for i, v := range job.Population {
			vals[i] = v+1
		}
		res := RRep{vals}
		job.results <- res
	}else {
		job.results <- nil
	}}
