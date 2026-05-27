"use client";

import React from 'react';
import DocsShell from '@/components/DocsShell';
import TerminalWindow from '@/components/TerminalWindow';
import { Terminal, Zap } from 'lucide-react';

export default function InitCommand() {
  return (
    <DocsShell 
      title="init" 
      subtitle="Initialize a repository for gitresolve deterministic tracking."
    >
      <div className="space-y-12">
        <section>
          <div className="flex items-center gap-3 mb-8">
            <div className="w-10 h-10 rounded-lg bg-blue-500/10 border border-blue-500/20 flex items-center justify-center">
              <Zap className="w-5 h-5 text-blue-500" />
            </div>
            <code className="text-2xl font-bold text-white bg-black px-3 py-1 rounded-lg border border-white/[0.05] tracking-tight">init</code>
          </div>
          
          <div className="docs-prose">
            <p className="text-[17px]">
              The <code>init</code> command prepares your local environment by creating the <code>.gitresolve</code> directory. This directory stores your local audit database, policy configurations, and temporary resolution buffers.
            </p>
            
            <TerminalWindow title="bash">
              <div className="space-y-1 text-[13px]">
                <div className="flex gap-3">
                  <span className="text-blue-500 font-bold">$</span>
                  <span className="text-white font-bold">gitresolve init</span>
                </div>
                <div className="text-[#888] mt-4">Creating .gitresolve directory...</div>
                <div className="text-[#888]">Initializing audit.db (SQLite)...</div>
                <div className="text-green-500 font-bold pt-2">✓ Repository initialized successfully.</div>
              </div>
            </TerminalWindow>

            <div className="grid grid-cols-1 gap-4 mt-8">
              <FlagItem flag="--force" desc="Overwrite any existing .gitresolve directory and reset the audit database." />
              <FlagItem flag="--template <name>" desc="Initialize with a specific policy template (e.g., 'strict', 'balanced')." />
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
