package collect

import (
    "encoding/json"
    "os"
)

type tfstate struct {
    Resources []struct {
        Type string `json:"type"`
        Name string `json:"name"`
        Instances []struct {
            Attributes map[string]any `json:"attributes"`
        } `json:"instances"`
    } `json:"resources"`
}

func LoadState(path string) (ResourceSet, error) {
    if path == "" {
        return ResourceSet{}, nil
    }
    b, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }
    var s tfstate
    if err := json.Unmarshal(b, &s); err != nil {
        return nil, err
    }
    out := ResourceSet{}
    for _, r := range s.Resources {
        attrs := map[string]string{}
        if len(r.Instances) > 0 {
            for k, v := range r.Instances[0].Attributes {
                switch vv := v.(type) {
                case string:
                    attrs[k] = vv
                default:
                    attrs[k] = ""
                }
            }
        }
        out[key(r.Type, r.Name)] = Resource{Type: r.Type, Name: r.Name, Attrs: attrs, Origin: "state"}
    }
    return out, nil
}