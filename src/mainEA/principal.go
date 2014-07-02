package main

import (
	"ea"
	"fmt"
	"math/rand"
	"runtime"
	"time"
	"encoding/json"
	//	"io/ioutil"
)

type TRes struct{
	CEvals                int
	InitPopSize           int
	ChromSize             int
}
type TSolution struct{
	Duration              int64
	BestFit               int
}
type SeqRes struct{
	TRes
	TSolution
}

func testSeqCEvals() {
	initialPopulationSize := 360
	chromosomeSize := 8
	pop := func() ea.TPopulation {
		return genPop(initialPopulationSize, chromosomeSize)
	}
	cEvals := 20
	obj := ea.SeqCEvals{ea.SeqConf{pop, ea.MaxOne, 0.3}, ea.CEvalsConf{cEvals}}
	initTime := time.Now()
	solution := obj.Run()
	endTime := time.Now()
	res := SeqRes{TRes{cEvals, initialPopulationSize, chromosomeSize}, TSolution{endTime.Sub(initTime).Nanoseconds(), solution.Fitness}}
	b, _ := json.Marshal(res)
	//	ioutil.WriteFile(fName, b, 0x777)
	//	fmt.Println("La mejor solución para 'SeqCEvals' es:", b)
	fmt.Println(string(b))
}

func testSeqFitnessQuality() {
	var qf ea.TQualityF = func(v int) bool { return v > 7 }
	var df ea.Tdo = func(i ea.TIndEval) {}
	initialPopulationSize := 360
	chromosomeSize := 8
	pop := func() ea.TPopulation {
		return genPop(initialPopulationSize, chromosomeSize)
	}
	obj := ea.SeqFitnessQuality{ea.SeqConf{pop, ea.MaxOne, 0.3}, ea.FitnessQualityConf{qf, df}}
	initTime := time.Now()
	solution, cEvals := obj.Run()
	endTime := time.Now()
	res := SeqRes{TRes{cEvals, initialPopulationSize, chromosomeSize}, TSolution{endTime.Sub(initTime).Nanoseconds(), solution.Fitness}}
	b, _ := json.Marshal(res)
	fmt.Println(string(b))
	//	fmt.Println("La mejor solución para 'SeqFitnessQuality' es: ", solution)
}

func genIndividual(n int) ea.TIndividual {
	res := make(ea.TIndividual, n)
	for i, _ := range res {
		if rand.Intn(100)%2 == 0 {
			res[i] = 1
		} else {
			res[i] = 0
		}
	}
	return res
}

func genPop(n int, m int) ea.TPopulation {
	res := make(ea.TPopulation, n)
	for i, _ := range res {
		res[i] = genIndividual(m)
	}
	return res
}

type ParRes struct{
	TRes
	MSizeEvals,
	MSizeReps,
	CEvaluators,
	CReproducers,
	CIslands int
	TSolution
}

func testParCEvals() {
	initialPopulationSize := 360
	chromosomeSize := 8
	pop := func() ea.TPopulation {
		return genPop(initialPopulationSize, chromosomeSize)
	}
	cEvals := 20000
	MSizeEvals := 50
	MSizeReps := 50
	CEvaluators := 7
	CReproducers := 5
	CIslands := 3
	obj := ea.ParCEvals{ea.ParConf{ea.SeqConf{pop, ea.MaxOne, 0.3}, MSizeEvals,
		MSizeReps, CEvaluators, CReproducers, CIslands}, ea.CEvalsConf{cEvals}}

	initTime := time.Now()
	solution := obj.Run()
	endTime := time.Now()
	res := ParRes{TRes{cEvals, initialPopulationSize, chromosomeSize}, MSizeEvals, MSizeReps,
		CEvaluators, CReproducers, CIslands, TSolution{endTime.Sub(initTime).Nanoseconds(), solution.Fitness}}
	b, _ := json.Marshal(res)
	fmt.Println(string(b))

	//	fmt.Println("La mejor solución para 'ParCEvals' es:", solution)
}

func testParFitnessQuality() {
	var qf ea.TQualityF = func(v int) bool { return v > 7 }
	var df ea.Tdo = func(i ea.TIndEval) {}
	initialPopulationSize := 360
	chromosomeSize := 8
	pop := func() ea.TPopulation {
		return genPop(initialPopulationSize, chromosomeSize)
	}
	MSizeEvals := 50
	MSizeReps := 50
	CEvaluators := 7
	CReproducers := 5
	CIslands := 3
	obj := ea.ParFitnessQuality{ea.ParConf{ea.SeqConf{pop, ea.MaxOne, 0.3},
		MSizeEvals, MSizeReps, CEvaluators, CReproducers, CIslands},
		ea.FitnessQualityConf{qf, df}}
	initTime := time.Now()
	solution, cEvals := obj.Run()
	endTime := time.Now()
	res := ParRes{TRes{cEvals, initialPopulationSize, chromosomeSize}, MSizeEvals, MSizeReps,
		CEvaluators, CReproducers, CIslands, TSolution{endTime.Sub(initTime).Nanoseconds(), solution.Fitness}}
	b, _ := json.Marshal(res)
	fmt.Println(string(b))

	//	fmt.Println("La mejor solución para 'ParFitnessQuality' es:", solution)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	testSeqCEvals()
	testSeqFitnessQuality()

	testParCEvals()
	testParFitnessQuality()
}
