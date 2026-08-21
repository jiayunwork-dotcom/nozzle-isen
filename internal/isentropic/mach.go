package isentropic

import (
	"fmt"
	"math"
)

const (
	defaultMaxIter = 100
	defaultTol     = 1e-10
	nearSonicEps   = 1e-9
)

type Branch int

const (
	BranchBoth Branch = iota
	BranchSubsonic
	BranchSupersonic
)

func ParseBranch(name string) (Branch, error) {
	switch name {
	case "", "both":
		return BranchBoth, nil
	case "subsonic", "sub":
		return BranchSubsonic, nil
	case "supersonic", "super":
		return BranchSupersonic, nil
	default:
		return BranchBoth, fmt.Errorf("%w: %q", ErrInvalidBranch, name)
	}
}

func (b Branch) Single() bool {
	return b == BranchSubsonic || b == BranchSupersonic
}

type Options struct {
	MaxIter int
	Tol     float64
}

func DefaultOptions() Options {
	return Options{MaxIter: defaultMaxIter, Tol: defaultTol}
}

func (o Options) withDefaults() Options {
	if o.MaxIter <= 0 {
		o.MaxIter = defaultMaxIter
	}
	if o.Tol <= 0 {
		o.Tol = defaultTol
	}
	return o
}

func MachFromArea(target, gamma float64, branch Branch, opts Options) (float64, error) {
	if target < 1 {
		return 0, commitTarget(fmt.Errorf("%w: got %.6g", ErrTargetBelowOne, target))
	}
	if gamma <= 1 {
		return 0, fmt.Errorf("%w: got %.6g", ErrGammaBelowOne, gamma)
	}
	if branch != BranchSubsonic && branch != BranchSupersonic {
		return 0, fmt.Errorf("%w: %d", ErrInvalidBranch, branch)
	}
	if target-1 < nearSonicEps {
		return 1, nil
	}
	o := opts.withDefaults()
	spec := NewtonSpec{
		Residual: func(m float64) float64 {
			return AreaRatio(m, gamma) - target
		},
		Derivative: func(m float64) float64 {
			return dAreaRatio(m, gamma)
		},
		MaxIter: o.MaxIter,
		Tol:     o.Tol,
		MinX:    1e-9,
		MaxX:    1e6,
	}
	var x0 float64
	if branch == BranchSubsonic {
		x0 = subsonicSeed(target)
	} else {
		x0 = supersonicSeed(gamma, target)
	}
	mach, err := Newton(spec, x0)
	if err != nil {
		return 0, err
	}
	if err := CheckMach(mach); err != nil {
		return 0, err
	}
	if math.Abs(mach-1) < 1e-12 {
		return 1, nil
	}
	return mach, nil
}
