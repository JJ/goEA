package seqEA

import (
	"goEA/pea"
	"sort"
)

func run(){

}

func runSeqEA(conf ConfSeqEA){
	p2Eval := make([]TIndividual, len(conf.Population))
	copy(p2Eval, conf.Population)

	for {
		tEvals := pea.Evaluate(ConfEvalSeq{conf.FEval, p2Eval})
		iEvals := sort.Sort(tEvals)
		p2Rep := pea.EnhanceParents(iEvals)
		lenSubp := len(iEvals)
		parents := pea.ParentsSelector(p2Rep, lenSubp / 2)
		nInds := make(TInds, 0)
		for _, ind := range parents {
			i1, i2 := Crossover(ind)
			nInds = append(nInds, i1, i2)
		}
		if lenSubp%2 == 1 {
			nInds = append(nInds, iEvals[0])
		}
	}

}
