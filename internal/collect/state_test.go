package collect

import (
    "os"
    "path/filepath"
    "testing"
)

func TestLoadState(t *testing.T) {
    p := filepath.Join("../../examples/simple/terraform.tfstate")
    rs, err := LoadState(p)
    if err != nil { t.Fatalf("err %v", err) }
    if _, ok := rs["aws_s3_bucket:logs"]; !ok { t.Fatalf("missing resource") }
}

func TestLoadDeployedFile(t *testing.T) {
    p := filepath.Join("../../examples/simple/deployed.json")
    rs, err := LoadDeployedFile(p)
    if err != nil { t.Fatalf("err %v", err) }
    if len(rs) < 2 { t.Fatalf("expected >=2") }
}

func TestLoadConfigDirInvalid(t *testing.T) {
    _, err := LoadConfig("/nonexistent")
    if err == nil { t.Fatalf("expected error") }
}

func TestLoadConfigDir(t *testing.T) {
    d := filepath.Join("../../examples/simple")
    rs, err := LoadConfig(d)
    if err != nil { t.Fatalf("err %v", err) }
    _ = os.Chmod(d, 0o755)
    if len(rs) == 0 { t.Fatalf("expected resources") }
}