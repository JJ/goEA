package main

import (
	"ea"
	"fmt"
	"runtime"
	"encoding/json"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
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

