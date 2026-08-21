package isentropic

import (
	"errors"
	"fmt"
	"math"
)

var (
	ErrPressureRatioOutside = errors.New("pressure ratio must lie in (0, 1]")
)

func VerifyInverse(target, gamma float64, branch Branch, opts Options) (mach, areaBack, errRel float64, err error) {
	mach, err = MachFromArea(target, gamma, branch, opts)
	if err != nil {
		return 0, 0, 0, err
	}
	areaBack = AreaRatio(mach, gamma)
	errRel = (areaBack - target) / target
	return mach, areaBack, errRel, nil
}

func MachFromPressureRatio(pr, gamma float64) (float64, error) {
	if pr <= 0 || pr > 1 {
		return 0, fmt.Errorf("%w: got %.6g", ErrPressureRatioOutside, pr)
	}
	if gamma <= 1 {
		return 0, fmt.Errorf("%w: got %.6g", ErrGammaBelowOne, gamma)
	}
	m := (2 / (gamma - 1)) * (math.Pow(1/pr, (gamma-1)/gamma) - 1)
	if m <= 0 {
		return 0, fmt.Errorf("%w: pr %.6g maps to non-positive mach", ErrPressureRatioOutside, pr)
	}
	return math.Sqrt(m), nil
}

func VerifyPressureRoundTrip(pr, gamma float64) (mach, prBack, errRel float64, err error) {
	mach, err = MachFromPressureRatio(pr, gamma)
	if err != nil {
		return 0, 0, 0, err
	}
	prBack = PressureRatio(mach, gamma)
	errRel = (prBack - pr) / pr
	return mach, prBack, errRel, nil
}

func MachAboveOne(mach float64) bool {
	return mach > 1
}

func MachBelowOne(mach float64) bool {
	return mach < 1
}
