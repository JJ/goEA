package main

import (
	"ea"
	"fmt"
	"math/rand"
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

	fmt.Println("La mejor solución es:", solution)

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

func testParCEvals() {
	pop := genPop(360, 8)
	obj := ea.ParCEvals{ea.ParConf{ea.SeqConf{pop,
		ea.MaxOne,
		0.3}, 50, 50, 7, 5},
		ea.CEvalsConf{20000}}

	solution := obj.Run()

	fmt.Println("La mejor solución es:", solution)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	testSeqCEvals()
	testSeqFitnessQuality()

	testParCEvals()

	//	res := make(chan int, 1)
	//	eatest.TestParAlg([]int{1, 3, 5, 7, 9, 11, 13, 15, 17, 19, 21, 23, 25, 27, 29, 31},
	//		1, 1,
	//		2, 2,
	//		10, res)
	//
	//	fmt.Println(<-res)

}
