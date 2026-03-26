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
Configuration conflicts are solved instantly. Whether it is `package.json` arrays or overlapping `.yaml` configurations, the engine unmarshals both divergent branches, performs a deep recursive map merge, and reserializes the output safely. 

### 3. Interactive Terminal Prompter
Safety is paramount. When encountering highly critical conflicts (like sensitive authentication modifications, mass code deletions, or complex logic), the engine strictly pauses operations. It isolates the conflict into a clean side-by-side terminal comparison and blocks for explicit human input: `[O]urs`, `[T]heirs`, or `[B]oth`.

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
| `gitresolve scan --target <branch>` | Predicts and previews conflict blocks against another branch **before** you trigger a merge. |
| `gitresolve status` | Inspects current conflicted files and prints per-block severity, type, and auto-resolution status. |
| `gitresolve merge [--dry-run]` | Strictly auto-resolves safe conflict blocks in existing conflicted files silently. |
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
* **TypeStructured:** Auto-resolved (JSON/YAML/TOML configurations deep-merged).
* **TypeDeleteModify:** Delegated to Interactive Prompt (high severity deletion protection).
* **TypeSignature:** Delegated to Prompt (AST-verified architecture changes).
* **TypeLogic:** Delegated to Prompt (core logic modifications or sensitive paths).

*(All state changes are locked atomically with `.gitresolve.lock` and backed up immediately as `<file>.gitresolve-orig`.)*
