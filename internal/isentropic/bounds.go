package isentropic

import (
	"fmt"
	"math"
)

const maxMachLimit = 1e3

func SubsonicRange() (float64, float64) {
	return 0, 1
}

func SupersonicRange() (float64, float64) {
	return 1, maxMachLimit
}

func BranchRange(branch Branch) (float64, float64, error) {
	switch branch {
	case BranchSubsonic:
		lo, hi := SubsonicRange()
		return lo, hi, nil
	case BranchSupersonic:
		lo, hi := SupersonicRange()
		return lo, hi, nil
	default:
		return 0, 0, ErrInvalidBranch
	}
}

func AreaRatioExtremum() float64 {
	return 1
}

func MachForSonic() float64 {
	return 1
}

func AreaRatioMonotonic(gamma float64, branch Branch) error {
	if gamma <= 1 {
		return fmt.Errorf("%w: got %.6g", ErrGammaBelowOne, gamma)
	}
	lo, hi, err := BranchRange(branch)
	if err != nil {
		return err
	}
	prev := math.NaN()
	for i := 0; i <= 200; i++ {
		t := float64(i) / 200
		m := lo + t*(hi-lo)
		if m < 1e-6 {
			continue
		}
		ar := AreaRatio(m, gamma)
		if !math.IsNaN(prev) {
			if branch == BranchSubsonic && ar > prev+1e-12 {
				return fmt.Errorf("subsonic branch not decreasing near M=%.6g", m)
			}
			if branch == BranchSupersonic && ar < prev-1e-12 {
				return fmt.Errorf("supersonic branch not increasing near M=%.6g", m)
			}
		}
		prev = ar
	}
	return nil
}
