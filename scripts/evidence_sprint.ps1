param(
    [string[]]$RepoUrls = @(
        "https://github.com/etcd-io/etcd.git",
        "https://github.com/spf13/cobra.git",
        "https://github.com/go-chi/chi.git"
    ),
    [string]$GitResolveBinary = "",
    [string]$EvidenceRoot = "",
    [int]$MaxMergeCandidates = 300,
    [int]$MaxConflictScenariosPerRepo = 3,
    [int]$MaxAttemptsPerRepo = 120,
    [int]$MaxSecondsPerRepo = 180,
    [int]$ProgressEvery = 25,
    [int]$MaxChangedFilesPerCandidate = 200,
    [ValidateSet("auto", "strict", "balanced", "aggressive")]
    [string]$PolicyProfile = "auto",
    [int]$TimeoutSeconds = 60,
    [switch]$SkipClone,
    [switch]$KeepRepos
)

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"
if (Get-Variable -Name PSNativeCommandUseErrorActionPreference -ErrorAction SilentlyContinue) {
    $PSNativeCommandUseErrorActionPreference = $false
}

# Support terminal invocations that pass multiple URLs as a single comma-separated argument.
if ($RepoUrls.Count -eq 1 -and $RepoUrls[0] -like "*,*") {
    $RepoUrls = @($RepoUrls[0].Split(',') | ForEach-Object { $_.Trim() } | Where-Object { $_ -ne "" })
}

function Invoke-NativeCommand {
    param(
        [string]$FilePath,
        [string[]]$ArgumentList,
        [string]$WorkingDirectory
    )

    $oldLocation = Get-Location
    $oldErrPref = $ErrorActionPreference
    try {
        Set-Location $WorkingDirectory
        $ErrorActionPreference = "Continue"
        $raw = & $FilePath @ArgumentList 2>&1
        $exitCode = $LASTEXITCODE
    }
    finally {
        $ErrorActionPreference = $oldErrPref
        Set-Location $oldLocation
    }

    $combined = ($raw | ForEach-Object { $_.ToString() }) -join "`n"

    return [PSCustomObject]@{
        ExitCode = $exitCode
        Output   = $combined.Trim()
    }
}

function Ensure-Directory {
    param([string]$Path)

    if (-not (Test-Path -Path $Path -PathType Container)) {
        New-Item -Path $Path -ItemType Directory | Out-Null
    }
}

function Resolve-RepoNameFromUrl {
    param([string]$Url)

    $name = Split-Path -Leaf $Url
    if ($name.EndsWith(".git")) {
        $name = $name.Substring(0, $name.Length - 4)
    }
    return $name
}

function Resolve-GitResolveBinary {
    param(
        [string]$Requested,
        [string]$WorkspaceRoot,
        [string]$ToolsDir
    )

    if ($Requested -and (Test-Path -Path $Requested -PathType Leaf)) {
        return (Resolve-Path $Requested).Path
    }

    $existing = Join-Path $WorkspaceRoot "gitresolve.exe"
    if (Test-Path -Path $existing -PathType Leaf) {
        return $existing
    }

    Ensure-Directory -Path $ToolsDir
    $builtBinary = Join-Path $ToolsDir "gitresolve-evidence.exe"
    $build = Invoke-NativeCommand -FilePath "go" -ArgumentList @("build", "-o", $builtBinary, ".") -WorkingDirectory $WorkspaceRoot
    if ($build.ExitCode -ne 0) {
        throw "go build failed while preparing gitresolve binary: $($build.Output)"
    }

    return $builtBinary
}

function Reset-RepoState {
    param([string]$RepoDir)

    [void](Invoke-NativeCommand -FilePath "git" -ArgumentList @("merge", "--abort") -WorkingDirectory $RepoDir)
    [void](Invoke-NativeCommand -FilePath "git" -ArgumentList @("reset", "--hard") -WorkingDirectory $RepoDir)
    [void](Invoke-NativeCommand -FilePath "git" -ArgumentList @("clean", "-fd") -WorkingDirectory $RepoDir)
}

