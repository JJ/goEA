package eatest

import (
	"fmt"
	//	"sync"
	"time"
	//	"syscall"
	"sync"
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
		fmt.Println("Extraidos", nSend2Eval, "elementos del", self.name)
	}
	return res
}

func (self *evalPool) append(nInds TPopulation) {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	self.pool = append(self.pool, nInds...)
	fmt.Println("Añadidos", len(nInds), "elementos al", self.name)
}

func (self *evalPool) Length()int {
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

	control := make(chan struct{}, 1)
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
	addJobsCEvals := func() {
		active := true
		for active {
			select {
			case <-control:
				active = false
			case eJobs <- EJob{p2Eval.extractElements(mSizeEval), eResults}:
			case rJobs <- RJob{p2Rep.extractElements(mSizeRep), rResults}:
			}
		}
		close(eJobs)
		close(rJobs)
	}
	bestSolution := -1
	waitAndProcessResults := func() {
		for ce := cEvals; ce > 0; {
			if cJobsdo < cantSee {
				cJobsdo++
				fmt.Println("Cant de p2Eval:", p2Eval.Length(), "Cant de p2Rep:", p2Rep.Length())
			}
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
					//					fmt.Println("R rep:", nInds)
					//					p2Eval = append(p2Eval, nInds...)
					p2Eval.append(nInds)
				}
			}
		}
		control <- struct{}{}

		fmt.Println("En la lista p2Rep hay:", p2Rep.Length(), "elementos.")
		res <- bestSolution
	}
	for i := 0; i < eCount; i++ {
		go doJobs(eJobs)
	}
	for i := 0; i < rCount; i++ {
		go doJobs(rJobs)
	}
	go addJobsCEvals()
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
