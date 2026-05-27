"use client";

import React from 'react';
import DocsShell from '@/components/DocsShell';
import TerminalWindow from '@/components/TerminalWindow';
import { Activity, BarChart3, ShieldCheck, Zap } from 'lucide-react';

export default function StatsCommand() {
  return (
    <DocsShell 
      title="stats" 
      subtitle="Data-driven visibility into your team's conflict resolution trends."
    >
      <div className="space-y-16">
        <section>
          <div className="flex items-center gap-3 mb-8">
            <div className="w-10 h-10 rounded-lg bg-blue-500/10 border border-blue-500/20 flex items-center justify-center">
              <BarChart3 className="w-5 h-5 text-blue-500" />
            </div>
            <code className="text-2xl font-bold text-white bg-black px-3 py-1 rounded-lg border border-white/[0.05] tracking-tight">stats</code>
          </div>

          <div className="docs-prose">
            <p className="text-[#a1a1aa] leading-relaxed text-[17px] font-medium max-w-4xl">
              Every decision made by gitresolve is persisted to a local SQLite database at <code>.gitresolve/audit.db</code>. This gives you a permanent audit trail of how code evolved during merges.
            </p>
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mt-8">
            <StatCard 
              icon={Activity}
              title="Audit Trail"
              desc="Permanent record of every conflict hash, strategy chosen, and engine confidence score."
            />
            <StatCard 
              icon={BarChart3}
              title="Trend Analysis"
              desc="Monitor escalation rates over time to identify high-risk areas in your codebase."
            />
          </div>
        </section>

        <section>
          <div className="mb-8">
            <h2 className="text-2xl font-bold text-white mb-4 tracking-tight">JSON Reporting for CI</h2>
            <p className="text-[#a1a1aa] text-[16px] font-medium">
              For automation addicts, the <code>--json</code> flag emits a structured object containing all critical health metrics of your merge environment.
            </p>
          </div>
          
          <TerminalWindow title="gitresolve stats --json">
            <div className="text-blue-500 font-mono text-[13px] whitespace-pre overflow-x-auto leading-relaxed">
{`{
  "total_decisions": 47,
  "auto_resolved": 31,
  "escalated_to_manual": 16,
  "escalation_rate": 0.34,
  "top_escalation_reasons": [
    { "reason": "semantic.field_type_conflict", "count": 9 },
    { "reason": "validation.go_syntax_failed", "count": 4 },
    { "reason": "parser.malformed_marker", "count": 3 }
  ]
}`}
            </div>
          </TerminalWindow>
        </section>

        <section>
          <div className="mb-8">
            <h2 className="text-2xl font-bold text-white mb-4 tracking-tight">Stable Reason Codes</h2>
            <p className="text-[#a1a1aa] text-[16px] font-medium">
              Reason codes follow a hierarchical namespace to help you identify <i>why</i> a conflict wasn't auto-resolved.
            </p>
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
             <ReasonItem namespace="parser.*" desc="Low-level failure to identify or recover balanced conflict markers." />
             <ReasonItem namespace="semantic.*" desc="High-level logic conflict, such as function signature or field type changes." />
             <ReasonItem namespace="strategy.*" desc="Decision blocked by active Policy Profile or explicit strategy constraints." />
             <ReasonItem namespace="validation.*" desc="Auto-resolution attempted but failed post-write syntax validation." />
          </div>
        </section>
      </div>
    </DocsShell>
  );
}

function StatCard({ icon: Icon, title, desc }) {
  return (
    <div className="p-6 rounded-xl bg-black border border-white/[0.05] hover-card group">
      <div className="w-10 h-10 rounded-lg bg-[#111] border border-[#222] flex items-center justify-center mb-6 group-hover:border-blue-500/50 transition-colors shadow-lg">
        <Icon className="w-5 h-5 text-white group-hover:text-blue-500 transition-colors" />
      </div>
      <h3 className="text-xl font-bold text-white mb-3 tracking-tight">{title}</h3>
      <p className="text-[14px] text-[#a1a1aa] font-medium leading-relaxed">{desc}</p>
    </div>
  );
}

function ReasonItem({ namespace, desc }) {
  return (
    <div className="p-6 rounded-xl bg-black border border-white/[0.05] hover-card group">
       <code className="text-blue-500 font-bold text-[14px] mb-2 block group-hover:text-white transition-colors tracking-tight">{namespace}</code>
       <p className="text-[14px] text-[#a1a1aa] font-medium leading-relaxed">{desc}</p>
    </div>
  );
}
