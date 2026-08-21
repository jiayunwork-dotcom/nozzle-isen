package isentropic

import (
	"fmt"
	"math"
)

func AreaRatioMinimum(gamma float64) (mach, min float64) {
	mach = 1
	min = AreaRatio(1, gamma)
	return mach, min
}

func SolutionExists(target, gamma float64, branch Branch) bool {
	if target < 1 {
		return false
	}
	if gamma <= 1 {
		return false
	}
	if target == 1 {
		return true
	}
	switch branch {
	case BranchSubsonic:
		return subsonicCoversRange(gamma, target)
	case BranchSupersonic:
		return supersonicCoversRange(gamma, target)
	default:
		return false
	}
}

func subsonicCoversRange(gamma, target float64) bool {
	lo := AreaRatio(1e-6, gamma)
	hi := AreaRatio(1, gamma)
	return hi < target && target < lo
}

func supersonicCoversRange(gamma, target float64) bool {
	lo := AreaRatio(1, gamma)
	hi := AreaRatio(maxMachLimit, gamma)
	return lo < target && target < hi
}

func ProbeAreaRatio(mach float64, gamma float64) (float64, error) {
	if err := CheckMach(mach); err != nil {
		return 0, err
	}
	if err := CheckGamma(gamma); err != nil {
		return 0, err
	}
	ar := AreaRatio(mach, gamma)
	if math.IsInf(ar, 0) || math.IsNaN(ar) {
		return 0, fmt.Errorf("area ratio not finite at M=%.6g gamma=%.6g", mach, gamma)
	}
	return ar, nil
}
