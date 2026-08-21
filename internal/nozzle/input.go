package nozzle

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"nozzle-isen/internal/gas"
)

type Case struct {
	T0           float64 `json:"t0"`
	P0           float64 `json:"p0"`
	Gamma        float64 `json:"gamma"`
	R            float64 `json:"r"`
	AreaRatio    float64 `json:"area_ratio"`
	ThroatArea   float64 `json:"throat_area"`
	Branch       string  `json:"branch"`
	BackPressure float64 `json:"back_pressure"`
}

func LoadCase(path string) (Case, error) {
	var c Case
	data, err := os.ReadFile(path)
	if err != nil {
		return c, fmt.Errorf("read case file: %w", err)
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&c); err != nil {
		return c, fmt.Errorf("parse case json: %w", err)
	}
	if c.R == 0 {
		c.R = gas.Rair
	}
	return c, nil
}

func (c Case) State() gas.State {
	return gas.NewState(c.T0, c.P0, c.Gamma, c.R)
}

func (c Case) HasBackPressure() bool {
	return c.BackPressure > 0
}

func (c Case) BranchName() string {
	if c.Branch == "" {
		return "both"
	}
	return c.Branch
}

func (c Case) Describe() string {
	return fmt.Sprintf("T0=%.3g p0=%.4g gamma=%.4g R=%.3g Ae/A*=%.6g A*=%.6g branch=%s",
		c.T0, c.P0, c.Gamma, c.R, c.AreaRatio, c.ThroatArea, c.BranchName())
}

var ErrThroatAreaNonPositive = errors.New("throat area A* must be > 0")
