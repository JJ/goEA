package main

import (
	"ea"
	"fmt"
	"runtime"
	"encoding/json"
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
	obj := ea.NewMaxSATProblem("configMaxSAT.json", "f:/Mis Documentos/PhD/src/sclEA/Pool-Island/problems/uf100-01.cnf")
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
	//	maxOneTests()
	maxSATTests()
}

