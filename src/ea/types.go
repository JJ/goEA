package ea

import (
	"sync"
)

type IDo interface {
	Do()
}

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
	Ind     TIndividual
	Fitness int
}

type TFitnessResult struct {
	TIndEval
	CEvals int
}

func NewIndEval() *TIndEval {
	return &TIndEval{nil, -1}
}

func (this TIndEval) Greater(that TIndEval) bool {
	return this.Fitness < that.Fitness
}

type TIndsEvaluated []TIndEval

func (inds TIndsEvaluated) Len() int { return len(inds) }
func (inds TIndsEvaluated) Less(i, j int) bool { return inds[j].Greater(inds[i]) }
func (inds TIndsEvaluated) Swap(i, j int) { inds[i], inds[j] = inds[j], inds[i] }

type SeqConf struct {
	GetPopulation func() TPopulation
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
	CReproducers,
	CIslands int
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
	PMutation  float32
	results    chan <- TPopulation
}

type evalPool struct {
	pool TPopulation
	mutex *sync.RWMutex
}

type repPool struct {
	pool TIndsEvaluated
	mutex *sync.RWMutex
}

func NewEvalPool() *evalPool {
	return &evalPool{make(TPopulation, 0), new(sync.RWMutex)}
}

func (self *evalPool) Assign(pop TPopulation) {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	self.pool = make(TPopulation, len(pop))
	copy(self.pool, pop)
}

func (self *evalPool) ExtractElements(mSizeEval int) TPopulation {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	nSend2Eval := mSizeEval
	res := TPopulation(nil)
	if len(self.pool) > 0 {
		if len(self.pool) < nSend2Eval {
			nSend2Eval = len(self.pool)
		}
		res = append(TPopulation{}, self.pool[:nSend2Eval]...)
		self.pool = self.pool[nSend2Eval:]
		//		fmt.Println(self.name, nSend2Eval, "-")
	}
	return res
}

func (self *evalPool) Append(nInds TPopulation) {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	self.pool = append(self.pool, nInds...)
	//	fmt.Println(self.name, len(nInds), "+")
}

func (self *evalPool) Length() int {
	self.mutex.RLock()
	defer self.mutex.RUnlock()
	return len(self.pool)
}

func NewRepPool() *repPool {
	return &repPool{make(TIndsEvaluated, 0), new(sync.RWMutex)}
}

func (self *repPool) Assign(pop TIndsEvaluated) {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	self.pool = make(TIndsEvaluated, len(pop))
	copy(self.pool, pop)
}

func (self *repPool) ExtractElements(mSize int) TIndsEvaluated {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	nSend2 := mSize
	res := TIndsEvaluated(nil)
	if len(self.pool) > 0 {
		if len(self.pool) < nSend2 {
			nSend2 = len(self.pool)
		}
		res = append(TIndsEvaluated{}, self.pool[:nSend2]...)
		self.pool = self.pool[nSend2:]
		//		fmt.Println(self.name, nSend2Eval, "-")
	}
	return res
}

func (self *repPool) RemoveWorstN(n int) {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	if n > len(self.pool) {
		n = len(self.pool)
	}
	self.pool = self.pool[:len(self.pool)-n]
	// TODO: report when n > len(self.pool)
}

// Append merge the TIndsEvaluated.
// nInds should be ordered.
func (self *repPool) Append(nInds TIndsEvaluated) {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	u := self.pool
	v := nInds
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
	self.pool = a
	//	fmt.Println(self.name, len(nInds), "+")
}

func (self *repPool) Length() int {
	self.mutex.RLock()
	defer self.mutex.RUnlock()
	return len(self.pool)
}
