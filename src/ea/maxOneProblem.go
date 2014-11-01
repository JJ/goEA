package ea

import (
	"math/rand"
	"encoding/json"
	"io/ioutil"
	"time"
	"strings"
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

func (self *MaxOneProblem) QualityFitnessFunction(v int) bool {
	return v > self.problemConfig.ChromosomeSize-2
}
func (self *MaxOneProblem) FitnessFunction(ind TIndividual) int {
	res := 0
	for _, e := range ind {
		if e == 1 {
			res++
		}
	}
	return res
}

func (self *Problem) runSeq(fitnessFunction TFitnessFunc, qf TQualityF) *SequentialResult {
	obj := SequentialProblem{self.GenInitPop, fitnessFunction, qf, self.problemConfig.Evaluations, self.problemConfig.PMutation}
	initTime := time.Now()
	solution, cEvals := obj.Run()
	endTime := time.Now()
	return &SequentialResult{endTime.Sub(initTime).Nanoseconds(), cEvals, solution.Fitness}
}

func (self *Problem) runConcurrent(fitnessFunction TFitnessFunc, qf TQualityF) *ConcurrentResult {
	obj := ConcurrentProblem{self.GenInitPop, fitnessFunction, self.problemConfig.PMutation,
		self.problemConfig.EvaluatorsCapacity, self.problemConfig.ReproducersCapacity,
		self.problemConfig.EvaluatorsCount, self.problemConfig.ReproducersCount,
		self.problemConfig.IslandsCount, qf, self.problemConfig.Evaluations}

	initTime := time.Now()
	solution := obj.Run()
	endTime := time.Now()

	return &ConcurrentResult{solution.CEvals, self.problemConfig.EvaluatorsCapacity, self.problemConfig.ReproducersCapacity,
		self.problemConfig.EvaluatorsCount, self.problemConfig.ReproducersCount, self.problemConfig.IslandsCount,
		solution.Emigrations, endTime.Sub(initTime).Nanoseconds(), solution.Fitness}
}
