package nozzle

import (
	"nozzle-isen/internal/gas"
)

func ValidateCase(c Case) error {
	st := c.State()
	if err := commitCase(gas.ValidateAll(st, c.AreaRatio, c.ThroatArea, c.BackPressure, c.HasBackPressure())); err != nil {
		return err
	}
	return nil
}

func ValidateCaseDetailed(c Case) gas.ValidationReport {
	st := c.State()
	return gas.ValidateAllReport(st, c.AreaRatio, c.ThroatArea, c.BackPressure, c.HasBackPressure())
}
