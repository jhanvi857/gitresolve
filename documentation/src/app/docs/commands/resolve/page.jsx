"use client";

import React from 'react';
import DocsShell from '@/components/DocsShell';
import TerminalWindow from '@/components/TerminalWindow';
import { Cpu } from 'lucide-react';

export default function ResolveCommand() {
  return (
    <DocsShell 
      title="resolve" 
      subtitle="Interactive orchestration and automatic conflict resolution with structural risk escalation."
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
            <TerminalWindow title="bash">
              <div className="flex gap-3 text-[13px] font-mono">
                <span className="text-blue-500 font-bold">$</span>
                <span className="text-white font-bold">gitresolve resolve</span>
              </div>
            </TerminalWindow>

            <p className="text-[17px]">
              The primary interactive interface. Use <code>resolve</code> to step through conflicted blocks. It integrates the History-Aware Escalation Layer to perform pre-resolve divergence checks, evaluates structural risks (blast radius, coupled files, import cycles), and automatically merges safe trivial blocks while guiding complex blocks.
            </p>

            <TerminalWindow title="output">
              <div className="space-y-3 text-[13px] font-mono">
                <div className="text-yellow-400 font-semibold">Warning: branch is 15 commits behind main — alice@example.com authored changes touching files you also modified</div>
                <div className="text-[#888] pl-2">suggested: git fetch && git rebase origin/main</div>
                <div className="pt-2 text-white font-bold">Conflict in internal/payments/charge.go (lines 45-58):</div>
                <div className="text-amber-300 pl-2">reason: ProcessPayment is called from 14 other locations — escalating for manual review</div>
                <div className="text-blue-400 pl-2">suggested: go test ./... (run full suite before committing)</div>
                <div className="pt-2 text-blue-500 font-bold">Select resolution [1: ours, 2: theirs, 3: manual, 4: abort]: _</div>
              </div>
            </TerminalWindow>

            <h3 className="text-lg font-bold text-white mt-8 mb-4">Command Options</h3>
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <FlagItem flag="--file <path>" desc="Resolve conflicts in a specific file only." />
              <FlagItem flag="--verbose, -v" desc="Print full factual evidence (AST diff, historical authors, file couplings, decision log row) on escalation." />
              <FlagItem flag="--skip-sync-check" desc="Bypass pre-resolve branch divergence check against remote default branch." />
              <FlagItem flag="--non-interactive" desc="Exit with status 1 if any conflict requires human input. Perfect for CI gates." />
              <FlagItem flag="--dry-run" desc="Preview resolutions and escalation warnings without writing any changes to disk." />
              <FlagItem flag="--strategy <type>" desc="Force a fixed strategy (ours/theirs/both/interactive) for all conflicts." />
              <FlagItem flag="--policy-profile <p>" desc="Override active policy profile (strict, balanced, aggressive, auto)." />
              <FlagItem flag="--timeout <duration>" desc="Timeout for interactive prompt (e.g. 30s). Auto-selects theirs if reached." />
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
