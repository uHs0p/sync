package collect

import (
    "os"
    "path/filepath"
    "testing"
)

func TestLoadConfigParsesAttrs(t *testing.T) {
    dir := t.TempDir()
    tf := `resource "aws_s3_bucket" "logs" {
  bucket = "example-logs"
  acl = private
}
`
    p := filepath.Join(dir, "main.tf")
    _ = os.WriteFile(p, []byte(tf), 0o644)
    rs, err := LoadConfig(dir)
    if err != nil { t.Fatalf("err %v", err) }
    r, ok := rs["aws_s3_bucket:logs"]
    if !ok { t.Fatalf("missing resource") }
    if r.Attrs["bucket"] != "example-logs" { t.Fatalf("bucket attr not parsed") }
}