package ea

import (
	"math/rand"
	"encoding/json"
	"io/ioutil"
	"time"
	"strings"
)

type TRes struct{
	Evaluations    int
	PopSize        int
	ChromosomeSize int
}
type TSolution struct{
	Duration int64
	BestFit  int
}
type SeqRes struct{
	TRes
	TSolution
}

func genIndividual(n int) TIndividual {
	res := make(TIndividual, n)
	for i, _ := range res {
		if rand.Intn(100)%2 == 0 {
			res[i] = 1
		} else {
			res[i] = 0
		}
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

type ParRes struct{
	TRes
	EvaluatorsCapacity,
	ReproducersCapacity,
	EvaluatorsCount,
	ReproducersCount,
	IslandsCount int
	TSolution
}

type Data struct{
	EvaluatorsCount,
	ReproducersCount,
	EvaluatorsCapacity,
	ReproducersCapacity,
	Evaluations,
	PopSize,
	ChromosomeSize,
	IslandsCount         int
	PMutation            float32
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

type Problem struct{
	problemConfig *Data
}

type IProblem interface {
	QualityFitnessFunction(v int) bool
	DoWhenQualityFitnessTrue(i TIndEval)
	FitnessFunction(ind TIndividual) int
}

func (self *Problem) GenInitPop() TPopulation {
	return genPop(self.problemConfig.PopSize, self.problemConfig.ChromosomeSize)
}

func (self *Problem) runSeqCEvals(fitnessFunction TFitnessFunc) *SeqRes {
	obj := SeqCEvals{SeqConf{self.GenInitPop, fitnessFunction, self.problemConfig.PMutation},
		CEvalsConf{self.problemConfig.Evaluations}}
	initTime := time.Now()
	solution := obj.Run()
	endTime := time.Now()
	return &SeqRes{TRes{self.problemConfig.Evaluations, self.problemConfig.PopSize, self.problemConfig.ChromosomeSize},
		TSolution{endTime.Sub(initTime).Nanoseconds(), solution.Fitness}}
}

func (self *Problem) runSeqFitnessQuality(fitnessFunction TFitnessFunc, qf TQualityF, df Tdo) *SeqRes {
	obj := SeqFitnessQuality{SeqConf{self.GenInitPop, fitnessFunction, self.problemConfig.PMutation}    ,
		FitnessQualityConf{qf, df}}
	initTime := time.Now()
	solution, cEvals := obj.Run()
	endTime := time.Now()
	return &SeqRes{TRes{cEvals, self.problemConfig.PopSize, self.problemConfig.ChromosomeSize},
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

func (self *MaxOneProblem) RunSeqFitnessQuality() *SeqRes {
	return self.Problem.runSeqFitnessQuality(self.FitnessFunction, self.QualityFitnessFunction, self.DoWhenQualityFitnessTrue)
}

func (self *MaxOneProblem) RunParCEvals() *ParRes {
	return self.Problem.runParCEvals(self.FitnessFunction)
}

func (self *MaxOneProblem) RunParFitnessQuality() *ParRes {
	return self.Problem.runParFitnessQuality(self.FitnessFunction, self.QualityFitnessFunction, self.DoWhenQualityFitnessTrue)
}

func (self *Problem) runParCEvals(fitnessFunction TFitnessFunc) *ParRes {
	obj := ParCEvals{ParConf{SeqConf{self.GenInitPop, fitnessFunction, self.problemConfig.PMutation}    ,
		self.problemConfig.EvaluatorsCapacity, self.problemConfig.ReproducersCapacity, self.problemConfig.EvaluatorsCount,
		self.problemConfig.ReproducersCount, self.problemConfig.IslandsCount}, CEvalsConf{self.problemConfig.Evaluations}}
	initTime := time.Now()
	solution := obj.Run()
	endTime := time.Now()
	return &ParRes{TRes{self.problemConfig.Evaluations, self.problemConfig.PopSize, self.problemConfig.ChromosomeSize}, self.problemConfig.EvaluatorsCapacity, self.problemConfig.ReproducersCapacity,
		self.problemConfig.EvaluatorsCount, self.problemConfig.ReproducersCount, self.problemConfig.IslandsCount, TSolution{endTime.Sub(initTime).Nanoseconds(), solution.Fitness}}
}

func (self *Problem) runParFitnessQuality(fitnessFunction TFitnessFunc, qf TQualityF, df Tdo) *ParRes {
	obj := ParFitnessQuality{ParConf{SeqConf{self.GenInitPop, fitnessFunction, self.problemConfig.PMutation}    ,
		self.problemConfig.EvaluatorsCapacity, self.problemConfig.ReproducersCapacity, self.problemConfig.EvaluatorsCount,
		self.problemConfig.ReproducersCount, self.problemConfig.IslandsCount},
		FitnessQualityConf{qf, df}}
	initTime := time.Now()
	solution, cEvals := obj.Run()
	endTime := time.Now()
	return &ParRes{TRes{cEvals, self.problemConfig.PopSize, self.problemConfig.ChromosomeSize}, self.problemConfig.EvaluatorsCapacity, self.problemConfig.ReproducersCapacity,
		self.problemConfig.EvaluatorsCount, self.problemConfig.ReproducersCount, self.problemConfig.IslandsCount, TSolution{endTime.Sub(initTime).Nanoseconds(), solution.Fitness}}
}
