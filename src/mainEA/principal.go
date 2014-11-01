package main

import (
	"ea"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	obj := ea.NewMaxSATProblem("configMaxSAT.json", "./problems/uf100-01.cnf")
	generalAction := func() {
		res := obj.RunSeq()
		b, _ := json.Marshal(*res)
		fmt.Println(string(b))
	}
	if len(os.Args) < 2 {
		generalAction()
	} else {
		switch os.Args[1] {

		case "seq":
			res := obj.RunSeq()
			b, _ := json.Marshal(*res)
			fmt.Println(string(b))

		case "par":
			res := obj.RunConcurrent()
			b, _ := json.Marshal(*res)
			fmt.Println(string(b))

		default:
			generalAction()
		}
	}
}
