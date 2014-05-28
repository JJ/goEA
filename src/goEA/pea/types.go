package pea

import (
	//"fmt"
)

// TIndividual represents a chromosome, normaly a char sequence.
type TIndividual []rune

// TInds is a sequence of TIndividual.
type TInds []TIndividual

// Pair is two individuals, named for reproduction.
type Pair struct{
	a TIndividual
	b TIndividual
}

// FitnessFunc is the type for function that evaluate a chromosome.
type FitnessFunc func(TIndividual) int

// IndEval: an individual and its fitness.
type IndEval struct {
	ind     TIndividual
	fitness int
}
type TIndsEvaluated []IndEval

func (this IndEval) Greater(that IndEval) bool {
	return this.fitness > that.fitness
}

func (inds TIndsEvaluated) Len() int { return len(inds) }
func (inds TIndsEvaluated) Less(i, j int) bool {
	return inds[j].Greater(inds[i])
}
func (inds TIndsEvaluated) Swap(i, j int) {
	inds[i], inds[j] = inds[j], inds[i]
}

// ConfIsland is the input of an island gorutine.
type ConfIsland struct {
	Control          chan string
	Population       []TIndividual
	ChSndEmigrant    chan <- TIndividual
	ChRcvEmigrant    <-chan TIndividual
	RCount, ECount,
	MSize, CEvals    int
}

// ConfIsland is the input of an evaluator gorutine.
type ConfEval struct {
	chRcvPop     <-chan []TIndividual
	chSndPopEval chan <- TIndsEvaluated
	fFEval       FitnessFunc
	mSize        int
}

// ConfIsland is the input of an reproducer gorutine.
type ConfRep struct {
	chRcvPop <-chan TIndsEvaluated
	chSndPop chan <- []TIndividual
	mSize    int
}
