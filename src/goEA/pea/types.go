package pea

import (
	//"fmt"
)

// TIndividual represents a chromosome, romally a sequence.
type TIndividual string

// FitnessFunc is the type for function that evaluate a chromosome.
type FitnessFunc func(TIndividual) int

// IndEval: an individual and its fitness.
type IndEval struct {
	ind   TIndividual
	value int
}

// ConfIsland is the input of an island gorutine.
type ConfIsland struct {
	Control       chan string
	Population    []TIndividual
	ChSndEmigrant chan <- TIndividual
	ChRcvEmigrant <-chan TIndividual
	RCount        int
	ECount        int
	MSize         int
	CEvals        int
}

// ConfIsland is the input of an evaluator gorutine.
type ConfEval struct {
	chRcvPop     <-chan []TIndividual
	chSndPopEval chan <- []IndEval
	fFEval       FitnessFunc
	mSize        int
}

// ConfIsland is the input of an reproducer gorutine.
type ConfRep struct {
	chRcvPop <-chan []IndEval
	chSndPop chan <- []TIndividual
	mSize    int
}
