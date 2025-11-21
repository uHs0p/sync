package compare

import (
    "sort"

    "terraformsync/internal/collect"
)

type Triplet struct {
    Key     string                        `json:"key" yaml:"key"`
    Code    *collect.Resource             `json:"code" yaml:"code"`
    State   *collect.Resource             `json:"state" yaml:"state"`
    Deployed *collect.Resource            `json:"deployed" yaml:"deployed"`
    Drift   string                        `json:"drift" yaml:"drift"`
    Severity string                       `json:"severity" yaml:"severity"`
}

type Result struct {
    Items []Triplet `json:"items" yaml:"items"`
    Summary map[string]int `json:"summary" yaml:"summary"`
    DurationMs int `json:"duration_ms" yaml:"duration_ms"`
}

func ThreeWay(code, state, deployed collect.ResourceSet) Result {
    keys := map[string]struct{}{}
    for k := range code { keys[k] = struct{}{} }
    for k := range state { keys[k] = struct{}{} }
    for k := range deployed { keys[k] = struct{}{} }
    all := make([]string, 0, len(keys))
    for k := range keys { all = append(all, k) }
    sort.Strings(all)

    summary := map[string]int{}
    items := []Triplet{}
    for _, k := range all {
        var c, s, d *collect.Resource
        if v, ok := code[k]; ok { vv := v; c = &vv }
        if v, ok := state[k]; ok { vv := v; s = &vv }
        if v, ok := deployed[k]; ok { vv := v; d = &vv }
        drift, severity := classify(c, s, d)
        items = append(items, Triplet{Key: k, Code: c, State: s, Deployed: d, Drift: drift, Severity: severity})
        summary[drift]++
    }
    return Result{Items: items, Summary: summary}
}

func classify(c, s, d *collect.Resource) (string, string) {
    if c != nil && s != nil && d != nil {
        if equalAttrs(c.Attrs, s.Attrs) && equalAttrs(s.Attrs, d.Attrs) {
            return "in_sync", "low"
        }
        return "three_way_conflict", "high"
    }
    if c != nil && s != nil && d == nil {
        if equalAttrs(c.Attrs, s.Attrs) { return "missing_deployed", "medium" }
        return "code_state_diff", "high"
    }
    if c != nil && s == nil && d != nil {
        if equalAttrs(c.Attrs, d.Attrs) { return "missing_state", "medium" }
        return "code_deployed_diff", "high"
    }
    if c == nil && s != nil && d != nil {
        if equalAttrs(s.Attrs, d.Attrs) { return "missing_code", "medium" }
        return "state_deployed_diff", "high"
    }
    if c != nil && s == nil && d == nil { return "only_in_code", "medium" }
    if c == nil && s != nil && d == nil { return "only_in_state", "medium" }
    if c == nil && s == nil && d != nil { return "only_in_deployed", "medium" }
    return "unknown", "low"
}

func equalAttrs(a, b map[string]string) bool {
    if len(a) != len(b) { return false }
    for k, v := range a { if b[k] != v { return false } }
    return true
}