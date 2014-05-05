package pea

import (
	//"fmt"
)

type TIndividual string
type FitnessFunc func(TIndividual) int
type TPool map[TIndividual]MValue

type MValue struct {
	fitness int
	tate    byte
}

type IndEval struct {
	ind   TIndividual
	value int
}

type RepSndMsg struct {
	popl  TPool
	mSize int
}

type ConfIsland struct {
	Control       chan string
	Population    []TIndividual
	ChSndEmigrant chan TIndividual
	ChRcvEmigrant chan TIndividual
	RCount        int
	ECount        int
	MSize         int
	CEvals        int
}

type ConfEval struct {
	chRcvPop     chan []TIndividual
	chSndPopEval chan []IndEval
	fFEval       FitnessFunc
	mSize        int
}

type ConfRep struct {
	chRcvPop chan RepSndMsg
	chSndPop chan []TIndividual
	mSize    int
}
