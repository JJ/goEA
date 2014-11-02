package ea

import (
	//	"math/rand"
	"sort"
	//	"fmt"
)

func (s *SequentialProblem) Run() (TIndEval, int) {
	population := s.GetPopulation()
	p2Eval := make(TPopulation, len(population))
	copy(p2Eval, population)
	alcanzadaSolucion := false
	bestSolution := NewIndEval()
	Do := func(ind TIndEval) {
		bestSolution = &ind
		alcanzadaSolucion = true
	}
	cEvals := 0
	for !alcanzadaSolucion && cEvals < s.Evaluations {
		IndEvals := Evaluate(p2Eval, s.FitnessF, s.QualityF, Do)
		sort.Sort(IndEvals)
		cEvals += len(IndEvals)
		p2Eval = Reproduce(IndEvals, s.PMutation)
	}
	return *bestSolution, cEvals
}
