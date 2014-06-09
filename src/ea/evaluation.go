package ea

func Evaluate(pop TPopulation, ff TFitnessFunc, qf TQualityF, df Tdo) (bool, TIndsEvaluated) {
	res := make(TIndsEvaluated, len(pop))
	i := 0
	mejorEncontrado := false
	for i < len(pop) && !mejorEncontrado {
		sol := ff(pop[i])
		indT := make(TIndividual, len(pop[i]))
		copy(indT, pop[i])
		nEntry := TIndEval{indT, sol}
		if qf(sol) {
			mejorEncontrado = true
			df(nEntry)
		}
		res[i] = nEntry
		i++
	}
	if mejorEncontrado {
		nRes := make(TIndsEvaluated, i)
		copy(nRes, res)
		res = nRes
	}
	return mejorEncontrado, res
}

func evaluator(chRcvPop chan TPopulation, chSndPopEval chan TIndsEvaluated, ff TFitnessFunc, qf TQualityF, df Tdo) {

	var active = true
	for active {

		select { // "select bloqueante" para garantizar el control continuo

		case work := <-chRcvPop:
			_, res := Evaluate(work, ff, qf, df)
			chSndPopEval <- res

		}

	}
}
