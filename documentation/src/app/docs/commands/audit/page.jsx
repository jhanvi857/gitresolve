"use client";

import React from 'react';
import DocsShell from '@/components/DocsShell';
import TerminalWindow from '@/components/TerminalWindow';
import { Search, Database } from 'lucide-react';

export default function AuditCommand() {
  return (
    <DocsShell 
      title="audit" 
      subtitle="Query the local decision database for historical evidence and traceability."
    >
      <div className="space-y-12">
        <section>
          <div className="flex items-center gap-3 mb-8">
            <div className="w-10 h-10 rounded-lg bg-blue-500/10 border border-blue-500/20 flex items-center justify-center">
              <Database className="w-5 h-5 text-blue-500" />
            </div>
            <code className="text-2xl font-bold text-white bg-black px-3 py-1 rounded-lg border border-white/[0.05] tracking-tight">audit</code>
          </div>
          
          <div className="docs-prose">
            <p className="text-[17px]">
              The <code>audit</code> command provides a searchable interface into the <code>.gitresolve/audit.db</code>. It is essential for security teams and lead engineers who need to understand exactly why a specific resolution was chosen weeks or months ago.
            </p>

            <TerminalWindow title="bash">
              <div className="space-y-4 text-[13px]">
                <div className="flex gap-3">
                  <span className="text-blue-500 font-bold">$</span>
                  <span className="text-white font-bold">gitresolve audit --file main.go</span>
                </div>
                <div className="font-mono pt-2">
                  <div className="text-[#888]">Found 2 historical resolutions for main.go</div>
                  <div className="mt-2 p-3 rounded bg-white/[0.03] border border-white/[0.05]">
                    <div className="text-white">ID: res_8f2a1</div>
                    <div className="text-white">Date: 2026-05-01 14:22:01</div>
                    <div className="text-white">Strategy: [O]urs</div>
                    <div className="text-white">Confidence: 0.98</div>
                  </div>
                </div>
              </div>
            </TerminalWindow>

            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 mt-8">
              <FlagItem flag="--file <path>" desc="Filter audit logs to only include decisions for a specific file." />
              <FlagItem flag="--id <hash>" desc="Look up a specific resolution event by its unique ID." />
              <FlagItem flag="--export <format>" desc="Export the audit database to CSV or JSON for external analysis." />
              <FlagItem flag="--verify-integrity" desc="Verify that the audit database has not been tampered with since initialization." />
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
