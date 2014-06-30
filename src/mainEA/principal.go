package main

import (
	"ea"
	"fmt"
	"math/rand"
	"runtime"
)

func testSeqCEvals() {
	pop := func() ea.TPopulation {
		return genPop(360, 8)
	}
	obj := ea.SeqCEvals{ea.SeqConf{pop,
		ea.MaxOne,
		0.3},
		ea.CEvalsConf{20}}

	solution := obj.Run()

	fmt.Println("La mejor solución para 'SeqCEvals' es:", solution)

}

func testSeqFitnessQuality() {
	var qf ea.TQualityF = func(v int) bool { return v > 7 }
	var df ea.Tdo = func(i ea.TIndEval) {}
	pop := func() ea.TPopulation {
		return genPop(360, 8)
	}
	obj := ea.SeqFitnessQuality{
		ea.SeqConf{pop,
			ea.MaxOne,
			0.3},
		ea.FitnessQualityConf{qf, df}}

	solution := obj.Run()

	fmt.Println("La mejor solución para 'SeqFitnessQuality' es: ", solution)

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
	pop := func() ea.TPopulation {
		return genPop(360, 8)
	}
	obj := ea.ParCEvals{ea.ParConf{ea.SeqConf{pop,
		ea.MaxOne,
		0.3}, 50, 50, 7, 5, 3},
		ea.CEvalsConf{20000}}

	solution := obj.Run()

	fmt.Println("La mejor solución para 'ParCEvals' es:", solution)
}

//func testParFitnessQuality() {
//	var qf ea.TQualityF = func(v int) bool { return v > 7 }
//	var df ea.Tdo = func(i ea.TIndEval) {}
//
//	pop := genPop(360, 8)
//	obj := ea.ParFitnessQuality{ea.ParConf{ea.SeqConf{pop,
//		ea.MaxOne,
//		0.3}, 50, 50, 7, 5},
//		ea.FitnessQualityConf{qf, df}}
//
//	solution := obj.Run()
//
//	fmt.Println("La mejor solución para 'ParFitnessQuality' es:", solution)
//}
func testParFitnessQuality() {
	var qf ea.TQualityF = func(v int) bool { return v > 7 }
	var df ea.Tdo = func(i ea.TIndEval) {}

	pop := func() ea.TPopulation {
		return genPop(360, 8)
	}
	obj := ea.ParFitnessQuality{ea.ParConf{ea.SeqConf{pop,
		ea.MaxOne,
		0.3}, 50, 50, 7, 5, 3},
		ea.FitnessQualityConf{qf, df}}

	solution := obj.Run()

	fmt.Println("La mejor solución para 'ParFitnessQuality' es:", solution)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	testSeqCEvals()
	testSeqFitnessQuality()

	testParCEvals()
	testParFitnessQuality()

}
