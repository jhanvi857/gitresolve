"use client";

import React from 'react';
import DocsShell from '@/components/DocsShell';

export default function StatsMetrics() {
  return (
    <DocsShell 
      title="Stats & Metrics" 
      subtitle="Data-driven visibility into your team's conflict resolution trends."
    >
      <div className="space-y-12">
        <section>
          <h2 className="text-xl font-semibold text-white mb-4">Local Observability</h2>
          <p className="text-gray-400 leading-relaxed mb-6">
            Every decision made by <code className="text-gray-300">gitresolve</code> is persisted to a local SQLite database (typically stored in <code className="text-gray-300">.git/gitresolve.db</code>). This gives you a permanent audit trail of how code evolved during merges.
          </p>
        </section>

        <section>
          <h2 className="text-xl font-semibold text-white mb-4">JSON Reporting for CI</h2>
          <p className="text-gray-400 mb-6">
            For automation addicts, the <code className="text-white">--json</code> flag emits a structured object containing all critical health metrics of your merge environment.
          </p>
          
          <div className="code-window border border-[#222]">
            <div className="code-header border-b border-[#111] px-4 py-2 flex gap-2">
              <div className="w-1.5 h-1.5 rounded-full bg-[#333]" /><div className="w-1.5 h-1.5 rounded-full bg-[#333]" /><div className="w-1.5 h-1.5 rounded-full bg-[#333]" />
              <span className="text-[10px] text-gray-600 font-mono ml-2">gitresolve stats --json</span>
            </div>
            <div className="code-content bg-black p-5 font-mono text-[13px] leading-relaxed text-blue-300 whitespace-pre">
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
          </div>
        </section>

        <section>
          <h2 className="text-xl font-semibold text-white mb-4">Stable Reason Codes</h2>
          <p className="text-gray-400 mb-6">
            Reason codes follow a hierarchical namespace to help you identify <i>why</i> a conflict wasn't auto-resolved.
          </p>
          
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
             <ReasonItem namespace="parser.*" desc="Low-level failure to identify or recover balanced conflict markers." />
             <ReasonItem namespace="semantic.*" desc="High-level logic conflict, such as function signature or field type changes." />
             <ReasonItem namespace="strategy.*" desc="Decision blocked by active Policy Profile or explicit strategy constraints." />
             <ReasonItem namespace="validation.*" desc="Auto-resolution attempted but failed post-write syntax validation." />
          </div>
        </section>

        <section className="mt-20 p-8 rounded-xl border border-[#222] bg-[#050505]">
          <h4 className="text-[12px] font-bold text-white mb-2 uppercase tracking-widest flex items-center gap-2 text-pink-500">
              <span className="w-1.5 h-1.5 rounded-full bg-pink-500 shrink-0" />
              CI Release Gating
          </h4>
          <p className="text-gray-500 text-[13px] leading-relaxed">
            Leading teams use these stats to gate deployments. If the <code className="text-gray-300">escalation_rate</code> is too high on 
            a stable release branch, it indicates the integration is too complex and may require a dedicated rebase period rather than a merge.
          </p>
        </section>
      </div>
    </DocsShell>
  );
}

function ReasonItem({ namespace, desc }) {
  return (
    <div className="p-4 rounded-lg bg-[#0a0a0a] border border-[#222]">
       <code className="text-white font-semibold text-xs mb-1 block">{namespace}</code>
       <p className="text-[11px] text-gray-500">{desc}</p>
    </div>
  );
}
