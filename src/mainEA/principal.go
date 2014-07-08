package main

import (
	"ea"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
)

func maxOneTests() {
	obj := ea.NewMaxOneProblem("configMaxOnes.json")
	fmt.Println("SeqCEvals")
	res := obj.RunSeqCEvals()
	b, _ := json.Marshal(*res)
	fmt.Println(string(b))
	fmt.Println()
	fmt.Println("SeqFitnessQuality")
	res = obj.RunSeqFitnessQuality()
	b, _ = json.Marshal(*res)
	fmt.Println(string(b))
	fmt.Println()
	fmt.Println("ParCEvals")
	pres := obj.RunParCEvals()
	b, _ = json.Marshal(*pres)
	fmt.Println(string(b))
	fmt.Println()
	fmt.Println("ParFitnessQuality")
	pres = obj.RunParFitnessQuality()
	b, _ = json.Marshal(*pres)
	fmt.Println(string(b))
}

func maxSATTests() {
	obj := ea.NewMaxSATProblem("configMaxSAT.json", "./problems/uf100-01.cnf")
	fmt.Println("SeqCEvals")
	res := obj.RunSeqCEvals()
	b, _ := json.Marshal(*res)
	fmt.Println(string(b))
	fmt.Println()
	fmt.Println("SeqFitnessQuality")
	res = obj.RunSeqFitnessQuality()
	b, _ = json.Marshal(*res)
	fmt.Println(string(b))
	fmt.Println()
	fmt.Println("ParCEvals")
	pres := obj.RunParCEvals()
	b, _ = json.Marshal(*pres)
	fmt.Println(string(b))
	fmt.Println()
	fmt.Println("ParFitnessQuality")
	pres = obj.RunParFitnessQuality()
	b, _ = json.Marshal(*pres)
	fmt.Println(string(b))
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	testing := false
	if testing {
		maxOneTests()
		maxSATTests()
		os.Exit(0)
	}
	obj := ea.NewMaxSATProblem("configMaxSAT.json", "./problems/uf100-01.cnf")
	//	fmt.Println(obj.ToString()) // To check the problem instance.
	generalAction := func() {
		res := obj.RunSeqCEvals()
		b, _ := json.Marshal(*res)
		fmt.Println(string(b))
	}
	if len(os.Args) < 2 {
		generalAction()
	} else {
		switch os.Args[1] {
		case "seqfq":
			res := obj.RunSeqFitnessQuality()
			b, _ := json.Marshal(*res)
			fmt.Println(string(b))
		case "parce":
			res := obj.RunParCEvals()
			b, _ := json.Marshal(*res)
			fmt.Println(string(b))
		case "parfq":
			res := obj.RunParFitnessQuality()
			b, _ := json.Marshal(*res)
			fmt.Println(string(b))
		default:
			generalAction()
		}
	}
}
