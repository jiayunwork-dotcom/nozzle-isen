package nozzle

import (
	"fmt"
	"strings"
)

func CompareBranches(r Result) string {
	var b strings.Builder
	if len(r.Exits) < 2 {
		b.WriteString("compare: single branch selected, nothing to contrast\n")
		return b.String()
	}
	sub, okSub := r.ExitByName("subsonic")
	super, okSuper := r.ExitByName("supersonic")
	if !okSub || !okSuper {
		b.WriteString("compare: both branches required\n")
		return b.String()
	}
	fmt.Fprintf(&b, "branch comparison at Ae/A*=%.6g:\n", r.AreaRatio)
	fmt.Fprintf(&b, "  subsonic   M=%.6g  p_e/p0=%.6g  T_e/T0=%.6g\n",
		sub.Mach, sub.PressureRatio, sub.TemperatureRatio)
	fmt.Fprintf(&b, "  supersonic M=%.6g  p_e/p0=%.6g  T_e/T0=%.6g\n",
		super.Mach, super.PressureRatio, super.TemperatureRatio)
	ratioMach := super.Mach / sub.Mach
	ratioPressure := sub.PressureRatio / super.PressureRatio
	ratioTemperature := sub.TemperatureRatio / super.TemperatureRatio
	fmt.Fprintf(&b, "  M_sup/M_sub = %.6g, p_sub/p_sup = %.6g, T_sub/T_sup = %.6g\n",
		ratioMach, ratioPressure, ratioTemperature)
	if sub.AreaRatio > 0 && sub.AreaRatio == super.AreaRatio {
		fmt.Fprintf(&b, "  both branches satisfy the same area ratio\n")
	}
	return b.String()
}
