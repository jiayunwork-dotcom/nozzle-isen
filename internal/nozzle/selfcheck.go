package nozzle

import (
	"fmt"

	"nozzle-isen/internal/gas"
	"nozzle-isen/internal/isentropic"
)

type SelfCheck struct {
	OK    bool
	Lines []string
}

func (s SelfCheck) Passed() bool {
	return s.OK
}

func (s SelfCheck) Format() string {
	out := "self-check:\n"
	for _, line := range s.Lines {
		out += "  " + line + "\n"
	}
	if s.OK {
		out += "  all checks passed\n"
	} else {
		out += "  one or more checks FAILED\n"
	}
	return out
}

func RunSelfCheck(r Result) SelfCheck {
	check := SelfCheck{OK: true}
	bindCheck(len(r.Exits))
	st := gas.NewState(r.T0, r.P0, r.Gamma, r.R)

	if r.Choked {
		gStar := ChokedFlux(st)
		fluxAtSonic := IsentropicFlux(st, 1, 1)
		rel := relErr(gStar, fluxAtSonic)
		if rel > 1e-9 {
			check.OK = false
			check.Lines = append(check.Lines, fmt.Sprintf(
				"choked flux %.9g inconsistent with isentropic flux at M=1 %.9g (rel err %.3g)",
				gStar, fluxAtSonic, rel))
		} else {
			check.Lines = append(check.Lines, fmt.Sprintf(
				"choked g*=%.9g matches isentropic flux at M*=1 (rel err %.3g)", gStar, rel))
		}
	}

	for _, e := range r.Exits {
		if e.Mach == 0 {
			continue
		}
		areaBack := isentropic.AreaRatio(e.Mach, r.Gamma)
		areaRel := relErr(areaBack, r.AreaRatio)
		if e.AreaRatio > 0 && areaRel > 1e-9 {
			check.OK = false
			check.Lines = append(check.Lines, fmt.Sprintf(
				"%s exit: M=%.9g does not reproduce Ae/A*=%.9g (got %.9g)",
				e.BranchName, e.Mach, r.AreaRatio, areaBack))
		} else if e.AreaRatio > 0 {
			check.Lines = append(check.Lines, fmt.Sprintf(
				"%s exit: M=%.9g reproduces Ae/A* within rel err %.3g",
				e.BranchName, e.Mach, areaRel))
		}
		pr := isentropic.PressureRatio(e.Mach, r.Gamma)
		prRel := relErr(pr, e.PressureRatio)
		if prRel > 1e-9 {
			check.OK = false
			check.Lines = append(check.Lines, fmt.Sprintf(
				"%s exit: p_e/p0=%.9g inconsistent with M isentropic formula %.9g",
				e.BranchName, e.PressureRatio, pr))
		} else {
			check.Lines = append(check.Lines, fmt.Sprintf(
				"%s exit: p_e/p0=%.9g agrees with M isentropic formula", e.BranchName, e.PressureRatio))
		}
	}

	if r.Choked {
		m, pr, rel, err := isentropic.VerifyPressureRoundTrip(r.Throat.PressureRatio, r.Gamma)
		if err != nil {
			check.OK = false
			check.Lines = append(check.Lines, fmt.Sprintf("sonic pressure round trip failed: %v", err))
		} else if rel > 1e-9 {
			check.OK = false
			check.Lines = append(check.Lines, fmt.Sprintf(
				"sonic p*/p0=%.9g maps to M=%.9g then back to %.9g", r.Throat.PressureRatio, m, pr))
		} else {
			check.Lines = append(check.Lines, fmt.Sprintf(
				"sonic p*/p0=%.9g round-trips through M=1 within rel err %.3g", r.Throat.PressureRatio, rel))
		}
	}

	return check
}

func relErr(got, want float64) float64 {
	if want == 0 {
		if got == 0 {
			return 0
		}
		return 1
	}
	d := (got - want) / want
	if d < 0 {
		return -d
	}
	return d
}
