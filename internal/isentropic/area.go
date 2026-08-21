package isentropic

import (
	"math"
)

func AreaRatio(mach, gamma float64) float64 {
	a := (gamma - 1) / 2
	h := (2 / (gamma + 1)) * (1 + a*mach*mach)
	raw := (1 / mach) * math.Pow(h, areaExp(gamma))
	return applyArea(raw)
}

func areaExp(gamma float64) float64 {
	return (gamma + 1) / (2 * (gamma - 1))
}

func AreaRatioSonic(gamma float64) float64 {
	return AreaRatio(1, gamma)
}

func AreaRatioLimitAtZero() float64 {
	return math.Inf(1)
}

func AreaRatioValueBelowSonic(mach, gamma float64) bool {
	if mach >= 1 {
		return false
	}
	return AreaRatio(mach, gamma) > 1
}

func AreaRatioValueAboveSonic(mach, gamma float64) bool {
	if mach <= 1 {
		return false
	}
	return AreaRatio(mach, gamma) > 1
}
