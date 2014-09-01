package ea

import (
	//	"fmt"
	//	"syscall"
	"math/rand"
	"sort"
)

func (s *ParCEvals) Run() TPoolCEvalsResult {
	cantsEvals := make([]int, s.CIslands)
	a := s.CEvals / s.CIslands
	for i, _ := range cantsEvals {
		cantsEvals[i] = a
	}
	rest := s.CEvals % s.CIslands
	for i := 0; i < rest; i++ {
		cantsEvals[i]++
	}
	migrantsChannels := make([]chan TIndEval, s.CIslands)
	for i, _ := range migrantsChannels {
		migrantsChannels[i] = make(chan TIndEval, 1)
	}
	res := make(chan TPoolCEvalsResult, 1)
	for i := 0; i < s.CIslands; i++ {
		go PoolManagerCEvals(
			s.GetPopulation(),
			s.CEvaluators, s.CReproducers,
			s.MSizeEvals, s.MSizeReps,
			s.PMutation, cantsEvals[i],
			s.FitnessF, res,
			migrantsChannels[i], migrantsChannels[(i+1)%s.CIslands])
	}
	bestSolution := NewIndEval()
	cEmigrations := 0
	for i := 0; i < s.CIslands; i++ {
		sol := <-res
		cEmigrations += sol.Emigrations
		if bestSolution.Fitness < sol.Fitness {
			bestSolution = &sol.TIndEval
		}
	}
	return TPoolCEvalsResult{*bestSolution, cEmigrations}
}

// PoolManager is the gorutine for control de workers. The island manager.
func PoolManagerCEvals(population TPopulation,
	eCount int, rCount int,
	mSizeEval int, mSizeRep int,
	pMutation float32, cEvals int,
	ff TFitnessFunc, res chan <- TPoolCEvalsResult,
	newMigrant <-chan TIndEval,
	migrantsDestination chan <- TIndEval) {

	//	workers := eCount + rCount
	eJobs := make(chan IDo, 1)
	rJobs := make(chan IDo, 1)
	eResults := make(chan TIndsEvaluated, 1)
	rResults := make(chan TPopulation, 1)

	control1 := make(chan struct{}, 1)
	control2 := make(chan struct{}, 1)
	cEmigrations := 0
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

	readingMigrants := func() {
		for nInd := range newMigrant {
			cEmigrations++
			p2Rep.Append(TIndsEvaluated{nInd})
			p2Rep.RemoveWorstN(1)
			// TODO: maintain pool size when p2Rep.Length() < 1
		}
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
					if rand.Intn(100)%2 == 0 && bestSolution.Ind != nil {
						migrantsDestination <- *bestSolution
					}
				}
			}
		}
		control1 <- struct{}{}
		control2 <- struct{}{}
		// fmt.Println("The End!")
		res <- TPoolCEvalsResult{*bestSolution, cEmigrations}
	}

	for i := 0; i < eCount; i++ {
		go doJobs(eJobs)
	}
	for i := 0; i < rCount; i++ {
		go doJobs(rJobs)
	}
	go addeJobs()
	go addrJobs()
	go readingMigrants()
	waitAndProcessResults()
}

func (s *ParFitnessQuality) Run() TPoolFitnessQualityResult {
	control := make(chan struct{}, s.CIslands-1)
	migrantsChannels := make([]chan TIndEval, s.CIslands)
	for i, _ := range migrantsChannels {
		migrantsChannels[i] = make(chan TIndEval, 1)
	}
	res := make(chan TPoolFitnessQualityResult, 1)
	for i := 0; i < s.CIslands; i++ {
		go PoolManagerFitnessQuality(
			s.GetPopulation(),
			s.CEvaluators, s.CReproducers,
			s.MSizeEvals, s.MSizeReps,
			s.PMutation, s.FitnessF,
			s.QualityF, s.Do,
			res,
			migrantsChannels[i], migrantsChannels[(i+1)%s.CIslands],
			control)
	}
	sol := <-res
	for i := 0; i < s.CIslands-1; i++ {
		control <- struct{}{}
	}
	return sol
}

// PoolManager is the gorutine for control de workers. The island manager.
func PoolManagerFitnessQuality(population TPopulation,
	eCount int, rCount int,
	mSizeEval int, mSizeRep int,
	pMutation float32, ff TFitnessFunc,
	qf TQualityF, df Tdo, res chan <- TPoolFitnessQualityResult,
	newMigrant <-chan TIndEval,
	migrantsDestination chan <- TIndEval,
	control chan struct{}) {

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
	bestSolution := NewIndEval()
	alcanzadaSolucion := false
	cEvals := 0
	cEmigrations := 0
	alcanzadaSolucionF := func(ind TIndEval) {
		// TODO: To precise "cEvals"
		df(ind)
		bestSolution = &ind
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

	readingMigrants := func() {
		for nInd := range newMigrant {
			cEmigrations++
			p2Rep.Append(TIndsEvaluated{nInd})
			p2Rep.RemoveWorstN(1)
			// TODO: maintain pool size when p2Rep.Length() < 1
		}
	}

	waitAndProcessResults := func() {
		for !alcanzadaSolucion {
			select { // Blocking
			case indEvals := <-eResults:
				if indEvals != nil && len(indEvals) > 0 {
					if bestSolution.Fitness < indEvals[0].Fitness {
						bestSolution = &indEvals[0]
					}
					p2Rep.Append(indEvals)
					cEvals += len(indEvals)
				}
			case nInds := <-rResults:
				if nInds != nil && len(nInds) > 0 {
					p2Eval.Append(nInds)
					if rand.Intn(100)%2 == 0 && bestSolution.Ind != nil {
						migrantsDestination <- *bestSolution
					}
				}
			case <-control:
			control1 <- struct{}{}
			control2 <- struct{}{}
				return
			}
		}
		control1 <- struct{}{}
		control2 <- struct{}{}
		// fmt.Println("The End!")
		res <- TPoolFitnessQualityResult{*bestSolution, cEvals, cEmigrations}
	}

	for i := 0; i < eCount; i++ {
		go doJobs(eJobs)
	}
	for i := 0; i < rCount; i++ {
		go doJobs(rJobs)
	}
	go addeJobs()
	go addrJobs()
	go readingMigrants()
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
