package main

import (
    "bufio"
    "encoding/json"
    "flag"
    "fmt"
    "os"
    "strings"
    "time"

    "terraformsync/internal/cliutil"
    "terraformsync/internal/collect"
    "terraformsync/internal/compare"
    "terraformsync/pkg/execx"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Fprintln(os.Stderr, "usage: terrainsync <plan|diff|resolve> [flags]")
        os.Exit(1)
    }
    cmd := os.Args[1]
    switch cmd {
    case "plan":
        fs := flag.NewFlagSet("plan", flag.ExitOnError)
        configDir := fs.String("config-dir", ".", "Path to Terraform configuration directory")
        _ = fs.Parse(os.Args[2:])
        if err := cliutil.ValidateInputs(*configDir, "", "mock", ""); err != nil {
            fmt.Fprintln(os.Stderr, err)
            os.Exit(1)
        }
        res := execx.RunTerraformPlan(*configDir)
        if !res.Found {
            fmt.Println("terraform not found; baseline plan unavailable")
            return
        }
        fmt.Println(res.Summary)
    case "diff":
        fs := flag.NewFlagSet("diff", flag.ExitOnError)
        configDir := fs.String("config-dir", ".", "Path to Terraform configuration directory")
        stateFile := fs.String("state-file", "", "Path to Terraform state file (.tfstate)")
        deployedSource := fs.String("deployed-source", "mock", "Deployed source: mock | file")
        deployedFile := fs.String("deployed-file", "", "Path to deployed data file when deployed-source=file")
        _ = fs.Parse(os.Args[2:])
        if err := cliutil.ValidateInputs(*configDir, *stateFile, *deployedSource, *deployedFile); err != nil {
            fmt.Fprintln(os.Stderr, err)
            os.Exit(1)
        }
        code, err := collect.LoadConfig(*configDir)
        if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
        state, err := collect.LoadState(*stateFile)
        if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
        var deployed collect.ResourceSet
        if *deployedSource == "mock" {
            deployed = collect.MockDeployed(code, state)
        } else {
            deployed, err = collect.LoadDeployedFile(*deployedFile)
            if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
        }
        start := nowMs()
        result := compare.ThreeWay(code, state, deployed)
        result.DurationMs = nowMs() - start
        b, err := toJSON(result)
        if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
        fmt.Println(string(b))
    case "resolve":
        fs := flag.NewFlagSet("resolve", flag.ExitOnError)
        configDir := fs.String("config-dir", ".", "Path to Terraform configuration directory")
        stateFile := fs.String("state-file", "", "Path to Terraform state file (.tfstate)")
        deployedSource := fs.String("deployed-source", "mock", "Deployed source: mock | file")
        deployedFile := fs.String("deployed-file", "", "Path to deployed data file when deployed-source=file")
        _ = fs.Parse(os.Args[2:])
        if err := cliutil.ValidateInputs(*configDir, *stateFile, *deployedSource, *deployedFile); err != nil {
            fmt.Fprintln(os.Stderr, err)
            os.Exit(1)
        }
        code, err := collect.LoadConfig(*configDir)
        if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
        state, err := collect.LoadState(*stateFile)
        if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
        var deployed collect.ResourceSet
        if *deployedSource == "mock" {
            deployed = collect.MockDeployed(code, state)
        } else {
            deployed, err = collect.LoadDeployedFile(*deployedFile)
            if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
        }
        start := nowMs()
        result := compare.ThreeWay(code, state, deployed)
        result.DurationMs = nowMs() - start
        runInteractive(result)
    default:
        fmt.Fprintln(os.Stderr, "unknown command")
        os.Exit(1)
    }
}

func toJSON(r compare.Result) ([]byte, error) {
    return json.MarshalIndent(r, "", "  ")
}

func runInteractive(r compare.Result) {
    selections := make(map[int]string)
    filterDrift := ""
    filterSeverity := ""
    reader := bufio.NewReader(os.Stdin)
    printList := func() {
        fmt.Println("Index | Key                          | Drift                 | Severity | Selection")
        for i, it := range r.Items {
            if filterDrift != "" && it.Drift != filterDrift { continue }
            if filterSeverity != "" && it.Severity != filterSeverity { continue }
            sel := selections[i]
            if sel == "" { sel = "-" }
            fmt.Printf("%5d | %-28s | %-20s | %-8s | %s\n", i+1, it.Key, it.Drift, it.Severity, sel)
        }
        fmt.Println("Commands: show <n> | select <n> code|state|deployed | filter drift=<v> | filter severity=<v> | clear | list | q")
    }
    printList()
    for {
        fmt.Print("> ")
        line, _ := reader.ReadString('\n')
        line = strings.TrimSpace(line)
        if line == "q" { return }
        if line == "list" { printList(); continue }
        if line == "clear" { filterDrift = ""; filterSeverity = ""; printList(); continue }
        if strings.HasPrefix(line, "filter ") {
            args := strings.TrimPrefix(line, "filter ")
            if strings.HasPrefix(args, "drift=") { filterDrift = strings.TrimPrefix(args, "drift="); printList(); continue }
            if strings.HasPrefix(args, "severity=") { filterSeverity = strings.TrimPrefix(args, "severity="); printList(); continue }
            fmt.Println("invalid filter; use drift=<v> or severity=<v>")
            continue
        }
        if strings.HasPrefix(line, "show ") {
            var idx int
            fmt.Sscanf(strings.TrimPrefix(line, "show "), "%d", &idx)
            if idx <= 0 || idx > len(r.Items) { fmt.Println("invalid index"); continue }
            it := r.Items[idx-1]
            fmt.Println("Key:", it.Key)
            fmt.Println("Drift:", it.Drift)
            if it.Code != nil { fmt.Println("Code:", it.Code.Attrs) } else { fmt.Println("Code:", "<nil>") }
            if it.State != nil { fmt.Println("State:", it.State.Attrs) } else { fmt.Println("State:", "<nil>") }
            if it.Deployed != nil { fmt.Println("Deployed:", it.Deployed.Attrs) } else { fmt.Println("Deployed:", "<nil>") }
            continue
        }
        if strings.HasPrefix(line, "select ") {
            parts := strings.Fields(line)
            if len(parts) != 3 { fmt.Println("usage: select <n> code|state|deployed"); continue }
            var idx int
            fmt.Sscanf(parts[1], "%d", &idx)
            if idx <= 0 || idx > len(r.Items) { fmt.Println("invalid index"); continue }
            choice := parts[2]
            if choice != "code" && choice != "state" && choice != "deployed" { fmt.Println("invalid choice"); continue }
            selections[idx-1] = choice
            fmt.Printf("selected %s for %s\n", choice, r.Items[idx-1].Key)
            continue
        }
        fmt.Println("unknown command")
    }
}

func nowMs() int {
    return int(timeNow().UnixNano() / 1e6)
}

var timeNow = func() time.Time { return time.Now() }