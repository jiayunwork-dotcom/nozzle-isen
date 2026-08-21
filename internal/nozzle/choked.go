package nozzle

import (
	"fmt"

	"nozzle-isen/internal/gas"
)

func IsChoked(c Case, sonicPressureRatio float64) (bool, string) {
	if !c.HasBackPressure() {
		return true, "design operation, throat reaches M*=1"
	}
	pbRatio := c.BackPressure / c.P0
	if pbRatio > sonicPressureRatio {
		return false, fmt.Sprintf(
			"back pressure pb/p0=%.6g exceeds sonic pressure ratio %.6g, throat cannot choke",
			pbRatio, sonicPressureRatio)
	}
	return true, "back pressure at or below sonic pressure ratio, throat chokes"
}

func BackPressureRegime(c Case, st gas.State) string {
	sonic := gas.CriticalPressureRatio(st.Gamma)
	pbRatio := c.BackPressure / st.P0
	switch {
	case !c.HasBackPressure():
		return "no back pressure specified"
	case pbRatio > sonic:
		return fmt.Sprintf("pb/p0=%.6g above sonic %.6g: unchoked subsonic flow", pbRatio, sonic)
	case pbRatio == sonic:
		return fmt.Sprintf("pb/p0=%.6g equals sonic %.6g: marginal choking", pbRatio, sonic)
	default:
		return fmt.Sprintf("pb/p0=%.6g below sonic %.6g: choked throat", pbRatio, sonic)
	}
}

func SonicPressureOf(c Case) float64 {
	return gas.CriticalPressureRatio(c.Gamma) * c.P0
}
