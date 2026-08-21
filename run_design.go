package main

import (
	"fmt"

	"nozzle-isen/internal/nozzle"
)

func runDesign(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: nozzle-isen design <input.json>")
	}
	c, err := nozzle.LoadCase(args[0])
	if err != nil {
		return fmt.Errorf("case %s: %w", args[0], err)
	}
	res, err := nozzle.Design(c)
	if err != nil {
		return fmt.Errorf("case %s: %w", args[0], err)
	}
	fmt.Print(nozzle.Report(res))
	return nil
}
