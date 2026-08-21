package gas

import (
	"math"
	"testing"
)

func TestValidateRejectsNonPositive(t *testing.T) {
	cases := []struct {
		name string
		t0   float64
		p0   float64
		g    float64
		r    float64
	}{
		{"zero total temperature", 0, 101325, 1.4, 287},
		{"zero total pressure", 300, 0, 1.4, 287},
		{"gamma at one", 300, 101325, 1.0, 287},
		{"negative gas constant", 300, 101325, 1.4, -1},
	}
	for _, c := range cases {
		s := NewState(c.t0, c.p0, c.g, c.r)
		if err := ValidateState(s); err == nil {
			t.Errorf("%s: want error, got nil", c.name)
		}
	}
}

func TestGammaNearOneRejected(t *testing.T) {
	for _, g := range []float64{1.0, 1.0002, 1.0005, 0.99} {
		if err := ValidateGamma(g); err == nil {
			t.Errorf("gamma=%.4g must be rejected, got nil error", g)
		}
	}
	if err := ValidateGamma(1.01); err != nil {
		t.Errorf("gamma=1.01 must pass, got %v", err)
	}
}

func TestCriticalRatiosAtGamma14(t *testing.T) {
	wantPR := 0.5282817877171742
	gotPR := CriticalPressureRatio(1.4)
	if math.Abs(gotPR-wantPR) > 1e-12 {
		t.Errorf("p*/p0 = %.15g, want %.15g", gotPR, wantPR)
	}
	wantTR := 2 / 2.4
	gotTR := CriticalTemperatureRatio(1.4)
	if math.Abs(gotTR-wantTR) > 1e-15 {
		t.Errorf("T*/T0 = %.15g, want %.15g", gotTR, wantTR)
	}
}

func TestChokedFluxMatchesIsentropicAtSonic(t *testing.T) {
	st := NewState(300, 101325, 1.4, 287)
	rel := FluxConsistency(st)
	if math.Abs(rel) > 1e-12 {
		t.Errorf("flux consistency rel error = %.3g, want ~0", rel)
	}
}

func TestStateDefaultsToAirConstant(t *testing.T) {
	s := NewState(300, 101325, 1.4, 0)
	if s.R != Rair {
		t.Errorf("R = %v, want %v when zero", s.R, Rair)
	}
	if s.Gamma != GammaAir {
		t.Errorf("gamma = %v, want %v", s.Gamma, GammaAir)
	}
}

func TestValidationReportAggregates(t *testing.T) {
	report := ValidateAllReport(NewState(0, 0, 1.0, -1), 0.5, 0, 200000, true)
	if report.Valid() {
		t.Fatal("all-invalid state must fail")
	}
	if report.T0Error == nil {
		t.Error("T0 error missing")
	}
	if report.P0Error == nil {
		t.Error("p0 error missing")
	}
	if report.GammaError == nil {
		t.Error("gamma error missing")
	}
	if report.RError == nil {
		t.Error("R error missing")
	}
	if report.AreaError == nil {
		t.Error("area ratio error missing")
	}
	if report.ThroatError == nil {
		t.Error("throat area error missing")
	}
	if report.PressureError == nil {
		t.Error("back pressure at or above p0 must error")
	}
}