function Get-ActionCount {
    param(
        [object]$Stats,
        [string[]]$ActionNames
    )

    $sum = 0
    foreach ($entry in $Stats.actions) {
        if ($ActionNames -contains $entry.action) {
            $sum += [int]$entry.count
        }
    }
    return $sum
}

function Get-ReasonCount {
    param(
        [object]$Stats,
        [string]$ReasonCode
    )

    $sum = 0
    foreach ($entry in $Stats.top_reason_codes) {
        if ($entry.reason_code -eq $ReasonCode) {
            $sum += [int]$entry.count
        }
    }
    return $sum
}

function Run-ConflictScenario {
    param(
        [string]$RepoDir,
        [string]$MergeCommit,
        [string]$P1,
        [string]$P2,
        [string]$GitResolvePath,
        [string]$Policy,
        [int]$TimeoutSec
    )

    Reset-RepoState -RepoDir $RepoDir

    $checkout = Invoke-NativeCommand -FilePath "git" -ArgumentList @("checkout", "--detach", $P1) -WorkingDirectory $RepoDir
    if ($checkout.ExitCode -ne 0) {
        return $null
    }

    $merge = Invoke-NativeCommand -FilePath "git" -ArgumentList @("merge", "--no-commit", "--no-ff", $P2) -WorkingDirectory $RepoDir

    $unmerged = Invoke-NativeCommand -FilePath "git" -ArgumentList @("diff", "--name-only", "--diff-filter=U") -WorkingDirectory $RepoDir
    $files = @($unmerged.Output -split "`r?`n" | Where-Object { $_.Trim() -ne "" })

    if ($files.Count -eq 0) {
        Reset-RepoState -RepoDir $RepoDir
        return $null
    }

    $dbPath = Join-Path $RepoDir ".git\gitresolve.db"
    if (Test-Path -Path $dbPath -PathType Leaf) {
        Remove-Item -Path $dbPath -Force
    }

    $resolveArgs = @(
        "resolve",
        "--non-interactive",
        "--policy-profile", $Policy,
        "--timeout", ("{0}s" -f $TimeoutSec)
    )
    $resolveRun = Invoke-NativeCommand -FilePath $GitResolvePath -ArgumentList $resolveArgs -WorkingDirectory $RepoDir

    $statsRun = Invoke-NativeCommand -FilePath $GitResolvePath -ArgumentList @("stats", "--json") -WorkingDirectory $RepoDir
    if ($statsRun.ExitCode -ne 0) {
        Reset-RepoState -RepoDir $RepoDir
        return [PSCustomObject]@{
            merge_commit      = $MergeCommit
            parent_1          = $P1
            parent_2          = $P2
            conflict_files    = $files
            conflict_count    = $files.Count
            resolve_exit_code = $resolveRun.ExitCode
            stats_error       = $statsRun.Output
        }
    }

    $stats = $statsRun.Output | ConvertFrom-Json

    $autoCount = Get-ActionCount -Stats $stats -ActionNames @("auto-resolve", "resolve")
    $manualCount = Get-ActionCount -Stats $stats -ActionNames @("manual", "manual-escalate")
    $blockedUnsafe = Get-ReasonCount -Stats $stats -ReasonCode "validation.syntax_failed"

    Reset-RepoState -RepoDir $RepoDir

    return [PSCustomObject]@{
        merge_commit              = $MergeCommit
        parent_1                  = $P1
        parent_2                  = $P2
        conflict_files            = $files
        conflict_count            = $files.Count
        merge_exit_code           = $merge.ExitCode
        resolve_exit_code         = $resolveRun.ExitCode
        total_decisions           = [int]$stats.total_decisions
        auto_resolved             = $autoCount
        manual_escalated          = $manualCount
        manual_escalation_rate    = [double]$stats.manual_escalation_rate
        validator_blocked_unsafe  = $blockedUnsafe
        stats_schema_version      = [string]$stats.schema_version
        stats_operation           = [string]$stats.operation
        top_reason_codes          = $stats.top_reason_codes
    }
}

