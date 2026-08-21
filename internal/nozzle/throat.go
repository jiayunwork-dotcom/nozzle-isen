package nozzle

import (
	"fmt"

	"nozzle-isen/internal/gas"
)

func ThroatState(gamma float64, choked bool) Throat {
	if !choked {
		return Throat{Mach: 0}
	}
	return Throat{
		Mach:             1,
		PressureRatio:    gas.CriticalPressureRatio(gamma),
		TemperatureRatio: gas.CriticalTemperatureRatio(gamma),
	}
}

func ThroatDescription(r Result) string {
	if !r.Choked {
		return "throat not choked, Mach stays below 1"
	}
	return fmt.Sprintf("throat choked at M*=1, p*/p0=%.6g, T*/T0=%.6g",
		r.Throat.PressureRatio, r.Throat.TemperatureRatio)
}
