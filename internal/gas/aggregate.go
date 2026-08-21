package gas

import (
	"errors"
	"fmt"
	"strings"
)

type ValidationReport struct {
	T0Error      error
	P0Error      error
	GammaError   error
	RError       error
	AreaError    error
	ThroatError  error
	PressureError error
}

func (r ValidationReport) Valid() bool {
	return r.T0Error == nil && r.P0Error == nil && r.GammaError == nil &&
		r.RError == nil && r.AreaError == nil && r.ThroatError == nil &&
		r.PressureError == nil
}

func (r ValidationReport) Error() string {
	var parts []string
	add := func(e error) {
		if e != nil {
			parts = append(parts, e.Error())
		}
	}
	add(r.T0Error)
	add(r.P0Error)
	add(r.GammaError)
	add(r.RError)
	add(r.AreaError)
	add(r.ThroatError)
	add(r.PressureError)
	if len(parts) == 0 {
		return ""
	}
	return fmt.Sprintf("%d validation error(s): %s", len(parts), strings.Join(parts, "; "))
}

func (r ValidationReport) Unwrap() error {
	if r.T0Error != nil {
		return r.T0Error
	}
	if r.P0Error != nil {
		return r.P0Error
	}
	if r.GammaError != nil {
		return r.GammaError
	}
	if r.RError != nil {
		return r.RError
	}
	if r.AreaError != nil {
		return r.AreaError
	}
	if r.ThroatError != nil {
		return r.ThroatError
	}
	return r.PressureError
}

var ErrValidationEmpty = errors.New("no validation error present")
