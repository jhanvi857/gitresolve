# gitresolve

[![Go Report Card](https://goreportcard.com/badge/github.com/jhanvi857/gitresolve)](https://goreportcard.com/report/github.com/jhanvi857/gitresolve)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

> A locally executed, deterministic Git merge conflict solver built on Abstract Syntax Trees and structured serialization.

Handling Git merge conflicts manually is error-prone, dangerous, and time-consuming. **gitresolve** intercepts your broken Git status, structurally evaluates the conflicted files, and safely resolves them without risking repository corruption.

## Why use gitresolve?

Standard Git performs blind text integrations. `gitresolve` actually parses your structures.

| Capability | Standard Git Merge | gitresolve Engine |
| :--- | :--- | :--- |
| **JSON/YAML/TOML** | Fails on text collisions | Merges data maps dynamically |
| **Logic Refactors** | Blindly overwrites text | Verifies AST signatures |
| **Formatting** | Flags entire block | Auto-heals whitespace |
| **Complex Deletions** | Overwrites or fails | Triggers Interactive Prompter |

---

## Core Engine Features

### 1. Abstract Syntax Tree (AST) Intelligence
Instead of analyzing raw text, `gitresolve` natively integrates `go-tree-sitter`. It compiles conflicting blocks into syntax trees to accurately pinpoint function signature modifications and logical refactors in **Go**, **JavaScript**, and **TypeScript**. If the syntax breaks, the merge is halted.

### 2. Structured Data Auto-Merger
Configuration conflicts are solved instantly. The engine performs a deep recursive map merge for **JSON**, **YAML**, and **TOML** using language-native parsers. 
- **Conservative Array Merging**: To avoid silent configuration corruption, overlapping array modifications (e.g., both branches adding different items to a server list) are marked as conflicts for human review.
- **Critical File Protection**: Files like `package.json`, `go.mod`, and `cargo.toml` are treated as high-severity. Auto-resolution is disabled by default for these files to ensure dependency integrity.

### 3. Interactive Terminal Prompter
Safety is paramount. When encountering highly critical conflicts (like sensitive authentication modifications, mass code deletions, or complex logic), the engine strictly pauses operations.
- **AST Fallback**: If parsing fails, it immediately delegates to the Interactive Prompter.
- **Lock Resilience**: Uses a multi-layered locking system with PID verification and signal handling (SIGINT/SIGTERM) to prevent stalled executions and allow recovery from crashes.

---

## Architectural Workflow

The operational flow prioritizes safety, executing natively without external API dependencies.

```mermaid
flowchart TD
    Start[User Triggers Resolve] --> LockRepo[Lock Repository]
    LockRepo --> ReadGit[Identify Conflicted Files]
    ReadGit --> LoopFiles{For Each File}
    
    LoopFiles --> Parsed[Parse Conflict Markers]
    Parsed --> Classify[AST & Heuristic Classification]
    
    Classify --> IsStructured{Is Config File?}
    IsStructured -- Yes --> StructuralMerge[Deep Map Serialization]
    StructuralMerge --> Output[Generate Clean File]
    
    IsStructured -- No --> IsSafeText{Is Safe Text Change?}
    IsSafeText -- Yes --> AutoResolve[Auto Merge Imports/Whitespace]
    AutoResolve --> Output
    
    IsSafeText -- No --> UserInput[Interactive Terminal Prompter]
    UserInput --> ReceiveInput[Receive Developer Decision]
    ReceiveInput --> Output
    
    Output --> Verify[Verify Syntax Validity]
    Verify --> Write[Atomic Write & Git Stage]
    Write --> NextFile{More Files?}
    
    NextFile -- Yes --> LoopFiles
    NextFile -- No --> Unlock[Unlock Repository]
    Unlock --> Finish[Cleanup & Exit]
```

---

## Command Reference

The CLI is designed to seamlessly sit on top of your standard Git workflow.

| Command | Description |
| :--- | :--- |
| `gitresolve resolve` | **Default Command**. Resolves all remaining conflicts interactively or automatically based on AST safety checks. |
| `gitresolve scan --target <branch>` | Non-destructively predicts conflicts against another branch using modern `git merge-tree`. |
| `gitresolve status` | Inspects current conflicted files and prints per-block severity, type, and auto-resolution status. |
| `gitresolve merge [--no-auto-structured]` | Auto-resolves safe blocks. Use `--no-auto-structured` to skip automated merging of JSON/YAML/TOML. |
| `gitresolve blame [--file <path>]` | Displays stored conflict-resolution history logged in the local SQLite database. |
| `gitresolve undo --steps N` | Resets the repository to a recorded snapshot SHA from recent GitResolve sessions. |

---

## Installation & Quick Start

Install directly via Go:

```bash
go install github.com/jhanvi857/gitresolve@latest
```

Or build the binary locally from source:

```bash
git clone https://github.com/jhanvi857/gitresolve
cd gitresolve
go build -o bin/gitresolve .
```

### Automatic Classification Ruleset
To ensure continuous repository integrity, conflicts are mapped to the following deterministic rules:

* **TypeWhitespace:** Auto-resolved (merges standard formatting differences).
* **TypeIdentical:** Auto-resolved (deduplicates exact parallel changes).
* **TypeImport:** Auto-resolved (deduplicates imports safely across languages).
* **TypeStructured:** Auto-resolved for non-critical files if no array/scalar overlaps exist.
* **TypeDeleteModify:** Delegated to Interactive Prompt (high severity deletion protection).
* **TypeSignature:** Delegated to Prompt (AST-verified architecture changes).
* **TypeLogic:** Delegated to Prompt (core logic modifications or sensitive paths).

*(All state changes are locked atomically with `.gitresolve.lock` and backed up immediately as `<file>.gitresolve-orig`.)*

---

## Technical Reliability
`gitresolve` is verified against a corpus of **18 real-world conflict scenarios** across multiple languages and formats, including:
- **TypeScript (TSX)** components with logical shifts.
- **Go** interface and function signature modifications.
- **Nested YAML/JSON** configuration overlaps.
- **Critical Dependency** files (`package.json`, `go.mod`, `Cargo.toml`).
- **Sensitive Security Paths** (auth, crypto, migrations).
- **TOML** deep map merges in Rust/Cargo projects.
