package pea

import (
	"fmt"
)

// PoolManager is the gorutine for control de workers. The island manager.
func PoolManager(conf ConfIsland) {
	p2Eval := make([]TIndividual, len(conf.Population))
	copy(p2Eval, conf.Population)
	// Siempre estarán ordenados: de mayor a menor.
	p2Rep := make(TIndsEvaluated, len(conf.Population))

	sndEvals := make(chan []TIndividual, conf.ECount)
	rcvEvals := make(chan TIndsEvaluated, conf.ECount)
	for i := 0; i < conf.ECount; i++ {
		go evaluator(ConfEval{sndEvals, rcvEvals, conf.FEval})
	}

	sndReps := make(chan TIndsEvaluated, conf.RCount)
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
				if len(p2Rep) < n2Send {
					n2Send = len(p2Rep)
				}
				// Mando los n2Send primeros (los mejores).
				sndReps <- append([]IndEval{}, p2Rep[:n2Send]...)
				p2Rep = p2Rep[n2Send:]
			}

			// Los individuos evaluados vienen ordenados por su fitness.
		case iEvals := <-rcvEvals:
			if iEvals != nil {
				p2Rep = Merge(p2Rep, iEvals)
				workDone += len(iEvals)
			}

		}
	}
}

// Merge is the mixer of two ordered sequences of individuals evaluated.
func Merge(u, v TIndsEvaluated) TIndsEvaluated {
	l := len(u) + len(v)
	a := make(TIndsEvaluated, l)
	i, j, k := 0, 0, 0
	for i < l {
		if j < len(v) && k < len(u) {
			if v[j].Greater(u[k]) {
				a[i] = v[j]
				j++
			}else {
				a[i] = u[k]
				k++
			}
		}else {
			if j >= len(v) {
				for k < len(u) {
					a[i] = u[k]
					i++; k++
				}
			} else {
				for j < len(v) {
					a[i] = v[j]
					i++; j++
				}
			}
		}
		i++
	}

	return a
}

//func updatePool(pool TPool, nI []TIndividual) TPool {
//	inds := make([]TIndividual, 0, len(pool))
//	for i := range pool {
//		inds = append(inds, i)
//	}
//
//	for _, k := range inds[len(nI):] {
//		delete(pool, k)
//	}
//
//	for _, e := range nI {
//		pool[e] = MValue{-1, 1}
//	}
//	return pool
//}
