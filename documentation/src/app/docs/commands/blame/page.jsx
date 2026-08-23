"use client";

import React from 'react';
import DocsShell from '@/components/DocsShell';
import TerminalWindow from '@/components/TerminalWindow';
import { History } from 'lucide-react';

export default function BlameCommand() {
  return (
    <DocsShell 
      title="blame" 
      subtitle="Inspect historical conflict resolutions and recurring conflict patterns."
    >
      <div className="space-y-12">
        <section>
          <div className="flex items-center gap-3 mb-8">
            <div className="w-10 h-10 rounded-lg bg-blue-500/10 border border-blue-500/20 flex items-center justify-center">
              <History className="w-5 h-5 text-blue-500" />
            </div>
            <code className="text-2xl font-bold text-white bg-black px-3 py-1 rounded-lg border border-white/[0.05] tracking-tight">blame</code>
          </div>
          
          <div className="docs-prose">
            <TerminalWindow title="bash">
              <div className="flex gap-3 text-[13px] font-mono">
                <span className="text-blue-500 font-bold">$</span>
                <span className="text-white font-bold">gitresolve blame</span>
              </div>
            </TerminalWindow>

            <p className="text-[17px]">
              The <code>blame</code> command queries the local SQLite decision history database to display past conflict events, resolution strategies used, severity levels, and recurring conflict pattern hotspots across your repository.
            </p>
            
            <TerminalWindow title="output">
              <div className="space-y-2 text-[13px] font-mono">
                <div className="text-[#888]">Conflict History:</div>
                <div className="text-[#555] font-bold">  TYPE            SEVERITY   STRATEGY   FILE</div>
                <div className="text-white">  func_decl       high       theirs     internal/payments/charge.go</div>
                <div className="text-white">  type_decl       medium     ours       models/user.go</div>
                <div className="text-white">  import_decl     trivial    both       pkg/auth/token.go</div>

                <div className="text-[#888] pt-4 font-bold">Conflict Pattern Analysis ($ gitresolve blame --patterns):</div>
                <div className="text-amber-400">  func_decl       8 occurrences</div>
                <div className="text-amber-400">  import_decl     12 occurrences</div>
              </div>
            </TerminalWindow>

            <h3 className="text-lg font-bold text-white mt-8 mb-4">Command Options</h3>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <FlagItem flag="--file <path>" desc="Show conflict history filtered to a specific file." />
              <FlagItem flag="--patterns" desc="Display aggregated conflict pattern analysis and occurrence frequencies." />
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
