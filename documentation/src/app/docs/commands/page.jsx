"use client";

import React from 'react';
import Link from 'next/link';
import DocsShell from '@/components/DocsShell';
import { Cpu, GitMerge, Shield, Activity, History, BarChart3, RotateCcw, Database, ArrowRight } from 'lucide-react';

const commandsList = [
  {
    name: "resolve",
    path: "/docs/commands/resolve",
    icon: Cpu,
    color: "text-blue-400",
    desc: "Interactive resolution orchestrator with pre-resolve divergence checks, blast radius analysis, and guided conflict selection.",
  },
  {
    name: "merge",
    path: "/docs/commands/merge",
    icon: GitMerge,
    color: "text-emerald-400",
    desc: "Automated batch merge engine utilizing precision AST heuristics, structured 3-way maps, and compiler verification gates.",
  },
  {
    name: "scan",
    path: "/docs/commands/scan",
    icon: Shield,
    color: "text-purple-400",
    desc: "Predictive conflict detection simulating merges via git merge-tree before modifying branches or staging files.",
  },
  {
    name: "status",
    path: "/docs/commands/status",
    icon: Activity,
    color: "text-cyan-400",
    desc: "Real-time index of unmerged conflict blocks, categorized by severity score (0-100) and auto-resolve eligibility.",
  },
  {
    name: "blame",
    path: "/docs/commands/blame",
    icon: History,
    color: "text-amber-400",
    desc: "Historical conflict resolution audit and recurring conflict pattern hotspot analysis from local SQLite database.",
  },
  {
    name: "policy check",
    path: "/docs/commands/policy",
    icon: Shield,
    color: "text-rose-400",
    desc: "Inspect active security policy profiles and path-based restriction rules for specific repository files.",
  },
  {
    name: "stats",
    path: "/docs/commands/stats",
    icon: BarChart3,
    color: "text-indigo-400",
    desc: "Data-driven visibility and JSON metrics reporting into team resolution trends and top escalation reason codes.",
  },
  {
    name: "undo",
    path: "/docs/commands/undo",
    icon: RotateCcw,
    color: "text-orange-400",
    desc: "Replay session history in reverse to safely roll back working tree and HEAD to snapshots recorded before the operation.",
  },
  {
    name: "db repair",
    path: "/docs/commands/db",
    icon: Database,
    color: "text-teal-400",
    desc: "Verify SQLite database integrity and re-initialize a clean audit store if corruption or WAL failures occur.",
  },
];

export default function Commands() {
  return (
    <DocsShell 
      title="CLI Command Reference" 
      subtitle="Complete reference for all 9 deterministic gitresolve subcommands."
    >
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {commandsList.map((cmd, i) => {
          const Icon = cmd.icon;
          return (
            <Link 
              key={i}
              href={cmd.path}
              className="p-6 rounded-xl bg-black border border-white/[0.05] hover-card group flex flex-col justify-between"
            >
              <div>
                <div className="flex items-center justify-between mb-4">
                  <div className="w-10 h-10 rounded-lg bg-[#111] border border-[#222] flex items-center justify-center group-hover:border-blue-500/50 transition-colors">
                    <Icon className={`w-5 h-5 ${cmd.color}`} />
                  </div>
                  <code className="font-mono text-xs text-[#888] bg-[#111] px-2 py-1 rounded">
                    gitresolve {cmd.name}
                  </code>
                </div>
                <h3 className="text-xl font-bold text-white mb-2 tracking-tight group-hover:text-blue-400 transition-colors">
                  {cmd.name}
                </h3>
                <p className="text-[14px] text-[#a1a1aa] leading-relaxed font-medium">
                  {cmd.desc}
                </p>
              </div>

              <div className="mt-6 pt-4 border-t border-white/[0.05] flex items-center justify-between text-xs font-bold text-[#666] group-hover:text-white transition-colors">
                <span>View Reference</span>
                <ArrowRight className="w-3.5 h-3.5 group-hover:translate-x-1 transition-transform" />
              </div>
            </Link>
          );
        })}
      </div>
    </DocsShell>
  );
}
