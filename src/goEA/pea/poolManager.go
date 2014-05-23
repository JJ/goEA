package pea

import (
	"fmt"
)

func chooseInds2Eval(mSize int, pool TPool) {

}

// PoolManager is the gorutine for control de workers. The island manager.
func PoolManager(conf ConfIsland) {
	p2Eval := make([]TIndividual, len(conf.Population))
	append(p2Eval, conf.Population...)
	// Siempre estarán ordenados: de mayor a menor.
	p2Rep := make([]IndEval, len(conf.Population))

	sndEvals := make(chan []TIndividual, conf.ECount)
	rcvEvals := make(chan []IndEval, conf.ECount)
	for i := 0; i < conf.ECount; i++ {
		go evaluator(ConfEval{sndEvals, rcvEvals, maxOne, conf.MSize})
	}

	sndReps := make(chan []IndEval, conf.RCount)
	rcvReps := make(chan []TIndividual, conf.RCount)
	for i := 0; i < conf.RCount; i++ {
		go reproducer(ConfRep{sndReps, rcvReps, conf.MSize})
	}

	var active = true
	var workDone = 0
	for active && workDone >= conf.CEvals {
		select { // "select bloqueante" para garantizar el control continuo
		case cmd := <-conf.Control:
			switch cmd {
			case "start":

			case "finalize":
				active = false
			default:
				fmt.Printf("Mensaje de control %v no entendido.\n", cmd)
			}

		case nInds := <-rcvReps:
			if nInds != nil {
				p2Eval = append(p2Eval, nInds...)
				n2Send := conf.MSize
				if len(p2Rep) < conf.MSize {
					n2Send = len(p2Rep)
				}
				// Mando los n2Send primeros (los mejores).
				sndReps <- p2Rep[:n2Send]
				p2Rep = p2Rep[n2Send:]
			}

			// Los individuos evaluados vienen ordenados por su fitness
		case iEvals := <-rcvEvals:
			if iEvals != nil {
				for _, par := range iEvals {
					pool1[par.ind] = MValue{par.value, 2}
				}

				workDone += len(iEvals)
			}

		}
	}

}

func updatePool(pool TPool, nI []TIndividual) TPool {
	inds := make([]TIndividual, 0, len(pool))
	for i := range pool {
		inds = append(inds, i)
	}

	for _, k := range inds[len(nI):] {
		delete(pool, k)
	}

	for _, e := range nI {
		pool[e] = MValue{-1, 1}
	}
	return pool
}
