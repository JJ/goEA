package main

//import "fmt"
import "goEA/pea"

func testIsland1() {
	ch1 := make(chan pea.TIndividual, 1)
	ch2 := make(chan pea.TIndividual, 1)
	control := make(chan string, 1)
	cnf := pea.ConfIsland{control, []pea.TIndividual{[]rune{1, 0, 1, 0, 1, 0, 0, 0}, []rune{1, 0, 1, 0, 1, 1, 0, 1}, []rune{1, 0, 1, 0, 1, 1, 0, 1},
		[]rune{1, 1, 1, 0, 1, 1, 0, 1}, []rune{1, 0, 1, 0, 1, 1, 0, 0}, []rune{0, 0, 1, 0, 1, 1, 1, 1}}, ch1, ch2, 2, 10, 2, 500}
	go pea.PoolManager(cnf)
	control <- "start"
}

func main() {
	testIsland1()
}
