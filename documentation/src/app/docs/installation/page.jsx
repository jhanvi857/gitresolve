"use client";

import React from 'react';
import DocsShell from '@/components/DocsShell';
import TerminalWindow from '@/components/TerminalWindow';

const Step = ({ number, title, children }) => (
  <div className="relative pl-12 pb-10 last:pb-0 group">
    <div className="absolute left-0 top-1 w-8 h-8 rounded-lg border border-white/[0.1] bg-[#0a0a0a] flex items-center justify-center text-white font-bold text-xs z-10 group-hover:border-blue-500 transition-all duration-300">
      {number}
    </div>
    <div className="absolute left-4 top-10 bottom-0 w-px bg-white/[0.05] group-last:hidden"></div>
    <h3 className="text-xl font-bold text-white mb-4 tracking-tight">{title}</h3>
    <div className="text-[#a1a1aa] text-[15px] leading-relaxed font-medium">
      {children}
    </div>
  </div>
);

export default function Installation() {
  return (
    <DocsShell 
      title="Installation" 
      subtitle="Get gitresolve up and running in your environment. Compiles to a single static binary for maximum portability."
    >
      <div className="flex flex-col gap-6 mt-8">
        <Step number="1" title="Prerequisites">
          <p>
            You will need the <span className="text-white font-bold underline underline-offset-4 decoration-blue-500/20">Go 1.20+</span> toolchain installed on your local machine to compile the engine.
          </p>
          <div className="mt-6 px-4 py-2 rounded-lg bg-black border border-white/[0.05] font-mono text-[13px] inline-flex items-center gap-3">
             <span className="text-blue-500 font-bold">$</span> 
             <span className="text-white">go version</span>
          </div>
        </Step>

        <Step number="2" title="Install the CLI">
          <p>
            Pull the latest version directly from the source repository using the Go toolchain.
          </p>
          <TerminalWindow title="bash">
            <div className="flex gap-3">
              <span className="text-blue-500 font-bold">$</span>
              <span className="text-white">go install github.com/jhanvi857/gitresolve@latest</span>
            </div>
          </TerminalWindow>
          <p className="mt-4 text-[13px] text-[#555] font-medium">
            Make sure your <code className="text-white">$GOPATH/bin</code> is included in your system <code className="text-white">PATH</code>.
          </p>
        </Step>

        <Step number="3" title="Verify Installation">
          <p>
            Test the installation by confirming the version and help outputs.
          </p>
          <TerminalWindow title="bash">
            <div className="flex gap-3">
              <span className="text-blue-500 font-bold">$</span>
              <span className="text-white">gitresolve --help</span>
            </div>
          </TerminalWindow>
        </Step>

        <Step number="4" title="Security & Integrity">
          <p>
            For production environments, verify the integrity of the binary using checksums and Cosign signatures.
          </p>
          <TerminalWindow title="verification">
            <div className="flex gap-3">
              <span className="text-blue-500 font-bold">$</span>
              <span className="text-white">cosign verify-blob --certificate checksums.txt.pem --signature checksums.txt.sig checksums.txt</span>
            </div>
          </TerminalWindow>
        </Step>

        <Step number="5" title="Quick Start Commands">
          <p>
            Try some basic commands to see gitresolve in action.
          </p>
          <TerminalWindow title="quick start">
            <div className="space-y-4">
              <div>
                <div className="text-[#555] mb-1"># View current conflicts with block-level severity</div>
                <div className="flex gap-3">
                  <span className="text-blue-500 font-bold">$</span>
                  <span className="text-white font-bold">gitresolve status</span>
                </div>
              </div>
              <div>
                <div className="text-[#555] mb-1"># Resolve interactively</div>
                <div className="flex gap-3">
                  <span className="text-blue-500 font-bold">$</span>
                  <span className="text-white font-bold">gitresolve resolve</span>
                </div>
              </div>
              <div>
                <div className="text-[#555] mb-1"># Predict conflicts before a merge</div>
                <div className="flex gap-3">
                  <span className="text-blue-500 font-bold">$</span>
                  <span className="text-white font-bold">gitresolve scan --target main</span>
                </div>
              </div>
            </div>
          </TerminalWindow>
        </Step>
      </div>

      <div className="mt-16 p-8 rounded-xl border border-white/[0.05] bg-black hover-card">
          <h4 className="text-[11px] font-bold text-white mb-4 uppercase tracking-[0.2em] flex items-center gap-2">
              <div className="w-2 h-2 rounded-full bg-blue-500 shadow-[0_0_8px_rgba(0,112,243,0.5)]" />
              Enterprise Support
          </h4>
          <p className="text-[#a1a1aa] text-[15px] font-medium leading-relaxed">
            For air-gapped environments or secure server deployments, gitresolve can be compiled into a static binary with <code className="text-white">CGO_ENABLED=0</code> for zero-dependency execution. 
            Contact our engineering team for specialized compliance support.
          </p>
      </div>
    </DocsShell>
  );
}

