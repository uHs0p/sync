# Terraform Sync MVP

## System Architecture Overview

- CLI with commands: `plan`, `diff`, `resolve` implemented using standard library.
- Data collection: parses `.tf` files via lightweight parser, reads `.tfstate` JSON, and loads deployed data (mock or external file).
- Comparison engine: three-way diff computes drift categories and severity; outputs JSON including `duration_ms` metrics.
- Interactive console: list resources and allow selection to view Code/State/Deployed details; non-interactive mode prints JSON.
- External integration: optional baseline `terraform plan -detailed-exitcode` if Terraform CLI is installed.

## Installation and Setup

- Prerequisites: Go 1.24+; optional Terraform CLI in PATH.
- Build binary: `go build ./cmd/terrainsync` (outputs `terrainsync.exe` in project root on Windows).
- Verify tests: `go test ./...`.
- Example assets: `examples/simple` includes `.tf`, `.tfstate`, and `deployed.json`.

## User Guide

- Diff: `./terrainsync.exe diff -config-dir examples/simple -state-file examples/simple/terraform.tfstate --deployed-source file --deployed-file examples/simple/deployed.json`.
- Plan: `./terrainsync.exe plan -config-dir examples/simple`.
- Resolve: `./terrainsync.exe resolve -config-dir examples/simple -state-file examples/simple/terraform.tfstate --deployed-source file --deployed-file examples/simple/deployed.json` and follow on-screen prompts; press `q` to quit.

## Presentation Highlights

- Core value proposition: reconcile drift across Code, State, and Deployed with clear categorization and a guided resolution workflow.
- Key user workflows: baseline plan, three-way diff (JSON), interactive inspection of per-resource triplets.
- Technical implementation: modular collectors, pure-Go CLI, deterministic diff, input validation, graceful fallback when Terraform is missing.
- Performance metrics: each `diff` execution reports `duration_ms` in the JSON result; use this to compare runs across environments and scales.
- Scalability considerations: map-based aggregation, O(N) merging across sources; collectors and comparisons are designed for future parallelization.

## Known Limitations and Future Roadmap

- Config parsing is a simplified resource/block parser; full HCL evaluation and module graph traversal will be added.
- Deployed source is mock/file-based; provider RPC integration and cloud queries are planned.
- Resolve currently provides inspection; per-resource selection and bulk strategies will be added in subsequent phases.
- Phases 6–9: resolution actions, backups/restore, reporting (JSON/Markdown/HTML), CI/CD integration, policies (OPA), packaging.

## Performance and Scalability

- Operates on in-memory sets; comparisons scale linearly with resource count.
- `duration_ms` helps track performance; for large inputs, enable future parallel collectors.