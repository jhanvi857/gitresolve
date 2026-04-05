"use client";

import React from 'react';

const FlowNode = ({ title, active, children, id }) => (
  <div id={id} className={`p-4 rounded-xl border flex flex-col items-center justify-center text-center transition-all min-w-[140px] ${active ? 'bg-[#111] border-white/40 text-white shadow-xl shadow-white/5' : 'bg-[#0a0a0a] border-[#222] text-gray-500 opacity-60'}`}>
    <span className="text-[9px] uppercase tracking-widest font-black mb-1.5 opacity-50">{title}</span>
    <div className="font-mono text-[11px] font-semibold tracking-tight px-1">
      {children}
    </div>
  </div>
);

const Line = () => <div className="h-10 w-px bg-gradient-to-b from-[#333] to-white/20 self-center" />;
const HorizontalLine = () => <div className="w-10 h-px bg-[#333] self-center shrink-0" />;

const SectionHeader = ({ title, desc }) => (
    <div className="mb-12 border-b border-[#222] pb-6 mt-20 first:mt-0">
        <h2 className="text-xl font-semibold text-white tracking-tight uppercase tracking-widest mb-3">{title}</h2>
        <p className="text-gray-500 text-[15px] leading-relaxed max-w-2xl">{desc}</p>
    </div>
);

