package cost

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEstimateProducesReport(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main\n\nfunc main() {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "main_test.go"), []byte("package main\n\nfunc TestX() {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	report, err := Estimate([]string{"--path", dir})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(report, "# Cleo Cost Estimate") {
		t.Fatalf("unexpected report:\n%s", report)
	}
	if !strings.Contains(report, "Rates Source: cached") {
		t.Fatalf("expected default rates source in report:\n%s", report)
	}
	if !strings.Contains(report, "| Metric") || !strings.Contains(report, "| Value |") {
		t.Fatalf("expected markdown table output:\n%s", report)
	}
}

func TestEstimateManualRate(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "a.py"), []byte("print('x')\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	report, err := Estimate([]string{"--path", dir, "--rates-source", "manual", "--hourly-rate", "150"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(report, "Rates Source: manual") {
		t.Fatalf("unexpected report:\n%s", report)
	}
}

func TestEstimatePlainFormat(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "a.go"), []byte("package a\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	report, err := Estimate([]string{"--path", dir, "--format", "plain"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(report, "Cleo Cost Estimate") {
		t.Fatalf("unexpected report:\n%s", report)
	}
	if strings.Contains(report, "# Cleo Cost Estimate") {
		t.Fatalf("expected non-markdown plain report:\n%s", report)
	}
}

func TestEstimateJSONFormat(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "a.ts"), []byte("export const a = 1\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	report, err := Estimate([]string{"--path", dir, "--format", "json"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(report, "\"title\": \"Cleo Cost Estimate\"") {
		t.Fatalf("unexpected json report:\n%s", report)
	}
}
