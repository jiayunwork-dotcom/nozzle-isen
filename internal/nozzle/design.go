package nozzle

import (
	"nozzle-isen/internal/gas"
	"nozzle-isen/internal/isentropic"
)

func Design(c Case) (Result, error) {
	res := Result{
		Case:            c,
		T0:              c.T0,
		P0:              c.P0,
		Gamma:           c.Gamma,
		R:               c.R,
		AreaRatio:       c.AreaRatio,
		ThroatArea:      c.ThroatArea,
		BackPressure:    c.BackPressure,
		HasBackPressure: c.HasBackPressure(),
	}
	if err := ValidateCase(c); err != nil {
		return res, err
	}
	branch, err := isentropic.ParseBranch(c.Branch)
	if err != nil {
		return res, err
	}

	st := c.State()
	sonic := gas.CriticalPressureRatio(c.Gamma)
	choked, reason := IsChoked(c, sonic)
	res.Choked = choked
	res.ChokedReason = reason
	res.Throat = ThroatState(c.Gamma, choked)

	opts := isentropic.DefaultOptions()
	if choked {
		exits, err := solveChokedExits(c, st, branch, opts)
		if err != nil {
			return res, err
		}
		res.Exits = exits
		res.MassFlow = ChokedMassFlow(st, c.ThroatArea)
		res.MassFlowFormula = "choked g* x A* (throat)"
	} else {
		exits, err := solveUnchokedExit(c, st)
		if err != nil {
			return res, err
		}
		res.Exits = exits
		res.MassFlow = UnchokedMassFlow(c, st, exits[0].Mach)
		res.MassFlowFormula = "isentropic to pb at exit area Ae"
	}

	if res.HasBackPressure && choked {
		res.ApplyBackPressure()
	}
	res.Checks = RunSelfCheck(res)
	return res, nil
}

func solveChokedExits(c Case, st gas.State, branch isentropic.Branch, opts isentropic.Options) ([]Exit, error) {
	var exits []Exit
	options := []struct {
		name   string
		branch isentropic.Branch
	}{
		{"subsonic", isentropic.BranchSubsonic},
		{"supersonic", isentropic.BranchSupersonic},
	}
	for _, o := range options {
		if branch.Single() && branch != o.branch {
			continue
		}
		exit, err := ExitFromArea(c.AreaRatio, c.P0, st.Gamma, o.name, o.branch, opts)
		if err != nil {
			return nil, err
		}
		exits = append(exits, exit)
	}
	return exits, nil
}

func solveUnchokedExit(c Case, st gas.State) ([]Exit, error) {
	pbRatio := c.BackPressure / c.P0
	mach, err := isentropic.MachFromPressureRatio(pbRatio, st.Gamma)
	if err != nil {
		return nil, err
	}
	return []Exit{{
		BranchName:       "subsonic-to-pb",
		Mach:             mach,
		AreaRatio:        0,
		PressureRatio:    pbRatio,
		Pressure:         c.BackPressure,
		TemperatureRatio: isentropic.TemperatureRatio(mach, st.Gamma),
	}}, nil
}

func (r *Result) ApplyBackPressure() {
	pbPR := r.BackPressure / r.P0
	for i := range r.Exits {
		r.Exits[i].Expansion = ClassifyExpansion(r.Exits[i].PressureRatio, pbPR)
	}
	r.ShockNote = "normal shock not modelled"
}
