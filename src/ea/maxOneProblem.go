package ea

import (
	"math/rand"
	"encoding/json"
	"io/ioutil"
	"time"
	"strings"
//	"fmt"
)

func genIndividual(n int) TIndividual {
	res := make(TIndividual, n)
	for i, _ := range res {
		res[i] = byte(rand.Intn(2))
	}
	return res
}

func genPop(n int, m int) TPopulation {
	res := make(TPopulation, n)
	for i, _ := range res {
		res[i] = genIndividual(m)
	}
	return res
}

type MaxOneProblem struct{
	Problem
}

func NewMaxOneProblem(configFile string) *MaxOneProblem {
	b, _ := ioutil.ReadFile(configFile)
	dec := json.NewDecoder(strings.NewReader(string(b)))
	var m Data
	dec.Decode(&m)
	return &MaxOneProblem{Problem{&m}}
}

func (self *Problem) GenInitPop() TPopulation {
	return genPop(self.problemConfig.PopSize, self.problemConfig.ChromosomeSize)
}

func (self *Problem) runSeqCEvals(fitnessFunction TFitnessFunc) *SeqRes {
	obj := SeqCEvals{SeqConf{self.GenInitPop, fitnessFunction, self.problemConfig.PMutation},
		CEvalsConf{self.problemConfig.Evaluations}}
	initTime := time.Now()
	solution, ce := obj.Run()
	endTime := time.Now()
	return &SeqRes{TRes{ce},
		TSolution{endTime.Sub(initTime).Nanoseconds(), solution.Fitness}}
}


func (self *MaxOneProblem) QualityFitnessFunction(v int) bool {
	return v > self.problemConfig.ChromosomeSize-2
}
func (self *MaxOneProblem) DoWhenQualityFitnessTrue(i TIndEval) {}
func (self *MaxOneProblem) FitnessFunction(ind TIndividual) int {
	res := 0
	for _, e := range ind {
		if e == 1 {
			res++
		}
	}
	return res
}

func (self *MaxOneProblem) RunSeqCEvals() *SeqRes {
	return self.Problem.runSeqCEvals(self.FitnessFunction)
}

func (self *MaxOneProblem) RunParCEvals() *ParRes {
	return self.Problem.runParCEvals(self.FitnessFunction)
}

func (self *Problem) runParCEvals(fitnessFunction TFitnessFunc) *ParRes {
	obj := ParCEvals{ParConf{SeqConf{self.GenInitPop, fitnessFunction, self.problemConfig.PMutation}    ,
		self.problemConfig.EvaluatorsCapacity, self.problemConfig.ReproducersCapacity, self.problemConfig.EvaluatorsCount,
		self.problemConfig.ReproducersCount}, CEvalsConf{self.problemConfig.Evaluations}}
	resChan := make(chan ParRes, 1)
	initTime := time.Now()
	obj.Run(func(res TPoolCEvalsResult){
		endTime := time.Now()
		resChan <- ParRes{TRes{res.Evaluations},
		self.problemConfig.EvaluatorsCapacity, self.problemConfig.ReproducersCapacity,
		self.problemConfig.EvaluatorsCount, self.problemConfig.ReproducersCount,
		TSolution{endTime.Sub(initTime).Nanoseconds(), res.Fitness}}
	})
	res := <-resChan
	return &res
}
