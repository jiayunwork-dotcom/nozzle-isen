package gas

import "fmt"

func RatioString(value float64) string {
	return fmt.Sprintf("%.6g", value)
}

func StateString(s State) string {
	return fmt.Sprintf("T0=%.3g K p0=%.4g Pa gamma=%.4g R=%.3g J/(kg.K)",
		s.T0, s.P0, s.Gamma, s.R)
}
