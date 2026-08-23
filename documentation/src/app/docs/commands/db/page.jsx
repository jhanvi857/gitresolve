"use client";

import React from 'react';
import DocsShell from '@/components/DocsShell';
import TerminalWindow from '@/components/TerminalWindow';
import { Database } from 'lucide-react';

export default function DbCommand() {
  return (
    <DocsShell 
      title="db" 
      subtitle="Manage and repair the local SQLite decision history database."
    >
      <div className="space-y-12">
        <section>
          <div className="flex items-center gap-3 mb-8">
            <div className="w-10 h-10 rounded-lg bg-blue-500/10 border border-blue-500/20 flex items-center justify-center">
              <Database className="w-5 h-5 text-blue-500" />
            </div>
            <code className="text-2xl font-bold text-white bg-black px-3 py-1 rounded-lg border border-white/[0.05] tracking-tight">db repair</code>
          </div>
          
          <div className="docs-prose">
            <TerminalWindow title="bash">
              <div className="flex gap-3 text-[13px] font-mono">
                <span className="text-blue-500 font-bold">$</span>
                <span className="text-white font-bold">gitresolve db repair</span>
              </div>
            </TerminalWindow>

            <p className="text-[17px]">
              The <code>db repair</code> command verifies the integrity of the local SQLite database at <code>.gitresolve/audit.db</code>. If corruption or unreadable WAL state is detected, it archives the damaged file and re-initializes a clean database.
            </p>
            
            <TerminalWindow title="output">
              <div className="space-y-1 text-[13px] font-mono">
                <div className="text-[#888]">Checking database at: .gitresolve/audit.db</div>
                <div className="text-emerald-400 font-bold">Database integrity check passed.</div>
              </div>
            </TerminalWindow>
          </div>
        </section>
      </div>
    </DocsShell>
  );
}
