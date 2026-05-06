"use client";

import React from 'react';
import DocsShell from '@/components/DocsShell';

export default function CommandsReference() {
  return (
    <DocsShell 
      title="Commands Reference" 
      subtitle="Complete documentation for the gitresolve CLI tool."
    >
      <div className="space-y-32">
        {/* SCAN COMMAND */}
        <section id="scan" className="scroll-mt-24">
          <div className="flex flex-col md:flex-row md:items-center gap-4 mb-8">
            <div className="flex items-center gap-3">
              <span className="w-8 h-8 rounded-lg bg-blue-600/10 border border-blue-500/20 flex items-center justify-center text-blue-500 font-bold text-[10px]">CMD</span>
              <code className="text-xl font-bold text-white bg-[#0a0a0a] px-3 py-1 rounded border border-[#222]">gitresolve scan</code>
            </div>
            <span className="text-gray-500 font-mono text-sm leading-none">Predictive Conflict Detection</span>
          </div>
          
          <div className="prose-layout">
            <p className="mb-6 text-gray-400">
              The <code className="text-white bg-[#111] px-1.5 py-0.5 rounded">scan</code> command is your early-warning system. By leveraging <code className="text-gray-300">git merge-tree</code>, it simulates a merge between your current HEAD and a target branch to find potential overlaps before you even run a merge command.
            </p>
            
            <div className="code-window mb-8 border border-[#222]">
              <div className="code-header border-b border-[#111] px-4 py-2 flex gap-2">
                <div className="w-1.5 h-1.5 rounded-full bg-[#333]" /><div className="w-1.5 h-1.5 rounded-full bg-[#333]" /><div className="w-1.5 h-1.5 rounded-full bg-[#333]" />
              </div>
              <div className="code-content bg-black p-5 font-mono text-[13px] leading-relaxed">
                <span className="text-blue-500">$</span> gitresolve scan --target main<br />
                <span className="text-gray-500">Scanning for potential conflicts against main...</span><br />
                <span className="text-white">Potential conflicts detected: 2 files/blocks.</span><br />
                <span className="text-gray-500">Conflict hints:</span><br />
                <span className="text-yellow-500"> - Conflict (modify/modify): internal/db/schema.go</span><br />
                <span className="text-yellow-500"> - Conflict (add/add): config/default.yaml</span>
              </div>
            </div>

            <div className="grid grid-cols-1 gap-4 mt-8">
              <FlagItem flag="--target <branch>" desc="The target branch to compare against (defaults to main)." />
            </div>
          </div>
        </section>

        {/* STATUS COMMAND */}
        <section id="status" className="scroll-mt-24">
          <div className="flex flex-col md:flex-row md:items-center gap-4 mb-8">
            <div className="flex items-center gap-3">
              <span className="w-8 h-8 rounded-lg bg-green-600/10 border border-green-500/20 flex items-center justify-center text-green-500 font-bold text-[10px]">CMD</span>
              <code className="text-xl font-bold text-white bg-[#0a0a0a] px-3 py-1 rounded border border-[#222]">gitresolve status</code>
            </div>
            <span className="text-gray-500 font-mono text-sm leading-none">Real-time Conflict Indexing</span>
          </div>
          
          <div className="prose-layout">
            <p className="mb-6 text-gray-400">
              Displays the current unmerged blocks, categorized by their <strong>Severity Score</strong> and auto-resolution eligibility. Unlike standard <code className="text-gray-300">git status</code>, this looks inside files to identify logical overlap patterns.
            </p>

            <div className="code-window mb-8 border border-[#222]">
              <div className="code-header border-b border-[#111] px-4 py-2 flex gap-2">
                <div className="w-1.5 h-1.5 rounded-full bg-[#333]" /><div className="w-1.5 h-1.5 rounded-full bg-[#333]" /><div className="w-1.5 h-1.5 rounded-full bg-[#333]" />
              </div>
              <div className="code-content bg-black p-5 font-mono text-[13px] leading-relaxed">
                <span className="text-blue-500">$</span> gitresolve status<br />
                <span className="text-gray-500 block mb-2">Conflict Status Check:</span>
                <span className="text-gray-600">  SCORE  TYPE            AUTO  FILE</span><br />
                <span className="text-blue-400">  5      whitespace      yes   internal/net/socket.go</span><br />
                <span className="text-green-400">  10     imports         yes   main.go</span><br />
                <span className="text-red-500">  99     logic-overlap   no    pkg/auth/handler.go</span><br />
                <br />
                <span className="text-white">Total conflict blocks: 3</span>
              </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mt-8">
               <div className="p-4 rounded-lg bg-[#0a0a0a] border border-[#222]">
                  <h4 className="text-white font-semibold text-xs mb-2 uppercase tracking-widest">Trivial (0-10)</h4>
                  <p className="text-[11px] text-gray-500">Safe for auto-merge. Whitespace or simple import deduplication.</p>
               </div>
               <div className="p-4 rounded-lg bg-[#0a0a0a] border border-[#222]">
                  <h4 className="text-white font-semibold text-xs mb-2 uppercase tracking-widest">Structural (11-50)</h4>
                  <p className="text-[11px] text-gray-500">Requires syntax verification gates or policy-based strategy.</p>
               </div>
               <div className="p-4 rounded-lg bg-[#0a0a0a] border border-[#222]">
                  <h4 className="text-white font-semibold text-xs mb-2 uppercase tracking-widest">High Risk (99+)</h4>
                  <p className="text-[11px] text-gray-500">Logical overlaps that default to manual review or strict policy escalation.</p>
               </div>
            </div>
          </div>
        </section>

        {/* MERGE COMMAND */}
        <section id="merge" className="scroll-mt-24">
          <div className="flex flex-col md:flex-row md:items-center gap-4 mb-8">
            <div className="flex items-center gap-3">
              <span className="w-8 h-8 rounded-lg bg-purple-600/10 border border-purple-500/20 flex items-center justify-center text-purple-500 font-bold text-[10px]">CMD</span>
              <code className="text-xl font-bold text-white bg-[#0a0a0a] px-3 py-1 rounded border border-[#222]">gitresolve merge</code>
            </div>
            <span className="text-gray-500 font-mono text-sm leading-none">Autonomous Smart Triage</span>
          </div>
          
          <div className="prose-layout">
            <p className="mb-6 text-gray-400">
              The <code className="text-white bg-[#111] px-1.5 py-0.5 rounded">merge</code> command performs non-interactive resolution. It only applies changes that the engine is confident in through deterministic rules or configured policy profiles.
            </p>

            <div className="code-window mb-8 border border-[#222]">
              <div className="code-header border-b border-[#111] px-4 py-2 flex gap-2">
                <div className="w-1.5 h-1.5 rounded-full bg-[#333]" /><div className="w-1.5 h-1.5 rounded-full bg-[#333]" /><div className="w-1.5 h-1.5 rounded-full bg-[#333]" />
              </div>
              <div className="code-content bg-black p-5 font-mono text-[13px] leading-relaxed">
                <span className="text-blue-500">$</span> gitresolve merge<br />
                <span className="text-gray-500">Scanning index. Found 2 unmerged conflicts...</span><br />
                <br />
                <span className="text-white block">--- Processing internal/net/socket.go ---</span>
                <span className="text-green-500"> {">"} Successfully auto-resolved 100% of conflicts and staged.</span><br />
                <br />
                <span className="text-white block">--- Processing pkg/auth/handler.go ---</span>
                <span className="text-yellow-500"> {">"} Escalating conflict [Severity 99] LogicOverlap</span><br />
                <br />
                <span className="text-gray-400 italic">Merge scan complete. 1 file resolved, 1 requires manual review.</span>
              </div>
            </div>

            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 mt-8">
              <FlagItem flag="--dry-run" desc="Preview auto-resolutions without modifying any files." />
              <FlagItem flag="--policy-profile <profile>" desc="Apply a risk posture (auto/strict/balanced/aggressive)." />
              <FlagItem flag="--max-file-bytes <bytes>" desc="Skip files larger than this limit (default 10MB). Set to -1 for unlimited." />
              <FlagItem flag="--no-auto-structured" desc="Disable auto-resolution for JSON/YAML/TOML files." />
              <FlagItem flag="--shadow" desc="Simulate resolution and record hash diffs to logs." />
            </div>
          </div>
        </section>

        {/* RESOLVE COMMAND */}
        <section id="resolve" className="scroll-mt-24">
          <div className="flex flex-col md:flex-row md:items-center gap-4 mb-8">
            <div className="flex items-center gap-3">
              <span className="w-8 h-8 rounded-lg bg-orange-600/10 border border-orange-500/20 flex items-center justify-center text-orange-500 font-bold text-[10px]">CMD</span>
              <code className="text-xl font-bold text-white bg-[#0a0a0a] px-3 py-1 rounded border border-[#222]">gitresolve resolve</code>
            </div>
            <span className="text-gray-500 font-mono text-sm leading-none">Interactive Conflict Orchestration</span>
          </div>
          
          <div className="prose-layout">
            <p className="mb-6 text-gray-400">
              The primary interactive interface. Use <code className="text-white bg-[#111] px-1.5 py-0.5 rounded">resolve</code> to step through conflicted blocks. It combines auto-resolution for trivial blocks with an interactive prompt for logical conflicts.
            </p>

            <div className="code-window mb-8 border border-[#222]">
              <div className="code-header border-b border-[#111] px-4 py-2 flex gap-2">
                <div className="w-1.5 h-1.5 rounded-full bg-[#333]" /><div className="w-1.5 h-1.5 rounded-full bg-[#333]" /><div className="w-1.5 h-1.5 rounded-full bg-[#333]" />
              </div>
              <div className="code-content bg-black p-5 font-mono text-[13px] leading-relaxed">
                <span className="text-blue-500">$</span> gitresolve resolve<br />
                <span className="text-white">[Scalar] main.go (L21-25)</span><br />
                <span className="text-gray-500"> [O]urs:   // this comment is different</span><br />
                <span className="text-gray-500"> [T]heirs: // this is a comment</span><br />
                <span className="text-gray-400"> Options: [O]urs [T]heirs [B]oth [M]anual edit [S]kip</span><br />
                <br />
                <span className="text-white font-bold">Select [O/T/B/M/S]: _</span>
              </div>
            </div>

            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 mt-8">
              <FlagItem flag="--non-interactive" desc="Exit with status 1 if any conflict requires human input. Perfect for CI gates." />
              <FlagItem flag="--timeout <duration>" desc="Auto-select 'theirs' after a timeout (e.g., 30s) during interactive prompts." />
              <FlagItem flag="--strategy <type>" desc="Force a fixed strategy (ours/theirs/both) for all automatable conflicts." />
              <FlagItem flag="--enforce-gates" desc="Apply release gates based on manual escalation rates." />
              <FlagItem flag="--manual-rate-gate <%>" desc="Maximum allowed percentage of manual escalations (default 60%)." />
              <FlagItem flag="--log-level <level>" desc="Set log level: error, warn, info, debug, trace (default: warn)." />
              <FlagItem flag="-v, --verbose" desc="Shorthand for --log-level info. Enables detailed structured logging." />
            </div>
          </div>
        </section>

        {/* BLAME COMMAND */}
        <section id="blame" className="scroll-mt-24">
          <div className="flex flex-col md:flex-row md:items-center gap-4 mb-8">
            <div className="flex items-center gap-3">
              <span className="w-8 h-8 rounded-lg bg-[#333] border border-[#444] flex items-center justify-center text-[#888] font-bold text-[10px]">CMD</span>
              <code className="text-xl font-bold text-white bg-[#0a0a0a] px-3 py-1 rounded border border-[#222]">gitresolve blame</code>
            </div>
            <span className="text-gray-500 font-mono text-sm leading-none">Conflict Auditability</span>
          </div>
          
          <div className="prose-layout">
            <p className="mb-6 text-gray-400">
              Queries the session log to output the history of conflicts and resolutions. Essential for post-mortem audits and identifying repetitive conflict patterns.
            </p>
            <div className="grid grid-cols-1 gap-4 mt-4">
              <FlagItem flag="--patterns" desc="Display an analysis of conflict pattern frequencies over time." />
              <FlagItem flag="--file <path>" desc="Filter history for a specific file." />
            </div>
          </div>
        </section>

        {/* STATS COMMAND */}
        <section id="stats" className="scroll-mt-24">
          <div className="flex flex-col md:flex-row md:items-center gap-4 mb-8">
            <div className="flex items-center gap-3">
              <span className="w-8 h-8 rounded-lg bg-pink-600/10 border border-pink-500/20 flex items-center justify-center text-pink-500 font-bold text-[10px]">CMD</span>
              <code className="text-xl font-bold text-white bg-[#0a0a0a] px-3 py-1 rounded border border-[#222]">gitresolve stats</code>
            </div>
            <span className="text-gray-500 font-mono text-sm leading-none">Observability & Metrics</span>
          </div>
          
          <div className="prose-layout text-gray-400">
            <p className="mb-6">
              Reports decision metrics and top reason codes from local observability logs. Use this to monitor your team&apos;s automation efficiency.
            </p>
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 mt-8">
              <FlagItem flag="--json" desc="Emit metrics as a machine-readable JSON object for CI pipes." />
              <FlagItem flag="--top <N>" desc="Show the top N escalation reason codes (default 8)." />
              <FlagItem flag="--operation <all/resolve/merge>" desc="Filter stats by specific operation type." />
            </div>
          </div>
        </section>

        {/* UNDO COMMAND */}
        <section id="undo" className="scroll-mt-24">
          <div className="flex flex-col md:flex-row md:items-center gap-4 mb-8">
            <div className="flex items-center gap-3">
              <span className="w-8 h-8 rounded-lg bg-red-600/10 border border-red-500/20 flex items-center justify-center text-red-500 font-bold text-[10px]">CMD</span>
              <code className="text-xl font-bold text-white bg-[#0a0a0a] px-3 py-1 rounded border border-[#222]">gitresolve undo</code>
            </div>
            <span className="text-gray-500 font-mono text-sm leading-none">Fail-safe Rollback</span>
          </div>
          
          <div className="prose-layout text-gray-400">
            <p className="mb-6">
              Restores the repository to a recorded snapshot SHA from a recent session. Every time a write occurs, a pre-resolution snapshot is captured.
            </p>
            <div className="grid grid-cols-1 gap-4 mt-4">
              <FlagItem flag="--steps <N>" desc="Undo the last N operations (default 1)." />
            </div>
          </div>
        </section>
      </div>
    </DocsShell>
  );
}

function FlagItem({ flag, desc }) {
  return (
    <div className="flex flex-col p-4 rounded-xl bg-[#050505] border border-[#1a1a1a] hover:border-[#333] transition-colors">
       <code className="text-blue-400 font-mono text-[13px] mb-1 font-semibold">{flag}</code>
       <p className="text-[12px] text-gray-500 leading-relaxed">{desc}</p>
    </div>
  );
}