function Should-SkipCandidate {
    param(
        [string]$RepoDir,
        [string]$MergeCommit,
        [int]$MaxChangedFiles
    )

    $changedFilesRun = Invoke-NativeCommand -FilePath "git" -ArgumentList @("diff-tree", "--no-commit-id", "-m", "--name-only", "-r", $MergeCommit) -WorkingDirectory $RepoDir
    if ($changedFilesRun.ExitCode -ne 0) {
        return $true
    }

    $files = @($changedFilesRun.Output -split "`r?`n" | Where-Object { $_.Trim() -ne "" })
    if ($files.Count -eq 0) {
        return $true
    }
    if ($files.Count -gt $MaxChangedFiles) {
        return $true
    }

    $interesting = $false
    foreach ($f in $files) {
        if ($f.EndsWith(".go") -or $f.EndsWith("go.mod") -or $f.EndsWith("go.sum")) {
            $interesting = $true
            break
        }
    }
    return -not $interesting
}

$workspaceRoot = Split-Path -Parent $PSScriptRoot
if (-not $EvidenceRoot) {
    $EvidenceRoot = Join-Path $workspaceRoot "evidence"
}

$reposDir = Join-Path $EvidenceRoot "_repos"
$rawDir = Join-Path $EvidenceRoot "raw"
$toolsDir = Join-Path $EvidenceRoot "tools"
Ensure-Directory -Path $EvidenceRoot
Ensure-Directory -Path $reposDir
Ensure-Directory -Path $rawDir
Ensure-Directory -Path $toolsDir

$gitresolvePath = Resolve-GitResolveBinary -Requested $GitResolveBinary -WorkspaceRoot $workspaceRoot -ToolsDir $toolsDir

$repoResults = New-Object System.Collections.Generic.List[object]

