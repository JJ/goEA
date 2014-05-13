package pea

import (
	//"fmt"
)

type TIndividual string
type FitnessFunc func(TIndividual) int
type TPool map[TIndividual]MValue

type MValue struct {
	fitness  int
	state    byte
}

type IndEval struct {
	ind   TIndividual
	value int
}

type TRepSndMsg struct{
	pool TPool
	size int
}

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

type ConfEval struct {
	chRcvPop     <-chan []TIndividual
	chSndPopEval chan <- []IndEval
	fFEval       FitnessFunc
	mSize        int
}

type ConfRep struct {
	chRcvPop <-chan TRepSndMsg
	chSndPop chan <- []TIndividual
	mSize    int
}
