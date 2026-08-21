package nozzle

import (
	"fmt"
	"math"
)

func ratio(v float64) string {
	return fmt.Sprintf("%.6g", v)
}

func value(v float64) string {
	return fmt.Sprintf("%.4g", v)
}

func pressureWithUnit(v float64) string {
	return fmt.Sprintf("%.4g Pa", v)
}

func pressureKpa(v float64) string {
	return fmt.Sprintf("%.4g kPa", v/1000)
}

func massFlowWithUnit(v float64) string {
	return fmt.Sprintf("%.4g kg/s", v)
}

func areaWithUnit(v float64) string {
	return fmt.Sprintf("%.4g m2", v)
}

func temperatureWithUnit(v float64) string {
	return fmt.Sprintf("%.4g K", v)
}

func significant(v float64) string {
	if math.Abs(v) < 1e-6 || math.Abs(v) > 1e6 {
		return fmt.Sprintf("%.3e", v)
	}
	return fmt.Sprintf("%.6g", v)
}
