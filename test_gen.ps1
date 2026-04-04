# test_gen.ps1
$t = "c:\Users\family\OneDrive\Desktop\gitresolve\tests"
if (-not (Test-Path $t)) { New-Item $t -ItemType Directory }

function Create-Conflict($name, $base, $ours, $theirs, $files) {
    $p = Join-Path $t $name
    if (Test-Path $p) { Remove-Item $p -Recurse -Force }
    New-Item $p -ItemType Directory | Out-Null
    Push-Location $p
    git init -b main -q
    git config user.email "t@e.com"; git config user.name "T"; git config core.autocrlf false
    $files | ForEach-Object { $base[$_] | Out-File $_ -Encoding utf8 }
    git add .; git commit -m "b" -q
    git checkout -b ours -q
    $files | ForEach-Object { $ours[$_] | Out-File $_ -Encoding utf8 }
    git add .; git commit -m "o" -q
    git checkout main -q
    git checkout -b theirs -q
    $files | ForEach-Object { $theirs[$_] | Out-File $_ -Encoding utf8 }
    git add .; git commit -m "t" -q
    git checkout ours -q
    git merge theirs -m "m" --no-edit 2>$null | Out-Null
    Pop-Location
}

# Already did E1-E5 but let's re-run or just do M1+
# Level 2
$m1_b = @{"d.json"='{"v":"1"}'}; $m1_o = @{"d.json"='{"v":"2","o":"1"}'}; $m1_t = @{"d.json"='{"v":"2","t":"1"}'}
Create-Conflict "M1_json" $m1_b $m1_o $m1_t @("d.json")

$m2_b = @{"l.yaml"="i:`n- a`n- b"}; $m2_o = @{"l.yaml"="i:`n- a`n- b2"}; $m2_t = @{"l.yaml"="i:`n- a2`n- b"}
Create-Conflict "M2_yaml" $m2_b $m2_o $m2_t @("l.yaml")

$m3_b = @{"t.ts"="interface U {n:s}"}; $m3_o = @{"t.ts"="interface U {n:s;id:n}"}; $m3_t = @{"t.ts"="interface U {n:s;e:s}"}
Create-Conflict "M3_ts" $m3_b $m3_o $m3_t @("t.ts")

# Level 3
$h2_b = @{"u.go"="package u`nfunc O(){}`nfunc K(){}"}; $h2_o = @{"u.go"="package u`nfunc K(){}"}; $h2_t = @{"u.go"="package u`nfunc O(){println(\"m\")}`nfunc K(){}"}
Create-Conflict "H2_del_mod" $h2_b $h2_o $h2_t @("u.go")

# Adding remaining ones quickly
$s4_b = @{"C.toml"="[f]`nd=[]"}; $s4_o = @{"C.toml"="[f]`nd=[]`na=[1]"}; $s4_t = @{"C.toml"="[f]`nd=[2]"}
Create-Conflict "S4_cargo" $s4_b $s4_o $s4_t @("C.toml")
