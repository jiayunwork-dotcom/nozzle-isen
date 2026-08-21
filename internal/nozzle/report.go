package nozzle

import (
	"bytes"
	"fmt"
)

func Report(r Result) string {
	var b bytes.Buffer
	b.WriteString("nozzle-isen design: isentropic converging-diverging nozzle audit\n")
	fmt.Fprintf(&b, "stagnation state: T0=%.3g K p0=%.4g Pa gamma=%.4g R=%.3g J/(kg.K)\n",
		r.T0, r.P0, r.Gamma, r.R)
	fmt.Fprintf(&b, "area ratio Ae/A* = %.6g, throat area A* = %.6g m2\n",
		r.AreaRatio, r.ThroatArea)
	fmt.Fprintf(&b, "selected expansion: %s\n", r.Case.BranchName())

	if r.Choked {
		b.WriteString("\nthroat (choked M*=1):\n")
		fmt.Fprintf(&b, "  p*/p0 = %s\n", ratio(r.Throat.PressureRatio))
		fmt.Fprintf(&b, "  T*/T0 = %s\n", ratio(r.Throat.TemperatureRatio))
	} else {
		b.WriteString("\nthroat: not choked, throat Mach stays below 1\n")
	}

	b.WriteString("\nexit branches:\n")
	for _, e := range r.Exits {
		fmt.Fprintf(&b, "  %-16s M=%-12s p_e/p0=%-12s p_e=%-12s T_e/T0=%s\n",
			e.BranchName+":", ratio(e.Mach), ratio(e.PressureRatio),
			value(e.Pressure), ratio(e.TemperatureRatio))
	}

	fmt.Fprintf(&b, "\nchoked: %v\n", r.Choked)
	fmt.Fprintf(&b, "  reason: %s\n", r.ChokedReason)
	fmt.Fprintf(&b, "mass flow qm = %s kg/s\n", ratio(r.MassFlow))
	fmt.Fprintf(&b, "  formula: %s\n", r.MassFlowFormula)

	if r.HasBackPressure {
		fmt.Fprintf(&b, "\nback pressure pb = %.4g Pa (pb/p0 = %s)\n",
			r.BackPressure, ratio(r.BackPressure/r.P0))
		if r.Choked {
			for _, e := range r.Exits {
				fmt.Fprintf(&b, "  %s exit: %s\n", e.BranchName,
					e.Expansion.Describe(e.PressureRatio, r.BackPressure/r.P0))
			}
			if r.ShockNote != "" {
				fmt.Fprintf(&b, "  note: %s\n", r.ShockNote)
			}
		} else {
			b.WriteString("  exit falls isentropically to pb, no design expansion ratio applies\n")
		}
	}
	b.WriteString("\n" + r.Checks.Format())
	return b.String()
}
