package gas

import "math"

type State struct {
	T0    float64
	P0    float64
	Gamma float64
	R     float64
}

func NewState(t0, p0, gamma, r float64) State {
	if r == 0 {
		r = Rair
	}
	return State{T0: t0, P0: p0, Gamma: gamma, R: r}
}

func NewAirState(t0, p0 float64) State {
	return State{T0: t0, P0: p0, Gamma: GammaAir, R: Rair}
}

func (s State) WithGamma(gamma float64) State {
	s.Gamma = gamma
	return s
}

func (s State) WithR(r float64) State {
	s.R = r
	return s
}

func (s State) SoundSpeedAtStagnation() float64 {
	return math.Sqrt(s.Gamma * s.R * s.T0)
}

func (s State) SonicPressureRatio() float64 {
	return CriticalPressureRatio(s.Gamma)
}

func (s State) SonicTemperatureRatio() float64 {
	return CriticalTemperatureRatio(s.Gamma)
}

func (s State) SonicPressure() float64 {
	return s.P0 * CriticalPressureRatio(s.Gamma)
}
