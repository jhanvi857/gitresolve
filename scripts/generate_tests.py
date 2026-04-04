import os
import subprocess
import shutil
import stat

ROOT = os.path.abspath(".")
TESTS_DIR = os.path.join(ROOT, "tests")

def remove_readonly(func, path, excinfo):
    os.chmod(path, stat.S_IWRITE)
    func(path)

if not os.path.exists(TESTS_DIR):
    os.makedirs(TESTS_DIR)

def run(cmd, cwd=None):
    subprocess.run(cmd, shell=True, check=True, cwd=cwd, capture_output=True)

def create_conflict(name, base, ours, theirs):
    p = os.path.join(TESTS_DIR, name)
    if os.path.exists(p):
        shutil.rmtree(p, onerror=remove_readonly)
    os.makedirs(p)
    
    run("git init -b main", cwd=p)
    run('git config user.email "t@e.com"', cwd=p)
    run('git config user.name "T"', cwd=p)
    run("git config core.autocrlf false", cwd=p)
    
    for f, content in base.items():
        fp = os.path.join(p, f)
        os.makedirs(os.path.dirname(fp), exist_ok=True)
        with open(fp, "w", encoding="utf-8") as file:
            file.write(content)
    run("git add .", cwd=p)
    run('git commit -m "base"', cwd=p)
    
    run("git checkout -b ours", cwd=p)
    for f, content in ours.items():
        fp = os.path.join(p, f)
        os.makedirs(os.path.dirname(fp), exist_ok=True)
        with open(fp, "w", encoding="utf-8") as file:
            file.write(content)
    run("git add .", cwd=p)
    run('git commit -m "ours"', cwd=p)
    
    run("git checkout main", cwd=p)
    run("git checkout -b theirs", cwd=p)
    for f, content in theirs.items():
        fp = os.path.join(p, f)
        os.makedirs(os.path.dirname(fp), exist_ok=True)
        with open(fp, "w", encoding="utf-8") as file:
            file.write(content)
    run("git add .", cwd=p)
    run('git commit -m "theirs"', cwd=p)
    
    run("git checkout ours", cwd=p)
    subprocess.run("git merge theirs -m 'm' --no-edit", shell=True, cwd=p, capture_output=True)
    
    if name == "E5_lock":
        with open(os.path.join(p, ".gitresolve.lock"), "w") as f:
            f.write("stale")

CASES = [
    ("E1_whitespace", {"m.go": "package m\n\nfunc m(){f.P(\"h\")}"}, {"m.go": "package m\n\nfunc m(){f.P(\"h\")   }"}, {"m.go": "package m\n\n\nfunc m(){f.P(\"h\")}"}),
    ("E2_identical", {"d.txt": "v1"}, {"d.txt": "v2"}, {"d.txt": "v2"}),
    ("E3_imports", {"a.go": "package a\nimport (\"c\")"}, {"a.go": "package a\nimport (\"c\"\n\"f\")"}, {"a.go": "package a\nimport (\"c\"\n\"f\")"}),
    ("E4_yaml", {"c.yaml": "a: 1\nb: 2"}, {"c.yaml": "a: 2\nb: 2"}, {"c.yaml": "a: 1\nb: 3"}),
    ("E5_lock", {"f.txt": "a"}, {"f.txt": "b"}, {"f.txt": "c"}),
    ("M1_json", {"d.json": '{"v":"1"}'}, {"d.json": '{"v":"2","o":"1"}'}, {"d.json": '{"v":"2","t":"1"}'}),
    ("M2_yaml", {"l.yaml": "i:\n- a\n- b"}, {"l.yaml": "i:\n- a\n- b2"}, {"l.yaml": "i:\n- a2\n- b"}),
    ("M3_ts", {"t.ts": "interface U {n:s}"}, {"t.ts": "interface U {n:s\nid:n}"}, {"t.ts": "interface U {n:s\ne:s}"}),
    ("M4_toml", {"c.toml": "[d]\nhost=\"l\""}, {"c.toml": "[d]\nhost=\"p\""}, {"c.toml": "[d]\nhost=\"l\"\ns=true"}),
    ("M5_go_mod", {"go.mod": "module t\nrequire (a v1)"}, {"go.mod": "module t\nrequire (a v1\nb v1)"}, {"go.mod": "module t\nrequire (a v1\nc v1)"}),
    ("H1_pkg", {"package.json": '{"v":"1.0","s":{"a":"b"}}'}, {"package.json": '{"v":"1.1","s":{"a":"b","t":"j"}}'}, {"package.json": '{"v":"1.0","s":{"a":"b","b":"t"}}'}),
    ("H2_del_mod", {"u.go": "package u\nfunc O(){}\nfunc K(){}"}, {"u.go": "package u\nfunc K(){}"}, {"u.go": "package u\nfunc O(){println(\"m\")}\nfunc K(){}"}),
    ("H3_security", {"p.json": '{"a":["/d"]}'}, {"p.json": '{"a":["/d","/e/p"]}'}, {"p.json": '{"a":["/d","/v/l"]}'}),
    ("H4_ctrl_c", {"l.txt": "1\n2\n3"}, {"l.txt": "1o\n2\n3"}, {"l.txt": "1\n2\n3t"}),
    ("H5_multi", {"f1": "a", "f2": "b"}, {"f1": "ao"}, {"f2": "bt"}),
    ("S1_ast", {"c.go": "package m\nfunc m(){if t{if t{p(\"b\")}}}}"}, {"c.go": "package m\nfunc m(){if t{if t{p(\"o\")}}}}"}, {"c.go": "package m\nfunc m(){if t{if t{p(\"t\")}}}}"}),
    ("S2_lock_contention", {"d.txt": "A"}, {"d.txt": "B"}, {"d.txt": "C"}),
    ("S3_db_migration", {"db/1.sql": "C1"}, {"db/2o.sql": "C2"}, {"db/2t.sql": "C3"}),
    ("S4_cargo", {"C.toml": "[f]\nd=[]"}, {"C.toml": "[f]\nd=[]\na=[1]"}, {"C.toml": "[f]\nd=[2]"}),
    ("S5_undo", {"s.txt": "p1"}, {"s.txt": "p2o"}, {"s.txt": "p2t"}),
]

for name, base, ours, theirs in CASES:
    print(f"Creating {name}...", flush=True)
    try:
        create_conflict(name, base, ours, theirs)
    except Exception as e:
        print(f"Failed {name}: {e}", flush=True)

print("Done.")
