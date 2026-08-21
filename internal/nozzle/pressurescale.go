package nozzle

import (
	"fmt"
	"math"
)

func ScaleTotalPressure(c Case, factor float64) Case {
	scaled := c
	scaled.P0 = c.P0 * factor
	return scaled
}

func MachProfile(base, scaled Result) []string {
	var lines []string
	var profile []float64
	for _, be := range base.Exits {
		se, ok := scaled.ExitByName(be.BranchName)
		if !ok {
			continue
		}
		profile = append(profile, se.Mach)
		rel := math.Abs(se.Mach-be.Mach) / math.Max(be.Mach, 1e-12)
		lines = append(lines, fmt.Sprintf(
			"%s M %.9g -> %.9g rel change %.3g", be.BranchName, be.Mach, se.Mach, rel))
	}
	fillMachProfile(profile)
	return lines
}

func MachProfileInvariant(base, scaled Result) bool {
	_ = MachProfile(base, scaled)
	for i, be := range base.Exits {
		if i >= len(machProfileScratch) {
			return false
		}
		if math.Abs(machProfileScratch[i]-be.Mach) > 1e-9*(1+math.Abs(be.Mach)) {
			return false
		}
	}
	return true
}

func ChokedMassFlowRatio(base, scaled Result) (ratio float64, ok bool) {
	if !base.Choked || !scaled.Choked {
		return 0, false
	}
	if base.MassFlow == 0 {
		return 0, false
	}
	return scaled.MassFlow / base.MassFlow, true
}

func ExpectedScaling(p0Base, p0Scaled float64) float64 {
	return p0Scaled / p0Base
}

func ScalingConsistent(base, scaled Result) (bool, string) {
	ratio, ok := ChokedMassFlowRatio(base, scaled)
	if !ok {
		return false, "both cases must be choked to compare mass flow scaling"
	}
	expected := ExpectedScaling(base.P0, scaled.P0)
	rel := math.Abs(ratio-expected) / expected
	if rel > 1e-9 {
		return false, fmt.Sprintf("mass flow scaled by %.9g, expected %.9g (rel err %.3g)", ratio, expected, rel)
	}
	if !MachProfileInvariant(base, scaled) {
		return false, "mach profile changed with total pressure"
	}
	return true, fmt.Sprintf("mass flow scales with p0 (factor %.9g) and Mach profile is invariant", ratio)
}
