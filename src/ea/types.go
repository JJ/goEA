package ea

// TIndividual represents a chromosome, normaly a char sequence.
type TIndividual []rune

// TPopulation is a sequence of TIndividual.
type TPopulation []TIndividual

// Pair is two individuals, named for reproduction.
type Pair struct {
	a TIndividual
	b TIndividual
}

// TFitnessFunc is the type for function that evaluate a chromosome.
type TFitnessFunc func(TIndividual) int

type TFSelPop2Eval func() TPopulation
type TFSelPop2Rep func() TIndsEvaluated

// IndEval: an individual and its fitness.
type TIndEval struct {
	ind     TIndividual
	fitness int
}

func (this TIndEval) Greater(that TIndEval) bool {
	return this.fitness < that.fitness
}

type TIndsEvaluated []TIndEval

func (inds TIndsEvaluated) Len() int { return len(inds) }
func (inds TIndsEvaluated) Less(i, j int) bool { return inds[j].Greater(inds[i]) }
func (inds TIndsEvaluated) Swap(i, j int) { inds[i], inds[j] = inds[j], inds[i] }

type SeqConf struct {
	Population TPopulation
	FitnessF   TFitnessFunc
	PMutation  float32
}

type CEvalsConf struct {
	CEvals int
}

type Tdo func(TIndEval)

type TQualityF func(int) bool

type FitnessQualityConf struct {
	QualityF TQualityF
	Do       Tdo
}

type SeqCEvals struct {
	SeqConf
	CEvalsConf
}

type SeqFitnessQuality struct {
	SeqConf
	FitnessQualityConf
}

type ParConf struct {
	SeqConf
	MSizeEvals,
	MSizeReps,
	CEvaluators,
	CReproducers int
}

type ParCEvals struct {
	ParConf
	CEvalsConf
}

type ParFitnessQuality struct {
	ParConf
	FitnessQualityConf
}

type EJob struct {
	Population     TPopulation
	FitnessF       TFitnessFunc
	QualityF       TQualityF
	DoFunc         Tdo
	results        chan <- TIndsEvaluated
}

type RJob struct {
	IndEvals   TIndsEvaluated
	PMutation float32
	results   chan <- TPopulation
}
