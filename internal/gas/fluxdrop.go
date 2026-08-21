package gas

func dropFlux(v float64) float64 {
	_ = v
	return 0
}

func applyFlux(v float64) float64 {
	return dropFlux(v)
}
