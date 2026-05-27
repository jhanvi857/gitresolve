"use client";

import React from 'react';
import Image from 'next/image';
import DocsShell from '@/components/DocsShell';
import { Cpu, Shield, Activity, Zap, ArrowRight } from 'lucide-react';

export default function Architecture() {
  return (
    <DocsShell
      title="Architecture"
      subtitle="How gitresolve handles conflicts with mathematical precision and local determinism."
    >
      <div className="space-y-16">
        <section>
          <div className="mb-8">
            <h2 className="text-2xl font-bold text-white mb-4 tracking-tight">Core Principles</h2>
            <p className="text-[#a1a1aa] leading-relaxed text-[17px] font-medium max-w-3xl">
              gitresolve is built on the principle of <strong>Local Determinism</strong>. Every resolution is produced by a fixed pipeline of rule-based engines that require zero network access and provide bit-identical results for the same input.
            </p>
          </div>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <FeatureCard 
               title="AST Integration" 
               icon={Cpu}
               desc="Instead of line diffs, we see syntax trees. This allows safe merging of unordered blocks like Go imports or Java annotations."
            />
            <FeatureCard 
               title="Rooted IO Sandbox" 
               icon={Shield}
               desc="Mandatory os.Root sandboxing for all file operations (CWE-22 mitigation). IO is mathematically restricted to the repository root."
            />
            <FeatureCard 
               title="Pre-write Validation" 
               icon={Activity}
               desc="No resolution is written to disk unless it passes the language parser. Syntax errors trigger immediate manual escalation."
            />
            <FeatureCard 
               title="Audit Persistence" 
               icon={Zap}
               desc="Every decision is recorded in a namespaced log, making it possible to query 'why' every resolution occurred months later."
            />
          </div>
        </section>

        <section>
          <div className="mb-12">
            <h2 className="text-2xl font-bold text-white mb-4 tracking-tight">The Resolution Pipeline</h2>
            <p className="text-[#a1a1aa] text-[16px] font-medium">A five-stage sequential processing engine.</p>
          </div>
          <div className="relative">
            <div className="absolute left-5 top-0 bottom-0 w-px bg-gradient-to-b from-blue-500/50 via-white/[0.05] to-transparent" />
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
                desc="Running the actual merge logic. For structured files, this is a deep recursive map merge. For code, it&apos;s an AST transformation." 
              />
              <PipeStep 
                num="05" 
                title="Post-Resolution Syntax Gate" 
                desc="The final check. We run the native language compiler/parser on the 'merged' result. If valid, we write. If invalid, we roll back." 
              />
            </div>
          </div>
        </section>

        <section>
          <div className="mb-8">
            <h2 className="text-2xl font-bold text-white mb-4 tracking-tight">Engine Deep Dive</h2>
            <p className="text-[#a1a1aa] text-[16px] font-medium">Visual breakdown of the deterministic resolution engine.</p>
          </div>
          <div className="p-1 rounded-2xl bg-white/[0.05] border border-white/[0.1] overflow-hidden shadow-2xl">
            <div className="bg-black rounded-xl overflow-hidden">
               <Image 
                 src="/gitresolve_architecture_diagram.png" 
                 alt="Gitresolve System Architecture" 
                 width={1200}
                 height={800}
                 className="w-full h-auto opacity-90"
               />
            </div>
          </div>
          <p className="mt-6 text-[11px] text-[#333] font-bold uppercase tracking-[0.2em] text-center">
            Marker Identification → Semantic Analysis → Syntax Validation
          </p>
        </section>

        <section className="pb-16">
           <div className="p-8 rounded-2xl bg-blue-500/5 border border-blue-500/10 relative overflow-hidden group hover-card">
              <div className="absolute top-0 right-0 w-96 h-96 bg-blue-500/10 blur-[120px] rounded-full -mr-48 -mt-48 transition-opacity duration-500 group-hover:opacity-100 opacity-50" />
              <h3 className="text-xl font-bold text-white mb-4 relative z-10">Structured Data Engine</h3>
              <p className="text-[#a1a1aa] leading-relaxed relative z-10 text-[17px] font-medium max-w-4xl">
                For JSON, YAML, and TOML, gitresolve bypasses line-based text merging entirely. It parses both sides into memory, performs a 3-way recursive object merge including <strong>conservative array unioning</strong>, and re-serializes the result. This prevents common errors where merging two objects results in an invalid JSON comma or duplicated keys.
              </p>
              <div className="mt-8 flex items-center gap-2 text-blue-500 font-bold text-[13px] uppercase tracking-widest cursor-pointer group-hover:gap-4 transition-all relative z-10">
                Read about structured merging <ArrowRight className="w-4 h-4" />
              </div>
           </div>
        </section>
      </div>
    </DocsShell>
  );
}

function FeatureCard({ title, desc, icon: Icon }) {
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

function PipeStep({ num, title, desc, active }) {
  return (
    <div className={`p-6 rounded-xl border ml-10 transition-all duration-300 relative ${active ? 'bg-blue-500/5 border-blue-500/20' : 'bg-transparent border-transparent hover:bg-white/[0.02] hover:border-white/[0.05]'}`}>
      <div className={`absolute -left-[2.75rem] top-7 w-3 h-3 rounded-full border-2 transition-all duration-300 ${active ? 'bg-blue-500 border-blue-500 shadow-[0_0_12px_rgba(59,130,246,0.5)]' : 'bg-black border-white/[0.1]'}`} />
      <div className="flex gap-6">
        <span className={`font-mono text-xs font-extrabold mt-1 tracking-tighter ${active ? 'text-blue-500' : 'text-[#333]'}`}>{num}</span>
        <div>
          <h4 className={`text-lg font-bold mb-2 ${active ? 'text-white' : 'text-[#a1a1aa]'}`}>{title}</h4>
          <p className="text-[15px] text-[#a1a1aa] font-medium leading-relaxed max-w-3xl">{desc}</p>
        </div>
      </div>
    </div>
  );
}
