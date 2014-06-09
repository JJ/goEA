package ea

import (
	//"fmt"
	"math/rand"
	"sort"
)

func (s *SeqCEvals) Run() TIndEval {
	p2Eval := make(TPopulation, len(s.Population))
	copy(p2Eval, s.Population)

	var qf TQualityF = func(v int) bool { return false }
	var df Tdo = func(i TIndEval) {}

	_, iEvals := Evaluate(p2Eval, s.FitnessF, qf, df)
	ce := len(iEvals)
	sort.Sort(iEvals)
	for ce < s.CEvals {

		p2Eval = Reproduce(iEvals, s.PMutation)

		_, iEvals := Evaluate(p2Eval, s.FitnessF, qf, df)

		sort.Sort(iEvals)
		ce += len(iEvals)

	}

	return iEvals[0]
}

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
