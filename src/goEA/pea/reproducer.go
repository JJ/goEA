package pea

// EnhanceParents get [n(n+1)/2] potentials parents. n = len(pop).
func EnhanceParents(pop TIndsEvaluated) TInds {
	n := len(pop)
	res := make(TInds, n*(n+1)/2)
	for i, indEval := range pop {
		for j := 0; j < n-i; j++ {
			if true { // TODO: usar probabilidad de crossover
				res = append(res, indEval.ind)
			}
		}
	}
	return res
}

// ParentsSelector gets n pairs for reproduction.
func ParentsSelector(pop TInds, n int) [] Pair {
	res := make([] Pair, n)
	for i := 0; i < n; i++ {
		//TODO:
	}
	return res
}

// Crossover function.
func Crossover(p Pair) (i1 TIndividual, i2 TIndividual) {
	//TODO:
	return p.a, p.b
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
			nInds := make(TInds, lenSubp)
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
	}}
