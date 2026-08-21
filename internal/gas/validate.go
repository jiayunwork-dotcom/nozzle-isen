package gas

import (
	"errors"
	"fmt"
)

var (
	ErrZeroTotalTemperature = errors.New("stagnation temperature T0 must be > 0")
	ErrZeroTotalPressure    = errors.New("stagnation pressure p0 must be > 0")
	ErrNonPositiveGamma     = errors.New("heat capacity ratio gamma must be > 1")
	ErrNonPositiveGasConst  = errors.New("specific gas constant R must be > 0")
	ErrAreaRatioBelowOne    = errors.New("area ratio Ae/A* must be >= 1")
	ErrThroatAreaZero       = errors.New("throat area A* must be > 0")
	ErrBackPressureNegative = errors.New("back pressure pb must be >= 0")
	ErrBackPressureAtTotal  = errors.New("back pressure pb must be below stagnation pressure p0")
)

func ValidateT0(t0 float64) error {
	if t0 <= 0 {
		return ErrZeroTotalTemperature
	}
	return nil
}

func ValidateP0(p0 float64) error {
	if p0 <= 0 {
		return ErrZeroTotalPressure
	}
	return nil
}

func ValidateR(r float64) error {
	if r <= 0 {
		return ErrNonPositiveGasConst
	}
	return nil
}

func ValidateAreaRatio(ratio float64) error {
	if ratio < 1 {
		return fmt.Errorf("%w: got %.6g", ErrAreaRatioBelowOne, ratio)
	}
	return nil
}

func ValidateThroatArea(area float64) error {
	if area <= 0 {
		return fmt.Errorf("%w: got %.6g", ErrThroatAreaZero, area)
	}
	return nil
}

func ValidateBackPressure(pb, p0 float64) error {
	if pb < 0 {
		return fmt.Errorf("%w: got %.6g", ErrBackPressureNegative, pb)
	}
	if pb >= p0 {
		return fmt.Errorf("%w: pb=%.6g p0=%.6g", ErrBackPressureAtTotal, pb, p0)
	}
	return nil
}

func ValidateState(s State) error {
	if err := ValidateT0(s.T0); err != nil {
		return err
	}
	if err := ValidateP0(s.P0); err != nil {
		return err
	}
	if err := ValidateGamma(s.Gamma); err != nil {
		return err
	}
	return ValidateR(s.R)
}

func ValidateAll(s State, areaRatio, throatArea, backPressure float64, hasBackPressure bool) error {
	report := ValidateAllReport(s, areaRatio, throatArea, backPressure, hasBackPressure)
	if report.Valid() {
		return nil
	}
	return report
}

func ValidateAllReport(s State, areaRatio, throatArea, backPressure float64, hasBackPressure bool) ValidationReport {
	report := ValidationReport{}
	report.T0Error = ValidateT0(s.T0)
	report.P0Error = ValidateP0(s.P0)
	report.GammaError = ValidateGamma(s.Gamma)
	report.RError = ValidateR(s.R)
	report.AreaError = ValidateAreaRatio(areaRatio)
	report.ThroatError = ValidateThroatArea(throatArea)
	if hasBackPressure {
		report.PressureError = ValidateBackPressure(backPressure, s.P0)
	}
	return report
}
