# scripts/create_tests.ps1
$t = (Join-Path (Resolve-Path .).Path "tests")
if (-not (Test-Path $t)) { New-Item $t -ItemType Directory -Force }

# Helper to avoid repetitive code
function Mk($n) {
    $p = Join-Path $t $n
    if (Test-Path $p) { Remove-Item $p -Recurse -Force }
    New-Item $p -ItemType Directory -Force | Out-Null
    Push-Location $p
    git init -b main -q
    git config user.email "t@e.com"; git config user.name "T"; git config core.autocrlf false
}

# --- Level 1 ---
Mk "E1_whitespace"
"package m`nfunc m(){f.P(\"h\")}" | Out-File m.go -Encoding utf8; git add .; git commit -m "b" -q
git checkout -b ours -q; "package m`nfunc m(){f.P(\"h\")   }" | Out-File m.go -Encoding utf8; git add .; git commit -m "o" -q
git checkout main -q; git checkout -b theirs -q; "package m`n`nfunc m(){f.P(\"h\")}" | Out-File m.go -Encoding utf8; git add .; git commit -m "t" -q
git checkout ours -q; git merge theirs 2>$null; Pop-Location

Mk "E2_identical"
"v1" | Out-File d.txt; git add .; git commit -m "b" -q
git checkout -b ours -q; "v2" | Out-File d.txt; git add .; git commit -m "o" -q
git checkout main -q; git checkout -b theirs -q; "v2" | Out-File d.txt; git add .; git commit -m "t" -q
git checkout ours -q; git merge theirs 2>$null; Pop-Location

Mk "E3_imports"
"package a`nimport (\"c\")" | Out-File a.go -Encoding utf8; git add .; git commit -m "b" -q
git checkout -b ours -q; "package a`nimport (\"c\"`n\"f\")" | Out-File a.go -Encoding utf8; git add .; git commit -m "o" -q
git checkout main -q; git checkout -b theirs -q; "package a`nimport (\"c\"`n\"f\")" | Out-File a.go -Encoding utf8; git add .; git commit -m "t" -q
git checkout ours -q; git merge theirs 2>$null; Pop-Location

Mk "E4_yaml_scalar"
"a: 1`nb: 2" | Out-File c.yaml; git add .; git commit -m "b" -q
git checkout -b ours -q; "a: 2`nb: 2" | Out-File c.yaml; git add .; git commit -m "o" -q
git checkout main -q; git checkout -b theirs -q; "a: 1`nb: 3" | Out-File c.yaml; git add .; git commit -m "t" -q
git checkout ours -q; git merge theirs 2>$null; Pop-Location

Mk "E5_lock_recovery"
"a" | Out-File f.txt; git add .; git commit -m "b" -q
git checkout -b ours -q; "b" | Out-File f.txt; git add .; git commit -m "o" -q
git checkout main -q; git checkout -b theirs -q; "c" | Out-File f.txt; git add .; git commit -m "t" -q
git checkout ours -q; git merge theirs 2>$null; "stale" | Out-File .gitresolve.lock; Pop-Location

# --- Level 2 ---
Mk "M1_json_merge"
'{"u":{"n":"j","s":{"t":"l"}}}' | Out-File d.json; git add .; git commit -m "b" -q
git checkout -b ours -q; '{"u":{"n":"j","s":{"t":"d","f":"i"}}}' | Out-File d.json; git add .; git commit -m "o" -q
git checkout main -q; git checkout -b theirs -q; '{"u":{"n":"j","s":{"t":"l","n":true}}}' | Out-File d.json; git add .; git commit -m "t" -q
git checkout ours -q; git merge theirs 2>$null; Pop-Location

Mk "M2_yaml_array"
"i:`n- a`n- b" | Out-File l.yaml; git add .; git commit -m "b" -q
git checkout -b ours -q; "i:`n- a`n- b2" | Out-File l.yaml; git add .; git commit -m "o" -q
git checkout main -q; git checkout -b theirs -q; "i:`n- a2`n- b" | Out-File l.yaml; git add .; git commit -m "t" -q
git checkout ours -q; git merge theirs 2>$null; Pop-Location

Mk "M3_ts_interface"
"interface U {n:s}" | Out-File t.ts; git add .; git commit -m "b" -q
git checkout -b ours -q; "interface U {n:s`nid:n}" | Out-File t.ts; git add .; git commit -m "o" -q
git checkout main -q; git checkout -b theirs -q; "interface U {n:s`ne:s}" | Out-File t.ts; git add .; git commit -m "t" -q
git checkout ours -q; git merge theirs 2>$null; Pop-Location

Mk "M4_toml_merge"
"[d]`nhost=\"l\"" | Out-File c.toml; git add .; git commit -m "b" -q
git checkout -b ours -q; "[d]`nhost=\"p\"`nssl=true" | Out-File c.toml; git add .; git commit -m "o" -q
git checkout main -q; git checkout -b theirs -q; "[d]`nhost=\"l\"`n[d.m]`ne=true" | Out-File c.toml; git add .; git commit -m "t" -q
git checkout ours -q; git merge theirs 2>$null; Pop-Location