export default function Architecture() {
  return (
    <div className="max-w-4xl mx-auto py-8">
      <header className="mb-16">
        <h1 className="text-3xl font-semibold tracking-tighter text-white mb-4">
          Documentation: System Architecture
        </h1>
        <p className="text-lg text-gray-500 leading-relaxed max-w-xl">
          A high-performance engine written in Go, mapping classic Git diff workflows to structural Abstract Syntax Tree (AST) analysis.
        </p>
      </header>

      {/* DETAILED ARCHITECTURE FLOW - Visual Representation */}
      <section className="mb-24 p-10 rounded-xl bg-[#030303] border border-[#222] relative overflow-hidden">
        <div className="absolute top-0 right-0 p-8 text-[40px] font-black opacity-[0.03] select-none uppercase tracking-tighter shrink-0 pointer-events-none">STRUCTURED</div>
        
        <h3 className="text-[10px] font-black text-gray-400 uppercase tracking-[0.3em] mb-12 text-center">Core Operational Flow</h3>
        
        <div className="flex flex-col items-center gap-0">
          <FlowNode title="Input Layer" active>Native .git Repository</FlowNode>
          <Line />
          <FlowNode title="Diagnostic" active>Conflict Stream Discovery</FlowNode>
          <Line />
          
          <div className="flex flex-col md:flex-row items-start gap-4 md:gap-8 justify-center">
            <div className="flex flex-col items-center">
                <FlowNode title="Structural Analysis" active>AST Parsing (tree-sitter)</FlowNode>
                <div className="h-6 w-px bg-[#222]" />
                <div className="text-[10px] text-gray-600 font-mono uppercase">Go / JS / TS / YAML</div>
            </div>
            
            <div className="hidden md:flex flex-col items-center justify-center pt-8">
                <div className="w-16 h-px bg-[#222]" />
            </div>

            <div className="flex flex-col items-center">
                <FlowNode title="Scoring Engine" active>Severity Score (1-10)</FlowNode>
                <div className="h-6 w-px bg-[#222]" />
                <div className="text-[10px] text-gray-600 font-mono uppercase">Predictive Risk</div>
            </div>
          </div>

          <Line />
          
          <div className="flex gap-4 items-center justify-center w-full max-w-md">
             <div className="flex-1 h-px bg-[#222]" />
             <div className="text-[10px] text-gray-600 font-black uppercase tracking-widest px-4">Evaluation</div>
             <div className="flex-1 h-px bg-[#222]" />
          </div>

          <div className="grid grid-cols-2 gap-4 w-full mt-10">
            <div className="flex flex-col items-center gap-4">
                <FlowNode title="Auto Path" active>Deterministic Resolution</FlowNode>
                <div className="h-6 w-px bg-[#222]" />
                <div className="p-3 rounded-lg border border-[#222] bg-black text-[10px] text-gray-400 font-mono">Atomic Disk Write</div>
            </div>
            <div className="flex flex-col items-center gap-4">
                <FlowNode title="Manual Path" active id="manual">Interactive Conflict SDK</FlowNode>
                <div className="h-6 w-px bg-[#222]" />
                <div className="p-3 rounded-lg border border-[#222] bg-black text-[10px] text-gray-400 font-mono">User UI Selection</div>
            </div>
          </div>
          
          <div className="mt-12 w-full pt-8 border-t border-[#111]">
             <div className="flex items-center justify-center gap-4">
                 <div className="p-2 rounded bg-blue-500/10 border border-blue-500/20 text-blue-400 font-bold text-[9px] uppercase tracking-widest">State Logged to SQLite</div>
                 <div className="p-2 rounded bg-green-500/10 border border-green-500/20 text-green-400 font-bold text-[9px] uppercase tracking-widest">FS Transaction Committed</div>
             </div>
          </div>
        </div>
      </section>

      {/* EXTENDED CONTENT - Longer scrolling documentation */}
      <div className="space-y-4">
        
        <SectionHeader 
            title="1. Conflict Diagnostic Layer" 
            desc="The engine starts by recursively scanning the Git index for unmerged files. Unlike standard git-merge, we don't just find the markers; we extract the context around them."
        />
        <div className="grid grid-cols-1 md:grid-cols-2 gap-8 text-[14px] text-gray-500 leading-relaxed">
            <p>
                The diagnostic layer identifies 'chunks' using the `Myers Diff` algorithm. For every chunk, it determines a <strong>Severity Tier</strong>. Scalar tiers (single-word changes) are immediately queued for auto-resolution, while structural markers (function signature shifts) require deep analysis.
            </p>
            <p>
                Each file is locked via a <strong>PID-verified file lock</strong> to prevent concurrent modification during the analysis phase, ensuring that the diagnostic results are 100% consistent with the on-disk state.
            </p>
        </div>

        <SectionHeader 
            title="2. Abstract Syntax Tree (AST) Intelligence" 
            desc="Instead of analyzing raw text, gitresolve integrates tree-sitter based grammars to convert code blocks into logical syntax trees."
        />
        <div className="p-8 rounded-xl bg-[#0a0a0a] border border-[#222] mb-10">
            <h4 className="text-white text-sm font-semibold mb-4 uppercase tracking-widest">How it resolves conflicts:</h4>
            <ul className="space-y-4 text-[13px] text-gray-400">
                <li className="flex gap-4">
                    <span className="w-1.5 h-1.5 rounded-full bg-blue-500 shrink-0 translate-y-1.5" />
                    <div><strong>Import Deduplication:</strong> Automatically detects when both branches added the same package but in different alphabetical order, resolving it seamlessly.</div>
                </li>
                <li className="flex gap-4">
                    <span className="w-1.5 h-1.5 rounded-full bg-blue-500 shrink-0 translate-y-1.5" />
                    <div><strong>Schema Shifts:</strong> For JSON/YAML files, the engine performs a recursive map merge. It identifies if a field was changed on one side and deleted on the other, providing an audit trail.</div>
                </li>
                <li className="flex gap-4">
                    <span className="w-1.5 h-1.5 rounded-full bg-blue-500 shrink-0 translate-y-1.5" />
                    <div><strong>Semantic Validation:</strong> If two function bodies were modified, gitresolve checks if the underlying semantic logic overlaps (e.g., both modified the same variable assignment).</div>
                </li>
            </ul>
        </div>

        <SectionHeader 
            title="3. Atomic Safety & Local Database" 
            desc="The engine prioritizes 'No Data Loss' as its primary directive. This is achieved via architectural persistence and atomic IO."
        />
        <div className="space-y-8 text-[14px] text-gray-500 leading-relaxed">
            <p>
                Every resolution operation is transactional. Before any on-disk write, the engine generates a <strong>Snapshot Signature</strong> of the current repository state. This signature is stored in a local SQLite session database located at `.gitresolve/history.db`.
            </p>
            <div className="grid grid-cols-3 gap-4">
                <div className="p-4 rounded border border-[#222] text-center">
                    <div className="text-white font-bold mb-1">POSIX Atomic</div>
                    <div className="text-[10px]">os.Rename swaps</div>
                </div>
                <div className="p-4 rounded border border-[#222] text-center">
                    <div className="text-white font-bold mb-1">Backup Points</div>
                    <div className="text-[10px]">.orig file copies</div>
                </div>
                <div className="p-4 rounded border border-[#222] text-center">
                    <div className="text-white font-bold mb-1">Capped Log</div>
                    <div className="text-[10px]">1000 history items</div>
                </div>
            </div>
            <p>
                In the event of a system crash, the database maintains the 'Locked' state of the repository. Upon restart, gitresolve detects the stale lock, verifies the PID, and offers a <strong>Resume or Revert</strong> option to the developer.
            </p>
        </div>

        <SectionHeader 
            title="4. Security Architecture" 
            desc="Designed for air-gapped environments where source code privacy is non-negotiable."
        />
        <p className="text-[14px] text-gray-500 leading-relaxed">
            Unlike many modern 'AI' solvers, gitresolve is <strong>100% offline</strong>. It never sends code fragments to remote APIs for resolution. The intelligence lives purely within the static Go binary, making it suitable for government, medical, and high-security enterprise environments that require full source data sovereignty.
        </p>

        <div className="pt-20">
             <div className="p-10 rounded-xl bg-gradient-to-br from-[#050505] to-black border border-[#222] text-center">
                 <h3 className="text-white font-semibold mb-2">Want to dive deeper into the code?</h3>
                 <p className="text-gray-500 text-sm mb-6">Review the implementation of the resolution engine on GitHub.</p>
                 <a href="https://github.com/jhanvi857/gitresolve" target="_blank" className="bg-white text-black font-bold px-6 py-2 rounded text-xs uppercase tracking-widest hover:bg-gray-200 transition-all">View source code branch</a>
             </div>
        </div>

      </div>
    </div>
  );
}
