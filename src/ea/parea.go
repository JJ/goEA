package ea

import (
//		"fmt"
	//	"syscall"
//	"math/rand"
	"sort"
)

func (s *ParCEvals) Run(reportSolution func(res TPoolCEvalsResult))  {
	//	workers := eCount + rCount
	eJobs := make(chan IDo, 1)
	rJobs := make(chan IDo, 1)
	eResults := make(chan TIndsEvaluated, 1)
	rResults := make(chan TPopulation, 1)

	control1 := make(chan struct{}, 1)
	control2 := make(chan struct{}, 1)
	p2Eval := NewEvalPool()
	p2Eval.Assign(s.GetPopulation())
	// Siempre estarán ordenados: de mayor a menor.
	p2Rep := NewRepPool()
	doJobs := func(jobs chan IDo) {
		for job := range jobs {
			job.Do()
		}
	}
	addeJobs := func() {
		var qf TQualityF = func(v int) bool { return false }
		var df Tdo = func(i TIndEval) {}
		active := true
		for active {
			select {
			case <-control1:
				active = false
			case eJobs <- EJob{p2Eval.ExtractElements(s.MSizeEvals), s.FitnessF, qf, df, eResults}:
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
			case rJobs <- RJob{p2Rep.ExtractElements(s.MSizeReps), s.PMutation, rResults}:
			}
		}
		close(rJobs)
	}

	bestSolution := NewIndEval()
	waitAndProcessResults := func() {
		ce := 0
		for ce < s.CEvals {
			select { // Blocking
			case indEvals := <-eResults:
				if indEvals != nil && len(indEvals) > 0 {
					if bestSolution.Fitness < indEvals[0].Fitness {
						bestSolution = &indEvals[0]
					}
					p2Rep.Append(indEvals)
					ce += len(indEvals)
				}

			case nInds := <-rResults:
				if nInds != nil && len(nInds) > 0 {
					p2Eval.Append(nInds)
				}
			}
		}
		control1 <- struct{}{}
		control2 <- struct{}{}
		// fmt.Println("The End!")
		reportSolution(TPoolCEvalsResult{*bestSolution, ce})
	}

	for i := 0; i < s.CEvaluators; i++ {
		go doJobs(eJobs)
	}
	for i := 0; i < s.CReproducers; i++ {
		go doJobs(rJobs)
	}
	go addeJobs()
	go addrJobs()
	waitAndProcessResults()
}

func (s *ParFitnessQuality) Run(reportSolution func(res TPoolFitnessQualityResult))  {
	//	workers := eCount + rCount
	eJobs := make(chan IDo, 1)
	rJobs := make(chan IDo, 1)
	eResults := make(chan TIndsEvaluated, 1)
	rResults := make(chan TPopulation, 1)

	control1 := make(chan struct{}, 1)
	control2 := make(chan struct{}, 1)
	p2Eval := NewEvalPool()
	p2Eval.Assign(s.GetPopulation())
	// Siempre estarán ordenados: de mayor a menor.
	p2Rep := NewRepPool()
	doJobs := func(jobs chan IDo) {
		for job := range jobs {
			job.Do()
		}
	}
	alcanzadaSolucion := false
	bestSolution := NewIndEval()
	Do := func(ind TIndEval) {
		bestSolution = &ind
		alcanzadaSolucion = true
	}
	ce := 0
	addeJobs := func() {
		active := true
		for active {
			select {
			case <-control1:
				active = false
			case eJobs <- EJob{p2Eval.ExtractElements(s.MSizeEvals), s.FitnessF, s.QualityF, Do, eResults}:
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
			case rJobs <- RJob{p2Rep.ExtractElements(s.MSizeReps), s.PMutation, rResults}:
			}
		}
		close(rJobs)
	}

	waitAndProcessResults := func() {
		for !alcanzadaSolucion {
			select { // Blocking
			case indEvals := <-eResults:
				if indEvals != nil && len(indEvals) > 0 {
					p2Rep.Append(indEvals)
					ce += len(indEvals)
				}

			case nInds := <-rResults:
				if nInds != nil && len(nInds) > 0 {
					p2Eval.Append(nInds)
				}
			}
		}
		control1 <- struct{}{}
		control2 <- struct{}{}
		// fmt.Println("The End!")
		reportSolution(TPoolFitnessQualityResult{*bestSolution, ce})
	}

	for i := 0; i < s.CEvaluators; i++ {
		go doJobs(eJobs)
	}
	for i := 0; i < s.CReproducers; i++ {
		go doJobs(rJobs)
	}
	go addeJobs()
	go addrJobs()
	waitAndProcessResults()
}

func (job EJob) Do() {
	if job.Population != nil && len(job.Population) > 0 {
		IndEvals := Evaluate(job.Population, job.FitnessF, job.QualityF, job.DoFunc)
		sort.Sort(IndEvals)
		job.results <- IndEvals
	} else {
		job.results <- nil
	}
}

func (job RJob) Do() {
	if job.IndEvals != nil && len(job.IndEvals) > 0 {
		reproductionResults := Reproduce(job.IndEvals, job.PMutation)
		job.results <- reproductionResults
	} else {
		job.results <- nil
	}
}
