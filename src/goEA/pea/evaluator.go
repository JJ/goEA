package pea

import (
	//"fmt"
)

func evaluator(conf ConfEval) {

	var active = true
	for active {
		select { // "select bloqueante" para garantizar el control continuo
		case work := <-conf.chRcvPop: {
			r := make(TIndsEvaluated, len(work))
			for i := range work {
				r[i] = IndEval{work[i], conf.fFEval(work[i])}
			}
			conf.chSndPopEval <- r
		}

		}
	}
}
