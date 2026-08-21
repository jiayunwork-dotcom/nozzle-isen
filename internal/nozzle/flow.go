package nozzle

import (
	"nozzle-isen/internal/gas"
)

func ChokedFlux(st gas.State) float64 {
	return gas.ChokedFlux(st)
}

func ChokedMassFlow(st gas.State, throatArea float64) float64 {
	return gas.ChokedMassFlow(st, throatArea)
}

func IsentropicFlux(st gas.State, mach, area float64) float64 {
	return gas.IsentropicFlux(st, mach, area)
}

func UnchokedMassFlow(c Case, st gas.State, exitMach float64) float64 {
	exitArea := c.ThroatArea * c.AreaRatio
	return IsentropicFlux(st, exitMach, exitArea)
}

func ExitArea(c Case) float64 {
	return c.ThroatArea * c.AreaRatio
}

func FluxConsistency(st gas.State) float64 {
	return gas.FluxConsistency(st)
}
