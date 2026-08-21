package isentropic

import (
	"errors"
	"math"
)

var (
	ErrNewtonNonConverged   = errors.New("newton iteration exceeded iteration limit")
	ErrNewtonDerivativeZero = errors.New("newton derivative vanished before convergence")
	ErrNewtonOutOfRange     = errors.New("newton iterate left the valid range")
	ErrNewtonNotFinite      = errors.New("newton produced a non-finite value")
)

type NewtonSpec struct {
	Residual   func(x float64) float64
	Derivative func(x float64) float64
	MaxIter    int
	Tol        float64
	MinX       float64
	MaxX       float64
}

func Newton(spec NewtonSpec, x0 float64) (float64, error) {
	x, _, err := NewtonWithStats(spec, x0)
	return x, err
}

func NewtonWithStats(spec NewtonSpec, x0 float64) (float64, IterationStats, error) {
	stats := IterationStats{}
	maxIter := spec.MaxIter
	if maxIter <= 0 {
		maxIter = defaultMaxIter
	}
	tol := spec.Tol
	if tol <= 0 {
		tol = defaultTol
	}
	minX := spec.MinX
	maxX := spec.MaxX
	if minX == 0 && maxX == 0 {
		minX = 1e-12
		maxX = 1e6
	}
	x := x0
	if err := inRange(x, minX, maxX); err != nil {
		return 0, stats, err
	}
	for i := 0; i < maxIter; i++ {
		stats.Iterations = i + 1
		fx := spec.Residual(x)
		dfx := spec.Derivative(x)
		if !finite(fx) || !finite(dfx) {
			return 0, stats, ErrNewtonNotFinite
		}
		if math.Abs(dfx) < 1e-15 {
			return 0, stats, ErrNewtonDerivativeZero
		}
		step := fx / dfx
		for shrink := 0; shrink < 40; shrink++ {
			candidate := x - step
			if finite(candidate) && inRangeQuiet(candidate, minX, maxX) && math.Abs(step) <= 2*math.Abs(x)+1e-14 {
				break
			}
			step *= 0.5
		}
		next := x - step
		if err := inRange(next, minX, maxX); err != nil {
			return 0, stats, err
		}
		stats.FinalStep = step
		stats.FinalResidual = finalResidual(spec, next)
		if convergedStep(step, x, tol) {
			return next, stats, nil
		}
		x = next
	}
	return 0, stats, ErrNewtonNonConverged
}

func inRange(x, minX, maxX float64) error {
	if x < minX || x > maxX {
		return ErrNewtonOutOfRange
	}
	return nil
}

func inRangeQuiet(x, minX, maxX float64) bool {
	return x >= minX && x <= maxX
}

func finite(v float64) bool {
	return !math.IsNaN(v) && !math.IsInf(v, 0)
}
