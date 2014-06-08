package ea

import (
	"sort"
	"math/rand"
)

func (s *SeqCEvals) Run() TIndEval{
	p2Eval := s.Population
	ce := 0

	var qf TQualityF = func(v int)bool {return false}
	var df Tdo = func(i TIndEval){}

	_, iEvals := Evaluate(p2Eval, s.FitnessF, qf, df)
	ce += len(iEvals)
	sort.Sort(iEvals)
	for ce < s.CEvals{
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
		_, iEvals := Evaluate(p2Eval, s.FitnessF, qf, df)
		sort.Sort(iEvals)
		ce += len(iEvals)
	}
	return iEvals[0]
}
