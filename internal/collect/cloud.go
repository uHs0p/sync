package collect

import (
    "encoding/json"
    "os"
)

func MockDeployed(code, state ResourceSet) ResourceSet {
    out := ResourceSet{}
    for k, v := range code {
        out[k] = Resource{Type: v.Type, Name: v.Name, Attrs: v.Attrs, Origin: "deployed"}
    }
    for k, v := range state {
    	if _, ok := out[k]; !ok {
            out[k] = Resource{Type: v.Type, Name: v.Name, Attrs: v.Attrs, Origin: "deployed"}
        }
    }
    return out
}

func LoadDeployedFile(path string) (ResourceSet, error) {
    b, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }
    var arr []Resource
    if err := json.Unmarshal(b, &arr); err != nil {
        return nil, err
    }
    out := ResourceSet{}
    for _, r := range arr {
        out[key(r.Type, r.Name)] = Resource{Type: r.Type, Name: r.Name, Attrs: r.Attrs, Origin: "deployed"}
    }
    return out, nil
}