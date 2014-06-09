package ea

import (
	"fmt"
)

// PoolManager is the gorutine for control de workers. The island manager.
func PoolManagerCEvals() {
	var qf TQualityF = func(v int) bool { return false }
	var df Tdo = func(i TIndEval) {}

	p2Eval := make(TPopulation, len(population))
	copy(p2Eval, population)
	// Siempre estarán ordenados: de mayor a menor.
	p2Rep := make(TIndsEvaluated, 0)

	sndEvals := make(chan TPopulation, eCount)
	rcvEvals := make(chan TIndsEvaluated, eCount)
	for i := 0; i < eCount; i++ {
		go evaluator(sndEvals, rcvEvals, ff, qf, df)
	}

	sndReps := make(chan TIndsEvaluated, rCount)
	rcvReps := make(chan TPopulation, rCount)
	for i := 0; i < rCount; i++ {
		go reproducer(sndReps, rcvReps, Mp, pMutation)
	}

	ce := 0
	for ce < cEvals {

		select { // "select bloqueante" para garantizar el control continuo

		case cmd := <-conf.Control:
			switch cmd {

			case "start":

			default:
				fmt.Printf("Mensaje de control %v no entendido.\n", cmd)
			}

			// Los individuos evaluados vienen ordenados por su fitness.
		case iEvals := <-rcvEvals:
			if iEvals != nil {
				p2Rep = Merge(p2Rep, iEvals)
				ce += len(iEvals)

				nSend2Eval := mSize
				if len(p2Eval) < nSend2Eval {
					nSend2Eval = len(p2Eval)
				}
				// Mando los nSend2Eval primeros (de los que han quedado).
				sndEvals <- append(TPopulation{}, p2Eval[:nSend2Eval]...)
				p2Eval = p2Eval[nSend2Eval:]
			}

		case nInds := <-rcvReps:
			if nInds != nil {
				p2Eval = append(p2Eval, nInds...)
				nSend2Rep := mSize
				if len(p2Rep) < nSend2Rep {
					nSend2Rep = len(p2Rep)
				}
				// Mando los nSend2Rep primeros (los mejores).
				sndReps <- append([]IndEval{}, p2Rep[:nSend2Rep]...)
				p2Rep = p2Rep[nSend2Rep:]
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
			} else {
				a[i] = u[k]
				k++
			}
		} else {
			if j >= len(v) {
				for k < len(u) {
					a[i] = u[k]
					i++
					k++
				}
			} else {
				for j < len(v) {
					a[i] = v[j]
					i++
					j++
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
