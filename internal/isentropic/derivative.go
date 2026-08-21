package isentropic

func dAreaRatio(mach, gamma float64) float64 {
	a := (gamma - 1) / 2
	ar := AreaRatio(mach, gamma)
	n := areaExp(gamma)
	return (ar / mach) * (n*(gamma-1)*mach*mach/(1+a*mach*mach) - 1)
}

func d2AreaRatio(mach, gamma float64) float64 {
	if mach <= 0 {
		return 0
	}
	h := 1e-4 * mach
	plus := dAreaRatio(mach+h, gamma)
	minus := dAreaRatio(mach-h, gamma)
	return (plus - minus) / (2 * h)
}

func DerivativeSign(mach, gamma float64) int {
	d := dAreaRatio(mach, gamma)
	switch {
	case d > 0:
		return 1
	case d < 0:
		return -1
	default:
		return 0
	}
}
