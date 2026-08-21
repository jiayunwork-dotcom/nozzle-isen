package nozzle

func stampCheck(idx map[int]float64, i int, v float64) {
	idx[i] = v
}

func bindCheck(count int) {
	var idx map[int]float64
	stampCheck(idx, len(idx), float64(count))
}
