package cliutil

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"
)

func ValidateInputs(configDir, stateFile, deployedSource, deployedFile string) error {
    if fi, err := os.Stat(configDir); err != nil || !fi.IsDir() {
        return fmt.Errorf("invalid config-dir: %s", configDir)
    }
    if stateFile != "" {
        if fi, err := os.Stat(stateFile); err != nil || fi.IsDir() {
            return fmt.Errorf("invalid state-file: %s", stateFile)
        }
    }
    ds := strings.ToLower(deployedSource)
    if ds != "mock" && ds != "file" {
        return fmt.Errorf("invalid deployed-source: %s", deployedSource)
    }
    if ds == "file" {
        if deployedFile == "" {
            return fmt.Errorf("deployed-file required when deployed-source=file")
        }
        if fi, err := os.Stat(deployedFile); err != nil || fi.IsDir() {
            return fmt.Errorf("invalid deployed-file: %s", deployedFile)
        }
    }
    abs, _ := filepath.Abs(configDir)
    _ = abs
    return nil
}