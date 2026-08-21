package nozzle

import (
	"math"
	"testing"
)

func baseCase() Case {
	return Case{
		T0:         300.0,
		P0:         101325.0,
		Gamma:      1.4,
		R:          287.0,
		AreaRatio:  2.0,
		ThroatArea: 0.01,
	}
}

func TestDesignGamma14Example(t *testing.T) {
	res, err := Design(baseCase())
	if err != nil {
		t.Fatalf("design failed: %v", err)
	}
	if !res.Choked {
		t.Error("throat must choke under design operation")
	}
	if res.Throat.Mach != 1 {
		t.Errorf("throat Mach = %v, want 1", res.Throat.Mach)
	}
	super, ok := res.ExitByName("supersonic")
	if !ok {
		t.Fatal("supersonic branch missing")
	}
	if super.Mach <= 1 {
		t.Errorf("supersonic M = %.6g, must be well above 1", super.Mach)
	}
	if math.Abs(res.MassFlow-2.364478) > 1e-3 {
		t.Errorf("qm = %.6g kg/s, want about 2.364478", res.MassFlow)
	}
	if math.Abs(res.Throat.PressureRatio-0.5282818) > 1e-6 {
		t.Errorf("p*/p0 = %.9g, want about 0.5282818", res.Throat.PressureRatio)
	}
}

func TestChokedMassFlowScalesWithPressure(t *testing.T) {
	base := baseCase()
	baseRes, err := Design(base)
	if err != nil {
		t.Fatalf("base design failed: %v", err)
	}
	scaled := ScaleTotalPressure(base, 2.0)
	scaledRes, err := Design(scaled)
	if err != nil {
		t.Fatalf("scaled design failed: %v", err)
	}
	ok, msg := ScalingConsistent(baseRes, scaledRes)
	if !ok {
		t.Errorf("scaling consistency: %s", msg)
	}
	if math.Abs(scaledRes.MassFlow-2*baseRes.MassFlow) > 1e-9*(1+baseRes.MassFlow) {
		t.Errorf("doubling p0 must double qm: base=%.9g scaled=%.9g",
			baseRes.MassFlow, scaledRes.MassFlow)
	}
}

func TestMachProfileInvariantWithPressure(t *testing.T) {
	base := baseCase()
	baseRes, err := Design(base)
	if err != nil {
		t.Fatalf("base design failed: %v", err)
	}
	for _, factor := range []float64{0.5, 1.5, 3.0} {
		scaled := ScaleTotalPressure(base, factor)
		scaledRes, err := Design(scaled)
		if err != nil {
			t.Fatalf("scaled design factor=%v failed: %v", factor, err)
		}
		if !MachProfileInvariant(baseRes, scaledRes) {
			t.Errorf("factor=%.2f: Mach profile must not depend on p0", factor)
		}
	}
}

func TestUnchokedHighBackPressureMassFlow(t *testing.T) {
	c := baseCase()
	c.BackPressure = 0.6 * c.P0
	res, err := Design(c)
	if err != nil {
		t.Fatalf("design failed: %v", err)
	}
	if res.Choked {
		t.Fatal("back pressure above sonic ratio must un-choke the throat")
	}
	sub, ok := res.ExitByName("subsonic-to-pb")
	if !ok {
		t.Fatal("subsonic-to-pb exit missing")
	}
	if sub.Mach >= 1 {
		t.Errorf("unchoked exit M = %.6g, must stay subsonic", sub.Mach)
	}
	if math.Abs(sub.PressureRatio-0.6) > 1e-10 {
		t.Errorf("exit p_e/p0 = %.10g, want back pressure ratio 0.6", sub.PressureRatio)
	}
	chokedFlux := ChokedMassFlow(c.State(), c.ThroatArea)
	if math.Abs(res.MassFlow-chokedFlux) < 1e-6 {
		t.Errorf("qm=%.9g must NOT equal the choked formula value %.9g",
			res.MassFlow, chokedFlux)
	}
	if res.MassFlow <= 0 {
		t.Errorf("qm = %.9g, must be positive", res.MassFlow)
	}
}

