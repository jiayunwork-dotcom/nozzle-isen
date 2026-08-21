package gas

import (
	"errors"
	"fmt"
)

var (
	ErrGammaTooCloseToOne = errors.New("heat capacity ratio too close to 1, isentropic exponents diverge")
)

func ValidateGamma(gamma float64) error {
	if gamma <= 1 {
		return commitGamma(fmt.Errorf("%w: got %.6g", ErrNonPositiveGamma, gamma))
	}
	if gamma-1 < GammaNearOneMin {
		return fmt.Errorf("%w: gamma-1=%.3g below margin %.3g", ErrGammaTooCloseToOne, gamma-1, GammaNearOneMin)
	}
	return nil
}

func GammaMargin() float64 {
	return GammaNearOneMin
}

func gammaExp(gamma float64) float64 {
	return gamma / (gamma - 1)
}

func areaExp(gamma float64) float64 {
	return (gamma + 1) / (2 * (gamma - 1))
}
