"use client";

import React from 'react';
import DocsShell from '@/components/DocsShell';

export default function MergeFlow() {
  return (
    <DocsShell 
      title="Merge Logic & Flow" 
      subtitle="The deterministic decision tree behind every automated resolution."
    >
      <div className="space-y-16">
        <section>
          <h2 className="text-xl font-semibold text-white mb-6">Automated Triage</h2>
          <p className="text-gray-400 mb-8 leading-relaxed">
            gitresolve processes every conflict through a multi-tier confidence engine. We only automate when the risk of logic corruption is statistically negligible or explicitly allowed by your Policy Profile.
          </p>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
             <TriageCard 
               title="Trivial Normalization" 
               color="text-blue-400"
               desc="Whitespace shifts, identical line changes, or purely comment updates. Handled via bitwise line comparison."
             />
             <TriageCard 
               title="Structured Native Merge" 
               color="text-green-400"
               desc="JSON, YAML, or TOML edits that don't overlap on the same key. Uses 3-way object tree merging."
             />
             <TriageCard 
               title="Semantic Import Deduplication" 
               color="text-orange-400"
               desc="Merging additions of new packages or libraries without duplicating the import block header."
             />
             <TriageCard 
               title="Logical Escalation" 
               color="text-red-400"
               desc="Function signature changes, deletion of used logic. These halt automation and require auditor review."
             />
          </div>
        </section>

        <section>
           <h2 className="text-xl font-semibold text-white mb-6">Symmetric Brace Recovery</h2>
           <div className="p-8 rounded-2xl bg-[#050505] border border-[#1a1a1a] border-l-4 border-l-blue-500">
              <p className="text-sm text-gray-500 leading-relaxed">
                Standard Git merge markers often cut through code structures, leaving unbalanced braces ({"{"} or {"}"}) which break the compiler. gitresolve uses a <strong>Lookahead Recovery</strong> algorithm: when it identifies a conflict inside a code block, it scans the surrounding context to ensure all braces remain balanced in the final resolution. 
              </p>
              <div className="mt-6 flex items-center gap-2 text-[10px] font-mono text-gray-700 uppercase tracking-widest">
                 <span className="w-1.5 h-1.5 rounded-full bg-blue-500"></span>
                 Available for Go, Java, and TypeScript
              </div>
           </div>
        </section>

        <section>
           <h2 className="text-xl font-semibold text-white mb-6">The 3-Way Merge Process</h2>
           <div className="prose-layout text-gray-400 text-sm leading-relaxed space-y-4">
              <p>
                Unlike simple &quot;Ours vs Theirs&quot;, gitresolve attempts to find a <strong>Base Ancestor</strong> for every block. This allows us to determine what actually changed on each side:
              </p>
              <ul className="list-disc pl-5 space-y-2">
                 <li>If Side A changed line 5 and Side B is identical to Base, we apply Side A automatically.</li>
                 <li>If both Side A and Side B changed from Base but in the same way, we deduplicate.</li>
                 <li>If both changed from Base differently, we escalate to the <strong>Interactive Orchestrator</strong>.</li>
              </ul>
           </div>
        </section>
      </div>
    </DocsShell>
  );
}

function TriageCard({ title, color, desc }) {
  return (
    <div className="p-5 rounded-xl border border-[#1a1a1a] bg-[#080808]">
      <h3 className={`text-xs font-bold mb-2 uppercase tracking-tight ${color}`}>{title}</h3>
      <p className="text-[12px] text-gray-500 leading-relaxed">{desc}</p>
    </div>
  );
}
