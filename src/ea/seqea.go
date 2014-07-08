package ea

import (
	"math/rand"
	"sort"
	//	"fmt"
)

// Run is the method of SeqCEvals to find the solution by the amount of evaluations criteria.
func (s *SeqCEvals) Run() TIndEval {
	population := s.GetPopulation()
	p2Eval := make(TPopulation, len(population))
	copy(p2Eval, population)

	var qf TQualityF = func(v int) bool { return false }
	var df Tdo = func(i TIndEval) {}
	//	fmt.Println("A evaluar")
	IndEvals := Evaluate(p2Eval, s.FitnessF, qf, df)
	ce := len(IndEvals)
	sort.Sort(IndEvals)

	for ce < s.CEvals {
		p2Eval = Reproduce(IndEvals, s.PMutation)
		IndEvals = Evaluate(p2Eval, s.FitnessF, qf, df)
		sort.Sort(IndEvals)
		ce += len(IndEvals)
		//		fmt.Println("Hechas", ce)
	}
	return IndEvals[0]
}

// Run is the method of SeqFitnessQuality to find the solution by the fitness quality criteria.
func (s *SeqFitnessQuality) Run() (TIndEval, int) {
	population := s.GetPopulation()
	p2Eval := make(TPopulation, len(population))
	copy(p2Eval, population)
	alcanzadaSolucion := false
	alcanzadaSolucionF := func(ind TIndEval) {
		s.Do(ind)
		alcanzadaSolucion = true
	}
	iEvals := Evaluate(p2Eval, s.FitnessF, s.QualityF, alcanzadaSolucionF)
	cEvals := len(iEvals)
	sort.Sort(iEvals)
	for !alcanzadaSolucion {
		lenSubPop := len(iEvals)
		p2Rep := EnhanceParents(iEvals[:lenSubPop])
		parents := ParentsSelector(p2Rep, lenSubPop/2)
		nInds := make(TPopulation, 0)
		for _, ind := range parents {
			i1, i2 := Crossover(ind)
			nInds = append(nInds, i1, i2)
		}
		if lenSubPop%2 == 1 {
			nInds = append(nInds, iEvals[0].Ind)
		}
		for _, ind := range nInds {
			if rand.Float32() < s.PMutation {
				Mutate(ind)
			}
		}
		p2Eval = nInds
		iEvals = Evaluate(p2Eval, s.FitnessF, s.QualityF, alcanzadaSolucionF)
		cEvals += len(iEvals)
		sort.Sort(iEvals)
	}

	return iEvals[0], cEvals
}
