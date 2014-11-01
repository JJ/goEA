package ea

import (
	"math/rand"
	"sort"
	//	"fmt"
)

func (s *SequentialProblem) Run() (TIndEval, int) {
	population := s.GetPopulation()
	p2Eval := make(TPopulation, len(population))
	copy(p2Eval, population)
	alcanzadaSolucion := false
	alcanzadaSolucionF := func(ind TIndEval) {
		alcanzadaSolucion = true
	}
	iEvals := Evaluate(p2Eval, s.FitnessF, s.QualityF, alcanzadaSolucionF)
	cEvals := len(iEvals)
	sort.Sort(iEvals)
	for !alcanzadaSolucion && cEvals < s.Evaluations {
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
