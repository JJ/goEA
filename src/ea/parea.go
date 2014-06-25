package ea

import (
	// "fmt"
	"sort"
	//	"syscall"
)

func (s *ParCEvals) Run() TIndEval {
	//ch1 := make(chan TIndividual, 1)
	//ch2 := make(chan TIndividual, 1)
	res := make(chan TIndEval)
	go PoolManagerCEvals(
		s.Population,
		s.CEvaluators, s.CReproducers,
		s.MSizeEvals, s.MSizeReps,
		s.PMutation, s.CEvals,
		s.FitnessF,
		res)

	return <-res
}
func (s *ParFitnessQuality) Run() TIndEval {
	//ch1 := make(chan TIndividual, 1)
	//ch2 := make(chan TIndividual, 1)
	res := make(chan TIndEval)
	go PoolManagerFitnessQuality(
		s.Population,
		s.CEvaluators, s.CReproducers,
		s.MSizeEvals, s.MSizeReps,
		s.PMutation, s.FitnessF,
		s.QualityF, s.Do,
		res)

	return <-res
}

// PoolManager is the gorutine for control de workers. The island manager.
func PoolManagerCEvals(population TPopulation,
	eCount int, rCount int,
	mSizeEval int, mSizeRep int,
	pMutation float32, cEvals int,
	ff TFitnessFunc, res chan<- TIndEval) {

	//	workers := eCount + rCount
	eJobs := make(chan IDo, 1)
	rJobs := make(chan IDo, 1)
	eResults := make(chan TIndsEvaluated, 1)
	rResults := make(chan TPopulation, 1)

	control1 := make(chan struct{}, 1)
	control2 := make(chan struct{}, 1)
	p2Eval := NewEvalPool()
	p2Eval.Assign(population)
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
			case eJobs <- EJob{p2Eval.ExtractElements(mSizeEval), ff, qf, df, eResults}:
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
			case rJobs <- RJob{p2Rep.ExtractElements(mSizeRep), pMutation, rResults}:
			}
		}
		close(rJobs)
	}

	bestSolution := NewIndEval()
	waitAndProcessResults := func() {
		for ce := cEvals; ce > 0; {
			select { // Blocking
			case indEvals := <-eResults:
				if indEvals != nil && len(indEvals) > 0 {
					if bestSolution.Fitness < indEvals[0].Fitness {
						bestSolution = &indEvals[0]
					}
					p2Rep.Append(indEvals)
					ce -= len(indEvals)
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
		res <- *bestSolution
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

// PoolManager is the gorutine for control de workers. The island manager.
func PoolManagerFitnessQuality(population TPopulation,
	eCount int, rCount int,
	mSizeEval int, mSizeRep int,
	pMutation float32, ff TFitnessFunc,
	qf TQualityF, df Tdo, res chan<- TIndEval) {

	//	workers := eCount + rCount
	eJobs := make(chan IDo, 1)
	rJobs := make(chan IDo, 1)
	eResults := make(chan TIndsEvaluated, 1)
	rResults := make(chan TPopulation, 1)

	control1 := make(chan struct{}, 1)
	control2 := make(chan struct{}, 1)
	p2Eval := NewEvalPool()
	p2Eval.Assign(population)
	// Siempre estarán ordenados: de mayor a menor.
	p2Rep := NewRepPool()
	doJobs := func(jobs chan IDo) {
		for job := range jobs {
			job.Do()
		}
	}
	alcanzadaSolucion := false
	alcanzadaSolucionF := func(ind TIndEval) {
		df(ind)
		alcanzadaSolucion = true
	}
	addeJobs := func() {
		active := true
		for active {
			select {
			case <-control1:
				active = false
			case eJobs <- EJob{p2Eval.ExtractElements(mSizeEval), ff, qf, alcanzadaSolucionF, eResults}:
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
			case rJobs <- RJob{p2Rep.ExtractElements(mSizeRep), pMutation, rResults}:
			}
		}
		close(rJobs)
	}

	bestSolution := NewIndEval()
	waitAndProcessResults := func() {

		for !alcanzadaSolucion {
			select { // Blocking
			case indEvals := <-eResults:
				if indEvals != nil && len(indEvals) > 0 {
					if bestSolution.Fitness < indEvals[0].Fitness {
						bestSolution = &indEvals[0]
					}
					p2Rep.Append(indEvals)
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
		res <- *bestSolution
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
	if job.Population != nil && len(job.Population) > 0 {
		IndEvals := Evaluate(job.Population, job.FitnessF, job.QualityF, job.DoFunc)
		sort.Sort(IndEvals)
		//		fmt.Println("Evaluados:", len(IndEvals))
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
