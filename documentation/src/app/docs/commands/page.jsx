"use client";

import React from 'react';

const MacCodeWindow = ({ children, title }) => (
  <div className="code-window mt-6 border-subtle overflow-hidden shadow-2xl">
    <div className="code-header">
      <div className="dot dot-red" />
      <div className="dot dot-yellow" />
      <div className="dot dot-green" />
      {title && <span className="ml-4 text-[9px] text-gray-500 font-mono uppercase tracking-widest">{title}</span>}
    </div>
    <div className="code-content whitespace-pre overflow-x-auto text-[11px] leading-6 py-6">
      {children}
    </div>
  </div>
);

const CommandSection = ({ id, name, description, usecase, flags, example, details }) => (
  <div id={id} className="py-20 border-b border-[#111] last:border-0 group px-4 md:px-0 scroll-mt-20">
    <div className="max-w-4xl mx-auto">
      <div className="flex flex-col md:flex-row md:items-baseline justify-between gap-2 mb-6">
        <h2 className="text-2xl font-semibold tracking-tight text-white font-mono flex items-center gap-2">
          {name}
        </h2>
        <div className="text-[10px] text-gray-500 font-black px-2 py-0.5 rounded border border-[#222] bg-[#111] uppercase tracking-[0.2em] shrink-0">
          CLI Command
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-12 gap-12">
        <div className="lg:col-span-12">
          <p className="text-gray-400 text-base leading-relaxed mb-8 max-w-3xl">
            {description}
          </p>
        </div>

        <div className="lg:col-span-5 space-y-10">
          <div>
            <h4 className="text-[10px] font-black text-gray-500 uppercase tracking-widest mb-4">When to use</h4>
            <p className="text-[15px] font-medium text-white italic border-l-2 border-white/10 pl-5 py-2">
              "{usecase}"
            </p>
          </div>

          {flags && flags.length > 0 && (
            <div className="pt-2">
              <h4 className="text-[10px] font-black text-gray-500 uppercase tracking-widest mb-4">Flags & Modifiers</h4>
              <div className="flex flex-col gap-3">
                {flags.map((flag, idx) => (
                  <div key={idx} className="flex flex-col gap-1.5 p-3 rounded-lg border border-[#222] bg-[#050505] hover:border-[#444] transition-colors">
                    <code className="text-[12px] font-bold text-white font-mono">
                      {flag.name}
                    </code>
                    <div className="text-[11px] text-gray-500 font-medium leading-relaxed">
                      {flag.desc}
                    </div>
                  </div>
                ))}
              </div>
            </div>
          )}
        </div>

        <div className="lg:col-span-7 relative">
          <div className="mb-4 text-[11px] font-black uppercase tracking-widest text-[#444]">Example Terminal Output</div>
          <MacCodeWindow title="terminal — demo">
            {example}
          </MacCodeWindow>
          {details && (
            <div className="mt-8 p-6 rounded-xl bg-blue-500/[0.03] border border-blue-500/10 text-[13px] text-gray-400 leading-relaxed">
              <div className="text-white text-[11px] font-black mb-2 uppercase tracking-widest">Internal Details</div>
              {details}
            </div>
          )}
        </div>
      </div>
    </div>
  </div>
);

