package nozzle

import "fmt"

func Summary(r Result) string {
	s := fmt.Sprintf("choked=%v qm=%.6g kg/s", r.Choked, r.MassFlow)
	for _, e := range r.Exits {
		s += fmt.Sprintf(" %s_M=%.6g", e.BranchName, e.Mach)
	}
	return s
}

func ChokedSummary(r Result) string {
	if r.Choked {
		return fmt.Sprintf("yes qm=%.6g pstar_p0=%.6g", r.MassFlow, r.Throat.PressureRatio)
	}
	return "no"
}

func ExitSummary(r Result, name string) string {
	e, ok := r.ExitByName(name)
	if !ok {
		return fmt.Sprintf("%s branch not present", name)
	}
	return fmt.Sprintf("%s M=%.6g p_e/p0=%.6g", name, e.Mach, e.PressureRatio)
}

func MassFlowSummary(r Result) string {
	return fmt.Sprintf("%.6g kg/s via %s", r.MassFlow, r.MassFlowFormula)
}
