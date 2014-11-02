package ea

import (
	//	"fmt"
	"time"
	"math/rand"
)

// Reproduce is the reproduction routine.
func Reproduce(iEvals TIndsEvaluated, pMutation float32) TPopulation {
	lenSubPop := len(iEvals)
	parents := ParentsSelector(iEvals, lenSubPop/2)
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

// ParentsSelector gets n pairs for reproduction.
func ParentsSelector(pop TIndsEvaluated, n int) []Pair {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	res := make([]Pair, n)
	nPar := len(pop)
	for i := 0; i < n; i++ {
		n1 := r.Intn(nPar)
		n2 := r.Intn(nPar)
		n3 := r.Intn(nPar)
		i1 := TIndEval{pop[n1].Ind, pop[n1].Fitness}
		i2 := TIndEval{pop[n2].Ind, pop[n2].Fitness}
		i3 := TIndEval{pop[n3].Ind, pop[n3].Fitness}

		if (i1.Fitness < i2.Fitness) {
			if (i1.Fitness < i3.Fitness) {
				res[i] = Pair{i2.Ind, i3.Ind}
			} else {
				res[i] = Pair{i2.Ind, i1.Ind}
			}
		} else {
			if (i2.Fitness < i3.Fitness) {
				res[i] = Pair{i1.Ind, i3.Ind}
			} else {
				res[i] = Pair{i2.Ind, i1.Ind}
			}
		}
	}
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
