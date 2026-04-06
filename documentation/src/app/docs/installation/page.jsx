"use client";

import React from 'react';
import DocsShell from '@/components/DocsShell';

const Step = ({ number, title, children }) => (
  <div className="relative pl-10 pb-12 last:pb-0 group text-sm">
    <div className="absolute left-0 top-0 w-6 h-6 rounded-full border border-[#333] bg-[#0a0a0a] flex items-center justify-center text-white font-mono text-[9px] font-bold z-10 group-hover:border-white transition-colors">
      {number}
    </div>
    <div className="absolute left-3 top-6 bottom-0 w-px bg-[#222] group-last:hidden"></div>
    <h3 className="text-[17px] font-semibold text-white mb-3 tracking-tight">{title}</h3>
    <div className="text-gray-500 leading-relaxed text-[14px]">
      {children}
    </div>
  </div>
);

const MacCodeWindow = ({ children, title }) => (
  <div className="code-window mt-4 border-subtle overflow-hidden">
    <div className="code-header">
      <div className="dot dot-red" />
      <div className="dot dot-yellow" />
      <div className="dot dot-green" />
      {title && <span className="ml-4 text-[9px] text-gray-500 font-mono uppercase tracking-widest">{title}</span>}
    </div>
    <div className="code-content whitespace-pre overflow-x-auto bg-black text-[11px]">
      {children}
    </div>
  </div>
);

export default function Installation() {
  return (
    <DocsShell 
      title="Installation & Setup" 
      subtitle="Get gitresolve up and running in your environment. Compiles to a single static binary for maximum portability."
    >
      <div className="flex flex-col gap-2">
        <Step number="1" title="Prerequisites">
          <p className="mb-4">
            You will need the <span className="text-white font-medium underline underline-offset-4 decoration-[#333]">Go 1.20+</span> toolchain installed on your local machine.
          </p>
          <div className="px-3 py-1.5 rounded-md bg-[#0a0a0a] border border-[#222] font-mono text-[11px] inline-block">
             <span className="text-gray-500 mr-3">$</span> 
             <span className="text-white">go version</span>
          </div>
        </Step>

        <Step number="2" title="Install the CLI">
          <p className="mb-4">
            Pull the latest version directly from the source repository.
          </p>
          <MacCodeWindow title="terminal">
            <span className="token-cmd">go install</span> <span className="token-string">github.com/jhanvi857/gitresolve@latest</span>
          </MacCodeWindow>
          <p className="mt-4 text-[11px] text-gray-500 italic">
            Make sure your $GOPATH/bin is included in your system PATH.
          </p>
        </Step>

        <Step number="3" title="Verify Installation">
          <p className="mb-4">
            Test the installation by confirming the version and help outputs.
          </p>
          <MacCodeWindow title="terminal">
            <span className="token-cmd">gitresolve</span> <span className="token-flag">--help</span>
          </MacCodeWindow>
        </Step>

        <Step number="4" title="Workspace Setup (Optional)">
          <p className="mb-4">
            Configure ownership and specialized rules in a <code className="text-white bg-[#111] px-1.5 py-0.5 rounded border border-[#222] text-xs">.gitresolve/owners.json</code> file for enterprise scalability.
          </p>
          <MacCodeWindow title=".gitresolve/owners.json">
<span className="token-output">{`{
  "frontend_team": ["web/**", "components/**"],
  "platform_team": ["pkg/**", "internal/cloud/**"],
  "security_team": ["internal/secrets/**"]
}`}</span>
          </MacCodeWindow>
        </Step>
      </div>

      <div className="mt-20 p-8 rounded-xl border border-[#222] bg-[#050505]">
          <h4 className="text-[12px] font-bold text-white mb-2 uppercase tracking-widest flex items-center gap-2">
              <span className="w-1.5 h-1.5 rounded-full bg-blue-500 shrink-0" />
              Enterprise Support
          </h4>
          <p className="text-gray-500 text-[13px] leading-relaxed">
            For air-gapped environments or secure server deployments, gitresolve can be compiled into a static binary with CGO_ENABLED=0 for zero-dependency execution. 
            Contact our engineering team for specialized compliance support.
          </p>
      </div>
    </DocsShell>
  );
}

