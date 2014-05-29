package pea

import (
	//"fmt"
)

func Evaluate(conf ConfEvalSeq) TIndsEvaluated {
	res := make(TIndsEvaluated, len(conf.Population))
	for i := range conf.Population {
		res[i] = IndEval{conf.Population[i], conf.FEval(conf.Population[i])}
	}
	return res
}

func evaluator(conf ConfEval) {

	var active = true
	for active {
		select { // "select bloqueante" para garantizar el control continuo

		case work := <-conf.chRcvPop:
		conf.chSndPopEval <- Evaluate(ConfEvalSeq{work, conf.FEval})

		}
	}
}
