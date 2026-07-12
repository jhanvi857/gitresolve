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

export default function QuickStart() {
  return (
    <DocsShell 
      title="Quick Start" 
      subtitle="Get up and running with gitresolve in under 5 minutes."
    >
      <div className="flex flex-col gap-6 mt-8">
        <Step number="1" title="Initialize the Repository">
          <p>
            Start by initializing your Git repository to support gitresolve configurations and audits.
          </p>
          <TerminalWindow title="bash">
            <div className="flex gap-3">
              <span className="text-blue-500 font-bold">$</span>
              <span className="text-white font-bold">gitresolve init</span>
            </div>
          </TerminalWindow>
        </Step>

        <Step number="2" title="Check Current Conflict Status">
          <p>
            List all outstanding conflicts and view local confidence ratings.
          </p>
          <TerminalWindow title="bash">
            <div className="flex gap-3">
              <span className="text-blue-500 font-bold">$</span>
              <span className="text-white font-bold">gitresolve status</span>
            </div>
          </TerminalWindow>
        </Step>

        <Step number="3" title="Perform a Precision Scan">
          <p>
            Identify potential merging conflicts before they happen by scanning a target branch.
          </p>
          <TerminalWindow title="bash">
            <div className="flex gap-3">
              <span className="text-blue-500 font-bold">$</span>
              <span className="text-white font-bold">gitresolve scan --target main</span>
            </div>
          </TerminalWindow>
        </Step>

        <Step number="4" title="Resolve Interactively">
          <p>
            Interactively walk through unresolved AST structural conflicts.
          </p>
          <TerminalWindow title="bash">
            <div className="flex gap-3">
              <span className="text-blue-500 font-bold">$</span>
              <span className="text-white font-bold">gitresolve resolve</span>
            </div>
          </TerminalWindow>
        </Step>
      </div>
    </DocsShell>
  );
}
