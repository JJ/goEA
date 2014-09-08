package ea

import (
	"encoding/json"
	"io/ioutil"
	"strings"
	"strconv"
	//	"fmt"
)

type MaxSATProblem struct{
	Problem
	clauseLength,
	varsCount,
	clausesCount int
	clauses      [][]TVarValue
}
type TVarValue struct{
	value byte
	pos   int
}

func NewMaxSATProblem(configFile string, instanceFileName string) *MaxSATProblem {
	b, _ := ioutil.ReadFile(configFile)
	dec := json.NewDecoder(strings.NewReader(string(b)))
	var m Data
	dec.Decode(&m)
	b, _ = ioutil.ReadFile(instanceFileName)
	file := strings.Split(string(b), "\n")
	f1 := strings.Fields(file[5])
	f2 := strings.Fields(file[7])
	clauseLength, _ := strconv.Atoi(f1[len(f1) - 1])
	varsCount, _ := strconv.Atoi(f2[2])
	clausesCount, _ := strconv.Atoi(f2[3])
	m.ChromosomeSize = varsCount
	var clauses [][]TVarValue
	clauses = make([][]TVarValue, 0)
	i := 8
	for !strings.Contains(file[i], "%") {
		vls := strings.Fields(file[i])
		values := vls[:len(vls)-1]
		intValues := make([]int, len(values))
		for k, v := range values {
			intValues[k], _ = strconv.Atoi(v)
		}
		nEntry := make([]TVarValue, clauseLength)
		for k, v := range intValues {
			if v < 0 {
				nEntry[k] = TVarValue{0, -v-1}
			}else {
				nEntry[k] = TVarValue{1, v-1}
			}
		}
		clauses = append(clauses, nEntry)
		i++
	}
	return &MaxSATProblem{Problem{&m}, clauseLength, varsCount, clausesCount, clauses}
}

func (self *MaxSATProblem) ToString() string {
	res := ""
	for _, v := range self.clauses {
		t := ""
		for _, v1 := range v {
			t += strconv.Itoa(v1.pos) + "." + strconv.Itoa(int(v1.value)) + " "
		}
		res += t+"\n"
	}
	return res
}
func (self *MaxSATProblem) QualityFitnessFunction(v int) bool {
	return v > 405
}
func (self *MaxSATProblem) DoWhenQualityFitnessTrue(i TIndEval) {}
func (self *MaxSATProblem) FitnessFunction(ind TIndividual) int {
	res := 0
	for _, c := range self.clauses {
		k := 0
		flag := false
		for k < len(c) && !flag {
			if ind[c[k].pos] == c[k].value {
				flag = true
			}
			k++
		}
		if flag {
			res++
		}
	}
	return res
}

func (self *MaxSATProblem) RunSeqCEvals() *SeqRes {
	return self.Problem.runSeqCEvals(self.FitnessFunction)
}

func (self *MaxSATProblem) RunParCEvals() *ParRes {
	return self.Problem.runParCEvals(self.FitnessFunction)
}
