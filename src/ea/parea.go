package ea

import (
//"fmt"
)

func (s *ParCEvals) Run() TIndEval {
	//ch1 := make(chan TIndividual, 1)
	//ch2 := make(chan TIndividual, 1)
	control := make(chan string, 1)

	//go pea.PoolManager(cnf)
	control <- "start"

	return TIndEval{[]rune{0}, 0}
}
