"use client";

import React from 'react';
import { GitBranch, Cpu, Shield, Database, CheckCircle2, AlertTriangle, ArrowDown, Activity, Network } from 'lucide-react';

export default function ArchitectureDiagram() {
  return (
    <div className="my-8 p-6 md:p-8 rounded-2xl bg-[#09090c] border border-white/[0.08] shadow-2xl overflow-x-auto">
      <div className="min-w-[700px] flex flex-col items-center gap-6">
        
        {/* Stage 0: Input */}
        <div className="w-full flex justify-center">
          <div className="px-6 py-3.5 rounded-xl bg-[#111116] border border-blue-500/30 flex items-center gap-3 shadow-[0_0_20px_rgba(59,130,246,0.1)]">
            <GitBranch className="w-5 h-5 text-blue-400" />
            <div>
              <div className="text-[11px] uppercase tracking-widest text-[#888] font-bold">Input Stage</div>
              <div className="text-white font-bold text-[14px]">Git Working Tree & Conflict Markers</div>
            </div>
          </div>
        </div>

        <ArrowDown className="w-5 h-5 text-blue-500/60 animate-bounce" />

        {/* Stage 1: Pre-Resolve Divergence Check */}
        <div className="w-full max-w-xl p-4 rounded-xl bg-purple-500/5 border border-purple-500/20 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="w-9 h-9 rounded-lg bg-purple-500/10 border border-purple-500/30 flex items-center justify-center">
              <Activity className="w-4 h-4 text-purple-400" />
            </div>
            <div>
              <div className="text-[11px] uppercase tracking-widest text-purple-400 font-bold">Stage 1: Pre-Flight</div>
              <div className="text-white font-bold text-[14px]">Remote Divergence & Author Overlap Check</div>
            </div>
          </div>
          <span className="px-2.5 py-1 rounded text-[11px] font-bold bg-purple-500/10 text-purple-400 border border-purple-500/20">
            git log / git rev-list
          </span>
        </div>

        <ArrowDown className="w-5 h-5 text-white/20" />

        {/* Stage 2: AST Extraction & Call Graph */}
        <div className="w-full max-w-xl p-4 rounded-xl bg-blue-500/5 border border-blue-500/20 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="w-9 h-9 rounded-lg bg-blue-500/10 border border-blue-500/30 flex items-center justify-center">
              <Cpu className="w-4 h-4 text-blue-400" />
            </div>
            <div>
              <div className="text-[11px] uppercase tracking-widest text-blue-400 font-bold">Stage 2: AST Parsing</div>
              <div className="text-white font-bold text-[14px]">Syntax Tree & Cross-File Symbol Call Graph</div>
            </div>
          </div>
          <span className="px-2.5 py-1 rounded text-[11px] font-bold bg-blue-500/10 text-blue-400 border border-blue-500/20">
            go/parser & AST
          </span>
        </div>

        <ArrowDown className="w-5 h-5 text-white/20" />

        {/* Stage 3: History-Aware Risk Scoring Matrix */}
        <div className="w-full max-w-2xl p-5 rounded-2xl bg-amber-500/5 border border-amber-500/20">
          <div className="flex items-center justify-between mb-4 pb-3 border-b border-amber-500/10">
            <div className="flex items-center gap-2.5">
              <Shield className="w-4 h-4 text-amber-400" />
              <span className="text-amber-400 font-bold text-[13px] uppercase tracking-wider">
                Stage 3: History-Aware Risk Scoring Layer
              </span>
            </div>
            <span className="text-[11px] font-mono text-[#888]">Threshold Rules</span>
          </div>

          <div className="grid grid-cols-2 sm:grid-cols-4 gap-3">
            <div className="p-3 rounded-lg bg-black/40 border border-white/[0.05]">
              <div className="text-[10px] text-[#888] uppercase font-bold">Blast Radius</div>
              <div className="text-amber-400 font-bold text-[13px] mt-1">&gt; 10 Callers</div>
            </div>
            <div className="p-3 rounded-lg bg-black/40 border border-white/[0.05]">
              <div className="text-[10px] text-[#888] uppercase font-bold">Co-Change</div>
              <div className="text-blue-400 font-bold text-[13px] mt-1">Strength &ge; 0.60</div>
            </div>
            <div className="p-3 rounded-lg bg-black/40 border border-white/[0.05]">
              <div className="text-[10px] text-[#888] uppercase font-bold">Import Cycles</div>
              <div className="text-rose-400 font-bold text-[13px] mt-1">Tarjan SCC</div>
            </div>
            <div className="p-3 rounded-lg bg-black/40 border border-white/[0.05]">
              <div className="text-[10px] text-[#888] uppercase font-bold">Author Weights</div>
              <div className="text-purple-400 font-bold text-[13px] mt-1">e^(-days/90)</div>
            </div>
          </div>
        </div>

        <ArrowDown className="w-5 h-5 text-white/20" />

        {/* Stage 4: Policy & Syntax Gate */}
        <div className="w-full max-w-xl p-4 rounded-xl bg-emerald-500/5 border border-emerald-500/20 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="w-9 h-9 rounded-lg bg-emerald-500/10 border border-emerald-500/30 flex items-center justify-center">
              <CheckCircle2 className="w-4 h-4 text-emerald-400" />
            </div>
            <div>
              <div className="text-[11px] uppercase tracking-widest text-emerald-400 font-bold">Stage 4: Policy Gate</div>
              <div className="text-white font-bold text-[14px]">Policy Profile Routing & Syntax Verification Gate</div>
            </div>
          </div>
          <span className="px-2.5 py-1 rounded text-[11px] font-bold bg-emerald-500/10 text-emerald-400 border border-emerald-500/20">
            policy.json
          </span>
        </div>

        <ArrowDown className="w-5 h-5 text-white/20" />

        {/* Stage 5: Outcomes & Persistence */}
        <div className="w-full max-w-2xl grid grid-cols-1 md:grid-cols-2 gap-4">
          <div className="p-4 rounded-xl bg-emerald-500/10 border border-emerald-500/30 flex items-start gap-3">
            <CheckCircle2 className="w-5 h-5 text-emerald-400 shrink-0 mt-0.5" />
            <div>
              <div className="text-emerald-400 font-bold text-[13px]">Auto-Resolution Path</div>
              <p className="text-[12px] text-[#a1a1aa] mt-1">
                Safe AST merge, post-write validation passed, bit-identical output written to disk.
              </p>
            </div>
          </div>

          <div className="p-4 rounded-xl bg-amber-500/10 border border-amber-500/30 flex items-start gap-3">
            <AlertTriangle className="w-5 h-5 text-amber-400 shrink-0 mt-0.5" />
            <div>
              <div className="text-amber-400 font-bold text-[13px]">Manual Escalation Path</div>
              <p className="text-[12px] text-[#a1a1aa] mt-1">
                2-line structural risk summary + parameter-templated suggested remediation command.
              </p>
            </div>
          </div>
        </div>

        {/* Bottom: SQLite Audit Log */}
        <div className="w-full max-w-xl p-3.5 rounded-xl bg-black border border-white/[0.08] flex items-center justify-between">
          <div className="flex items-center gap-2.5">
            <Database className="w-4 h-4 text-blue-400" />
            <span className="text-white text-[13px] font-bold">Permanent Audit Trail</span>
          </div>
          <span className="text-[11px] font-mono text-[#888]">
            SQLite: .gitresolve/audit.db (WAL mode)
          </span>
        </div>

      </div>
    </div>
  );
}
