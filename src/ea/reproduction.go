package ea

import (
	//	"fmt"
	"time"
	"math/rand"
)

// Reproduce is the reproduction routine.
func Reproduce(iEvals TIndsEvaluated, pMutation float32) TPopulation {
	lenSubPop := len(iEvals)
	p2Rep := EnhanceParents(iEvals)
	parents := ParentsSelector(p2Rep, lenSubPop/2)
	nInds := make(TPopulation, 0)
	for _, ind := range parents {
		i1, i2 := Crossover(ind)
		nInds = append(nInds, i1, i2)
	}
	if lenSubPop%2 == 1 {
		nInds = append(nInds, iEvals[0].Ind)
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for _, ind := range nInds {
		if r.Float32() < pMutation {
			Mutate(ind)
		}
	}
	return nInds
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
				rInd := make(TIndividual, len(indEval.Ind))
				copy(rInd, indEval.Ind)
				res[indx] = rInd
				indx++
			}
		}
	}
	return res
}

// ParentsSelector gets n pairs for reproduction.
func ParentsSelector(pop TPopulation, n int) []Pair {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	res := make([]Pair, n)
	nPar := len(pop)
	for i := 0; i < n; i++ {
		m1 := r.Intn(nPar)
		m2 := r.Intn(nPar)
		for pop[m1].Equals(pop[m2]) {
			m2 = r.Intn(nPar)
		}
		i1 := make(TIndividual, len(pop[m1]))
		copy(i1, pop[m1])
		i2 := make(TIndividual, len(pop[m2]))
		copy(i2, pop[m2])
		res[i] = Pair{i1, i2}
	}
	//	fmt.Println("yes")
	return res
}

// Crossover function.
func Crossover(p Pair) (a TIndividual, b TIndividual) {
	indLength := len(p.a)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	cPoint := r.Intn(indLength - 1) + 1
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
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	pos := r.Intn(len(ind))
	ind[pos] = changeGen(ind[pos])
}

func changeGen(i byte) byte {
	var res byte = 0
	if i == 0 {
		res = 1
	}
	return res
}