Mk "M5_go_mod"
"module t`nrequire (a v1)" | Out-File go.mod; git add .; git commit -m "b" -q
git checkout -b ours -q; "module t`nrequire (a v1`nb v1)" | Out-File go.mod; git add .; git commit -m "o" -q
git checkout main -q; git checkout -b theirs -q; "module t`nrequire (a v1`nc v1)" | Out-File go.mod; git add .; git commit -m "t" -q
git checkout ours -q; git merge theirs 2>$null; Pop-Location

# --- Level 3 ---
Mk "H1_pkg_scripts"
'{"v":"1.0","s":{"a":"b"}}' | Out-File package.json; git add .; git commit -m "b" -q
git checkout -b ours -q; '{"v":"1.1","s":{"a":"b","t":"j"}}' | Out-File package.json; git add .; git commit -m "o" -q
git checkout main -q; git checkout -b theirs -q; '{"v":"1.0","s":{"a":"b","b":"t"}}' | Out-File package.json; git add .; git commit -m "t" -q
git checkout ours -q; git merge theirs 2>$null; Pop-Location

Mk "H2_del_mod"
"package u`nfunc O(){}`nfunc K(){}" | Out-File u.go -Encoding utf8; git add .; git commit -m "b" -q
git checkout -b ours -q; "package u`nfunc K(){}" | Out-File u.go -Encoding utf8; git add .; git commit -m "o" -q
git checkout main -q; git checkout -b theirs -q; "package u`nfunc O(){println(\"m\")}`nfunc K(){}" | Out-File u.go -Encoding utf8; git add .; git commit -m "t" -q
git checkout ours -q; git merge theirs 2>$null; Pop-Location

Mk "H3_security_path"
'{"a":["/d"]}' | Out-File p.json; git add .; git commit -m "b" -q
git checkout -b ours -q; '{"a":["/d","/e/p"]}' | Out-File p.json; git add .; git commit -m "o" -q
git checkout main -q; git checkout -b theirs -q; '{"a":["/d","/v/l"]}' | Out-File p.json; git add .; git commit -m "t" -q
git checkout ours -q; git merge theirs 2>$null; Pop-Location

Mk "H4_ctrl_c"
"1`n2`n3" | Out-File l.txt; git add .; git commit -m "b" -q
git checkout -b ours -q; "1o`n2`n3" | Out-File l.txt; git add .; git commit -m "o" -q
git checkout main -q; git checkout -b theirs -q; "1`n2`n3t" | Out-File l.txt; git add .; git commit -m "t" -q
git checkout ours -q; git merge theirs 2>$null; Pop-Location

Mk "H5_multi"
"a" | Out-File f1; "b" | Out-File f2; git add .; git commit -m "b" -q
git checkout -b ours -q; "ao" | Out-File f1; "bo" | Out-File f2; git add .; git commit -m "o" -q
git checkout main -q; git checkout -b theirs -q; "at" | Out-File f1; "bt" | Out-File f2; git add .; git commit -m "t" -q
git checkout ours -q; git merge theirs 2>$null; Pop-Location

# --- Level 4 ---
Mk "S1_ast_fail"
"package m`nfunc m(){if t{if t{p(\"b\")}}}}" | Out-File c.go -Encoding utf8; git add .; git commit -m "b" -q
git checkout -b ours -q; "package m`nfunc m(){if t{if t{p(\"o\")}}}}" | Out-File c.go -Encoding utf8; git add .; git commit -m "o" -q
git checkout main -q; git checkout -b theirs -q; "package m`nfunc m(){if t{if t{p(\"t\")}}}}" | Out-File c.go -Encoding utf8; git add .; git commit -m "t" -q
git checkout ours -q; git merge theirs 2>$null; Pop-Location

Mk "S2_lock_contention"
"A" | Out-File d.txt; git add .; git commit -m "b" -q
git checkout -b ours -q; "B" | Out-File d.txt; git add .; git commit -m "o" -q
git checkout main -q; git checkout -b theirs -q; "C" | Out-File d.txt; git add .; git commit -m "t" -q
git checkout ours -q; git merge theirs 2>$null; Pop-Location

Mk "S3_db_migration"
"C1" | Out-File 1.sql; git add .; git commit -m "b" -q
git checkout -b ours -q; "C2" | Out-File 2o.sql; git add .; git commit -m "o" -q
git checkout main -q; git checkout -b theirs -q; "C3" | Out-File 2t.sql; git add .; git commit -m "t" -q
# This might not conflict on content, but it's okay for testing multiple files or naming.
Pop-Location

Mk "S4_cargo_feature"
"[f]`nd=[]" | Out-File C.toml; git add .; git commit -m "b" -q
git checkout -b ours -q; "[f]`nd=[]`na=[1]" | Out-File C.toml; git add .; git commit -m "o" -q
git checkout main -q; git checkout -b theirs -q; "[f]`nd=[2]" | Out-File C.toml; git add .; git commit -m "t" -q
git checkout ours -q; git merge theirs 2>$null; Pop-Location

Mk "S5_undo_integrity"
"p1" | Out-File s.txt; git add .; git commit -m "b" -q
git checkout -b ours -q; "p2o" | Out-File s.txt; git add .; git commit -m "o" -q
git checkout main -q; git checkout -b theirs -q; "p2t" | Out-File s.txt; git add .; git commit -m "t" -q
git checkout ours -q; git merge theirs 2>$null; Pop-Location

Write-Host "All tests created."
