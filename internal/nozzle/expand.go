package nozzle

import (
	"fmt"
)

type BackPressurePoint struct {
	Pb       float64
	PbRatio  float64
	Choked   bool
	MassFlow float64
	ExitMach float64
	ExitName string
}

func (p BackPressurePoint) String() string {
	return fmt.Sprintf("pb=%.4g (%.6g p0) choked=%v qm=%.6g M=%.6g",
		p.Pb, p.PbRatio, p.Choked, p.MassFlow, p.ExitMach)
}

func OperatingLine(c Case, pressures []float64) ([]BackPressurePoint, error) {
	points := make([]BackPressurePoint, 0, len(pressures))
	for _, pb := range pressures {
		pt, err := pointAt(c, pb)
		if err != nil {
			return nil, err
		}
		points = append(points, pt)
	}
	return points, nil
}

func OperatingLineAuto(c Case, n int) ([]BackPressurePoint, error) {
	if n < 2 {
		return nil, fmt.Errorf("operating line needs at least 2 points, got %d", n)
	}
	sonic := SonicPressureOf(c)
	upper := c.P0 * (1 - 1e-4)
	points := make([]BackPressurePoint, 0, n)
	for i := 0; i < n; i++ {
		t := float64(i) / float64(n-1)
		pb := sonic + t*(upper-sonic)
		pt, err := pointAt(c, pb)
		if err != nil {
			return nil, err
		}
		points = append(points, pt)
	}
	return points, nil
}

func pointAt(c Case, pb float64) (BackPressurePoint, error) {
	run := c
	run.BackPressure = pb
	res, err := Design(run)
	if err != nil {
		return BackPressurePoint{}, err
	}
	pt := BackPressurePoint{
		Pb:       pb,
		PbRatio:  pb / c.P0,
		Choked:   res.Choked,
		MassFlow: res.MassFlow,
	}
	if len(res.Exits) > 0 {
		pt.ExitMach = res.Exits[0].Mach
		pt.ExitName = res.Exits[0].BranchName
	}
	return pt, nil
}

func ChokedFraction(points []BackPressurePoint) float64 {
	if len(points) == 0 {
		return 0
	}
	count := 0
	for _, p := range points {
		if p.Choked {
			count++
		}
	}
	return float64(count) / float64(len(points))
}

func MaxMassFlow(points []BackPressurePoint) (float64, float64) {
	max := 0.0
	pbAtMax := 0.0
	for _, p := range points {
		if p.MassFlow > max {
			max = p.MassFlow
			pbAtMax = p.Pb
		}
	}
	return max, pbAtMax
}
