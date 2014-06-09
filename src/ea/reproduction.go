package ea

import (
	//"fmt"
	"math/rand"
)

// reproducer is the working gorutine for reproduce the individuals.
func reproducer(chRcvPop chan TIndsEvaluated, chSndPop chan TPopulation, Mp int, pMutation float32) {

	var active = true
	for active {
		select { // "select bloqueante" para garantizar el control continuo
		case subp := <-chRcvPop:

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

			chSndPop <- nInds

		}
	}
}

// EnhanceParents get [n(n+1)/2] potentials parents. n = len(pop).
// Repeate n times the best individual, n-1 times the second one, ...
// A simple strategy, to analyze a better one.
func EnhanceParents(pop TIndsEvaluated) TPopulation {
	n := len(pop)
	res := make(TPopulation, n*(n+1)/2)
	indx := 0
	for i, indEval := range pop {
		for j := 0; j < n-i; j++ {
			if true { // TODO: usar probabilidad de crossover
				rInd := make(TIndividual, len(indEval.ind))
				copy(rInd, indEval.ind)
				res[indx] = rInd
				indx++
			}
		}
	}
	return res
}

// ParentsSelector gets n pairs for reproduction.
func ParentsSelector(pop TPopulation, n int) []Pair {
	res := make([]Pair, n)
	nPar := len(pop)
	for i := 0; i < n; i++ {
		m1 := rand.Intn(nPar)
		m2 := rand.Intn(nPar)

		i1 := make(TIndividual, len(pop[m1]))
		copy(i1, pop[m1])
		i2 := make(TIndividual, len(pop[m2]))
		copy(i2, pop[m2])

		res[i] = Pair{i1, i2}
	}
	return res
}

// Crossover function.
func Crossover(p Pair) (a TIndividual, b TIndividual) {
	indLength := len(p.a)

	cPoint := rand.Intn(indLength-1) + 1
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

// Mutate one chromosome of the individual.
func Mutate(ind TIndividual) {
	pos := rand.Intn(len(ind))
	ind[pos] = changeGen(ind[pos])
}

func changeGen(i rune) rune {
	var res rune = 0
	if i == 0 {
		res = 1
	}
	return res
}
