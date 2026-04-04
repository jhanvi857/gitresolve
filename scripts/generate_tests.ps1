# scripts/generate_tests.ps1
param([string]$TargetDir = (Join-Path $PSScriptRoot "..\tests"))
if (-not (Test-Path $TargetDir)) { New-Item -ItemType Directory -Path $TargetDir -Force | Out-Null }

function Create-Conflict {
    param($Name, $Base, $Ours, $Theirs, $Extra)
    $path = Join-Path $TargetDir $Name
    if (Test-Path $path) { Remove-Item -Recurse -Force $path }
    New-Item -ItemType Directory -Path $path -Force | Out-Null
    Push-Location $path
    try {
        git init -q
        git config user.email "test@example.com"
        git config user.name "Test User"
        git config commit.gpgSign false
        "init" | Out-File .gitignore -Encoding utf8
        git add .
        git commit -m "init" -q
        git checkout -b base -q
        $Base.GetEnumerator() | % { 
            $p = Split-Path $_.Key -Parent
            if ($p) { New-Item -Directory $p -Force | Out-Null }
            $_.Value | Out-File $_.Key -Encoding utf8
        }
        git add .
        git commit -m "base" -q
        git checkout -b ours -q
        $Ours.GetEnumerator() | % { 
            $p = Split-Path $_.Key -Parent
            if ($p) { New-Item -Directory $p -Force | Out-Null }
            $_.Value | Out-File $_.Key -Encoding utf8 
        }
        git add .
        git commit -m "ours" -q
        git checkout base -q
        git checkout -b theirs -q
        $Theirs.GetEnumerator() | % { 
            $p = Split-Path $_.Key -Parent
            if ($p) { New-Item -Directory $p -Force | Out-Null }
            $_.Value | Out-File $_.Key -Encoding utf8 
        }
        git add .
        git commit -m "theirs" -q
        if ($Extra) { $Extra | % { Invoke-Expression $_; git add .; git commit -m "extra" -q } }
        git checkout ours -q
        git merge theirs 2>$null
        Write-Host "Done: $Name"
    } finally { Pop-Location }
}

Write-Host "Level 1"
Create-Conflict "E1_whitespace" @{"main.go"="package main`n`nfunc main() {`n    fmt.Println(\"hello\")`n}"} @{"main.go"="package main`n`nfunc main() {`n    fmt.Println(\"hello\")    `n}"} @{"main.go"="package main`n`nfunc main() {`n`n    fmt.Println(\"hello\")`n}"}
Create-Conflict "E2_identical" @{"data.txt"="v1`n"} @{"data.txt"="v2`n"} @{"data.txt"="v2`n"}
Create-Conflict "E3_imports" @{"api.go"="package api`nimport (\"context\")"} @{"api.go"="package api`nimport (\"context\"`n\"fmt\")"} @{"api.go"="package api`nimport (\"context\"`n\"fmt\")"}
Create-Conflict "E4_yaml_scalar" @{"config.yaml"="a: 1`nb: 2"} @{"config.yaml"="a: 2`nb: 2"} @{"config.yaml"="a: 1`nb: 3"}
Create-Conflict "E5_lock_recovery" @{"f.txt"="a"} @{"f.txt"="b"} @{"f.txt"="c"} @("New-Item .gitresolve.lock -ItemType File -Value 'stale'")

Write-Host "Level 2"
$m1_b = '{"u":{"n":"J","s":{"t":"l"}}}'; $m1_o = '{"u":{"n":"J","s":{"t":"d","f":"i"}}}'; $m1_t = '{"u":{"n":"J","s":{"t":"l","n":true}}}'
Create-Conflict "M1_json_merge" @{"d.json"=$m1_b} @{"d.json"=$m1_o} @{"d.json"=$m1_t}
Create-Conflict "M2_yaml_array" @{"l.yaml"="i:`n - a`n - b"} @{"l.yaml"="i:`n - a`n - b2"} @{"l.yaml"="i:`n - a2`n - b"}
Create-Conflict "M3_ts_interface" @{"t.ts"="interface U {n:s}"} @{"t.ts"="interface U {n:s`nid:n}"} @{"t.ts"="interface U {n:s`ne:s}"}
Create-Conflict "M4_toml_merge" @{"c.toml"="[d]`nhost=\"l\""} @{"c.toml"="[d]`nhost=\"p\"`nssl=true"} @{"c.toml"="[d]`nhost=\"l\"`n[d.m]`ne=true"}
Create-Conflict "M5_go_mod_conflict" @{"go.mod"="module t`ngo 1.21`nrequire (a v1)"} @{"go.mod"="module t`ngo 1.21`nrequire (a v1`nb v2)"} @{"go.mod"="module t`ngo 1.21`nrequire (a v1`nc v3)"}

Write-Host "Level 3"
Create-Conflict "H1_pkg_scripts" @{"package.json"='{"v":"1.0","s":{"a":"b"}}'} @{"package.json"='{"v":"1.1","s":{"a":"b","t":"j"}}'} @{"package.json"='{"v":"1.0","s":{"a":"b","b":"t"}}'}
Create-Conflict "H2_del_mod" @{"u.go"="package u`nfunc O(){}`nfunc K(){}"} @{"u.go"="package u`nfunc K(){}"} @{"u.go"="package u`nfunc O(){println(\"m\")}`nfunc K(){}"}
Create-Conflict "H3_security_path" @{"p.json"='{"a":["/d"]}'} @{"p.json"='{"a":["/d","/e/p"]}'} @{"p.json"='{"a":["/d","/v/l"]}'}
Create-Conflict "H4_ctrl_c_cleanup" @{"l.txt"="1`n2`n3"} @{"l.txt"="1o`n2`n3"} @{"l.txt"="1`n2`n3t"}
Create-Conflict "H5_multi_file" @{"f1"="a";"f2"="b"} @{"f1"="ao";"f2"="bo"} @{"f1"="at";"f2"="bt"}

Write-Host "Level 4"
Create-Conflict "S1_ast_parse_fail" @{"c.go"="package m`nfunc m(){if t{if t{p(\"b\")}}}}"} @{"c.go"="package m`nfunc m(){if t{if t{p(\"o\")}}}}"} @{"c.go"="package m`nfunc m(){if t{if t{p(\"t\")}}}}"}
Create-Conflict "S2_lock_contention" @{"d.txt"="A"} @{"d.txt"="B"} @{"d.txt"="C"}
Create-Conflict "S3_db_migration" @{"db/1.sql"="C1"} @{"db/2o.sql"="C2"} @{"db/2t.sql"="C3"}
Create-Conflict "S4_cargo_toml_flag" @{"C.toml"="[f]`nd=[]"} @{"C.toml"="[f]`nd=[]`na=[1]"} @{"C.toml"="[f]`nd=[2]"}
Create-Conflict "S5_undo_integrity" @{"s.txt"="p1"} @{"s.txt"="p2o"} @{"s.txt"="p2t"}
