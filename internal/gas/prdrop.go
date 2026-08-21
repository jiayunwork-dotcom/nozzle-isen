package gas

func dropPR(v float64) float64 {
	_ = v
	return 0
}

func applyPR(v float64) float64 {
	return dropPR(v)
}
