"use client";

import React from 'react';
import DocsShell from '@/components/DocsShell';

export default function PolicyProfiles() {
  return (
    <DocsShell 
      title="Policy Profiles" 
      subtitle="Configure risk posture and automation behavior per path and team."
    >
      <div className="space-y-12">
        <section>
          <h2 className="text-xl font-semibold text-white mb-4">Risk Management</h2>
          <p className="text-gray-400 leading-relaxed mb-6">
            Not all code carries the same risk. A conflict in a documentation file is trivial, while a conflict in a payment processing handler is critical. Policy profiles allow you to tune <code className="text-gray-300">gitresolve</code>&apos;s automation threshold based on the file&apos;s importance.
          </p>
          
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6 my-10">
            <ProfileCard 
              name="strict" 
              color="text-red-500"
              desc="Maximum escalation. Blocks 'Both' strategy for all source files. Forced manual review for sensitive paths." 
            />
            <ProfileCard 
              name="balanced" 
              color="text-blue-500"
              desc="Default posture. Auto-resolves trivial blocks but escalates on semantic structural changes." 
            />
            <ProfileCard 
              name="aggressive" 
              color="text-green-500"
              desc="Maximum automation. Suitable for generated code, logs, or documentation where safety gates can be relaxed." 
            />
            <ProfileCard 
              name="auto" 
              color="text-purple-500"
              desc="Dynamic resolution. Matches files against .gitresolve/policy.json for path-specific overrides." 
            />
          </div>
        </section>

        <section>
          <h2 className="text-xl font-semibold text-white mb-4">Configuration: policy.json</h2>
          <p className="text-gray-400 mb-6">
            Create a <code className="text-white">.gitresolve/policy.json</code> at your repository root to define fine-grained rules.
          </p>
          
          <div className="code-window border border-[#222]">
            <div className="code-header border-b border-[#111] px-4 py-2 flex gap-2">
              <div className="w-1.5 h-1.5 rounded-full bg-[#333]" /><div className="w-1.5 h-1.5 rounded-full bg-[#333]" /><div className="w-1.5 h-1.5 rounded-full bg-[#333]" />
              <span className="text-[10px] text-gray-600 font-mono ml-2">.gitresolve/policy.json</span>
            </div>
            <div className="code-content bg-black p-5 font-mono text-[13px] leading-relaxed text-blue-300">
{`{
  "default": "balanced",
  "path_profiles": {
    "internal/auth/": "strict",
    "internal/payments/": "strict",
    "docs/": "aggressive",
    "scripts/": "aggressive"
  },
  "team_profiles": {
    "security-ops": "strict",
    "dx-team": "aggressive"
  }
}`}
            </div>
          </div>
        </section>

        <section>
          <h2 className="text-xl font-semibold text-white mb-4">Resolution Logic</h2>
          <div className="space-y-4">
            <div className="p-4 rounded-lg bg-[#0a0a0a] border border-[#222]">
               <h4 className="text-white font-medium text-sm mb-1">Path Overrides</h4>
               <p className="text-xs text-gray-500">The engine uses a <strong>longest-path match</strong>. If a file is in <code className="text-gray-400">internal/auth/utils.go</code>, it receives the &quot;strict&quot; profile even if the rest of <code className="text-gray-400">internal/</code> is balanced.</p>
            </div>
            <div className="p-4 rounded-lg bg-[#0a0a0a] border border-[#222]">
               <h4 className="text-white font-medium text-sm mb-1">Previewing Policy</h4>
               <p className="text-xs text-gray-500">You can preview which profile is active for a file by running:</p>
               <code className="text-[11px] text-blue-400 mt-2 block bg-black border border-[#111] p-2 rounded">gitresolve resolve --policy-profile auto --dry-run path/to/file.go</code>
            </div>
          </div>
        </section>
      </div>
    </DocsShell>
  );
}

function ProfileCard({ name, color, desc }) {
  return (
    <div className="p-6 rounded-xl bg-[#080808] border border-[#1a1a1a] flex flex-col gap-3">
      <div className="flex items-center gap-2">
        <div className={`w-2 h-2 rounded-full ${color.replace('text', 'bg')}`} />
        <h3 className={`font-mono text-sm font-bold tracking-widest ${color}`}>{name.toUpperCase()}</h3>
      </div>
      <p className="text-xs text-gray-500 font-medium leading-relaxed">{desc}</p>
    </div>
  );
}
