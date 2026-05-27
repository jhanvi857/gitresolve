"use client";

import React from 'react';
import DocsShell from '@/components/DocsShell';
import TerminalWindow from '@/components/TerminalWindow';
import { Shield } from 'lucide-react';

export default function ScanCommand() {
  return (
    <DocsShell 
      title="scan" 
      subtitle="Predictive conflict detection before the merge happens."
    >
      <div className="space-y-12">
        <section>
          <div className="flex items-center gap-3 mb-8">
            <div className="w-10 h-10 rounded-lg bg-blue-500/10 border border-blue-500/20 flex items-center justify-center">
              <Shield className="w-5 h-5 text-blue-500" />
            </div>
            <code className="text-2xl font-bold text-white bg-black px-3 py-1 rounded-lg border border-white/[0.05] tracking-tight">scan</code>
          </div>
          
          <div className="docs-prose">
            <p className="text-[17px]">
              The <code>scan</code> command is your early-warning system. By leveraging <code>git merge-tree</code>, it simulates a merge between your current HEAD and a target branch to find potential overlaps before you even run a merge command.
            </p>
            
            <TerminalWindow title="bash">
              <div className="space-y-1 text-[13px]">
                <div className="flex gap-3">
                  <span className="text-blue-500 font-bold">$</span>
                  <span className="text-white font-bold">gitresolve scan --target main</span>
                </div>
                <div className="text-[#888] mt-4">Scanning for potential conflicts against main...</div>
                <div className="text-white font-bold pt-2">Potential conflicts detected: 2 files/blocks.</div>
                <div className="space-y-0.5 pt-2">
                  <div className="text-yellow-500 font-medium">⚠ Conflict (modify/modify): internal/db/schema.go</div>
                  <div className="text-yellow-500 font-medium">⚠ Conflict (add/add): config/default.yaml</div>
                </div>
              </div>
            </TerminalWindow>

            <div className="grid grid-cols-1 gap-4 mt-8">
              <FlagItem flag="--target <branch>" desc="The target branch to compare against (defaults to main)." />
              <FlagItem flag="--depth <int>" desc="How many commits deep to look for potential diverging paths." />
              <FlagItem flag="--ignore-whitespace" desc="Skip conflicts that are purely whitespace changes." />
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
