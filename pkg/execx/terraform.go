package execx

import (
    "bytes"
    "os/exec"
)

type PlanResult struct {
    Found   bool
    ExitCode int
    Output  string
    Summary string
}

func RunTerraformPlan(dir string) PlanResult {
    path, err := exec.LookPath("terraform")
    if err != nil { return PlanResult{Found: false} }
    cmd := exec.Command(path, "plan", "-detailed-exitcode")
    cmd.Dir = dir
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &out
    err = cmd.Run()
    ec := 0
    if err != nil {
        if exitErr, ok := err.(*exec.ExitError); ok {
            ec = exitErr.ExitCode()
        } else {
            ec = 1
        }
    }
    sum := "no changes"
    if ec == 2 { sum = "changes present" }
    if ec == 1 { sum = "error running plan" }
    return PlanResult{Found: true, ExitCode: ec, Output: out.String(), Summary: sum}
}