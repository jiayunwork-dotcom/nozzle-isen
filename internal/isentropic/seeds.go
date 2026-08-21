package isentropic

import "math"

func subsonicSeed(target float64) float64 {
	scale := 1 + 0.4*math.Log(target)
	if scale < 1 {
		scale = 1
	}
	return 0.45 / scale
}

func supersonicSeed(gamma, target float64) float64 {
	logT := math.Log(target)
	if logT < 0 {
		logT = 0
	}
	return 1.2 + 0.45*logT
}
