"use client";

import React from 'react';
import DocsShell from '@/components/DocsShell';
import { Shield, Lock, Zap, ShieldCheck, Terminal, ArrowRight, ShieldAlert } from 'lucide-react';

export default function Security() {
  return (
    <DocsShell 
      title="Security & Privacy" 
      subtitle="How we protect your codebase by staying purely offline and mathematically secure."
    >
      <div className="space-y-16">
        <section>
          <div className="mb-8">
            <h2 className="text-2xl font-bold text-white mb-4 tracking-tight">Zero Data Leakage</h2>
            <p className="text-[#a1a1aa] leading-relaxed text-[17px] font-medium max-w-4xl">
              In an era of AI-driven tools, your source code is often treated as training data. gitresolve takes a different path. We believe your code is your most valuable asset and should never leave your machine.
            </p>
          </div>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <SecurityCard 
              icon={Lock}
              title="100% Offline"
              desc="The tool contains no networking code. It cannot send your code to a remote server because it doesn't know how to talk to the internet."
            />
            <SecurityCard 
              icon={Shield}
              title="CWE-22 Sandboxing"
              desc="Every file operation uses os.Root to ensure no read/write can escape the repository root. Path traversal is mathematically impossible."
            />
            <SecurityCard 
              icon={Zap}
              title="DoS Protection"
              desc="A mandatory 10MB gate prevents memory exhaustion. Maliciously oversized conflict files are skipped and escalated to manual review."
            />
            <SecurityCard 
              icon={ShieldCheck}
              title="Advisory Locking"
              desc="Uses native flock(2) and LockFileEx. Safe from PID-reuse attacks and race conditions in concurrent CI pipes."
            />
          </div>
        </section>

         <section>
          <div className="mb-12">
            <h2 className="text-2xl font-bold text-white mb-4 tracking-tight">Integrity & Verification</h2>
            <p className="text-[#a1a1aa] text-[16px] font-medium">Multi-layered security protocols for every resolution.</p>
          </div>
          <div className="relative">
            <div className="absolute left-5 top-0 bottom-0 w-px bg-gradient-to-b from-blue-500/50 via-white/[0.05] to-transparent" />
            <div className="space-y-4">
              <SecurityStep 
                num="01" 
                title="PII Privacy (Hashing)" 
                desc="Sensitive file content or conflict blocks are never stored in plain text in debug logs. We use 12-char SHA-256 hashes for event correlation." 
                active
              />
              <SecurityStep 
                num="02" 
                title="Syntax Validation" 
                desc="The merged code is passed through a language-specific syntax checker. If the merge creates invalid syntax, the operation is rolled back." 
              />
              <SecurityStep 
                num="03" 
                title="Supply Chain Security" 
                desc="All releases are signed via Cosign (OIDC) and include a CycloneDX SBOM. Binaries are verifiable against the public Rekor transparency log." 
              />
            </div>
          </div>
        </section>

        <section className="pb-16">
          <div className="p-8 rounded-xl border border-red-500/10 bg-red-500/5 hover-card relative overflow-hidden group">
            <div className="absolute top-0 right-0 w-64 h-64 bg-red-500/10 blur-[100px] rounded-full -mr-32 -mt-32 opacity-50" />
            <h4 className="text-[11px] font-bold text-white mb-4 uppercase tracking-[0.2em] flex items-center gap-2 text-red-500">
                <ShieldAlert className="w-4 h-4 shadow-[0_0_10px_rgba(239,68,68,0.5)]" />
                Security First
            </h4>
            <p className="text-[#a1a1aa] text-[16px] font-medium leading-relaxed max-w-4xl">
              Security issues should be reported privately according to the process in <code>SECURITY.md</code> at the repository root. Avoid opening public issues for unpatched vulnerabilities to protect the ecosystem.
            </p>
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

function SecurityStep({ num, title, desc, active }) {
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