export default function CommandsPage() {
  const commands = [
    {
      id: "scan",
      name: "gitresolve scan",
      description: "Predict potential conflicts before performing a destructive merge. By simulating a merge using the git merge-tree engine, the scan command identifies files that will require manual intervention.",
      usecase: "Use this during your build/test phase to proactively catch merge issues on your feature branch before they hit the main environment.",
      flags: [
        { name: "--target <branch>", desc: "Specifies the target branch to simulate a merge against (default: main)." },
        { name: "--no-fetch", desc: "Ignore remote branch states and only use local references." },
        { name: "--json", desc: "Output resolution summary as a machine-readable JSON object." }
      ],
      example: (
        <>
          <span className="token-output">$ </span><span className="token-cmd">gitresolve</span> <span className="token-arg">scan</span> <span className="token-flag">--target</span> <span className="token-string">develop</span><br />
          <span className="token-comment mt-4 block"># Simulating merge develop...feature-A</span>
          <span className="token-severity-high mt-4 block">! Predictive Conflict Detected</span>
          <span className="token-output block">- src/core/auth.go: Logical Overlap (Score 7)</span>
          <span className="token-output block">- pkg/config.yaml: YAML Conflict (Score 2 — Resolvable)</span>
          <br />
          <span className="token-output block mt-4 text-red-500 font-bold underline decoration-red-500/20 underline-offset-4">Scan failed: 1 high-risk file found.</span>
        </>
      ),
      details: "The engine uses the --write-tree flag of git merge-tree to calculate the resulting tree without modifying your worktree. This makes it completely safe to run even during active development."
    },
    {
      id: "status",
      name: "gitresolve status",
      description: "Provides a diagnostic view of current conflicts in your working directory. It grades each conflict by complexity using its built-in scoring engine.",
      usecase: "Run this after a failed git merge to see which files can be automatically resolved and which require your manual attention.",
      flags: [
        { name: "--detailed", desc: "Shows specific line numbers for each conflict block." },
        { name: "--short", desc: "Gives a high-level count of conflicted files only." }
      ],
      example: (
        <>
          <span className="token-output">$ </span><span className="token-cmd">gitresolve</span> <span className="token-arg">status</span><br />
          <span className="token-output block mt-4 font-bold border-b border-[#222] pb-2 px-2">DIAGNOSTIC REPORT — 2 FILES</span>
          <span className="token-output mt-4 block px-2">SCORE  AUTO  LOC   FILE</span>
          <span className="block px-2 mt-2"><span className="token-number">2 </span>     YES   4     pkg/utils/data.go</span>
          <span className="block px-2"><span className="token-severity-high font-bold px-1 rounded-sm bg-red-900/30">9 </span>     NO    22    internal/db.sql</span>
          <span className="token-comment mt-6 block px-2">Legend: Score 1-4 (Resolvable), 5-10 (Manual Override Needed)</span>
        </>
      ),
      details: "Severity is calculated based on the file extension (structured vs text), the presence of AST overlaps, and the line count of the conflicted block."
    },
    {
      id: "merge",
      name: "gitresolve merge",
      description: "Applies automated resolution logic to any safe conflict patterns. This includes whitespace normalization, import deduplication, and structured data folding.",
      usecase: "Use this as your first step after a conflict. It will automatically clean up the trivial parts of the merge for you.",
      flags: [
        { name: "--dry-run", desc: "Shows the proposed changes in the terminal without writing to disk." },
        { name: "--auto-structured", desc: "Enables deep merging for JSON, YAML, and TOML files (on by default)." },
        { name: "--verify", desc: "Runs syntax validation checks after writing the resolved file." }
      ],
      example: (
        <>
          <span className="token-output">$ </span><span className="token-cmd">gitresolve</span> <span className="token-arg">merge</span> <span className="token-flag">--verify</span><br />
          <span className="token-comment mt-4 block"># Batch resolving safe patterns...</span>
          <span className="token-output block mt-4">--- main.go (imports resolved)</span>
          <span className="token-output block">--- config.json (deep-map merge resolved)</span>
          <span className="token-severity-low mt-4 block font-bold text-green-500">✓ Syntax Verification Passed: internal/app_test.go</span>
          <span className="token-output mt-4 block">Total files modified: 3</span>
        </>
      ),
      details: "The merge command utilizes POSIX atomic writes. A temporary file is written and then atomically swapped with the original worktree file using os.Rename to prevent partial file corruption."
    },
    {
      id: "resolve",
      name: "gitresolve resolve",
      description: "Launches the interactive conflict SDK for remaining high-risk files. This guided experience allows you to jump between conflict blocks and select strategies.",
      usecase: "Run this to fix complex logical conflicts that the 'merge' command flagged as high-risk or ambiguous.",
      flags: [
        { name: "--file <path>", desc: "Force interactive resolution on a specific file only." },
        { name: "--strategy <choice>", desc: "Pre-set a strategy (ours/theirs/both) to bypass the prompt." },
        { name: "--timeout <time>", desc: "Automatic fallback to 'theirs' if the prompt isn't answered in time (useful for automation)." }
      ],
      example: (
        <>
          <span className="token-output">$ </span><span className="token-cmd">gitresolve</span> <span className="token-arg">resolve</span> <span className="token-flag">--file</span> <span className="token-string">core/system.go</span><br />
          <span className="token-output mt-4 block underline">INTERACTIVE BLOCK 1 of 4</span>
          <span className="token-variable mt-2 block">Found duplicate function definition.</span>
          <span className="token-output mt-2 block text-white/40">L122: func (s *Sys) Start() { }</span>
          <span className="token-output block">? Action: [O]urs / [T]heirs / [B]oth / [E]dit ? <span className="animate-pulse shadow-sm shadow-white">_</span></span>
        </>
      ),
      details: "The resolution engine tracks your choices in a local SQLite database, allowing you to resume interrupted resolution sessions later without losing your work."
    },
    {
      id: "blame",
      name: "gitresolve blame",
      description: "Audit the history of conflict resolutions in your repository. It provides patterns and insights into which files are frequently conflicting.",
      usecase: "Use this during team post-mortems to identify files that are becoming 'hotspots' for developer friction.",
      flags: [
        { name: "--patterns", desc: "Aggregates history into a statistical report of common conflict types." },
        { name: "--limit <N>", desc: "Only show the N most recent resolutions." }
      ],
      example: (
        <>
          <span className="token-output">$ </span><span className="token-cmd">gitresolve</span> <span className="token-arg">blame</span> <span className="token-flag">--patterns</span><br />
          <span className="token-output block mt-4 font-bold border-b border-[#222] pb-2">CONFLICT PATTERN ANALYSIS — LAST 30 DAYS</span>
          <span className="token-output mt-4 block">1.  <span className="token-variable font-bold">ImportConflict</span>  (42%) — Resolved via Merge</span>
          <span className="token-output block">2.  <span className="token-severity-high font-bold">LogicOverlap</span>    (26%) — Resolved via Interactive</span>
          <span className="token-output block">3.  <span className="token-string font-bold">SchemaShift</span>     (18%) — Resolved via Structured</span>
        </>
      ),
      details: "Blame data is derived from the .gitresolve/history.db file. This file is local and never pushed to remote servers."
    },
    {
      id: "undo",
      name: "gitresolve undo",
      description: "Roll back your repository to the exact state it was in before a gitresolve operation. This is safer than manual git resets as it preserves session history.",
      usecase: "If an automated resolution resulted in unexpected behavior, run undo to restore the conflicting markers perfectly.",
      flags: [
        { name: "--steps <N>", desc: "Number of operations to roll back in time (default: 1)." },
        { name: "--force", desc: "Ignore uncommitted worktree changes and force the rollback." }
      ],
      example: (
        <>
          <span className="token-output">$ </span><span className="token-cmd">gitresolve</span> <span className="token-arg">undo</span> <span className="token-flag">--steps</span> <span className="token-number">1</span><br />
          <span className="token-comment mt-4 block"># Fetching previous snapshot: eec91a1...</span>
          <span className="token-output block mt-2 text-green-500 font-semibold">Done: Restored markers in 3 files. State is now: CONFLICTING.</span>
        </>
      ),
      details: "The undo command verifies the SHA signature of your worktree before proceeding to ensure that manual edits aren't silently overwritten."
    }
  ];

  return (
    <div className="pb-32">
      <header className="pt-8 pb-16 border-b border-[#222]">
        <div className="max-w-4xl mx-auto px-4 md:px-0">
          <div className="text-white text-[10px] uppercase tracking-[0.4em] font-bold mb-6 flex items-center gap-2">
            <span className="w-1.5 h-1.5 rounded-full bg-blue-500" />
            CLI SDK v1.2 Reference
          </div>
          <h1 className="text-4xl font-semibold tracking-tighter text-white mb-6">
            Command Reference
          </h1>
          <p className="text-[17px] text-gray-500 max-w-2xl leading-relaxed">
            The gitresolve CLI provides a complete sdk for predicting, debugging, and resolving version control conflicts. Each command is designed to be highly secure and deterministic.
          </p>
        </div>
      </header>

      <div>
        {commands.map((cmd, idx) => (
          <CommandSection key={idx} {...cmd} />
        ))}
      </div>

      <div className="max-w-4xl mx-auto px-4 md:px-0 mt-32">
        {/* <div className="p-12 rounded-2xl border border-[#222] bg-gradient-to-br from-[#050505] to-black flex flex-col md:flex-row items-center justify-between gap-12 text-center md:text-left">
          <div className="max-w-md">
            <h3 className="text-xl font-semibold text-white mb-3 tracking-tight">Enterprise Scaling</h3>
            <p className="text-gray-500 text-sm leading-relaxed">Need help integrating gitresolve into a large-scale CI/CD pipeline or air-gapped infrastructure? Our engineering team provides specialized deployment guides.</p>
          </div>
          <button className="px-8 py-3 rounded-lg bg-white text-black font-bold text-xs hover:bg-gray-200 transition-all uppercase tracking-widest shrink-0 shadow-lg">
            Request Enterprise Guide
          </button>
        </div> */}
      </div>
    </div>
  );
}
