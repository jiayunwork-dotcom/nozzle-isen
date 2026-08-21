package nozzle

import (
	"fmt"
)

func ClassifyExpansion(designPressureRatio, backPressureRatio float64) ExpansionClass {
	switch {
	case backPressureRatio > designPressureRatio:
		return ExpansionOver
	case backPressureRatio < designPressureRatio:
		return ExpansionUnder
	default:
		return ExpansionMatched
	}
}

func (e ExpansionClass) Describe(exitPressureRatio, backPressureRatio float64) string {
	switch e {
	case ExpansionOver:
		return fmt.Sprintf("overexpanded (p_e/p0=%.6g < pb/p0=%.6g)",
			exitPressureRatio, backPressureRatio)
	case ExpansionUnder:
		return fmt.Sprintf("underexpanded (p_e/p0=%.6g > pb/p0=%.6g)",
			exitPressureRatio, backPressureRatio)
	default:
		return fmt.Sprintf("matched (p_e/p0=%.6g = pb/p0=%.6g)",
			exitPressureRatio, backPressureRatio)
	}
}

func ExpansionSummary(r Result) []string {
	if !r.HasBackPressure || !r.Choked {
		return nil
	}
	pbPR := r.BackPressure / r.P0
	lines := make([]string, 0, len(r.Exits))
	for _, e := range r.Exits {
		lines = append(lines, fmt.Sprintf("%s: %s", e.BranchName,
			e.Expansion.Describe(e.PressureRatio, pbPR)))
	}
	if r.ShockNote != "" {
		lines = append(lines, r.ShockNote)
	}
	return lines
}
