"use client";

import React from 'react';
import DocsShell from '@/components/DocsShell';
import TerminalWindow from '@/components/TerminalWindow';
import { Shield } from 'lucide-react';

export default function PolicyCommand() {
  return (
    <DocsShell 
      title="policy" 
      subtitle="Inspect policy profile resolution and rules for specific repository paths."
    >
      <div className="space-y-12">
        <section>
          <div className="flex items-center gap-3 mb-8">
            <div className="w-10 h-10 rounded-lg bg-blue-500/10 border border-blue-500/20 flex items-center justify-center">
              <Shield className="w-5 h-5 text-blue-500" />
            </div>
            <code className="text-2xl font-bold text-white bg-black px-3 py-1 rounded-lg border border-white/[0.05] tracking-tight">policy check</code>
          </div>
          
          <div className="docs-prose">
            <TerminalWindow title="bash">
              <div className="flex gap-3 text-[13px] font-mono">
                <span className="text-blue-500 font-bold">$</span>
                <span className="text-white font-bold">gitresolve policy check internal/payments/charge.go</span>
              </div>
            </TerminalWindow>

            <p className="text-[17px]">
              The <code>policy check &lt;file&gt;</code> command resolves the active policy profile for a given file. It reports whether the profile originated from an explicit flag, a path rule in <code>.gitresolve/policy.json</code>, a team rule, or the default profile.
            </p>
            
            <TerminalWindow title="output">
              <div className="space-y-2 text-[13px] font-mono">
                <div className="text-[#888]">Policy Check</div>
                <div className="text-white">  file: internal/payments/charge.go</div>
                <div className="text-white">  requested_profile: auto</div>
                <div className="text-emerald-400 font-bold">  resolved_profile: strict</div>
                <div className="text-white">  source: path_rule</div>
                <div className="text-[#888]">  matched_path: internal/payments/</div>
                <div className="text-rose-400">  strict_blocks_both_for_file: true</div>
              </div>
            </TerminalWindow>

            <h3 className="text-lg font-bold text-white mt-8 mb-4">Command Options</h3>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <FlagItem flag="--policy-profile <profile>" desc="Override policy evaluation with an explicit profile (auto, strict, balanced, aggressive)." />
              <FlagItem flag="--json" desc="Emit policy check resolution details in structured JSON format." />
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
