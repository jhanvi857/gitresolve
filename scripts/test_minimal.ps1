# scripts/test_minimal.ps1
$root = (Resolve-Path .).Path
$tests = Join-Path $root "tests"
if (-not (Test-Path $tests)) { New-Item $tests -ItemType Directory -Force }

$name = "E1_minimal"
$p = Join-Path $tests $name
if (Test-Path $p) { Remove-Item $p -Recurse -Force }
New-Item $p -ItemType Directory -Force | Out-Null
Push-Location $p
git init -b main -q
git config user.email "test@example.com"
git config user.name "Test User"
git config core.autocrlf false
"base" | Out-File f1 -Encoding utf8
git add .
git commit -m "base" -q
git checkout -b ours -q
"ours" | Out-File f1 -Encoding utf8
git add .
git commit -m "ours" -q
git checkout main -q
git checkout -b theirs -q
"theirs" | Out-File f1 -Encoding utf8
git add .
git commit -m "theirs" -q
git checkout ours -q
git merge theirs 2>$null
Pop-Location
Write-Host "Minimal Done"
