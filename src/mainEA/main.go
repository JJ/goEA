package main

import (
	"ea"
	"fmt"
	"runtime"
)

func testSeqCEvals() {
	obj := ea.SeqCEvals{ea.SeqConf{[]ea.TIndividual{
		[]rune{1, 0, 1, 0, 1, 0, 0, 0},
		[]rune{1, 0, 1, 0, 1, 1, 0, 1},
		[]rune{1, 0, 1, 0, 1, 1, 0, 1},
		[]rune{1, 1, 1, 0, 1, 1, 0, 1},
		[]rune{1, 0, 1, 0, 1, 1, 0, 0},
		[]rune{0, 0, 1, 0, 1, 1, 1, 1}},
		ea.MaxOne,
		0.3},
		ea.CEvalsConf{20}}

	solution := obj.Run()

	fmt.Println("La mejor solución es: ", solution)

}

func testSeqFitnessQuality() {

	var qf ea.TQualityF = func(v int) bool { return v > 7 }
	var df ea.Tdo = func(i ea.TIndEval) {}

	obj := ea.SeqFitnessQuality{
		ea.SeqConf{[]ea.TIndividual{
			[]rune{1, 0, 1, 0, 1, 0, 0, 0},
			[]rune{1, 0, 1, 0, 1, 1, 0, 1},
			[]rune{1, 0, 1, 0, 1, 1, 0, 1},
			[]rune{1, 1, 1, 0, 1, 1, 0, 1},
			[]rune{1, 0, 1, 0, 1, 1, 0, 0},
			[]rune{0, 0, 1, 0, 1, 1, 1, 1}},
			ea.MaxOne,
			0.3},
		ea.FitnessQualityConf{qf, df}}

	solution := obj.Run()

	fmt.Println("La mejor solución es: ", solution)

}

func testParCEvals() {
	obj := ea.ParCEvals{ea.ParConf{ea.SeqConf{[]ea.TIndividual{
		[]rune{1, 0, 1, 0, 1, 0, 0, 0},
		[]rune{1, 0, 1, 0, 1, 1, 0, 1},
		[]rune{1, 0, 1, 0, 1, 1, 0, 1},
		[]rune{1, 1, 1, 0, 1, 1, 0, 1},
		[]rune{1, 0, 1, 0, 1, 1, 0, 0},
		[]rune{0, 0, 1, 0, 1, 1, 1, 1}},
		ea.MaxOne,
		0.3}, 2, 2, 2, 2},
		ea.CEvalsConf{20}}

	solution := obj.Run()

	fmt.Println("La mejor solución es: ", solution)

}

func main() {
	//testSeqCEvals()
	//testSeqFitnessQuality()

	runtime.GOMAXPROCS(runtime.NumCPU())
	testParCEvals()

}
