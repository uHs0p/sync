package cliutil

import (
    "os"
    "path/filepath"
    "testing"
)

func TestValidateInputs(t *testing.T) {
    d := filepath.Join("../../examples/simple")
    s := filepath.Join(d, "terraform.tfstate")
    if err := ValidateInputs(d, s, "file", filepath.Join(d, "deployed.json")); err != nil {
        t.Fatalf("expected valid, got %v", err)
    }
    if err := ValidateInputs("/nope", s, "mock", ""); err == nil {
        t.Fatalf("expected invalid config dir")
    }
    tmp := t.TempDir()
    pf := filepath.Join(tmp, "x.txt")
    _ = os.WriteFile(pf, []byte("hi"), 0o644)
    if err := ValidateInputs(tmp, pf, "file", ""); err == nil {
        t.Fatalf("expected missing deployed-file error")
    }
    if err := ValidateInputs(tmp, pf, "bad", ""); err == nil {
        t.Fatalf("expected invalid deployed-source error")
    }
}