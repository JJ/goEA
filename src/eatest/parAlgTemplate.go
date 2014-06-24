package eatest

import (
	//	"fmt"
	//	"sync"
//	"time"
	//	"syscall"
	"sync"
)

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

type evalPool struct {
	pool TPopulation
	mutex *sync.RWMutex
	name string
}

type repPool struct {
	pool TPopulation
	mutex *sync.RWMutex
}

func NewEvalPool(name string) *evalPool {
	return &evalPool{make(TPopulation, 0), new(sync.RWMutex), name}
}

func (self *evalPool) assign(pop TPopulation) {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	self.pool = make(TPopulation, len(pop))
	copy(self.pool, pop)
}

func (self *evalPool) extractElements(mSizeEval int) TPopulation {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	nSend2Eval := mSizeEval
	res := TPopulation(nil)
	if len(self.pool) > 0 {
		if len(self.pool) < nSend2Eval {
			nSend2Eval = len(self.pool)
		}
		res = make(TPopulation, nSend2Eval)
		copy(res, self.pool)
		//			res = append(TPopulation{}, p2Eval[:nSend2Eval]...)
		self.pool = self.pool[nSend2Eval:]
		//		fmt.Println(self.name, nSend2Eval, "-")
	}
	return res
}

func (self *evalPool) append(nInds TPopulation) {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	self.pool = append(self.pool, nInds...)
	//	fmt.Println(self.name, len(nInds), "+")
}

func (self *evalPool) Length() int {
	self.mutex.RLock()
	defer self.mutex.RUnlock()
	return len(self.pool)
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

	control1 := make(chan struct{}, 1)
	control2 := make(chan struct{}, 1)
	//	p2Eval := make(TPopulation, len(population))
	//	copy(p2Eval, population)
	//
	//	p2Rep := make(TPopulation, 0)
	p2Rep := NewEvalPool("repPool")
	p2Eval := NewEvalPool("evalPool")
	p2Eval.assign(population)
	doJobs := func(jobs chan IDo) {
		for job := range jobs {
			job.Do()
		}
	}
	addeJobs := func() {
		active := true
		for active {
			select {
			case <-control1:
				active = false
			case eJobs <- EJob{p2Eval.extractElements(mSizeEval), eResults}:
			}
		}
		close(eJobs)
	}
	addrJobs := func() {
		active := true
		for active {
			select {
			case <-control2:
				active = false
			case rJobs <- RJob{p2Rep.extractElements(mSizeRep), rResults}:
			}
		}
		close(rJobs)
	}
	bestSolution := -1
	waitAndProcessResults := func() {
		for ce := cEvals; ce > 0; {
			select { // Blocking
			case indEvals := <-eResults:
				if indEvals != nil && len(indEvals) > 0 {
					if bestSolution < indEvals[0] {
						bestSolution = indEvals[0]
					}
					p2Rep.append(indEvals)
					ce -= len(indEvals)
				}
			case nInds := <-rResults:
				if nInds != nil && len(nInds) > 0 {
					p2Eval.append(nInds)
				}
			}
		}
		control1 <- struct{}{}
		control2 <- struct{}{}

		//		fmt.Println("En la lista p2Rep hay:", p2Rep.Length(), "elementos.")
		res <- bestSolution
	}
	for i := 0; i < eCount; i++ {
		go doJobs(eJobs)
	}
	for i := 0; i < rCount; i++ {
		go doJobs(rJobs)
	}
	go addeJobs()
	go addrJobs()
	waitAndProcessResults()
}

func (job EJob) Do() {
	if job.Population != nil {
		//		fmt.Println("E Job done:", len(job.Population))
		res := make(TPopulation, len(job.Population))
		for i, v := range job.Population {
			res[i] = v*2
		}
		job.results <- res
	}else {
//		finish := time.After(time.Duration(5))
//		<-finish
		job.results <- nil
	}
}

func (job RJob) Do() {
	if job.Population != nil {
		//		fmt.Println("R Job done:", len(job.Population))
		res := make(TPopulation, len(job.Population))
		for i, v := range job.Population {
			res[i] = v+1
		}
		job.results <- res
	}else {
//		finish := time.After(time.Duration(3))
//		<-finish
		job.results <- nil
	}}