foreach ($repoUrl in $RepoUrls) {
    $repoName = Resolve-RepoNameFromUrl -Url $repoUrl
    $repoDir = Join-Path $reposDir $repoName

    if (-not $SkipClone) {
        if (-not (Test-Path -Path $repoDir -PathType Container)) {
            Write-Host "Cloning $repoUrl ..."
            $clone = Invoke-NativeCommand -FilePath "git" -ArgumentList @("clone", "--filter=blob:none", $repoUrl, $repoDir) -WorkingDirectory $reposDir
            if ($clone.ExitCode -ne 0) {
                Write-Host "Clone failed for $repoUrl"
                $repoResults.Add([PSCustomObject]@{
                    repo                       = $repoName
                    url                        = $repoUrl
                    error                      = "clone_failed"
                    error_detail               = $clone.Output
                    scenarios_attempted        = 0
                    scenarios_with_conflicts   = 0
                    total_decisions            = 0
                    auto_resolved              = 0
                    manual_escalated           = 0
                    manual_escalation_rate     = 0
                    validator_blocked_unsafe   = 0
                    scenario_details           = @()
                })
                continue
            }
        }
        else {
            [void](Invoke-NativeCommand -FilePath "git" -ArgumentList @("fetch", "--all", "--prune") -WorkingDirectory $repoDir)
        }
    }
    elseif (-not (Test-Path -Path $repoDir -PathType Container)) {
        throw "SkipClone was set, but local repo path does not exist: $repoDir"
    }

    Write-Host "\nAnalyzing merge history in $repoName ..."
    $revList = Invoke-NativeCommand -FilePath "git" -ArgumentList @("rev-list", "--merges", ("--max-count={0}" -f $MaxMergeCandidates), "HEAD") -WorkingDirectory $repoDir
    if ($revList.ExitCode -ne 0) {
        $repoResults.Add([PSCustomObject]@{
            repo                       = $repoName
            url                        = $repoUrl
            error                      = "rev_list_failed"
            error_detail               = $revList.Output
            scenarios_attempted        = 0
            scenarios_with_conflicts   = 0
            total_decisions            = 0
            auto_resolved              = 0
            manual_escalated           = 0
            manual_escalation_rate     = 0
            validator_blocked_unsafe   = 0
            scenario_details           = @()
        })
        continue
    }

    $mergeCommits = @($revList.Output -split "`r?`n" | Where-Object { $_.Trim() -ne "" })
    $scenarios = New-Object System.Collections.Generic.List[object]
    $attempted = 0
    $considered = 0
    $skippedPrefilter = 0
    $repoStart = Get-Date

    foreach ($mergeCommit in $mergeCommits) {
        if ($scenarios.Count -ge $MaxConflictScenariosPerRepo) {
            break
        }

        $elapsed = (Get-Date) - $repoStart
        if ($elapsed.TotalSeconds -ge $MaxSecondsPerRepo) {
            Write-Host ("  stopping after time budget: {0}s" -f $MaxSecondsPerRepo)
            break
        }
        if ($attempted -ge $MaxAttemptsPerRepo) {
            Write-Host ("  stopping after attempt budget: {0}" -f $MaxAttemptsPerRepo)
            break
        }

        $considered++
        if ((($considered % $ProgressEvery) -eq 0) -and $ProgressEvery -gt 0) {
            Write-Host ("  progress: considered={0} attempted={1} captured={2} skipped={3}" -f $considered, $attempted, $scenarios.Count, $skippedPrefilter)
        }

        if (Should-SkipCandidate -RepoDir $repoDir -MergeCommit $mergeCommit -MaxChangedFiles $MaxChangedFilesPerCandidate) {
            $skippedPrefilter++
            continue
        }

        $parents = Invoke-NativeCommand -FilePath "git" -ArgumentList @("rev-list", "--parents", "-n", "1", $mergeCommit) -WorkingDirectory $repoDir
        if ($parents.ExitCode -ne 0) {
            continue
        }

        $parts = @($parents.Output -split "\s+" | Where-Object { $_.Trim() -ne "" })
        if ($parts.Count -lt 3) {
            continue
        }

        $attempted++
        $scenario = Run-ConflictScenario -RepoDir $repoDir -MergeCommit $mergeCommit -P1 $parts[1] -P2 $parts[2] -GitResolvePath $gitresolvePath -Policy $PolicyProfile -TimeoutSec $TimeoutSeconds
        if ($null -ne $scenario) {
            $scenarios.Add($scenario)
            Write-Host ("  captured conflict scenario {0}/{1}: {2}" -f $scenarios.Count, $MaxConflictScenariosPerRepo, $mergeCommit.Substring(0, 12))
        }
    }

    Reset-RepoState -RepoDir $repoDir

    $totalDecisions = 0
    $autoResolved = 0
    $manualEscalated = 0
    $blockedUnsafe = 0

    foreach ($s in $scenarios) {
        if ($s.PSObject.Properties.Name -contains "total_decisions") {
            $totalDecisions += [int]$s.total_decisions
        }
        if ($s.PSObject.Properties.Name -contains "auto_resolved") {
            $autoResolved += [int]$s.auto_resolved
        }
        if ($s.PSObject.Properties.Name -contains "manual_escalated") {
            $manualEscalated += [int]$s.manual_escalated
        }
        if ($s.PSObject.Properties.Name -contains "validator_blocked_unsafe") {
            $blockedUnsafe += [int]$s.validator_blocked_unsafe
        }
    }

    $manualRate = 0.0
    if ($totalDecisions -gt 0) {
        $manualRate = [math]::Round(($manualEscalated * 100.0) / $totalDecisions, 2)
    }

    $repoSummary = [PSCustomObject]@{
        repo                     = $repoName
        url                      = $repoUrl
        candidates_considered    = $considered
        candidates_skipped       = $skippedPrefilter
        scenarios_attempted      = $attempted
        scenarios_with_conflicts = $scenarios.Count
        total_decisions          = $totalDecisions
        auto_resolved            = $autoResolved
        manual_escalated         = $manualEscalated
        manual_escalation_rate   = $manualRate
        validator_blocked_unsafe = $blockedUnsafe
        scenario_details         = $scenarios
    }

    $repoResults.Add($repoSummary)

    $rawPath = Join-Path $rawDir ("{0}.json" -f $repoName)
    $repoSummary | ConvertTo-Json -Depth 8 | Set-Content -Path $rawPath -Encoding UTF8

    if (-not $KeepRepos) {
        try {
            Remove-Item -Path $repoDir -Recurse -Force -ErrorAction Stop
        }
        catch {
            Write-Warning ("Could not remove temporary repo clone {0}: {1}" -f $repoDir, $_.Exception.Message)
        }
    }
}

