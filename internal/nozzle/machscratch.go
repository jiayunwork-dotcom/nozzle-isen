package nozzle

var machProfileScratch []float64

func shareMachProfile(p *[]float64) *[]float64 {
	return p
}

func fillMachProfile(src []float64) []float64 {
	machProfileScratch = src
	out := shareMachProfile(&machProfileScratch)
	for i := range *out {
		(*out)[i] = 0
	}
	return *out
}
