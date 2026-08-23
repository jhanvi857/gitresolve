"use client";

import React from 'react';
import DocsShell from '@/components/DocsShell';
import { ShieldCheck, Lock, FileCode, CheckCircle2, AlertOctagon } from 'lucide-react';

export default function Security() {
  return (
    <DocsShell 
      title="Security & Privacy Standards" 
      subtitle="Architectural security guarantees built specifically for enterprise repositories."
    >
      <div className="space-y-16">
        {/* Zero LLM Guarantee Banner */}
        <section>
          <div className="p-8 rounded-2xl bg-gradient-to-r from-blue-500/10 via-blue-500/5 to-transparent border border-blue-500/20 relative overflow-hidden group hover-card">
            <div className="flex items-start gap-4">
              <div className="w-12 h-12 rounded-xl bg-blue-500/20 border border-blue-500/40 flex items-center justify-center shrink-0">
                <ShieldCheck className="w-6 h-6 text-blue-400" />
              </div>
              <div className="space-y-3">
                <h2 className="text-2xl font-bold text-white tracking-tight">Zero-LLM, Deterministic Guarantee</h2>
                <p className="text-[#a1a1aa] leading-relaxed text-[16px] font-medium">
                  <code>gitresolve</code> strictly does <strong className="text-white">NOT</strong> use Large Language Models (LLMs), probabilistic neural networks, or external AI APIs. Conflict resolution is performed entirely locally through Abstract Syntax Tree (AST) grammar checks, deterministic type unification, and formal policy rules.
                </p>
                <div className="flex flex-wrap gap-3 pt-2">
                  <span className="px-3 py-1 rounded-full text-xs font-bold bg-blue-500/10 text-blue-400 border border-blue-500/20">
                    100% Offline
                  </span>
                  <span className="px-3 py-1 rounded-full text-xs font-bold bg-emerald-500/10 text-emerald-400 border border-emerald-500/20">
                    Zero Data Leakage
                  </span>
                  <span className="px-3 py-1 rounded-full text-xs font-bold bg-purple-500/10 text-purple-400 border border-purple-500/20">
                    Bit-Identical Reproducibility
                  </span>
                </div>
              </div>
            </div>
          </div>
        </section>

        {/* Security Pillars */}
        <section>
          <div className="mb-8">
            <h2 className="text-2xl font-bold text-white mb-4 tracking-tight">Enterprise Hardening Mechanisms</h2>
            <p className="text-[#a1a1aa] text-[16px] font-medium">
              Every operation follows defence-in-depth isolation principles.
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <SecurityCard 
              icon={Lock}
              title="CWE-22 Rooted Path Sandboxing"
              desc="Mandatory os.Root directory confinement guarantees that all file reading and resolution writes are mathematically restricted within the repository boundaries, preventing path traversal exploits (CVE / CWE-22 mitigation)."
            />
            <SecurityCard 
              icon={FileCode}
              title="POSIX Atomic File Writes"
              desc="All auto-resolutions are written first to isolated temporary staging buffers and flushed with fsync() before an atomic os.Rename(). A sudden system crash or kernel panic will never leave files half-written or corrupted."
            />
            <SecurityCard 
              icon={CheckCircle2}
              title="Post-Write Syntax Verification Gate"
              desc="No auto-resolved code file is ever finalized on disk unless it passes native language compiler AST verification (Go syntax, JSON/YAML parser). If syntax fails, changes roll back immediately to manual review."
            />
            <SecurityCard 
              icon={AlertOctagon}
              title="Immutable SQLite Audit Ledger"
              desc="Every decision made by the engine is recorded with timestamp, commit hashes, matched rules, and confidence metrics in .gitresolve/audit.db in WAL mode for forensic auditability."
            />
          </div>
        </section>
      </div>
    </DocsShell>
  );
}

function SecurityCard({ icon: Icon, title, desc }) {
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
