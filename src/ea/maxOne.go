package ea

func MaxOne(ind TIndividual) int {
	var res = 0
	for _, e := range ind {
		if e == 1 {
			res++
		}
	}
	return res
}
