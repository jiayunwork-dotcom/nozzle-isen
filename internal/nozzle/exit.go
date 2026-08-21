package nozzle

import (
	"nozzle-isen/internal/isentropic"
)

func ExitFromArea(target, p0, gamma float64, name string, branch isentropic.Branch, opts isentropic.Options) (Exit, error) {
	mach, err := isentropic.MachFromArea(target, gamma, branch, opts)
	if err != nil {
		return Exit{}, err
	}
	pr := isentropic.PressureRatio(mach, gamma)
	return Exit{
		BranchName:       name,
		Mach:             mach,
		AreaRatio:        target,
		PressureRatio:    pr,
		Pressure:         pr * p0,
		TemperatureRatio: isentropic.TemperatureRatio(mach, gamma),
	}, nil
}

func (e Exit) ExitPressure(p0 float64) float64 {
	return e.PressureRatio * p0
}

func (e Exit) ExitTemperature(t0 float64) float64 {
	return e.TemperatureRatio * t0
}

func (e Exit) MachToPressureRatio() float64 {
	return e.PressureRatio
}

func (e Exit) IsSupersonic() bool {
	return e.Mach > 1
}

func (e Exit) IsSubsonic() bool {
	return e.Mach < 1
}
