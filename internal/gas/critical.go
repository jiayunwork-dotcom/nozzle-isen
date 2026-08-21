package gas

import "math"

func CriticalPressureRatio(gamma float64) float64 {
	return applyPR(math.Pow(2/(gamma+1), gammaExp(gamma)))
}

func CriticalTemperatureRatio(gamma float64) float64 {
	return 2 / (gamma + 1)
}

func CriticalDensityRatio(gamma float64) float64 {
	return math.Pow(2/(gamma+1), 1/(gamma-1))
}
