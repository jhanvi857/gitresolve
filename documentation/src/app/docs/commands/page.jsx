"use client";

import React from 'react';
import DocsShell from '@/components/DocsShell';

export default function CommandsReference() {
  return (
    <DocsShell 
      title="Commands Reference" 
      subtitle="Complete documentation for the gitresolve CLI tool."
    >
      <div className="space-y-24">
        {/* SCAN COMMAND */}
        <section id="scan" className="scroll-mt-24">
          <div className="flex flex-col md:flex-row md:items-center gap-4 mb-8">
            <div className="flex items-center gap-3">
              <span className="w-8 h-8 rounded-lg bg-blue-600/10 border border-blue-500/20 flex items-center justify-center text-blue-500 font-bold text-xs">CMD</span>
              <code className="text-xl font-bold text-white bg-[#0a0a0a] px-3 py-1 rounded border border-[#222]">gitresolve scan</code>
            </div>
            <span className="text-gray-500 font-mono text-sm">Predictive Conflict Detection</span>
          </div>
          
          <div className="prose-layout">
            <p className="mb-6">
              The <code className="text-white bg-[#111] px-1.5 py-0.5 rounded">scan</code> command is your early-warning system. By leveraging <code className="text-gray-300">git merge-tree</code>, it simulates a merge between your current HEAD and a target branch to find potential overlaps before you even run a merge command.
            </p>
            
            <div className="code-window mb-8 border border-[#222]">
              <div className="code-header border-b border-[#111] px-4 py-2 flex gap-2">
                <div className="w-2 h-2 rounded-full bg-[#333]" /><div className="w-2 h-2 rounded-full bg-[#333]" /><div className="w-2 h-2 rounded-full bg-[#333]" />
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

            <div className="my-8 p-6 rounded-xl bg-[#070707] border border-[#1a1a1a] border-l-2 border-l-blue-500">
               <h4 className="text-white font-semibold mb-4 text-sm uppercase tracking-wider">Available Flags</h4>
               <div className="space-y-4">
                 <div className="flex gap-4">
                   <code className="text-blue-400 shrink-0">--target</code>
                   <p className="text-xs text-gray-500">The branch to compare against (defaults to <span className="text-gray-300">main</span>).</p>
                 </div>
               </div>
            </div>
          </div>
        </section>

        {/* STATUS COMMAND */}
        <section id="status" className="scroll-mt-24">
          <div className="flex flex-col md:flex-row md:items-center gap-4 mb-8">
            <div className="flex items-center gap-3">
              <span className="w-8 h-8 rounded-lg bg-green-600/10 border border-green-500/20 flex items-center justify-center text-green-500 font-bold text-xs">CMD</span>
              <code className="text-xl font-bold text-white bg-[#0a0a0a] px-3 py-1 rounded border border-[#222]">gitresolve status</code>
            </div>
            <span className="text-gray-500 font-mono text-sm">Real-time Conflict Indexing</span>
          </div>
          
          <div className="prose-layout">
            <p className="mb-6">
              Displays the current unmerged blocks, categorized by their <strong>Severity Score</strong> and auto-resolution eligibility.
            </p>

            <div className="code-window mb-8 border border-[#222]">
              <div className="code-header border-b border-[#111] px-4 py-2 flex gap-2">
                <div className="w-2 h-2 rounded-full bg-[#333]" /><div className="w-2 h-2 rounded-full bg-[#333]" /><div className="w-2 h-2 rounded-full bg-[#333]" />
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

            <ul className="space-y-4 mt-6">
              <li className="flex gap-3 text-sm">
                <span className="text-white font-bold w-24 shrink-0">Score 0-10:</span> Trivial conflicts (whitespace, simple imports). Safe for auto-merge.
              </li>
              <li className="flex gap-3 text-sm">
                <span className="text-white font-bold w-24 shrink-0">Score 10-50:</span> Structural changes requiring verification gates.
              </li>
              <li className="flex gap-3 text-sm">
                <span className="text-white font-bold w-24 shrink-0">Score 99:</span> High-risk logical overlaps that default to manual review.
              </li>
            </ul>
          </div>
        </section>

        {/* MERGE COMMAND */}
        <section id="merge" className="scroll-mt-24">
          <div className="flex flex-col md:flex-row md:items-center gap-4 mb-8">
            <div className="flex items-center gap-3">
              <span className="w-8 h-8 rounded-lg bg-purple-600/10 border border-purple-500/20 flex items-center justify-center text-purple-500 font-bold text-xs">CMD</span>
              <code className="text-xl font-bold text-white bg-[#0a0a0a] px-3 py-1 rounded border border-[#222]">gitresolve merge</code>
            </div>
            <span className="text-gray-500 font-mono text-sm">Autonomous Smart Triage</span>
          </div>
          
          <div className="prose-layout">
            <p className="mb-6">
              The <code className="text-white bg-[#111] px-1.5 py-0.5 rounded">merge</code> command performs non-interactive resolution. It only applies changes that the engine is 100% confident in through deterministic rules.
            </p>

            <div className="code-window mb-8 border border-[#222]">
              <div className="code-header border-b border-[#111] px-4 py-2 flex gap-2">
                <div className="w-2 h-2 rounded-full bg-[#333]" /><div className="w-2 h-2 rounded-full bg-[#333]" /><div className="w-2 h-2 rounded-full bg-[#333]" />
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
                <span className="text-gray-500">   reason: semantic mismatch in Auth handler branch</span><br />
                <br />
                <span className="text-gray-400 italic">Merge scan complete. 1 file resolved, 1 requires manual review.</span>
              </div>
            </div>

            <div className="my-8 p-6 rounded-xl bg-[#070707] border border-[#1a1a1a] border-l-2 border-l-purple-500">
               <h4 className="text-white font-semibold mb-4 text-sm uppercase tracking-wider">Operational Mode</h4>
               <div className="space-y-4">
                 <div className="flex gap-4">
                   <code className="text-purple-400 shrink-0">--dry-run</code>
                   <p className="text-xs text-gray-500">Preview auto-resolutions without modifying any files.</p>
                 </div>
                 <div className="flex gap-4">
                   <code className="text-purple-400 shrink-0">--no-auto-structured</code>
                   <p className="text-xs text-gray-500">Disable auto-resolution for structured files (JSON/YAML/TOML).</p>
                 </div>
               </div>
            </div>
          </div>
        </section>

        {/* RESOLVE COMMAND */}
        <section id="resolve" className="scroll-mt-24">
          <div className="flex flex-col md:flex-row md:items-center gap-4 mb-8">
            <div className="flex items-center gap-3">
              <span className="w-8 h-8 rounded-lg bg-orange-600/10 border border-orange-500/20 flex items-center justify-center text-orange-500 font-bold text-xs">CMD</span>
              <code className="text-xl font-bold text-white bg-[#0a0a0a] px-3 py-1 rounded border border-[#222]">gitresolve resolve</code>
            </div>
            <span className="text-gray-500 font-mono text-sm">Interactive Conflict Orchestration</span>
          </div>
          
          <div className="prose-layout">
            <p className="mb-6">
              The primary interactive interface. Use <code className="text-white bg-[#111] px-1.5 py-0.5 rounded">resolve</code> to take manual actions on high-severity blocks through an interactive interface.
            </p>

            <div className="code-window mb-8 border border-[#222]">
              <div className="code-header border-b border-[#111] px-4 py-2 flex gap-2">
                <div className="w-2 h-2 rounded-full bg-[#333]" /><div className="w-2 h-2 rounded-full bg-[#333]" /><div className="w-2 h-2 rounded-full bg-[#333]" />
              </div>
              <div className="code-content bg-black p-5 font-mono text-[13px] leading-relaxed">
                <span className="text-blue-500">$</span> gitresolve resolve<br />
                <span className="text-white">Active Conflict: pkg/auth/handler.go [Line 124]</span><br />
                <span className="text-gray-500 italic">Type: LogicOverlap | Severity: 99</span><br />
                <br />
                <span className="text-white font-bold underline">Choose strategy:</span><br />
                <span className="text-blue-400">1. Keep Ours (Head)</span><br />
                <span className="text-blue-400">2. Keep Theirs (Incoming)</span><br />
                <span className="text-blue-400">3. Merge Both (Structural)</span><br />
                <span className="text-blue-400">4. Open in Editor</span><br />
                <br />
                <span className="text-white">Selection: _</span>
              </div>
            </div>
          </div>
        </section>


        {/* BLAME COMMAND */}
        <section id="blame" className="scroll-mt-24">
          <div className="flex flex-col md:flex-row md:items-center gap-4 mb-8">
            <div className="flex items-center gap-3">
              <span className="w-8 h-8 rounded-lg bg-[#333] border border-[#444] flex items-center justify-center text-[#888] font-bold text-xs">CMD</span>
              <code className="text-xl font-bold text-white bg-[#0a0a0a] px-3 py-1 rounded border border-[#222]">gitresolve blame</code>
            </div>
            <span className="text-gray-500 font-mono text-sm">Conflict Genealogy</span>
          </div>
          
          <div className="prose-layout">
            <p>
              Analyzes who contributed to the conflicting lines on both sides. Useful for identifying which cross-functional team members need to synchronize.
            </p>
          </div>
        </section>

        {/* UNDO COMMAND */}
        <section id="undo" className="scroll-mt-24">
          <div className="flex flex-col md:flex-row md:items-center gap-4 mb-8">
            <div className="flex items-center gap-3">
              <span className="w-8 h-8 rounded-lg bg-red-600/10 border border-red-500/20 flex items-center justify-center text-red-500 font-bold text-xs">CMD</span>
              <code className="text-xl font-bold text-white bg-[#0a0a0a] px-3 py-1 rounded border border-[#222]">gitresolve undo</code>
            </div>
            <span className="text-gray-500 font-mono text-sm">Fail-safe Rollback</span>
          </div>
          
          <div className="prose-layout">
            <p>
              If a resolution session feels unsafe, <code className="text-white bg-[#111] px-1.5 py-0.5 rounded">undo</code> restores everything to the moment before <code className="text-gray-300">gitresolve</code> touched the index.
            </p>
          </div>
        </section>
      </div>


    </DocsShell>
  );
}
