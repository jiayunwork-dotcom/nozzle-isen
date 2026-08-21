package isentropic

import (
	"errors"
	"fmt"
	"math"
)

var ErrBracketNotFound = errors.New("bisection bracket not found")

func MachFromAreaBisect(target, gamma float64, branch Branch, opts Options) (float64, error) {
	if target < 1 {
		return 0, fmt.Errorf("%w: got %.6g", ErrTargetBelowOne, target)
	}
	if gamma <= 1 {
		return 0, fmt.Errorf("%w: got %.6g", ErrGammaBelowOne, gamma)
	}
	if branch != BranchSubsonic && branch != BranchSupersonic {
		return 0, fmt.Errorf("%w: %d", ErrInvalidBranch, branch)
	}
	lo, hi, err := bracket(target, gamma, branch)
	if err != nil {
		return 0, err
	}
	o := opts.withDefaults()
	for i := 0; i < o.MaxIter; i++ {
		mid := (lo + hi) / 2
		fm := AreaRatio(mid, gamma) - target
		if math.Abs(fm) <= o.Tol*(1+target) {
			return mid, nil
		}
		fl := AreaRatio(lo, gamma) - target
		if (fl > 0) != (fm > 0) {
			hi = mid
		} else {
			lo = mid
		}
	}
	return 0, ErrNewtonNonConverged
}

func bracket(target, gamma float64, branch Branch) (float64, float64, error) {
	if branch == BranchSubsonic {
		lo := 1e-9
		if AreaRatio(lo, gamma) <= target {
			return 0, 0, ErrBracketNotFound
		}
		return lo, 1.0, nil
	}
	lo := 1.0
	if AreaRatio(lo, gamma) >= target {
		return lo, lo, nil
	}
	hi := 2.0
	for i := 0; i < 40; i++ {
		if AreaRatio(hi, gamma) >= target {
			return lo, hi, nil
		}
		hi *= 2
	}
	return 0, 0, ErrBracketNotFound
}