$grandTotal = 0
$grandAuto = 0
$grandManual = 0
$grandBlockedUnsafe = 0
$grandScenarios = 0

foreach ($repo in $repoResults) {
    $grandTotal += [int]$repo.total_decisions
    $grandAuto += [int]$repo.auto_resolved
    $grandManual += [int]$repo.manual_escalated
    $grandBlockedUnsafe += [int]$repo.validator_blocked_unsafe
    $grandScenarios += [int]$repo.scenarios_with_conflicts
}

$grandManualRate = 0.0
if ($grandTotal -gt 0) {
    $grandManualRate = [math]::Round(($grandManual * 100.0) / $grandTotal, 2)
}

$grandAutoRate = 0.0
if ($grandTotal -gt 0) {
    $grandAutoRate = [math]::Round(($grandAuto * 100.0) / $grandTotal, 2)
}

$report = [PSCustomObject]@{
    generated_at_utc = [DateTime]::UtcNow.ToString("o")
    config = [PSCustomObject]@{
        repo_urls                       = $RepoUrls
        max_merge_candidates            = $MaxMergeCandidates
        max_conflict_scenarios_per_repo = $MaxConflictScenariosPerRepo
        policy_profile                  = $PolicyProfile
        timeout_seconds                 = $TimeoutSeconds
        gitresolve_binary               = $gitresolvePath
    }
    aggregate = [PSCustomObject]@{
        repos_evaluated          = $repoResults.Count
        conflict_scenarios       = $grandScenarios
        total_decisions          = $grandTotal
        auto_resolved            = $grandAuto
        manual_escalated         = $grandManual
        auto_resolve_rate        = $grandAutoRate
        manual_escalation_rate   = $grandManualRate
        validator_blocked_unsafe = $grandBlockedUnsafe
    }
    repositories = $repoResults
}

$reportJsonPath = Join-Path $EvidenceRoot "evidence_report.json"
$report | ConvertTo-Json -Depth 10 | Set-Content -Path $reportJsonPath -Encoding UTF8

$md = New-Object System.Collections.Generic.List[string]
$md.Add("# GitResolve Evidence Sprint Report")
$md.Add("")
$md.Add(("Generated (UTC): {0}" -f $report.generated_at_utc))
$md.Add("")
$md.Add("## Aggregate")
$md.Add("")
$md.Add(("- Conflict scenarios: {0}" -f $report.aggregate.conflict_scenarios))
$md.Add(("- Total decisions: {0}" -f $report.aggregate.total_decisions))
$md.Add(("- Auto-resolved: {0} ({1}%)" -f $report.aggregate.auto_resolved, $report.aggregate.auto_resolve_rate))
$md.Add(("- Escalated to manual: {0} ({1}%)" -f $report.aggregate.manual_escalated, $report.aggregate.manual_escalation_rate))
$md.Add(("- Validator blocked unsafe writes: {0}" -f $report.aggregate.validator_blocked_unsafe))
$md.Add("")
$md.Add("## Repository Breakdown")
$md.Add("")
$md.Add("| Repository | Scenarios | Total Decisions | Auto | Manual | Manual Rate | Validator Blocked |")
$md.Add("| :--- | ---: | ---: | ---: | ---: | ---: | ---: |")
foreach ($repo in $report.repositories) {
    $md.Add(("| {0} | {1} | {2} | {3} | {4} | {5}% | {6} |" -f $repo.repo, $repo.scenarios_with_conflicts, $repo.total_decisions, $repo.auto_resolved, $repo.manual_escalated, $repo.manual_escalation_rate, $repo.validator_blocked_unsafe))
}

$reportMdPath = Join-Path $EvidenceRoot "evidence_report.md"
$md -join "`n" | Set-Content -Path $reportMdPath -Encoding UTF8

if ($report.aggregate.conflict_scenarios -eq 0) {
    Write-Warning "No conflict scenarios were captured. Increase -MaxMergeCandidates or add different repositories."
}

Write-Host "\nEvidence sprint complete."
Write-Host ("JSON report: {0}" -f $reportJsonPath)
Write-Host ("Markdown report: {0}" -f $reportMdPath)
