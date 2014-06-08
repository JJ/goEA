package ea

func Evaluate(pop TPopulation, ff TFitnessFunc, qf TQualityF, df Tdo)(bool, TIndsEvaluated){
	res := make(TIndsEvaluated, len(pop))
	i := 0
	mejorEncontrado := false
	for i < len(pop) && ! mejorEncontrado {
		sol := ff(pop[i])
		nEntry := TIndEval{pop[i], sol}
		if (! qf(sol)){
			res[i] = nEntry
		}else{
			mejorEncontrado = true
			df(nEntry)
		}
		i++
	}
	return mejorEncontrado, res
}
