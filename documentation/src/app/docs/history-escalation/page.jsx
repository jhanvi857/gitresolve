"use client";

import React from 'react';
import DocsShell from '@/components/DocsShell';
import TerminalWindow from '@/components/TerminalWindow';
import { Activity, GitBranch, Cpu, Network, Users, Terminal, ShieldAlert, Zap } from 'lucide-react';

export default function HistoryEscalation() {
  return (
    <DocsShell 
      title="History-Aware Escalation" 
      subtitle="Local, deterministic structural risk scoring and suggested remediation commands."
    >
      <div className="space-y-16">
        {/* Overview */}
        <section>
          <div className="mb-8">
            <h2 className="text-2xl font-bold text-white mb-4 tracking-tight">Zero-AI, 100% Deterministic Risk Scoring</h2>
            <p className="text-[#a1a1aa] leading-relaxed text-[17px] font-medium max-w-4xl">
              <code>gitresolve</code> incorporates a local, dependency-free escalation layer. Before touching any conflict file, the engine mines repository Git history and inspects AST call trees to answer two questions:
            </p>
            <ul className="list-disc list-inside space-y-2 mt-4 text-[#d4d4d8] text-[15px] font-medium">
              <li><strong>Why did this conflict happen?</strong> (e.g. stale branch divergence, simultaneous edits by teammates).</li>
              <li><strong>What is the blast radius?</strong> (e.g. widely referenced symbols, un-modified coupled partner files, dependency cycles).</li>
            </ul>
            <p className="text-[#a1a1aa] leading-relaxed text-[15px] font-medium mt-4">
              When structural risk exceeds policy thresholds, <code>gitresolve</code> escalates the block for human review with an informative 2-line explanation and a suggested command rendered from parameterized templates filled strictly from verified repository facts.
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <RuleCard 
              name="Blast Radius Analysis"
              code="semantic.high_blast_radius"
              icon={Cpu}
              color="text-amber-400"
              desc="Detects when a modified function or symbol is referenced by more than 10 other locations across the codebase."
            />
            <RuleCard 
              name="Missing Coupled Files"
              code="semantic.missing_coupled_file"
              icon={Network}
              color="text-blue-400"
              desc="Detects when historical co-change strength >= 0.60 indicates a paired file should have been modified, but was not."
            />
            <RuleCard 
              name="Import Cycle Detection"
              code="semantic.import_cycle"
              icon={ShieldAlert}
              color="text-red-400"
              desc="Evaluates Go package dependency graphs using Tarjan's SCC algorithm to flag conflicts sitting inside cyclic dependencies."
            />
            <RuleCard 
              name="Stale Branch Divergence"
              code="strategy.stale_branch_divergence"
              icon={GitBranch}
              color="text-purple-400"
              desc="Pre-flight check running before resolution to detect branches falling behind remote default branch with overlapping authors."
            />
          </div>
        </section>

        {/* CLI Terminal Views */}
        <section>
          <div className="mb-8">
            <h2 className="text-2xl font-bold text-white mb-4 tracking-tight">CLI Terminal Experience</h2>
            <p className="text-[#a1a1aa] text-[16px] font-medium">
              Clean 2-line summary in default mode; comprehensive evidence audit in verbose mode.
            </p>
          </div>

          <div className="space-y-6">
            <div>
              <h4 className="text-white font-bold text-sm mb-3 tracking-tight">1. Default View (Interactive Escalation)</h4>
              <TerminalWindow title="bash">
                <div className="space-y-3 font-mono text-[13px] leading-relaxed">
                  <div className="flex gap-3">
                    <span className="text-green-500 font-bold">$</span>
                    <span className="text-white font-bold">gitresolve resolve</span>
                  </div>
                  <div className="text-yellow-400 font-semibold">Warning: branch is 15 commits behind main — alice@example.com, bob@example.com authored changes touching files you also modified</div>
                  <div className="text-[#888] pl-2">suggested: git fetch && git rebase origin/main</div>
                  <div className="pt-2 text-white font-bold">Conflict in internal/payments/charge.go (lines 45-58):</div>
                  <div className="text-amber-300 pl-2">reason: ProcessPayment is called from 14 other locations — escalating for manual review</div>
                  <div className="text-blue-400 pl-2">suggested: go test ./... (run full suite before committing)</div>
                  <div className="pt-2 text-blue-500 font-bold">Select resolution [1: ours, 2: theirs, 3: manual, 4: abort]: _</div>
                </div>
              </TerminalWindow>
            </div>

            <div>
              <h4 className="text-white font-bold text-sm mb-3 tracking-tight">2. Verbose Evidence View (<code>--verbose</code> / <code>-v</code>)</h4>
              <TerminalWindow title="bash – verbose mode">
                <div className="space-y-2 font-mono text-[12px] leading-relaxed text-[#d4d4d8]">
                  <div className="flex gap-3">
                    <span className="text-green-500 font-bold">$</span>
                    <span className="text-white font-bold">gitresolve resolve -v</span>
                  </div>
                  <div className="text-white font-bold">Conflict in internal/payments/charge.go (lines 45-58):</div>
                  <div className="text-amber-300 pl-2">reason: ProcessPayment is called from 14 other locations — escalating for manual review</div>
                  <div className="text-blue-400 pl-2">suggested: go test ./... (run full suite before committing)</div>
                  <div className="text-[#888] pl-2 font-bold">[verbose evidence]</div>
                  <div className="text-[#a1a1aa] pl-4">AST / Conflict block diff:</div>
                  <div className="text-emerald-400 pl-6">--- ours ---</div>
                  <div className="text-emerald-400 pl-6">+ func ProcessPayment(amount int) error {'{'} ... {'}'}</div>
                  <div className="text-rose-400 pl-6">--- theirs ---</div>
                  <div className="text-rose-400 pl-6">- func ProcessPayment(amount int, currency string) error {'{'} ... {'}'}</div>
                  <div className="text-[#a1a1aa] pl-4 pt-1">Historical author contributions:</div>
                  <div className="text-[#888] pl-6">* alice@company.com (weight: 2.45, last touched: 2026-08-21)</div>
                  <div className="text-[#888] pl-6">* bob@company.com   (weight: 1.10, last touched: 2026-08-15)</div>
                  <div className="text-[#a1a1aa] pl-4 pt-1">Historical file couplings:</div>
                  <div className="text-[#888] pl-6">* handlers/checkout.go (count: 18, strength: 0.90)</div>
                  <div className="text-[#a1a1aa] pl-4 pt-1">Decision log row:</div>
                  <div className="text-[#888] pl-6">reason_code: semantic.high_blast_radius | confidence: 0.20 | type: func_decl | severity: high</div>
                </div>
              </TerminalWindow>
            </div>
          </div>
        </section>

        {/* Reason Codes Table */}
        <section>
          <div className="mb-8">
            <h2 className="text-2xl font-bold text-white mb-4 tracking-tight">Namespaced Reason Codes</h2>
            <p className="text-[#a1a1aa] text-[16px] font-medium">
              Stable reason codes written to the SQLite <code>decision_logs</code> table.
            </p>
          </div>

          <div className="overflow-x-auto rounded-xl border border-white/[0.05]">
            <table className="w-full text-left text-[14px]">
              <thead className="bg-white/[0.02] border-b border-white/[0.05] text-[#888] uppercase text-[11px] font-extrabold tracking-wider">
                <tr>
                  <th className="py-4 px-6">Reason Code</th>
                  <th className="py-4 px-6">Category</th>
                  <th className="py-4 px-6">Severity</th>
                  <th className="py-4 px-6">Trigger Condition</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-white/[0.05] text-[#d4d4d8] font-medium">
                <tr>
                  <td className="py-4 px-6 font-mono text-blue-400 text-[13px]">semantic.high_blast_radius</td>
                  <td className="py-4 px-6">Semantic</td>
                  <td className="py-4 px-6"><span className="px-2 py-0.5 rounded bg-red-500/10 text-red-400 text-[11px] font-bold">High</span></td>
                  <td className="py-4 px-6 text-[#a1a1aa]">Callee count &gt; <code>max_callers</code> (default 10)</td>
                </tr>
                <tr>
                  <td className="py-4 px-6 font-mono text-blue-400 text-[13px]">semantic.missing_coupled_file</td>
                  <td className="py-4 px-6">Semantic</td>
                  <td className="py-4 px-6"><span className="px-2 py-0.5 rounded bg-amber-500/10 text-amber-400 text-[11px] font-bold">Medium</span></td>
                  <td className="py-4 px-6 text-[#a1a1aa]">Co-change strength &ge; 0.60, partner file untouched</td>
                </tr>
                <tr>
                  <td className="py-4 px-6 font-mono text-blue-400 text-[13px]">semantic.import_cycle</td>
                  <td className="py-4 px-6">Semantic</td>
                  <td className="py-4 px-6"><span className="px-2 py-0.5 rounded bg-red-500/10 text-red-400 text-[11px] font-bold">High</span></td>
                  <td className="py-4 px-6 text-[#a1a1aa]">File package is member of a Tarjan SCC cycle</td>
                </tr>
                <tr>
                  <td className="py-4 px-6 font-mono text-purple-400 text-[13px]">strategy.stale_branch_divergence</td>
                  <td className="py-4 px-6">Strategy</td>
                  <td className="py-4 px-6"><span className="px-2 py-0.5 rounded bg-amber-500/10 text-amber-400 text-[11px] font-bold">Medium</span></td>
                  <td className="py-4 px-6 text-[#a1a1aa]">Behind commits &gt; <code>max_divergence_commits</code></td>
                </tr>
                <tr>
                  <td className="py-4 px-6 font-mono text-purple-400 text-[13px]">strategy.multi_author_conflict</td>
                  <td className="py-4 px-6">Strategy</td>
                  <td className="py-4 px-6"><span className="px-2 py-0.5 rounded bg-blue-500/10 text-blue-400 text-[11px] font-bold">Info</span></td>
                  <td className="py-4 px-6 text-[#a1a1aa]">Conflicted sides have distinct historical contributors</td>
                </tr>
              </tbody>
            </table>
          </div>
        </section>

        {/* Performance & Benchmarks */}
        <section className="pb-16">
          <div className="mb-8">
            <h2 className="text-2xl font-bold text-white mb-4 tracking-tight">Performance Guarantees</h2>
            <p className="text-[#a1a1aa] text-[16px] font-medium">
              Graph build + rule evaluation is designed to add negligible overhead to the resolve pipeline.
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            <div className="p-6 rounded-xl bg-black border border-white/[0.05] hover-card">
              <div className="text-[11px] font-extrabold uppercase tracking-widest text-[#888] mb-2">Target Budget</div>
              <div className="text-3xl font-extrabold text-white mb-1">&lt; 200 ms</div>
              <p className="text-[13px] text-[#a1a1aa]">Upper bound ceiling for full graph build on large repositories.</p>
            </div>
            <div className="p-6 rounded-xl bg-black border border-white/[0.05] hover-card">
              <div className="text-[11px] font-extrabold uppercase tracking-widest text-emerald-400 mb-2">Benchmark Result</div>
              <div className="text-3xl font-extrabold text-emerald-400 mb-1">~55 ms</div>
              <p className="text-[13px] text-[#a1a1aa]">Measured over 5,000 commits and 2,000 files in synthetic suite.</p>
            </div>
            <div className="p-6 rounded-xl bg-black border border-white/[0.05] hover-card">
              <div className="text-[11px] font-extrabold uppercase tracking-widest text-blue-400 mb-2">Incremental Sync</div>
              <div className="text-3xl font-extrabold text-blue-400 mb-1">~10 ms</div>
              <p className="text-[13px] text-[#a1a1aa]">Cached commit SHAs in SQLite enable near-instant re-runs.</p>
            </div>
          </div>
        </section>
      </div>
    </DocsShell>
  );
}

function RuleCard({ name, code, desc, color, icon: Icon }) {
  return (
    <div className="p-6 rounded-xl bg-black border border-white/[0.05] hover-card flex flex-col gap-4">
      <div className="flex items-center justify-between">
        <div className="w-10 h-10 rounded-lg bg-white/[0.03] border border-white/[0.08] flex items-center justify-center">
          <Icon className={`w-5 h-5 ${color}`} />
        </div>
        <code className="text-[11px] font-mono text-[#888] px-2 py-1 rounded bg-white/[0.03] border border-white/[0.05]">
          {code}
        </code>
      </div>
      <div>
        <h3 className="text-lg font-bold text-white mb-1 tracking-tight">{name}</h3>
        <p className="text-[14px] text-[#a1a1aa] font-medium leading-relaxed">{desc}</p>
      </div>
    </div>
  );
}
