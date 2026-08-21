package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}
	switch os.Args[1] {
	case "design":
		if err := runDesign(os.Args[2:]); err != nil {
			fail("design: %v", err)
		}
	case "help", "-h", "--help":
		fmt.Print(usage)
	default:
		fmt.Fprintf(os.Stderr, "nozzle-isen: unknown command %q\n\n%s", os.Args[1], usage)
		os.Exit(2)
	}
}

func fail(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "nozzle-isen: %s\n", fmt.Sprintf(format, args...))
	os.Exit(1)
}
