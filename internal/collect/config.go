package collect

import (
    "bufio"
    "os"
    "path/filepath"
    "regexp"
    "strings"
)

type Resource struct {
    Type   string            `json:"type" yaml:"type"`
    Name   string            `json:"name" yaml:"name"`
    Attrs  map[string]string `json:"attrs" yaml:"attrs"`
    Origin string            `json:"origin" yaml:"origin"`
}

type ResourceSet map[string]Resource

func key(t, n string) string { return t + ":" + n }

var resHeader = regexp.MustCompile(`^\s*resource\s+"([^"]+)"\s+"([^"]+)"\s*\{`)
var attrLine = regexp.MustCompile(`^\s*([A-Za-z0-9_]+)\s*=\s*"?([^"\n]+)"?`)

func LoadConfig(dir string) (ResourceSet, error) {
    if dir == "" {
        return ResourceSet{}, nil
    }
    abs, _ := filepath.Abs(dir)
    _, err := os.Stat(abs)
    if err != nil { return nil, err }
    out := ResourceSet{}
    entries, err := os.ReadDir(abs)
    if err != nil { return nil, err }
    for _, e := range entries {
        if e.IsDir() { continue }
        if !strings.HasSuffix(e.Name(), ".tf") { continue }
        p := filepath.Join(abs, e.Name())
        f, err := os.Open(p)
        if err != nil { return nil, err }
        scanner := bufio.NewScanner(f)
        var curType, curName string
        attrs := map[string]string{}
        inRes := false
        for scanner.Scan() {
            line := scanner.Text()
            if !inRes {
                if m := resHeader.FindStringSubmatch(line); len(m) == 3 {
                    curType, curName = m[1], m[2]
                    attrs = map[string]string{}
                    inRes = true
                }
                continue
            }
            if strings.Contains(line, "}") {
                if curType != "" && curName != "" {
                    out[key(curType, curName)] = Resource{Type: curType, Name: curName, Attrs: attrs, Origin: "code"}
                }
                inRes = false
                curType, curName = "", ""
                attrs = map[string]string{}
                continue
            }
            if m := attrLine.FindStringSubmatch(line); len(m) == 3 {
                attrs[m[1]] = strings.TrimSpace(m[2])
            }
        }
        _ = f.Close()
    }
    return out, nil
}