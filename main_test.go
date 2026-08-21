package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTemp(t *testing.T, name, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	return path
}

func captureStdout(f func()) string {
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	os.Stdout = w
	f()
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestCLIDesignOutput(t *testing.T) {
	out := captureStdout(func() {
		if err := runDesign([]string{"example/gamma14.json"}); err != nil {
			t.Fatalf("runDesign failed: %v", err)
		}
	})
	for _, key := range []string{"supersonic", "subsonic", "choked: true", "qm"} {
		if !strings.Contains(out, key) {
			t.Errorf("output must contain %q, got:\n%s", key, out)
		}
	}
	m := parseValue(t, out, "M=", 1)
	if m <= 1 {
		t.Errorf("printed supersonic M=%v, must be well above 1", m)
	}
}

func TestCLIRejectsInvalidInput(t *testing.T) {
	badArea := writeTemp(t, "badarea.json", `{"t0":300,"p0":101325,"gamma":1.4,"r":287,"area_ratio":0.5,"throat_area":0.01}`)
	if err := runDesign([]string{badArea}); err == nil {
		t.Fatal("area ratio below 1 must be rejected")
	}
	badGamma := writeTemp(t, "badgamma.json", `{"t0":300,"p0":101325,"gamma":1.0,"r":287,"area_ratio":2,"throat_area":0.01}`)
	if err := runDesign([]string{badGamma}); err == nil {
		t.Fatal("gamma at one must be rejected")
	}
	badP0 := writeTemp(t, "badp0.json", `{"t0":300,"p0":0,"gamma":1.4,"r":287,"area_ratio":2,"throat_area":0.01}`)
	if err := runDesign([]string{badP0}); err == nil {
		t.Fatal("zero stagnation pressure must be rejected")
	}
}

func TestCLIMissingCaseFile(t *testing.T) {
	if err := runDesign([]string{"no-such-file.json"}); err == nil {
		t.Fatal("missing case file must be rejected")
	}
	if err := runDesign([]string{}); err == nil {
		t.Fatal("missing argument must be rejected")
	}
}

func parseValue(t *testing.T, out, key string, skip int) float64 {
	t.Helper()
	found := 0
	for _, line := range strings.Split(out, "\n") {
		if strings.Contains(line, key) {
			found++
			if found <= skip {
				continue
			}
			rest := line[strings.Index(line, key)+len(key):]
			var v float64
			if _, err := fmt.Sscanf(rest, "%f", &v); err != nil {
				t.Fatalf("parse %q in line %q: %v", key, line, err)
			}
			return v
		}
	}
	t.Fatalf("no line with %q beyond skip=%d in:\n%s", key, skip, out)
	return 0
}
