"use client";

import React from 'react';
import DocsShell from '@/components/DocsShell';
import TerminalWindow from '@/components/TerminalWindow';
import { Download, Terminal, Shield, Zap } from 'lucide-react';

export default function Installation() {
  return (
    <DocsShell 
      title="Installation" 
      subtitle="Install the gitresolve standalone deterministic binary across Linux, macOS, and Windows."
    >
      <div className="space-y-16">
        {/* Quick Install */}
        <section>
          <div className="mb-8">
            <h2 className="text-2xl font-bold text-white mb-4 tracking-tight">Quick Install via Go Toolchain</h2>
            <p className="text-[#a1a1aa] leading-relaxed text-[17px] font-medium max-w-3xl">
              <code>gitresolve</code> is built in pure Go with zero external C-dependencies. If you have Go 1.22+ installed, you can compile and install directly to your <code>$GOPATH/bin</code>:
            </p>
          </div>

          <TerminalWindow title="bash">
            <div className="flex gap-3 text-[13px] font-mono">
              <span className="text-blue-500 font-bold">$</span>
              <span className="text-white font-bold">go install github.com/jhanvi857/gitresolve/cmd/gitresolve@latest</span>
            </div>
          </TerminalWindow>
        </section>

        {/* Build from Source */}
        <section>
          <div className="mb-8">
            <h2 className="text-2xl font-bold text-white mb-4 tracking-tight">Build from Source</h2>
            <p className="text-[#a1a1aa] text-[16px] font-medium">
              Clone the repository to build the optimized production release binary locally:
            </p>
          </div>

          <TerminalWindow title="bash">
            <div className="space-y-2 text-[13px] font-mono">
              <div className="flex gap-3">
                <span className="text-blue-500 font-bold">$</span>
                <span className="text-white font-bold">git clone https://github.com/jhanvi857/gitresolve.git</span>
              </div>
              <div className="flex gap-3">
                <span className="text-blue-500 font-bold">$</span>
                <span className="text-white font-bold">cd gitresolve</span>
              </div>
              <div className="flex gap-3">
                <span className="text-blue-500 font-bold">$</span>
                <span className="text-white font-bold">go build -o gitresolve ./cmd/gitresolve</span>
              </div>
              <div className="flex gap-3 pt-2">
                <span className="text-blue-500 font-bold">$</span>
                <span className="text-white font-bold">./gitresolve --version</span>
              </div>
              <div className="text-emerald-400 font-bold">gitresolve version 1.4.0 (pure offline engine)</div>
            </div>
          </TerminalWindow>
        </section>

        {/* Verification */}
        <section>
          <div className="mb-8">
            <h2 className="text-2xl font-bold text-white mb-4 tracking-tight">Verify Installation</h2>
            <p className="text-[#a1a1aa] text-[16px] font-medium">
              Check that <code>gitresolve</code> is available in your shell <code>$PATH</code>:
            </p>
          </div>

          <TerminalWindow title="bash">
            <div className="space-y-1 text-[13px] font-mono">
              <div className="flex gap-3">
                <span className="text-blue-500 font-bold">$</span>
                <span className="text-white font-bold">gitresolve --help</span>
              </div>
              <div className="text-[#888] pt-2">Offline, deterministic AST-powered Git conflict resolution engine</div>
              <div className="text-[#888]">Available Commands: resolve, merge, scan, status, blame, policy, stats, undo, db</div>
            </div>
          </TerminalWindow>
        </section>

        {/* Binary Properties */}
        <section>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            <InstallCard 
              icon={Shield} 
              title="100% Offline" 
              desc="Zero network telemetry, zero external API dependencies, zero LLMs."
            />
            <InstallCard 
              icon={Zap} 
              title="Sub-Millisecond" 
              desc="Optimized Go tree analysis resolves conflicts in milliseconds."
            />
            <InstallCard 
              icon={Terminal} 
              title="Git Native" 
              desc="Integrates seamlessly with existing git merge and rebase workflows."
            />
          </div>
        </section>
      </div>
    </DocsShell>
  );
}

function InstallCard({ icon: Icon, title, desc }) {
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