func TestDesignBranchSelection(t *testing.T) {
	sub := baseCase()
	sub.Branch = "subsonic"
	subRes, err := Design(sub)
	if err != nil {
		t.Fatalf("subsonic design failed: %v", err)
	}
	if len(subRes.Exits) != 1 {
		t.Errorf("subsonic selection must yield exactly one exit, got %d", len(subRes.Exits))
	}
	super := baseCase()
	super.Branch = "supersonic"
	superRes, err := Design(super)
	if err != nil {
		t.Fatalf("supersonic design failed: %v", err)
	}
	if len(superRes.Exits) != 1 {
		t.Errorf("supersonic selection must yield exactly one exit, got %d", len(superRes.Exits))
	}
}

func TestBackPressureOverExpandedUnderExpanded(t *testing.T) {
	low := baseCase()
	low.BackPressure = 0.3 * low.P0
	lowRes, err := Design(low)
	if err != nil {
		t.Fatalf("design low pb failed: %v", err)
	}
	super, _ := lowRes.ExitByName("supersonic")
	if super.Expansion != ExpansionOver {
		t.Errorf("pb/p0 below design exit ratio must be overexpanded, got %v", super.Expansion)
	}
	if lowRes.ShockNote != "normal shock not modelled" {
		t.Errorf("shock note = %q, want explicit not-modelled note", lowRes.ShockNote)
	}
	high := baseCase()
	high.BackPressure = 0.05 * high.P0
	highRes, err := Design(high)
	if err != nil {
		t.Fatalf("design high pb failed: %v", err)
	}
	super, _ = highRes.ExitByName("supersonic")
	if super.Expansion != ExpansionUnder {
		t.Errorf("pb/p0 above design exit ratio must be underexpanded, got %v", super.Expansion)
	}
}

func TestExitPressureConsistentWithMach(t *testing.T) {
	res, err := Design(baseCase())
	if err != nil {
		t.Fatalf("design failed: %v", err)
	}
	for _, e := range res.Exits {
		tr := e.TemperatureRatio
		wantTr := math.Pow(1+0.2*e.Mach*e.Mach, -1)
		if math.Abs(tr-wantTr) > 1e-10 {
			t.Errorf("%s: T_e/T0=%.12g inconsistent with M=%.9g formula %.12g",
				e.BranchName, tr, e.Mach, wantTr)
		}
		if math.Abs(e.Pressure- res.P0*e.PressureRatio) > 1e-9*(1+e.Pressure) {
			t.Errorf("%s: p_e=%.9g inconsistent with p0*p_e/p0=%.9g",
				e.BranchName, e.Pressure, res.P0*e.PressureRatio)
		}
	}
}

func TestDesignSelfCheckPasses(t *testing.T) {
	res, err := Design(baseCase())
	if err != nil {
		t.Fatalf("design failed: %v", err)
	}
	if !res.Checks.Passed() {
		t.Errorf("self-check failed:\n%s", res.Checks.Format())
	}
}

func TestOperatingLineChokedFraction(t *testing.T) {
	c := baseCase()
	points, err := OperatingLineAuto(c, 8)
	if err != nil {
		t.Fatalf("operating line failed: %v", err)
	}
	if len(points) != 8 {
		t.Fatalf("expected 8 points, got %d", len(points))
	}
	if !points[0].Choked {
		t.Error("lowest back pressure point must be choked")
	}
	if points[len(points)-1].Choked {
		t.Error("highest back pressure point must be unchoked")
	}
	frac := ChokedFraction(points)
	if frac <= 0 || frac >= 1 {
		t.Errorf("choked fraction = %.2f, want strictly between 0 and 1", frac)
	}
}

func TestValidateCaseRejectsBadInput(t *testing.T) {
	bad := []Case{
		{300, 101325, 1.4, 287, 0.5, 0.01, "", 0},
		{300, 101325, 1.4, 287, 2.0, 0, "", 0},
		{300, 101325, 1.4, 287, 2.0, 0.01, "", 150000},
	}
	for i, c := range bad {
		if err := ValidateCase(c); err == nil {
			t.Errorf("case %d must be rejected", i)
		}
	}
}
