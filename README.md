# gitresolve

gitresolve is a local Go CLI for handling Git merge conflicts with deterministic, rule-based logic.

## What It Does 

- Runs locally with no network/API calls.
- Detects conflicted files from Git status using go-git.
- Locks the repository root with `.gitresolve.lock` while running.
- For each conflicted file during `merge`:
    - creates a backup as `<file>.gitresolve-orig` (non-dry-run only),
    - parses inline conflict markers,
    - classifies each conflict,
    - auto-resolves only safe categories,
        - verifies resulting content (no conflict markers, JSON/YAML validity),
    - writes updates atomically (temp file + fsync + rename),
    - stages the file if all conflicts in that file were auto-resolved.

## Commands 

- `gitresolve merge [--dry-run]`
    - auto-resolves safe conflict blocks in existing conflicted files.
    - records session and conflict history in `.git/gitresolve.db`.
- `gitresolve status`
    - inspects current conflicted files and prints per-block severity/type/auto-resolve status.
- `gitresolve scan --target <branch>`
    - predicts conflict blocks against another branch using `git merge-base` + `git merge-tree`.
- `gitresolve resolve [--file <path>] [--strategy ours|theirs|both] [--dry-run]`
    - resolves remaining conflict blocks using explicit strategy and stages result.
- `gitresolve blame [--file <path>]`
    - shows stored conflict-resolution history from SQLite.
- `gitresolve undo --steps N`
    - resets repository to a recorded snapshot SHA from recent gitresolve sessions.

## Real End-to-End Flow (`gitresolve merge`)

This is the concrete runtime flow in the current implementation:

```mermaid
flowchart TD
    A[Open repo and create .gitresolve.lock] --> B[List conflicted files from worktree status]
    B --> C{For each conflicted file}
    C --> D[Backup original to .gitresolve-orig\n(skip when --dry-run)]
    D --> E[Parse conflict markers\n<<<<<<< ======= >>>>>>>]
    E --> F[Classify conflict blocks]
    F --> G{Auto-resolvable?}
    G -- Yes --> H[Compile resolved content]
    H --> I[Verify output\n(no markers + JSON/YAML checks)]
    I --> J[Atomic write\n(temp + fsync + rename)]
    J --> K{All blocks in file resolved?}
    K -- Yes --> L[Stage file with git add]
    K -- No --> M[Leave remaining blocks for manual resolution]
    G -- No --> M
    L --> N[Record conflict history in SQLite]
    M --> N
    N --> O[Release lock and finish]
```

1. Open repository (`.`) and create lock file `.gitresolve.lock`.
2. Query conflicted files from worktree status (`UpdatedButUnmerged`).
3. For each file:
     - backup original file to `<file>.gitresolve-orig` (unless `--dry-run`),
     - read file content from disk,
     - parse conflict blocks marked by `<<<<<<<`, `=======`, `>>>>>>>`,
     - classify each block,
     - auto-resolve if allowed,
     - rebuild full file content,
     - write atomically,
     - stage file with `git add` if all conflicts in that file were auto-resolved.
4. Release lock file at exit.

Additionally, `merge` records:

- one session snapshot (`operation=merge`, commit SHA before changes),
- one conflict-history row per conflict block (auto or manual-required).

### Conflict Classification Rules Actually Used

Classification is line/rule-based in current `merge` execution:

- `TypeWhitespace` -> auto-resolve.
- `TypeIdentical` -> auto-resolve.
- `TypeImport` -> auto-resolve (dedupe by exact line string).
- `TypeStructured` -> detected by extension (`.json/.yaml/.yml/.toml`) but not auto-resolved in merge path.
- `TypeDeleteModify`, `TypeSignature`, sensitive-path `TypeLogic` -> not auto-resolved.
- default logic conflicts -> not auto-resolved.

## Concrete Example

Suppose a file contains this conflict:

```txt
<<<<<<< HEAD
import "fmt"
import "net/http"
=======
import "fmt"
import "os"
>>>>>>> feature
```

`gitresolve merge` will:

- classify this block as import conflict,
- merge both sides,
- dedupe duplicate lines,
- produce:

```txt
import "fmt"
import "net/http"
import "os"
```

If every block in the file is auto-resolved, the file is staged.
If any block is not safely auto-resolvable, merge leaves those markers in place for manual review and prints escalation info.

### Dry-Run Note

`--dry-run` skips backup creation and does not persist file changes. In the current code path, attempted writes in dry-run return a dry-run error signal, so you may see an "atomic write failed" message even though no write is intended.

## Installation

```bash
go install github.com/jhanvi857/gitresolve@latest
```

Or local build:

```bash
git clone https://github.com/jhanvi857/gitresolve
cd gitresolve
go build -o bin/gitresolve .
```

## Development Notes

- Lock file path: `.gitresolve.lock` at repo root.
- Backup convention: `<file>.gitresolve-orig`.
- Atomic write strategy: temp file in same directory + `Sync` + `Rename`.
- If no conflicted files are found, `merge` exits with a status message indicating no conflicts.

## Scope and Practical Limits

- `merge` is deterministic and conservative: it auto-resolves safe classes only.
- For complex or semantic conflicts, use `resolve --strategy ...` (or manual editing) and then rerun status.
- `scan` depends on local Git support for `merge-tree`.