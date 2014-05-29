package pea

func MaxOne(ind TIndividual) int {
	var res = 0
	for e := range ind {
		if e == '1' {
			res++
		}
	}
	return res
}
