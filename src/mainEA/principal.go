package main

import (
	"ea"
	"encoding/json"
	"fmt"
	"runtime"
	"os"
)


func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	obj := ea.NewMaxSATProblem("configMaxSAT.json", "./problems/uf100-01.cnf")
//	obj := ea.NewMaxOneProblem("configMaxOnes.json")
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
		case "pce":
			res := obj.RunParCEvals()
			b, _ := json.Marshal(*res)
			fmt.Println(string(b))
		case "sfq":
			res := obj.RunSeqFitnessQuality()
			b, _ := json.Marshal(*res)
			fmt.Println(string(b))
		case "pfq":
			res := obj.RunParFitnessQuality()
			b, _ := json.Marshal(*res)
			fmt.Println(string(b))
		default:// case "sce":
			generalAction()
		}
	}
}
