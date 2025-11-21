package compare

import (
    "testing"
    "terraformsync/internal/collect"
)

func TestThreeWayClassification(t *testing.T) {
    code := collect.ResourceSet{"t:n": {Type: "t", Name: "n", Attrs: map[string]string{"a":"1"}}}
    state := collect.ResourceSet{"t:n": {Type: "t", Name: "n", Attrs: map[string]string{"a":"1"}}}
    deployed := collect.ResourceSet{"t:n": {Type: "t", Name: "n", Attrs: map[string]string{"a":"1"}}}
    r := ThreeWay(code, state, deployed)
    if len(r.Items) != 1 || r.Items[0].Drift != "in_sync" {
        t.Fatalf("expected in_sync")
    }
    deployed["t:n"] = collect.Resource{Type: "t", Name: "n", Attrs: map[string]string{"a":"2"}}
    r = ThreeWay(code, state, deployed)
    if r.Items[0].Drift != "three_way_conflict" {
        t.Fatalf("expected conflict")
    }
}