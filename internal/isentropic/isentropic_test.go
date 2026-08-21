package isentropic

import (
	"errors"
	"math"
	"testing"
)

func TestAreaRatioUnityAtMachOne(t *testing.T) {
	for _, g := range []float64{1.4, 1.2, 1.67, 1.1} {
		got := AreaRatio(1, g)
		if math.Abs(got-1) > 1e-12 {
			t.Errorf("gamma=%.2f: A/A* at M=1 = %.15g, want 1", g, got)
		}
	}
}

func TestMachFromAreaGamma14RatioTwo(t *testing.T) {
	sub, err := MachFromArea(2, 1.4, BranchSubsonic, DefaultOptions())
	if err != nil {
		t.Fatalf("subsonic solve failed: %v", err)
	}
	if math.Abs(sub-0.30590) > 1e-4 {
		t.Errorf("subsonic M = %.6g, want about 0.30590", sub)
	}
	super, err := MachFromArea(2, 1.4, BranchSupersonic, DefaultOptions())
	if err != nil {
		t.Fatalf("supersonic solve failed: %v", err)
	}
	if super < 2.0 {
		t.Errorf("supersonic M = %.6g, must be well above 1 for Ae/A*=2", super)
	}
	if math.Abs(super-2.1972) > 1e-3 {
		t.Errorf("supersonic M = %.6g, want about 2.1972", super)
	}
}

func TestMachFromAreaRecoversRatio(t *testing.T) {
	for _, branch := range []Branch{BranchSubsonic, BranchSupersonic} {
		mach, back, rel, err := VerifyInverse(3.5, 1.4, branch, DefaultOptions())
		if err != nil {
			t.Fatalf("branch %d inverse failed: %v", branch, err)
		}
		if rel > 1e-10 {
			t.Errorf("branch %d: area rel error = %.3g, want <= 1e-10 (M=%.9g back=%.9g)",
				branch, rel, mach, back)
		}
	}
}

func TestMachFromAreaRejectsBelowUnity(t *testing.T) {
	_, err := MachFromArea(0.5, 1.4, BranchSubsonic, DefaultOptions())
	if !errors.Is(err, ErrTargetBelowOne) {
		t.Errorf("want ErrTargetBelowOne, got %v", err)
	}
}

func TestMachFromAreaRejectsBadGamma(t *testing.T) {
	for _, g := range []float64{1.0, 0.9, -2} {
		_, err := MachFromArea(2, g, BranchSubsonic, DefaultOptions())
		if !errors.Is(err, ErrGammaBelowOne) {
			t.Errorf("gamma=%v: want ErrGammaBelowOne, got %v", g, err)
		}
	}
}

func TestNewtonConvergenceLimitExceeded(t *testing.T) {
	opts := Options{MaxIter: 2, Tol: 1e-10}
	_, err := MachFromArea(2, 1.4, BranchSupersonic, opts)
	if !errors.Is(err, ErrNewtonNonConverged) {
		t.Errorf("want ErrNewtonNonConverged when iteration limit is exhausted, got %v", err)
	}
}

func TestMachFromPressureRatioRoundTrip(t *testing.T) {
	for _, pr := range []float64{0.95, 0.8, 0.6, 0.3, 0.1} {
		mach, back, rel, err := VerifyPressureRoundTrip(pr, 1.4)
		if err != nil {
			t.Fatalf("pr=%.3g: %v", pr, err)
		}
		if rel > 1e-10 {
			t.Errorf("pr=%.3g: pressure round trip rel error %.3g (M=%.9g back=%.9g)",
				pr, rel, mach, back)
		}
	}
}

func TestSonicConsistency(t *testing.T) {
	if err := SonicConsistency(1.4); err != nil {
		t.Errorf("gamma=1.4: %v", err)
	}
	if err := SonicConsistency(1.67); err != nil {
		t.Errorf("gamma=1.67: %v", err)
	}
}

func TestBisectMatchesNewton(t *testing.T) {
	for _, branch := range []Branch{BranchSubsonic, BranchSupersonic} {
		want, err := MachFromArea(2, 1.4, branch, DefaultOptions())
		if err != nil {
			t.Fatalf("newton branch %d: %v", branch, err)
		}
		got, err := MachFromAreaBisect(2, 1.4, branch, DefaultOptions())
		if err != nil {
			t.Fatalf("bisect branch %d: %v", branch, err)
		}
		if math.Abs(got-want) > 1e-8 {
			t.Errorf("branch %d: bisect M=%.9g vs newton M=%.9g", branch, got, want)
		}
	}
}

func TestAreaRatioMonotonic(t *testing.T) {
	if err := AreaRatioMonotonic(1.4, BranchSubsonic); err != nil {
		t.Errorf("subsonic branch: %v", err)
	}
	if err := AreaRatioMonotonic(1.4, BranchSupersonic); err != nil {
		t.Errorf("supersonic branch: %v", err)
	}
}

func TestMachFromAreaNearSonicDegenerate(t *testing.T) {
	for _, branch := range []Branch{BranchSubsonic, BranchSupersonic} {
		mach, err := MachFromArea(1+1e-12, 1.4, branch, DefaultOptions())
		if err != nil {
			t.Fatalf("branch %d near-sonic: %v", branch, err)
		}
		if math.Abs(mach-1) > 1e-12 {
			t.Errorf("branch %d: near-sonic M=%.12g, want 1", branch, mach)
		}
	}
}
