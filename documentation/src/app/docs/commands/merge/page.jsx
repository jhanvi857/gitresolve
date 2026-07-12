"use client";

import React from 'react';
import DocsShell from '@/components/DocsShell';
import TerminalWindow from '@/components/TerminalWindow';
import { GitMerge } from 'lucide-react';

export default function MergeCommand() {
  return (
    <DocsShell 
      title="merge" 
      subtitle="Run smart merge on current conflicted files using precision AST heuristics."
    >
      <div className="space-y-12">
        <section>
          <div className="flex items-center gap-3 mb-8">
            <div className="w-10 h-10 rounded-lg bg-blue-500/10 border border-blue-500/20 flex items-center justify-center">
              <GitMerge className="w-5 h-5 text-blue-500" />
            </div>
            <code className="text-2xl font-bold text-white bg-black px-3 py-1 rounded-lg border border-white/[0.05] tracking-tight">merge</code>
          </div>
          
          <div className="docs-prose">
            <p className="text-[17px]">
              The <code>merge</code> command processes files containing Git merge conflict markers. It parses the blocks, runs language-specific AST validation, checks policy boundaries, and auto-applies safe resolutions. Complex conflicts are safely escalated for manual inspection.
            </p>
            
            <TerminalWindow title="bash">
              <div className="space-y-1 text-[13px]">
                <div className="flex gap-3">
                  <span className="text-blue-500 font-bold">$</span>
                  <span className="text-white font-bold">gitresolve merge</span>
                </div>
                <div className="text-[#888] mt-4">Engine Bootup: Initializing gitresolve in directory &apos;.&apos;</div>
                <div className="text-[#888]">Scanning index. Found 2 unmerged conflicts...</div>
                <div className="text-[#888] pt-2">--- Processing main.go ---</div>
                <div className="text-green-500 font-bold">✓ Resolution applied automatically to main.go.</div>
                <div className="text-green-500 font-bold pt-4">Merge complete. auto_resolved: 2, manual_escalations: 0.</div>
              </div>
            </TerminalWindow>

            <h3 className="text-lg font-bold text-white mt-8 mb-4">Command Options</h3>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <FlagItem flag="--dry-run" desc="Show what resolutions would be applied without modifying any files." />
              <FlagItem flag="--shadow" desc="Simulate resolutions and record decisions without writing to the disk." />
              <FlagItem flag="--policy-profile <name>" desc="The security policy profile: auto, strict, balanced, or aggressive." />
              <FlagItem flag="--no-auto-structured" desc="Disable auto-resolution for structured files (JSON/YAML/TOML)." />
              <FlagItem flag="--enforce-gates" desc="Enforce release gate thresholds (manual rate and validation failures)." />
              <FlagItem flag="--manual-rate-gate <%>" desc="Maximum allowed manual escalation rate (default: 60%)." />
            </div>
          </div>
        </section>
      </div>
    </DocsShell>
  );
}

function FlagItem({ flag, desc }) {
  return (
    <div className="p-6 rounded-xl bg-black border border-white/[0.05] hover-card group">
       <code className="text-blue-500 font-bold text-[14px] mb-2 block group-hover:text-white transition-colors tracking-tight">{flag}</code>
       <p className="text-[14px] text-[#a1a1aa] font-medium leading-relaxed">{desc}</p>
    </div>
  );
}
