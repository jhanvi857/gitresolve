param(
    [string]$TestsRoot = (Join-Path $PSScriptRoot "tests"),
    [string]$GitResolveBinary = "",
    [switch]$ResolveRemaining = $true,
    [ValidateSet("interactive", "ours", "theirs", "both")]
    [string]$ResolveStrategy = "theirs",
    [switch]$StopOnFirstFailure
)

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"
if (Get-Variable -Name PSNativeCommandUseErrorActionPreference -ErrorAction SilentlyContinue) {
    $PSNativeCommandUseErrorActionPreference = $false
}

function Resolve-GitResolveBinary {
    param([string]$Requested)

    if ($Requested -and (Test-Path -Path $Requested -PathType Leaf)) {
        return (Resolve-Path $Requested).Path
    }

    $candidates = @(
        (Join-Path $PSScriptRoot "gitresolve-windows-amd64.exe"),
        (Join-Path $PSScriptRoot "gitresolve-windows-arm64.exe"),
        (Join-Path $PSScriptRoot "gitresolve.exe")
    )

    foreach ($candidate in $candidates) {
        if (Test-Path -Path $candidate -PathType Leaf) {
            return (Resolve-Path $candidate).Path
        }
    }

    throw "No usable gitresolve binary found. Checked: $($candidates -join ', ')"
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
        StdOut   = $combined
        StdErr   = ""
        Combined = $combined.Trim()
    }
}

function Test-JsonFile {
    param([string]$Path)
    try {
        Get-Content -Raw -Path $Path | ConvertFrom-Json -ErrorAction Stop | Out-Null
        return $null
    }
    catch {
        return "JSON parse failed: $($_.Exception.Message)"
    }
}

function Test-YamlFile {
    param([string]$Path)
    if (-not (Get-Command ConvertFrom-Yaml -ErrorAction SilentlyContinue)) {
        return "YAML parser unavailable in this PowerShell version"
    }

    try {
        Get-Content -Raw -Path $Path | ConvertFrom-Yaml -ErrorAction Stop | Out-Null
        return $null
    }
    catch {
        return "YAML parse failed: $($_.Exception.Message)"
    }
}

function Test-GoFile {
    param([string]$Path)
    if (-not (Get-Command gofmt -ErrorAction SilentlyContinue)) {
        return "gofmt unavailable for syntax check"
    }

    $run = Invoke-NativeCommand -FilePath "gofmt" -ArgumentList @("-e", "-l", $Path) -WorkingDirectory (Split-Path -Parent $Path)
    if ($run.ExitCode -ne 0) {
        return "Go parse failed (gofmt returned non-zero)"
    }

    return $null
}

function Find-MarkerFiles {
    param([string]$RepoPath)

    return Get-ChildItem -Path $RepoPath -Recurse -File -Force |
        Where-Object {
            $_.FullName -notmatch "\\.git\\" -and
            $_.Name -notlike "*.gitresolve-orig" -and
            $_.Name -ne ".gitresolve.lock"
        } |
        Where-Object {
            try {
                Select-String -Path $_.FullName -Pattern "^<<<<<<<|^=======|^>>>>>>>" -SimpleMatch:$false -Quiet
            }
            catch {
                $false
            }
        }
}

