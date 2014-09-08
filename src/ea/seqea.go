package ea

import (
//	"math/rand"
	"sort"
	//	"fmt"
)

// Run is the method of SeqCEvals to find the solution by the amount of evaluations criteria.
func (s *SeqCEvals) Run() (*TIndEval, int) {
	population := s.GetPopulation()
	p2Eval := make(TPopulation, len(population))
	copy(p2Eval, population)

	var qf TQualityF = func(v int) bool { return false }
	var df Tdo = func(i TIndEval) {}
	bestSolution := NewIndEval()
	ce := 0
	for ce < s.CEvals {
		IndEvals := Evaluate(p2Eval, s.FitnessF, qf, df)
		sort.Sort(IndEvals)
		ce += len(IndEvals)
		p2Eval = Reproduce(IndEvals, s.PMutation)
		if bestSolution.Fitness < IndEvals[0].Fitness {
			bestSolution = &IndEvals[0]
		}
	}
	return bestSolution, ce
}
