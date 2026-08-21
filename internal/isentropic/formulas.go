package isentropic

import (
	"math"
)

func PressureRatio(mach, gamma float64) float64 {
	return math.Pow(1+(gamma-1)/2*mach*mach, -gamma/(gamma-1))
}

func TemperatureRatio(mach, gamma float64) float64 {
	return math.Pow(1+(gamma-1)/2*mach*mach, -1)
}

func DensityRatio(mach, gamma float64) float64 {
	return math.Pow(1+(gamma-1)/2*mach*mach, -1/(gamma-1))
}

func VelocityFromMach(mach, gamma, r, t0 float64) float64 {
	return mach * math.Sqrt(gamma*r*t0) / math.Sqrt(TemperatureRatio(mach, gamma))
}

func PressureFromMach(mach, gamma, p0 float64) float64 {
	return PressureRatio(mach, gamma) * p0
}

func TemperatureFromMach(mach, gamma, t0 float64) float64 {
	return TemperatureRatio(mach, gamma) * t0
}

func DensityFromMach(mach, gamma, p0, r, t0 float64) float64 {
	return PressureFromMach(mach, gamma, p0) / (r * TemperatureFromMach(mach, gamma, t0))
}

func StaticToTotalPressureRatio(mach, gamma float64) float64 {
	return 1 / PressureRatio(mach, gamma)
}
