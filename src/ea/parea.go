package ea

import (
//"fmt"
)

func (s *ParCEvals) Run() TIndEval {
	//ch1 := make(chan TIndividual, 1)
	//ch2 := make(chan TIndividual, 1)
	control := make(chan string, 1)

	go PoolManagerCEvals(
		s.Population,
		s.CEvaluators, s.CReproducers,
		s.MSizeEvals, s.MSizeReps,
		s.PMutation, s.CEvals,
		s.FitnessF,
		control)

	control <- "start"

	return TIndEval{[]rune{0}, 0}
}