function Test-Fixture {
    param(
        [string]$FixturePath,
        [string]$BinaryPath,
        [bool]$DoResolve,
        [string]$Strategy
    )

    $fixtureName = Split-Path -Leaf $FixturePath
    Write-Host ""
    Write-Host "=== Running $fixtureName ==="

    Push-Location $FixturePath
    try {
        $mergeRun = Invoke-NativeCommand -FilePath $BinaryPath -ArgumentList @("merge") -WorkingDirectory $FixturePath
        $mergeOutput = $mergeRun.Combined
        $mergeExit = $mergeRun.ExitCode

        $resolveOutput = $null
        $resolveExit = 0
        if ($DoResolve) {
            $resolveRun = Invoke-NativeCommand -FilePath $BinaryPath -ArgumentList @("resolve", "--strategy", $Strategy, "--non-interactive") -WorkingDirectory $FixturePath
            $resolveOutput = $resolveRun.Combined
            $resolveExit = $resolveRun.ExitCode
        }

        $gitRun = Invoke-NativeCommand -FilePath "git" -ArgumentList @("diff", "--name-only", "--diff-filter=U") -WorkingDirectory $FixturePath
        $unmerged = @($gitRun.Combined -split "`r?`n" | Where-Object { $_.Trim() -ne "" })
        $markerFiles = @(Find-MarkerFiles -RepoPath $FixturePath)

        $syntaxErrors = New-Object System.Collections.Generic.List[string]
        $filesToCheck = Get-ChildItem -Path $FixturePath -Recurse -File |
            Where-Object {
                $_.FullName -notmatch "\\.git\\" -and
                $_.Name -notlike "*.gitresolve-orig" -and
                $_.Name -ne ".gitresolve.lock"
            }

        foreach ($file in $filesToCheck) {
            $ext = $file.Extension.ToLowerInvariant()
            $err = $null
            switch ($ext) {
                ".json" { $err = Test-JsonFile -Path $file.FullName }
                ".yaml" { $err = Test-YamlFile -Path $file.FullName }
                ".yml"  { $err = Test-YamlFile -Path $file.FullName }
                ".go"   { $err = Test-GoFile -Path $file.FullName }
                ".toml" { $err = "TOML parser check not available in PowerShell; rely on gitresolve verification" }
                default { }
            }

            if ($err) {
                $relative = $file.FullName.Substring($FixturePath.Length).TrimStart('\\')
                $syntaxErrors.Add("$relative -> $err")
            }
        }

        $unmergedCount = @($unmerged | Where-Object { $_ -and $_.Trim() -ne "" }).Count
        $markerCount = @($markerFiles).Count
        $hardSyntaxCount = @($syntaxErrors | Where-Object { $_ -notlike "*TOML parser check*" }).Count

        $pass = ($mergeExit -eq 0) -and ($unmergedCount -eq 0) -and ($markerCount -eq 0) -and ($hardSyntaxCount -eq 0)

        return [PSCustomObject]@{
            Fixture       = $fixtureName
            MergeExit     = $mergeExit
            ResolveExit   = $resolveExit
            UnmergedCount = $unmergedCount
            MarkerCount   = $markerCount
            SyntaxErrors  = $syntaxErrors
            Passed        = $pass
            MergeOutput   = $mergeOutput
            ResolveOutput = if ($resolveOutput) { $resolveOutput } else { "" }
        }
    }
    finally {
        Pop-Location
    }
}

if (-not (Test-Path -Path $TestsRoot -PathType Container)) {
    throw "Tests root not found: $TestsRoot"
}
$GitResolveBinary = Resolve-GitResolveBinary -Requested $GitResolveBinary

$fixtures = Get-ChildItem -Path $TestsRoot -Directory | Sort-Object Name
if ($fixtures.Count -eq 0) {
    throw "No fixture folders found in $TestsRoot"
}

$results = New-Object System.Collections.Generic.List[object]

foreach ($fixture in $fixtures) {
    $result = Test-Fixture -FixturePath $fixture.FullName -BinaryPath $GitResolveBinary -DoResolve:$ResolveRemaining -Strategy $ResolveStrategy
    $results.Add($result)

    if ($result.Passed) {
        Write-Host "PASS: $($result.Fixture)"
    }
    else {
        Write-Host "FAIL: $($result.Fixture)"
        Write-Host "  mergeExit=$($result.MergeExit) unmerged=$($result.UnmergedCount) markers=$($result.MarkerCount)"

        if ($result.SyntaxErrors.Count -gt 0) {
            Write-Host "  syntax findings:"
            foreach ($e in $result.SyntaxErrors) {
                Write-Host "   - $e"
            }
        }

        if ($StopOnFirstFailure) {
            break
        }
    }
}

$passCount = @($results | Where-Object Passed).Count
$failCount = $results.Count - $passCount

Write-Host ""
Write-Host "=== Summary ==="
Write-Host "Total: $($results.Count)"
Write-Host "Passed: $passCount"
Write-Host "Failed: $failCount"

if ($failCount -gt 0) {
    Write-Host ""
    Write-Host "Failed fixtures:"
    $results | Where-Object { -not $_.Passed } | ForEach-Object { Write-Host " - $($_.Fixture)" }
    exit 1
}

exit 0
