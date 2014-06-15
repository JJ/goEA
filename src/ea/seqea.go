package ea

import (
	//"fmt"
	"math/rand"
	"sort"
)

// Run is the method of SeqCEvals to find the solution by the amount of evaluations criteria.
func (s *SeqCEvals) Run() TIndEval {
	p2Eval := make(TPopulation, len(s.Population))
	copy(p2Eval, s.Population)

	var qf TQualityF = func(v int) bool { return false }
	var df Tdo = func(i TIndEval) {}

	_, IndEvals := Evaluate(p2Eval, s.FitnessF, qf, df)
	ce := len(IndEvals)
	sort.Sort(IndEvals)

	for ce < s.CEvals {
		reproductionResults := Reproduce(IndEvals, s.PMutation)
		p2Eval = reproductionResults
		_, IndEvals = Evaluate(p2Eval, s.FitnessF, qf, df)
		sort.Sort(IndEvals)
		ce += len(IndEvals)
	}
	return IndEvals[0]
}

// Run is the method of SeqFitnessQuality to find the solution by the fitness quality criteria.
func (s *SeqFitnessQuality) Run() TIndEval {

	p2Eval := make(TPopulation, len(s.Population))
	copy(p2Eval, s.Population)

	alcanzadaSolucion, iEvals := Evaluate(p2Eval, s.FitnessF, s.QualityF, s.Do)
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
			nInds = append(nInds, iEvals[0].ind)
		}
		for _, ind := range nInds {
			if rand.Float32() < s.PMutation {
				Mutate(ind)
			}
		}
		p2Eval = nInds
		alcanzadaSolucion, iEvals = Evaluate(p2Eval, s.FitnessF, s.QualityF, s.Do)
		sort.Sort(iEvals)
	}

	return iEvals[0]
}
