"use client";

import React from 'react';
import DocsShell from '@/components/DocsShell';
import TerminalWindow from '@/components/TerminalWindow';
import { Activity } from 'lucide-react';

export default function StatusCommand() {
  return (
    <DocsShell 
      title="status" 
      subtitle="Real-time indexing of current unmerged blocks and logical overlaps."
    >
      <div className="space-y-12">
        <section>
          <div className="flex items-center gap-3 mb-8">
            <div className="w-10 h-10 rounded-lg bg-green-500/10 border border-green-500/20 flex items-center justify-center">
              <Activity className="w-5 h-5 text-green-500" />
            </div>
            <code className="text-2xl font-bold text-white bg-black px-3 py-1 rounded-lg border border-white/[0.05] tracking-tight">status</code>
          </div>
          
          <div className="docs-prose">
            <TerminalWindow title="bash">
              <div className="flex gap-3 text-[13px] font-mono">
                <span className="text-blue-500 font-bold">$</span>
                <span className="text-white font-bold">gitresolve status</span>
              </div>
            </TerminalWindow>

            <p className="text-[17px]">
              Displays current unmerged conflict blocks across your working tree, categorized by their <strong>Severity Score</strong> and auto-resolution eligibility. Unlike standard <code>git status</code>, this inspects inside conflict blocks to classify syntax patterns.
            </p>

            <TerminalWindow title="output">
              <div className="space-y-4 font-mono text-[13px]">
                <div className="flex justify-between border-b border-white/[0.05] pb-2 mb-2">
                  <span className="text-[#555] font-bold">SCORE</span>
                  <span className="text-[#555] font-bold">TYPE</span>
                  <span className="text-[#555] font-bold">AUTO</span>
                  <span className="text-[#555] font-bold">FILE</span>
                </div>
                <div className="flex justify-between py-1.5">
                  <span className="text-blue-500 font-bold">5</span>
                  <span className="text-white">whitespace</span>
                  <span className="text-emerald-400 font-bold">YES</span>
                  <span className="text-[#888]">internal/net/socket.go</span>
                </div>
                <div className="flex justify-between py-1.5">
                  <span className="text-blue-500 font-bold">10</span>
                  <span className="text-white">imports</span>
                  <span className="text-emerald-400 font-bold">YES</span>
                  <span className="text-[#888]">main.go</span>
                </div>
                <div className="flex justify-between py-1.5">
                  <span className="text-red-500 font-bold">99</span>
                  <span className="text-white font-bold">logic-overlap</span>
                  <span className="text-rose-400 font-bold">NO</span>
                  <span className="text-[#888]">pkg/auth/handler.go</span>
                </div>
              </div>
            </TerminalWindow>

            <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mt-8">
               <SeverityCard 
                 title="Trivial (0-10)" 
                 color="text-blue-500" 
                 desc="Safe for auto-merge. Whitespace or simple import deduplication." 
               />
               <SeverityCard 
                 title="Structural (11-50)" 
                 color="text-yellow-500" 
                 desc="Requires syntax verification gates or policy-based strategy." 
               />
               <SeverityCard 
                 title="High Risk (99+)" 
                 color="text-red-500" 
                 desc="Logical overlaps that default to manual review or strict policy escalation." 
               />
            </div>
          </div>
        </section>
      </div>
    </DocsShell>
  );
}

function SeverityCard({ title, color, desc }) {
  return (
    <div className="p-6 rounded-xl bg-black border border-white/[0.05] hover-card">
       <h4 className={`font-bold text-[10px] mb-3 uppercase tracking-[0.2em] ${color}`}>{title}</h4>
       <p className="text-[14px] text-[#a1a1aa] font-medium leading-relaxed">{desc}</p>
    </div>
  );
}
