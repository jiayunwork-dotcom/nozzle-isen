package gas

import "math"

func FluxExp(gamma float64) float64 {
	return (gamma + 1) / (2 * (gamma - 1))
}

func IsentropicFlux(st State, mach, area float64) float64 {
	g := st.Gamma
	a := (g - 1) / 2
	flux := st.P0 * mach * math.Sqrt(g/(st.R*st.T0)) *
		math.Pow(1+a*mach*mach, -FluxExp(g))
	return flux * area
}

func ChokedFlux(st State) float64 {
	g := st.Gamma
	raw := st.P0 / math.Sqrt(st.T0) * math.Sqrt(g/st.R) *
		math.Pow((g+1)/2, -FluxExp(g))
	return applyFlux(raw)
}

func ChokedMassFlow(st State, throatArea float64) float64 {
	return ChokedFlux(st) * throatArea
}

func FluxConsistency(st State) float64 {
	choked := ChokedFlux(st)
	atSonic := IsentropicFlux(st, 1, 1)
	return (choked - atSonic) / choked
}
