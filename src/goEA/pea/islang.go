package pea

import (
	"fmt"
)

func Island(conf ConfIsland) {
	pool := make(TPool, len(conf.Population))

	for _, e := range conf.Population {
		pool[e] = MValue{-1, 1}
	}
	sndEvals := make(chan []TIndividual, conf.ECount)
	rcvEvals := make(chan []IndEval, conf.ECount)
	for i := 0; i < conf.ECount; i++ {
		go evaluator(ConfEval{sndEvals, rcvEvals, maxOne, conf.MSize})
	}
	sndReps := make(chan RepSndMsg, conf.RCount)
	rcvReps := make(chan []TIndividual, conf.RCount)
	for i := 0; i < conf.RCount; i++ {
		go reproducer(ConfRep{sndReps, rcvReps, conf.MSize})
	}

	var active = true
	for active {
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
				pool = updatePool(pool, nInds)
				sndReps <- RepSndMsg{pool, conf.MSize}
			}

		}
	}

}

func updatePool(pool TPool, nI []TIndividual) TPool {
	inds := make([]TIndividual, 0, len(pool))
	for i := range pool{
		inds = append(inds, i)
	}
	for _, k := range inds[len(nI):]{
		delete(pool, k)
	}

	for _, e := range nI {
		pool[e] = MValue{-1, 1}
	}
	return pool
}
