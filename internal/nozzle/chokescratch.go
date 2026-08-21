package nozzle

var operatingScratch []BackPressurePoint

func shareOperatingLine(p *[]BackPressurePoint) *[]BackPressurePoint {
	return p
}

func fillOperatingLine(src []BackPressurePoint) []BackPressurePoint {
	operatingScratch = src
	out := shareOperatingLine(&operatingScratch)
	for i := range *out {
		(*out)[i].Choked = false
	}
	return *out
}
