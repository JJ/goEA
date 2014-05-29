package pea

import (
	"math/rand"
)

// EnhanceParents get [n(n+1)/2] potentials parents. n = len(pop).
// Repeate n times the best individual, n-1 times the second one, ...
// A simple strategy, to analyze a better one.
func EnhanceParents(pop TIndsEvaluated) TInds {
	n := len(pop)
	res := make(TInds, n*(n+1)/2)
	indx := 0
	for i, indEval := range pop {
		for j := 0; j < n-i-1; j++ {
			if true { // TODO: usar probabilidad de crossover
				res[indx] = indEval.ind
				indx++
			}
		}
	}
	return res
}

// ParentsSelector gets n pairs for reproduction.
func ParentsSelector(pop TInds, n int) [] Pair {
	res := make([] Pair, n)
	nPar := len(pop)
	for i := 0; i < n; i++ {
		m1 := rand.Intn(nPar)
		m2 := rand.Intn(nPar)
		res[i] = Pair{pop[m1], pop[m2]}
	}
	return res
}

// Crossover function.
func Crossover(p Pair) (a TIndividual, b TIndividual) {
	indLength := len(p.a)
	cPoint := rand.Intn(indLength)
	res1 := make(TIndividual, indLength)
	res2 := make(TIndividual, indLength)

	for i := 0; i < cPoint; i++ {
		res1[i] = p.a[i]
		res2[i] = p.b[i]
	}
	for i := cPoint; i < indLength; i++ {
		res1[i] = p.b[i]
		res2[i] = p.a[i]
	}

	return res1, res2
}

// reproducer is the working gorutine for reproduce the individuals.
func reproducer(conf ConfRep) {

	var active = true
	for active {
		select { // "select bloqueante" para garantizar el control continuo
		case subp := <-conf.chRcvPop: {
			fparents := EnhanceParents(subp)
			lenSubp := len(subp)
			n := lenSubp / 2
			parents := ParentsSelector(fparents, n)
			nInds := make(TInds, 0)
			for _, ind := range parents {
				i1, i2 := Crossover(ind)
				nInds = append(nInds, i1, i2)
			}
			if lenSubp%2 == 1 {
				nInds = append(nInds, subp[0].ind)
			}

			// TODO: mutar sobre nInds

			conf.chSndPop <- nInds

		}

		}
	}
}
