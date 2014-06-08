package main

import (
	"ea"
	"fmt"
)

func testSeqCEvals() {
	conf := ea.SeqCEvals{ea.SeqConf{[]ea.TIndividual{
		[]rune{1, 0, 1, 0, 1, 0, 0, 0},
		[]rune{1, 0, 1, 0, 1, 1, 0, 1},
		[]rune{1, 0, 1, 0, 1, 1, 0, 1},
		[]rune{1, 1, 1, 0, 1, 1, 0, 1},
		[]rune{1, 0, 1, 0, 1, 1, 0, 0},
		[]rune{0, 0, 1, 0, 1, 1, 1, 1}},
		ea.MaxOne,
		0.3},
		ea.CEvalsConf{20}}

	solution := conf.Run()

	fmt.Println("La mejor solución es: ", solution)

}

func testSeqFitnessQuality() {

	var qf ea.TQualityF = func(v int) bool { return v > 7 }
	var df ea.Tdo = func(i ea.TIndEval) {}

	conf := ea.SeqFitnessQuality{
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

	solution := conf.Run()

	fmt.Println("La mejor solución es: ", solution)

}

func main() {
	//testSeqCEvals()
	testSeqFitnessQuality()
}
