# scripts/generate_tests_final.ps1
$t = "c:\Users\family\OneDrive\Desktop\gitresolve\tests"
if (-not (Test-Path $t)) { New-Item $t -ItemType Directory }

$cases = @(
    # L1
    @{ Name="E1_whitespace"; Base=@{"m.go"="package m`nfunc m(){f.P(\"h\")}"}; Ours=@{"m.go"="package m`nfunc m(){f.P(\"h\")   }"}; Theirs=@{"m.go"="package m`n`nfunc m(){f.P(\"h\")}"} },
    @{ Name="E2_identical"; Base=@{"d.txt"="v1"}; Ours=@{"d.txt"="v2"}; Theirs=@{"d.txt"="v2"} },
    @{ Name="E3_imports"; Base=@{"a.go"="package a`nimport (\"c\")"}; Ours=@{"a.go"="package a`nimport (\"c\"`n\"f\")"}; Theirs=@{"a.go"="package a`nimport (\"c\"`n\"f\")"} },
    @{ Name="E4_yaml"; Base=@{"c.yaml"="a: 1`nb: 2"}; Ours=@{"c.yaml"="a: 2`nb: 2"}; Theirs=@{"c.yaml"="a: 1`nb: 3"} },
    @{ Name="E5_lock"; Base=@{"f.txt"="a"}; Ours=@{"f.txt"="b"}; Theirs=@{"f.txt"="c"} },
    # L2
    @{ Name="M1_json"; Base=@{"d.json"='{"v":"1"}'}; Ours=@{"d.json"='{"v":"2","o":"1"}'}; Theirs=@{"d.json"='{"v":"2","t":"1"}'} },
    @{ Name="M2_yaml"; Base=@{"l.yaml"="i:`n- a`n- b"}; Ours=@{"l.yaml"="i:`n- a`n- b2"}; Theirs=@{"l.yaml"="i:`n- a2`n- b"} },
    @{ Name="M3_ts"; Base=@{"t.ts"="interface U {n:s}"}; Ours=@{"t.ts"="interface U {n:s;id:n}"}; Theirs=@{"t.ts"="interface U {n:s;e:s}"} },
    @{ Name="M4_toml"; Base=@{"c.toml"="[d]`nhost=\"l\""}; Ours=@{"c.toml"="[d]`nhost=\"p\""}; Theirs=@{"c.toml"="[d]`nhost=\"l\"`ns=true"} },
    @{ Name="M5_go_mod"; Base=@{"go.mod"="module t`nrequire (a v1)"}; Ours=@{"go.mod"="module t`nrequire (a v1`nb v1)"}; Theirs=@{"go.mod"="module t`nrequire (a v1`nc v1)"} },
    # L3
    @{ Name="H1_pkg"; Base=@{"package.json"='{"v":"1.0","s":{"a":"b"}}'}; Ours=@{"package.json"='{"v":"1.1","s":{"a":"b","t":"j"}}'}; Theirs=@{"package.json"='{"v":"1.0","s":{"a":"b","b":"t"}}'} },
    @{ Name="H2_del_mod"; Base=@{"u.go"="package u`nfunc O(){}`nfunc K(){}"}; Ours=@{"u.go"="package u`nfunc K(){}"}; Theirs=@{"u.go"="package u`nfunc O(){println(\"m\")}`nfunc K(){}"} },
    @{ Name="H3_security"; Base=@{"p.json"='{"a":["/d"]}'}; Ours=@{"p.json"='{"a":["/d","/e/p"]}'}; Theirs=@{"p.json"='{"a":["/d","/v/l"]}'} },
    @{ Name="H4_ctrl_c"; Base=@{"l.txt"="1`n2`n3"}; Ours=@{"l.txt"="1o`n2`n3"}; Theirs=@{"l.txt"="1`n2`n3t"} },
    @{ Name="H5_multi"; Base=@{"f1"="a";"f2"="b"}; Ours=@{"f1"="ao"}; Theirs=@{"f2"="bt"} },
    # L4
    @{ Name="S1_ast"; Base=@{"c.go"="package m`nfunc m(){if t{if t{p(\"b\")}}}}"}; Ours=@{"c.go"="package m`nfunc m(){if t{if t{p(\"o\")}}}}"}; Theirs=@{"c.go"="package m`nfunc m(){if t{if t{p(\"t\")}}}}"} },
    @{ Name="S2_lock_contention"; Base=@{"d.txt"="A"}; Ours=@{"d.txt"="B"}; Theirs=@{"d.txt"="C"} },
    @{ Name="S3_db_migration"; Base=@{"db/1.sql"="C1"}; Ours=@{"db/2o.sql"="C2"}; Theirs=@{"db/2t.sql"="C3"} },
    @{ Name="S4_cargo"; Base=@{"C.toml"="[f]`nd=[]"}; Ours=@{"C.toml"="[f]`nd=[]`na=[1]"}; Theirs=@{"C.toml"="[f]`nd=[2]"} },
    @{ Name="S5_undo"; Base=@{"s.txt"="p1"}; Ours=@{"s.txt"="p2o"}; Theirs=@{"s.txt"="p2t"} }
)

foreach ($c in $cases) {
    Write-Host "Setting up $($c.Name)..."
    $p = Join-Path $t $c.Name
    if (Test-Path $p) { Remove-Item $p -Recurse -Force }
    New-Item $p -ItemType Directory | Out-Null
    Push-Location $p
    git init -b main -q
    git config user.email "t@e.com"; git config user.name "T"; git config core.autocrlf false
    
    $fset = $c.Base.Keys
    foreach ($f in $fset) { $c.Base[$f] | Out-File $f -Encoding utf8 }
    git add .; git commit -m "b" -q
    
    git checkout -b ours -q
    $fset = $c.Ours.Keys
    foreach ($f in $fset) { $c.Ours[$f] | Out-File $f -Encoding utf8 }
    git add .; git commit -m "o" -q
    
    git checkout main -q
    git checkout -b theirs -q
    $fset = $c.Theirs.Keys
    foreach ($f in $fset) { $c.Theirs[$f] | Out-File $f -Encoding utf8 }
    git add .; git commit -m "t" -q
    
    git checkout ours -q
    git merge theirs -m "m" --no-edit 2>$null | Out-Null
    
    if ($c.Name -eq "E5_lock") { "stale" | Out-File .gitresolve.lock }
    
    Pop-Location
}
Write-Host "Done setting up 20 tests."
