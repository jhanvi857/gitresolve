"use client";

import React from 'react';
import DocsShell from '@/components/DocsShell';
import TerminalWindow from '@/components/TerminalWindow';
import { RotateCcw } from 'lucide-react';

export default function UndoCommand() {
  return (
    <DocsShell 
      title="undo" 
      subtitle="Replay session log in reverse to safely undo conflict resolutions."
    >
      <div className="space-y-12">
        <section>
          <div className="flex items-center gap-3 mb-8">
            <div className="w-10 h-10 rounded-lg bg-blue-500/10 border border-blue-500/20 flex items-center justify-center">
              <RotateCcw className="w-5 h-5 text-blue-500" />
            </div>
            <code className="text-2xl font-bold text-white bg-black px-3 py-1 rounded-lg border border-white/[0.05] tracking-tight">undo</code>
          </div>
          
          <div className="docs-prose">
            <TerminalWindow title="bash">
              <div className="flex gap-3 text-[13px] font-mono">
                <span className="text-blue-500 font-bold">$</span>
                <span className="text-white font-bold">gitresolve undo --steps 1</span>
              </div>
            </TerminalWindow>

            <p className="text-[17px]">
              The <code>undo</code> command rolls back resolutions by inspecting the local session history in SQLite. It restores repository HEAD and working tree state back to the exact commit snapshot recorded before the operation ran.
            </p>
            
            <TerminalWindow title="output">
              <div className="space-y-1 text-[13px] font-mono">
                <div className="text-[#888]">Undoing last 1 operation(s) -&gt; resetting to 7b4c9e8f...</div>
                <div className="text-emerald-400 font-bold">✓ Undo successful.</div>
              </div>
            </TerminalWindow>

            <h3 className="text-lg font-bold text-white mt-8 mb-4">Command Options</h3>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <FlagItem flag="--steps <N>" desc="Number of historical resolution operations to roll back (default: 1)." />
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
