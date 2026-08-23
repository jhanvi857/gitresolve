"use client";

import React from 'react';
import Link from 'next/link';
import DocsShell from '@/components/DocsShell';
import TerminalWindow from '@/components/TerminalWindow';
import { Sparkles, ArrowRight, ShieldCheck, Cpu } from 'lucide-react';

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
      subtitle="Get up and running with gitresolve in under 2 minutes."
    >
      <div className="space-y-8 mt-4">
        {/* Top Feature Spotlight Callout */}
        <Link 
          href="/docs/history-escalation"
          className="block p-5 rounded-2xl bg-gradient-to-r from-amber-500/10 via-amber-500/5 to-transparent border border-amber-500/30 hover:border-amber-500/60 transition-all group shadow-xl"
        >
          <div className="flex items-center justify-between gap-4">
            <div className="flex items-center gap-3.5">
              <div className="w-10 h-10 rounded-xl bg-amber-500/20 border border-amber-500/40 flex items-center justify-center shrink-0">
                <Sparkles className="w-5 h-5 text-amber-400" />
              </div>
              <div>
                <div className="flex items-center gap-2">
                  <span className="text-white font-bold text-[15px]">New Feature: History-Aware Escalation Layer</span>
                  <span className="px-2 py-0.5 rounded-full text-[9px] font-extrabold uppercase tracking-wider bg-amber-500/20 text-amber-300 border border-amber-500/30">
                    v1.4
                  </span>
                </div>
                <p className="text-[13px] text-[#a1a1aa] mt-0.5">
                  Pre-resolve divergence warnings, symbol blast radius scoring, and co-change coupling risk detection.
                </p>
              </div>
            </div>
            <div className="flex items-center gap-1.5 text-amber-400 font-bold text-xs shrink-0 group-hover:translate-x-1 transition-transform">
              <span>Explore Guide</span>
              <ArrowRight className="w-4 h-4" />
            </div>
          </div>
        </Link>

        {/* Quick Steps */}
        <div className="flex flex-col gap-6 pt-4">
          <Step number="1" title="Scan for Predictive Overlaps">
            <p>
              Simulate a merge against your target branch before merging to detect upcoming conflict hotspots early.
            </p>
            <TerminalWindow title="bash">
              <div className="flex gap-3 text-[13px] font-mono">
                <span className="text-blue-500 font-bold">$</span>
                <span className="text-white font-bold">gitresolve scan --target main</span>
              </div>
            </TerminalWindow>
          </Step>

          <Step number="2" title="Check Current Working Tree Status">
            <p>
              Index existing conflict blocks in your working tree, view severity scores, and check auto-resolve eligibility.
            </p>
            <TerminalWindow title="bash">
              <div className="flex gap-3 text-[13px] font-mono">
                <span className="text-blue-500 font-bold">$</span>
                <span className="text-white font-bold">gitresolve status</span>
              </div>
            </TerminalWindow>
          </Step>

          <Step number="3" title="Resolve with History-Aware Risk Scoring">
            <p>
              Interactively step through conflicts with instant divergence checks, blast radius analysis, and suggested remediation commands.
            </p>
            <TerminalWindow title="bash">
              <div className="flex gap-3 text-[13px] font-mono">
                <span className="text-blue-500 font-bold">$</span>
                <span className="text-white font-bold">gitresolve resolve</span>
              </div>
            </TerminalWindow>
          </Step>

          <Step number="4" title="Inspect Resolution Blame & Audit">
            <p>
              Query the local SQLite database to inspect historical conflict decisions and recurring conflict patterns.
            </p>
            <TerminalWindow title="bash">
              <div className="flex gap-3 text-[13px] font-mono">
                <span className="text-blue-500 font-bold">$</span>
                <span className="text-white font-bold">gitresolve blame --patterns</span>
              </div>
            </TerminalWindow>
          </Step>
        </div>
      </div>
    </DocsShell>
  );
}
