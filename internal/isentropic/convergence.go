package isentropic

import (
	"math"
)

type IterationStats struct {
	Iterations int
	FinalStep  float64
	FinalResidual float64
}

func residualScale(x float64) float64 {
	if math.Abs(x) < 1 {
		return 1
	}
	return math.Abs(x)
}

func convergedStep(step, x, tol float64) bool {
	return math.Abs(step) <= tol*residualScale(x)
}

func finalResidual(spec NewtonSpec, x float64) float64 {
	return math.Abs(spec.Residual(x))
}
