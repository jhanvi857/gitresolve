"use client";

import React from 'react';
import Image from 'next/image';
import DocsShell from '@/components/DocsShell';

export default function Architecture() {
  return (
    <DocsShell
      title="Architecture"
      subtitle="How gitresolve handles conflicts without a centralized brain."
    >
      <div className="space-y-16">
        <section>
          <h2 className="text-xl font-semibold text-white mb-4">Core Principles</h2>
          <p className="text-gray-400 leading-relaxed mb-6">
            gitresolve is built on the principle of <strong>Local Determinism</strong>. Every resolution is produced by a fixed pipeline of rule-based engines that require zero network access and provide bit-identical results for the same input.
          </p>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            <FeatureCard 
               title="AST Integration" 
               desc="Instead of line diffs, we see syntax trees. This allows safe merging of unordered blocks like Go imports or Java annotations."
            />
            <FeatureCard 
               title="Rooted IO Sandbox" 
               desc="Mandatory os.Root sandboxing for all file operations (CWE-22 mitigation). IO is mathematically restricted to the repository root."
            />
            <FeatureCard 
               title="Pre-write Validation" 
               desc="No resolution is written to disk unless it passes the language parser. Syntax errors trigger immediate manual escalation."
            />
            <FeatureCard 
               title="Audit Persistence" 
               desc="Every decision is recorded in a namespaced log, making it possible to query 'why' every resolution occurred months later."
            />
          </div>
        </section>

        <section>
          <div className="flex items-center gap-4 mb-8">
            <h2 className="text-xl font-semibold text-white">The Resolution Pipeline</h2>
            <div className="h-px flex-1 bg-[#222]"></div>
          </div>
          
          <div className="space-y-4">
            <PipeStep 
              num="01" 
              title="Symmetric Marker Identification" 
              desc="Scanning for conflict markers with active brace-balance recovery. If markers are malformed or nested improperly, we escalate immediately." 
              active
            />
            <PipeStep 
              num="02" 
              title="Heuristic Classification" 
              desc="The engine categorizes the block: TypeIdentical, TypeWhitespace, TypeImport, TypeStructured, or TypeScalar." 
            />
            <PipeStep 
              num="03" 
              title="Policy-Injected Routing" 
              desc="Active Policy Profiles (strict/aggressive) modify the confidence threshold required to proceed with automation." 
            />
            <PipeStep 
              num="04" 
              title="Semantic Strategy Execution" 
              desc="Running the actual merge logic. For structured files, this is a deep recursive map merge. For code, it's an AST transformation." 
            />
            <PipeStep 
              num="05" 
              title="Post-Resolution Syntax Gate" 
              desc="The final check. We run the native language compiler/parser on the 'merged' result. If valid, we write. If invalid, we roll back." 
            />
          </div>
        </section>

        <section className="pt-10">
           <div className="flex items-center gap-4 mb-8">
              <h2 className="text-xl font-semibold text-white">Engine Deep Dive</h2>
              <div className="h-px flex-1 bg-[#222]"></div>
           </div>
           <div className="p-2 rounded-2xl bg-[#0a0a0a] border border-[#222] overflow-hidden group">
              <Image 
                src="/architecture.png" 
                alt="Gitresolve System Architecture" 
                width={1200}
                height={800}
                className="w-full h-auto rounded-xl opacity-90 group-hover:opacity-100 transition-opacity"
              />
              <div className="p-6 bg-black">
                <p className="text-xs text-gray-500 leading-relaxed italic text-center">
                  Visual breakdown of the deterministic resolution engine: from marker identification to post-resolution syntax validation.
                </p>
              </div>
           </div>
        </section>

        <section className="pt-10 pb-20">
           <h2 className="text-xl font-semibold text-white mb-6">Structured Data Engine</h2>
           <div className="p-8 rounded-2xl bg-[#050505] border border-[#1a1a1a] relative overflow-hidden group">
              <div className="absolute top-0 right-0 w-64 h-64 bg-blue-500/5 blur-[100px] rounded-full -mr-32 -mt-32"></div>
              <p className="text-sm text-gray-500 leading-relaxed relative z-10">
                For JSON, YAML, and TOML, gitresolve bypasses line-based text merging entirely. It parses both sides into memory, performs a 3-way recursive object merge including <strong>conservative array unioning</strong>, and re-serializes the result. This prevents common errors where merging two objects results in an invalid JSON comma or duplicated keys.
              </p>
           </div>
        </section>
      </div>
    </DocsShell>
  );
}

function FeatureCard({ title, desc }) {
  return (
    <div className="p-6 rounded-xl bg-[#0a0a0a] border border-[#222] hover:border-[#333] transition-colors">
      <h3 className="text-white font-medium mb-2 text-sm">{title}</h3>
      <p className="text-[12px] text-gray-500 leading-relaxed">{desc}</p>
    </div>
  );
}

function PipeStep({ num, title, desc, active }) {
  return (
    <div className={`p-6 rounded-xl border transition-all ${active ? 'bg-blue-600/5 border-blue-500/20 shadow-[0_0_30px_rgba(59,130,246,0.05)]' : 'bg-[#080808] border-[#1a1a1a]'}`}>
      <div className="flex gap-6">
        <span className={`font-mono text-[10px] font-bold mt-1 ${active ? 'text-blue-500' : 'text-gray-700'}`}>{num}</span>
        <div>
          <h4 className={`text-sm font-semibold mb-2 ${active ? 'text-white' : 'text-gray-400'}`}>{title}</h4>
          <p className="text-xs text-gray-500 leading-relaxed max-w-2xl">{desc}</p>
        </div>
      </div>
    </div>
  );
}
