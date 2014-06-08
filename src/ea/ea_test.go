package ea

import (
	"testing"
	"fmt"
	"sort"
)

func ExampleEnhanceLists() {
	pop := []int{5, 4, 3, 2}
	n := len(pop)
	res := make([]int, n*(n+1)/2)
	indx := 0
	for i, indEval := range pop {
		for j := 0; j < n-i; j++ {
			res[indx] = indEval
			indx++
		}
	}
	fmt.Println(res)

	// Output: [5 5 5 5 4 4 4 3 3 2]
}

func ExampleMerge() {
	v := []int{5, 4, 3, 3, -1, -2, -10}
	u := []int{11, 9, 7, 2, 1, -3, -5}
	l := len(u) + len(v)
	a := make([]int, l)
	i, j, k := 0, 0, 0
	for i < l {
		if j < len(v) && k < len(u) {
			if v[j] > u[k] {
				a[i] = v[j]
				j++
			}else {
				a[i] = u[k]
				k++
			}
		}else {
			if j >= len(v) {
				for k < len(u) {
					a[i] = u[k]
					i++
					k++
				}
			} else {
				for j < len(v) {
					a[i] = v[j]
					i++
					j++
				}
			}
		}

		i++
	}
	fmt.Println(a)

	// Output: [11 9 7 5 4 3 3 2 1 -1 -2 -3 -5 -10]
}

func ExampleCrossover() {
	p := Pair{[]rune{1, 0, 0, 1, 1, 0},
		[]rune{1, 1, 0, 0, 1, 1}}

	indLength := len(p.a)
	cPoint := 2
	res1 := make(TIndividual, indLength)
	res2 := make(TIndividual, indLength)

	for i := 0; i < cPoint; i++ {
		res1[i] = p.a[i]
		res2[i] = p.b[i]
	}
	for i := cPoint; i < indLength; i++ {
		res1[i] = p.b[i]
		res2[i] = p.a[i]
	}

	fmt.Println(res1)
	fmt.Println(res2)

	// Output:
	// [1 0 0 0 1 1]
	// [1 1 0 1 1 0]
}

func ExampleSortIndEval() {
	p := TIndsEvaluated{TIndEval{[]rune{1,0,1,1}, 3}, TIndEval{[]rune{1,0,0,0}, 1},
		TIndEval{[]rune{1,1,1,1}, 4}, TIndEval{[]rune{0,1,0,1}, 2},
		TIndEval{[]rune{1,0,0,0}, 1}}

	sort.Sort(p)
	fmt.Println(p)

	// Output:
	// [{[1 1 1 1] 4} {[1 0 1 1] 3} {[0 1 0 1] 2} {[1 0 0 0] 1} {[1 0 0 0] 1}]
}

func Test2(t *testing.T) {
	//	t.Fail()
	t.SkipNow()
}

func Test_SeqCEvals(t *testing.T) {
	conf := SeqCEvals{SeqConf{[]TIndividual{[]rune{1, 0, 1, 0, 1, 0, 0, 0}, []rune{1, 0, 1, 0, 1, 1, 0, 1}, []rune{1, 0, 1, 0, 1, 1, 0, 1},
		[]rune{1, 1, 1, 0, 1, 1, 0, 1}, []rune{1, 0, 1, 0, 1, 1, 0, 0}, []rune{0, 0, 1, 0, 1, 1, 1, 1}},
	MaxOne,
	0.3},
	CEvalsConf{20}}

	solution := conf.Run()

//fmt.Println("La mejor solución es: ", solution)

	if solution.fitness < 4{
		t.Fail()
	}

}
