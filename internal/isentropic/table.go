package isentropic

import (
	"fmt"
)

type TableRow struct {
	Mach             float64
	AreaRatio        float64
	PressureRatio    float64
	TemperatureRatio float64
}

func RowAt(mach, gamma float64) (TableRow, error) {
	if err := CheckMach(mach); err != nil {
		return TableRow{}, err
	}
	if gamma <= 1 {
		return TableRow{}, fmt.Errorf("%w: got %.6g", ErrGammaBelowOne, gamma)
	}
	return TableRow{
		Mach:             mach,
		AreaRatio:        AreaRatio(mach, gamma),
		PressureRatio:    PressureRatio(mach, gamma),
		TemperatureRatio: TemperatureRatio(mach, gamma),
	}, nil
}

func AreaTable(gamma float64, machs []float64) ([]TableRow, error) {
	rows := make([]TableRow, 0, len(machs))
	for _, m := range machs {
		row, err := RowAt(m, gamma)
		if err != nil {
			return nil, err
		}
		rows = append(rows, row)
	}
	return rows, nil
}

func MachLine(gamma, machMin, machMax float64, n int) ([]TableRow, error) {
	if n < 2 {
		return nil, fmt.Errorf("table needs at least 2 points, got %d", n)
	}
	if machMax <= machMin {
		return nil, fmt.Errorf("table range empty: [%.6g, %.6g]", machMin, machMax)
	}
	rows := make([]TableRow, 0, n)
	for i := 0; i < n; i++ {
		t := float64(i) / float64(n-1)
		m := machMin + t*(machMax-machMin)
		if m == 0 {
			m = 1e-9
		}
		row, err := RowAt(m, gamma)
		if err != nil {
			return nil, err
		}
		rows = append(rows, row)
	}
	return rows, nil
}

func TableSatisfiesSonicMinimum(gamma float64, rows []TableRow) error {
	for _, row := range rows {
		if row.Mach < 1 {
			continue
		}
		if row.AreaRatio < 1-1e-12 {
			return fmt.Errorf("area ratio %.6g below sonic minimum at M=%.6g", row.AreaRatio, row.Mach)
		}
	}
	return nil
}

func TableMaxArea(rows []TableRow) float64 {
	max := 0.0
	for _, row := range rows {
		if row.AreaRatio > max {
			max = row.AreaRatio
		}
	}
	return max
}
