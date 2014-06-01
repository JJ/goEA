package seqEA

import (
	"goEA/pea"
	"sort"
	"math/rand"
)

func run() {

}

func runSeqEA(conf ConfSeqEA) {

	p2Eval := make([]TIndividual, len(conf.Population))
	copy(p2Eval, conf.Population)

	for {
		iEvals := pea.Evaluate(ConfEvalSeq{conf.FEval, p2Eval})
		sort.Sort(iEvals)
		lenSubPop := len(iEvals)
		p2Rep := pea.EnhanceParents(iEvals[:lenSubPop])
		parents := pea.ParentsSelector(p2Rep, lenSubPop/2)
		nInds := make(TInds, 0)
		for _, ind := range parents {
			i1, i2 := Crossover(ind)
			nInds = append(nInds, i1, i2)
		}

		if lenSubPop%2 == 1 {
			nInds = append(nInds, iEvals[0])
		}

		for _, ind := range nInds {
			if rand.float32() < conf.PMut {
				pea.Mutate(ind)
			}
		}
		p2Eval = nInds
	}

}
