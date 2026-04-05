"use client";

import React from 'react';

const ShieldIcon = () => (
    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="text-white opacity-40 group-hover:opacity-100 transition-opacity">
        <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"></path>
    </svg>
);

const SecurityCard = ({ title, icon, children }) => (
    <div className="p-10 rounded-xl border border-[#222] bg-[#0a0a0a] hover:border-white transition-colors group">
        <div className="w-10 h-10 rounded bg-[#111] border border-[#222] flex items-center justify-center mb-10 text-white text-[11px]">
            {icon}
        </div>
        <h3 className="text-[17px] font-semibold text-white mb-3 tracking-tight">{title}</h3>
        <div className="text-gray-500 leading-relaxed text-[14px]">
            {children}
        </div>
    </div>
);

export default function Security() {
  return (
    <div className="max-w-4xl mx-auto py-8">
      <header className="mb-16">
        <div className="inline-flex items-center gap-2 px-3 py-1 rounded bg-white text-black text-[10px] font-black uppercase tracking-widest mb-6">
            Privacy First
        </div>
        <h1 className="text-3xl font-semibold tracking-tighter text-white mb-4 uppercase">
            Security Standard
        </h1>
        <p className="text-[17px] text-gray-500 leading-relaxed max-w-xl">
            The gitresolve engine is built for mission-critical engineering where proprietary source code leaves the local machine under zero circumstances.
        </p>
      </header>

      <section className="mb-12">
        <div className="p-10 rounded-xl border border-[#222] bg-[#0a0a0a] relative overflow-hidden">
            <h2 className="text-xl font-semibold text-white mb-6 tracking-tight uppercase tracking-[0.1em]">
                Zero LLM Integrations
            </h2>
            <div className="space-y-4 text-gray-500 leading-relaxed relative z-10 max-w-2xl text-[15px]">
                <p>
                    <strong className="text-white font-medium">gitresolve strictly does not use LLMs (Large Language Models) or probabilistic AI networks to resolve your source code.</strong> 
                </p>
                <p>
                    Sending proprietary intellectual property to third-party APIs acts as a major security and compliance leak. Mathematical trees are deterministic; LLMs are not.
                </p>
            </div>
            <div className="mt-12 flex flex-wrap gap-10 py-6 border-t border-[#222]">
                <div className="flex items-center gap-2">
                    <span className="w-1 h-1 rounded-full bg-blue-500" />
                    <span className="text-[10px] font-bold text-white uppercase tracking-[0.15em]">100% Privacy</span>
                </div>
                <div className="flex items-center gap-2">
                    <span className="w-1 h-1 rounded-full bg-blue-500" />
                    <span className="text-[10px] font-bold text-white uppercase tracking-[0.15em]">0 API Calls</span>
                </div>
                <div className="flex items-center gap-2">
                    <span className="w-1 h-1 rounded-full bg-blue-500" />
                    <span className="text-[10px] font-bold text-white uppercase tracking-[0.15em]">Verifiable Logs</span>
                </div>
            </div>
        </div>
      </section>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-10 mb-20">
          <SecurityCard title="Atomic Reliability" icon={<ShieldIcon />}>
              <p>
                  Writing to files is safe. We use POSIX atomic renaming logic via os.Rename to ensure repository state is never corrupted during power loss or system panics.
              </p>
          </SecurityCard>

          <SecurityCard title="Reversible Auditing" icon={<ShieldIcon />}>
              <p>
                  Every operation is tracked in a local SQLite session database. Use gitresolve undo to restore your repository post-merge to its exact pre-operation state perfectly.
              </p>
          </SecurityCard>
      </div>

      <footer className="mt-20 p-10 rounded-xl border border-[#222] bg-[#050505] text-center max-w-lg mx-auto">
          <h4 className="text-[13px] font-bold text-white mb-2 uppercase tracking-[0.2em] opacity-80">Enterprise Compliance</h4>
          <p className="text-gray-500 text-[12px] leading-relaxed mb-8 px-4">
              Designed to pass SOC2 and ISO 27001 audits by proving that no sensitive material is exported from the environment.
          </p>
          <button className="px-5 py-2.5 rounded-md bg-white text-black font-bold text-[10px] hover:bg-gray-200 transition-all uppercase tracking-widest shrink-0">
              Request Security Audit Log
          </button>
      </footer>
    </div>
  );
}
