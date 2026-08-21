package nozzle

import (
	"fmt"
)

type ExpansionClass int

const (
	ExpansionMatched ExpansionClass = iota
	ExpansionUnder
	ExpansionOver
)

func (e ExpansionClass) String() string {
	switch e {
	case ExpansionUnder:
		return "underexpanded"
	case ExpansionOver:
		return "overexpanded"
	default:
		return "matched"
	}
}

type Throat struct {
	Mach             float64
	PressureRatio    float64
	TemperatureRatio float64
}

type Exit struct {
	BranchName       string
	Mach             float64
	AreaRatio        float64
	PressureRatio    float64
	Pressure         float64
	TemperatureRatio float64
	Expansion        ExpansionClass
}

type Result struct {
	Case            Case
	T0              float64
	P0              float64
	Gamma           float64
	R               float64
	AreaRatio       float64
	ThroatArea      float64
	Choked          bool
	ChokedReason    string
	Throat          Throat
	Exits           []Exit
	MassFlow        float64
	MassFlowFormula string
	BackPressure    float64
	HasBackPressure bool
	ShockNote       string
	Checks          SelfCheck
}

func (r Result) ExitByName(name string) (Exit, bool) {
	for _, e := range r.Exits {
		if e.BranchName == name {
			return e, true
		}
	}
	return Exit{}, false
}

func (r Result) StateString() string {
	return fmt.Sprintf("T0=%.3g K p0=%.4g Pa gamma=%.4g R=%.3g J/(kg.K)",
		r.T0, r.P0, r.Gamma, r.R)
}

func (t Throat) String() string {
	return fmt.Sprintf("M=%.6g p/p0=%.6g T/T0=%.6g", t.Mach, t.PressureRatio, t.TemperatureRatio)
}

func (e Exit) String() string {
	return fmt.Sprintf("%s M=%.6g p_e/p0=%.6g", e.BranchName, e.Mach, e.PressureRatio)
}
