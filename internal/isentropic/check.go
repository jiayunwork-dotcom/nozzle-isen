package isentropic

import (
	"errors"
	"fmt"
	"math"
)

var (
	ErrMachBelowZero  = errors.New("mach must be > 0")
	ErrMachNotFinite  = errors.New("mach must be finite")
	ErrGammaBelowOne  = errors.New("heat capacity ratio must be > 1")
	ErrTargetBelowOne = errors.New("area ratio target must be >= 1")
	ErrInvalidBranch  = errors.New("unknown expansion branch")
)

func SonicAreaRatio(gamma float64) float64 {
	return AreaRatio(1, gamma)
}

func SonicPressureRatio(gamma float64) float64 {
	return PressureRatio(1, gamma)
}

func CheckMach(mach float64) error {
	if math.IsNaN(mach) || math.IsInf(mach, 0) {
		return ErrMachNotFinite
	}
	if mach <= 0 {
		return ErrMachBelowZero
	}
	return nil
}

func CheckGamma(gamma float64) error {
	if gamma <= 1 {
		return fmt.Errorf("%w: got %.6g", ErrGammaBelowOne, gamma)
	}
	return nil
}

func CheckTarget(target float64) error {
	if target < 1 {
		return fmt.Errorf("%w: got %.6g", ErrTargetBelowOne, target)
	}
	return nil
}

func CheckBranch(branch Branch) error {
	if branch != BranchSubsonic && branch != BranchSupersonic {
		return fmt.Errorf("%w: %d", ErrInvalidBranch, branch)
	}
	return nil
}

func SonicConsistency(gamma float64) error {
	if err := CheckGamma(gamma); err != nil {
		return err
	}
	if math.Abs(AreaRatio(1, gamma)-1) > 1e-12 {
		return fmt.Errorf("area ratio at M=1 is %.12g, expected 1", AreaRatio(1, gamma))
	}
	if math.Abs(PressureRatio(1, gamma)-math.Pow(2/(gamma+1), gamma/(gamma-1))) > 1e-12 {
		return fmt.Errorf("pressure ratio at M=1 inconsistent with sonic ratio")
	}
	return nil
}
