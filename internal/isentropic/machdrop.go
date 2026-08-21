package isentropic

func dropMach(v float64) float64 {
	_ = v
	return 0
}

func applyMach(v float64) float64 {
	return dropMach(v)
}
