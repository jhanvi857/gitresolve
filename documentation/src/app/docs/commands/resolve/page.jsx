"use client";

import React from 'react';
import DocsShell from '@/components/DocsShell';
import TerminalWindow from '@/components/TerminalWindow';
import { Cpu } from 'lucide-react';

export default function ResolveCommand() {
  return (
    <DocsShell 
      title="resolve" 
      subtitle="Interactive orchestration and automatic conflict resolution."
    >
      <div className="space-y-12">
        <section>
          <div className="flex items-center gap-3 mb-8">
            <div className="w-10 h-10 rounded-lg bg-blue-500/10 border border-blue-500/20 flex items-center justify-center">
              <Cpu className="w-5 h-5 text-blue-500" />
            </div>
            <code className="text-2xl font-bold text-white bg-black px-3 py-1 rounded-lg border border-white/[0.05] tracking-tight">resolve</code>
          </div>
          
          <div className="docs-prose">
            <p className="text-[17px]">
              The primary interactive interface. Use <code>resolve</code> to step through conflicted blocks. It combines auto-resolution for trivial blocks with an interactive prompt for logical conflicts.
            </p>

            <TerminalWindow title="bash">
              <div className="space-y-4 text-[13px]">
                <div className="flex gap-3">
                  <span className="text-blue-500 font-bold">$</span>
                  <span className="text-white font-bold">gitresolve resolve</span>
                </div>
                <div className="pt-4 space-y-1">
                  <div className="text-white font-bold">[Scalar] main.go (L21-25)</div>
                  <div className="text-[#888]"> [O]urs:   // this comment is different</div>
                  <div className="text-[#888]"> [T]heirs: // this is a comment</div>
                  <div className="pt-4 text-blue-500 font-bold">Options: [O]urs [T]heirs [B]oth [M]anual [S]kip</div>
                  <div className="text-white">Select action: _</div>
                </div>
              </div>
            </TerminalWindow>

            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 mt-8">
              <FlagItem flag="--non-interactive" desc="Exit with status 1 if any conflict requires human input. Perfect for CI gates." />
              <FlagItem flag="--timeout <duration>" desc="Auto-select 'theirs' after a timeout (e.g., 30s) during interactive prompts." />
              <FlagItem flag="--strategy <type>" desc="Force a fixed strategy (ours/theirs/both) for all conflicts." />
              <FlagItem flag="--enforce-gates" desc="Apply release gates based on manual escalation rates." />
              <FlagItem flag="--path <glob>" desc="Only resolve conflicts in files matching the specified glob pattern." />
              <FlagItem flag="--dry-run" desc="Preview resolutions without writing any changes to disk." />
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
