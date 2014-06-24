package ea

import (
	"fmt"
	"sync"
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

// PoolManager is the gorutine for control de workers. The island manager.
func PoolManagerCEvals(population TPopulation,
	eCount int, rCount int,
	mSizeEval int, mSizeRep int,
	pMutation float32, cEvals int,
	ff TFitnessFunc, res chan <- TIndEval) {

	//	workers := eCount + rCount
	eJobs := make(chan EJob, eCount)
	rJobs := make(chan RJob, rCount)
	eResults := make(chan TIndsEvaluated, 1)
	rResults := make(chan TPopulation, 1)

	control := make(chan struct {}, 0)
	p2Eval := make(TPopulation, len(population))
	copy(p2Eval, population)

	// Siempre estarán ordenados: de mayor a menor.
	p2Rep := make(TIndsEvaluated, 0)

	var mp2Eval sync.Mutex
	var mp2Rep sync.Mutex

	var cont1 = 0
	var cont2 = 0

	selPop2Eval := func() TPopulation {
		nSend2Eval := mSizeEval

		mp2Eval.Lock()

		if len(p2Eval) < nSend2Eval {
			nSend2Eval = len(p2Eval)
		}
		// Mando los nSend2Eval primeros (de los que han quedado).
		res := append(TPopulation{}, p2Eval[:nSend2Eval]...)

		//		syscall.Exit(1)

		p2Eval = p2Eval[nSend2Eval:]
		cont1++
		if cont1 <= cantSee {
		}
		if len(res) > 0 {
			fmt.Println("Creando EJob, len(pop):", len(res))
			//			fmt.Println("selPop2Eval, p2Eval len:", len(p2Eval))
		}


		mp2Eval.Unlock()

		return res
	}
	selPop2Rep := func() []TIndEval {
		nSend2Rep := mSizeRep

		mp2Rep.Lock()
		if len(p2Rep) < nSend2Rep {
			nSend2Rep = len(p2Rep)
		}
		// Mando los nSend2Rep primeros (los mejores).
		//		res := append([]TIndEval{}, p2Rep[:nSend2Rep]...)
		//		p2Rep = p2Rep[nSend2Rep:]

		//		cont2++
		//		if cont2 <= cantSee {
		//			fmt.Println("selPop2Rep, p2Rep len:", len(p2Rep))
		//			fmt.Println("selPop2Rep, res:", len(res))
		//		}

		mp2Rep.Unlock()

		//		return res
		return []TIndEval{}
	}

	for i := 0; i < eCount; i++ {
		go doEvalJobs(eJobs)
	}
	//	for i := 0; i < rCount; i++ {
	//		go doRepJobs(rJobs)
	//	}
	go addJobsCEvals(control, eJobs, rJobs, selPop2Eval, selPop2Rep, ff, pMutation, eResults, rResults)

	waitAndProcessResults := func() {
		for ce := cEvals; ce > 0; {
			//			fmt.Println(ce)
			select { // Blocking
			case indEvals := <-eResults:
				if indEvals != nil && len(indEvals) > 0 {
					mp2Rep.Lock()

					p2Rep = Merge(p2Rep, indEvals)
					//					fmt.Println("Nuevos evaluados:", len(indEvals))
					//					fmt.Println("Para rep:", len(p2Rep))
					//					logPools()
					mp2Rep.Unlock()
					ce -= len(indEvals)
				}

				mp2Rep.Lock()
				if len(indEvals) > 0 {
					//					fmt.Println("Nuevos evaluados:", len(indEvals))
					//					fmt.Println("Para rep:        ", len(p2Rep))
				}
				if cont2 <= cantSee {
					cont2++
				}
				//					logPools()
				mp2Rep.Unlock()

			case nInds := <-rResults:
				if nInds != nil && len(nInds) > 0 {
					mp2Eval.Lock()
					p2Eval = append(p2Eval, nInds...)

					//					fmt.Println("Nuevos individuos:", len(nInds))
					//					fmt.Println("Para eval:", len(p2Eval))

					mp2Eval.Unlock()
				}
			}

		}

		control <- struct{}{}
		//		fmt.Println("SE ACABOOOOOOOOOOOOOOOOOOOOO:")
		//		fmt.Println(p2Rep)
		res <- p2Rep[0]
	}

	waitAndProcessResults()
}

var cJobsdo = 0
var cantSee = 10

func addJobsCEvals(control chan struct{}, eJobs chan <- EJob, rJobs chan <- RJob,
	selPop2Eval func() TPopulation, selPop2Rep func() []TIndEval,
	FitnessF TFitnessFunc, PMutation float32,
	reportEvalResults chan TIndsEvaluated,
	reportRepResults chan TPopulation) {

	var qf TQualityF = func(v int) bool { return false }
	var df Tdo = func(i TIndEval) {}

	active := true
	for active {
		select {

		case <-control:
			active = false

		case eJobs <- EJob{selPop2Eval(), FitnessF, qf, df, reportEvalResults}:

		case rJobs <- RJob{selPop2Rep(), PMutation, reportRepResults}:

		}
	}
	close(eJobs)
	close(rJobs)
}

func doEvalJobs(jobs <-chan EJob) {
	for job := range jobs {
		if len(job.Population) > 0 {
			fmt.Println("job.Do con:", len(job.Population))
		}else if cJobsdo < cantSee {
			cJobsdo++
			fmt.Println("job.Do vacio.")
		}
		job.Do()
	}
}

// TODO: Cambiar el chanel de jods por uno de referencias a ibjetos con metodo "Do"
func doRepJobs(jobs <-chan RJob) {
	for job := range jobs {
		job.Do()
	}
}


func (job EJob) Do() {
	if job.Population != nil && len(job.Population) > 0 {
		_, IndEvals := Evaluate(job.Population, job.FitnessF, job.QualityF, job.DoFunc)
		sort.Sort(IndEvals)
		//		fmt.Println("Evaluados:", len(IndEvals))
		job.results <- IndEvals
	}else {
		job.results <- nil
	}
}

func (job RJob) Do() {
	if job.IndEvals != nil && len(job.IndEvals) > 0 {
		reproductionResults := Reproduce(job.IndEvals, job.PMutation)
		job.results <- reproductionResults
	}else {
		job.results <- nil
	}
}

// Merge is the mixer of two ordered sequences of individuals evaluated.
func Merge(u, v TIndsEvaluated) TIndsEvaluated {
	l := len(u) + len(v)
	a := make(TIndsEvaluated, l)
	i, j, k := 0, 0, 0
	for i < l {
		if j < len(v) && k < len(u) {
			if v[j].Greater(u[k]) {
				a[i] = v[j]
				j++
			} else {
				a[i] = u[k]
				k++
			}
		} else {
			if j >= len(v) {
				for k < len(u) {
					a[i] = u[k]
					i++
					k++
				}
			} else {
				for j < len(v) {
					a[i] = v[j]
					i++
					j++
				}
			}
		}
		i++
	}

	return a
}


