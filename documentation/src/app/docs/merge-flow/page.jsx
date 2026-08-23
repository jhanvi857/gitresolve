"use client";

import React from 'react';
import DocsShell from '@/components/DocsShell';
import { GitMerge, Layers, ShieldCheck, Cpu } from 'lucide-react';

export default function MergeFlow() {
  return (
    <DocsShell 
      title="Merge Logic & Flow" 
      subtitle="The deterministic decision tree behind every automated resolution."
    >
      <div className="space-y-16">
        <section>
          <div className="mb-8">
            <h2 className="text-2xl font-bold text-white mb-4 tracking-tight">Automated Conflict Triage</h2>
            <p className="text-[#a1a1aa] leading-relaxed text-[17px] font-medium max-w-4xl">
              <code>gitresolve</code> evaluates each conflict block against language-specific grammar and policy constraints. We only auto-resolve when the probability of semantic logic distortion is provably zero.
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
             <TriageCard 
               title="Trivial Normalization" 
               color="text-blue-400" 
               desc="Whitespace shifts, indentation alignment, identical code lines, or comments. Handled via bitwise line deduplication."
             />
             <TriageCard 
               title="Structured 3-Way Object Merge" 
               color="text-emerald-400" 
               desc="Non-overlapping JSON, YAML, or TOML key edits. Deep recursive tree merging with conservative array unioning."
             />
             <TriageCard 
               title="Semantic Import Deduplication" 
               color="text-amber-400" 
               desc="Merging additions of new packages or libraries across multiple branches without duplicating import declarations."
             />
             <TriageCard 
               title="Structural Logic Escalation" 
               color="text-rose-400" 
               desc="Function signature changes, type alterations, or deletion of active call sites. Escalated for human review with suggested fix."
             />
          </div>
        </section>

        <section>
          <div className="mb-8">
            <h2 className="text-2xl font-bold text-white mb-4 tracking-tight">Symmetric Brace Recovery</h2>
            <p className="text-[#a1a1aa] text-[16px] font-medium">
              Preventing syntax fractures caused by naive line cutting.
            </p>
          </div>

          <div className="p-8 rounded-2xl bg-black border border-white/[0.08] relative overflow-hidden group hover-card">
            <p className="text-[#a1a1aa] text-[15px] leading-relaxed max-w-3xl">
              Standard Git merge markers often bisect code structures, leaving unclosed brackets (<code>{"{"}</code> or <code>{"}"}</code>) that produce immediate compiler failure. <code>gitresolve</code> applies lookahead AST scanning: when a conflict block is detected inside a scope, it calculates surrounding delimiter balance to guarantee valid syntax closure.
            </p>
            <div className="mt-6 flex items-center gap-2 text-[11px] font-mono text-[#666] uppercase tracking-widest">
               <span className="w-2 h-2 rounded-full bg-blue-500 shadow-[0_0_8px_rgba(59,130,246,0.8)]"></span>
               Supported Languages: Go, Java, TypeScript, JavaScript, Python
            </div>
          </div>
        </section>

        <section>
          <div className="mb-8">
            <h2 className="text-2xl font-bold text-white mb-4 tracking-tight">The 3-Way Merge Process</h2>
            <p className="text-[#a1a1aa] text-[16px] font-medium">
              Comparing against common ancestor commits to identify genuine intent.
            </p>
          </div>

          <div className="space-y-4">
             <div className="p-6 rounded-xl bg-black border border-white/[0.05] hover-card">
               <h4 className="text-white font-bold text-base mb-2">1. One-Sided Modification</h4>
               <p className="text-[14px] text-[#a1a1aa]">If Side A modified a block and Side B is identical to the Base ancestor commit, Side A is auto-applied seamlessly.</p>
             </div>
             <div className="p-6 rounded-xl bg-black border border-white/[0.05] hover-card">
               <h4 className="text-white font-bold text-base mb-2">2. Identical Modifications</h4>
               <p className="text-[14px] text-[#a1a1aa]">If both sides made identical edits relative to the Base ancestor, they are deduplicated automatically.</p>
             </div>
             <div className="p-6 rounded-xl bg-black border border-white/[0.05] hover-card">
               <h4 className="text-white font-bold text-base mb-2">3. Divergent Modifications</h4>
               <p className="text-[14px] text-[#a1a1aa]">If both sides made different structural edits to the same logic, the engine runs the History-Aware Escalation scoring matrix and presents an interactive prompt.</p>
             </div>
          </div>
        </section>
      </div>
    </DocsShell>
  );
}

function TriageCard({ title, color, desc }) {
  return (
    <div className="p-6 rounded-xl bg-black border border-white/[0.05] hover-card">
      <h3 className={`text-sm font-bold mb-3 uppercase tracking-wider ${color}`}>{title}</h3>
      <p className="text-[14px] text-[#a1a1aa] leading-relaxed">{desc}</p>
    </div>
  );
}
