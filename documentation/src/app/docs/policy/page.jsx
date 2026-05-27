"use client";

import React from 'react';
import DocsShell from '@/components/DocsShell';
import TerminalWindow from '@/components/TerminalWindow';
import { ShieldAlert, ShieldCheck, Shield, Zap, ArrowRight } from 'lucide-react';

export default function PolicyProfiles() {
  return (
    <DocsShell 
      title="Policy Profiles" 
      subtitle="Configure risk posture and automation behavior per path and team."
    >
      <div className="space-y-16">
        <section>
          <div className="mb-8">
            <h2 className="text-2xl font-bold text-white mb-4 tracking-tight">Risk Management</h2>
            <p className="text-[#a1a1aa] leading-relaxed text-[17px] font-medium max-w-4xl">
              Not all code carries the same risk. A conflict in a documentation file is trivial, while a conflict in a payment processing handler is critical. Policy profiles allow you to tune gitresolve&apos;s automation threshold based on the file&apos;s importance.
            </p>
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <ProfileCard 
              name="strict" 
              icon={ShieldAlert}
              color="text-red-500"
              desc="Maximum escalation. Blocks 'Both' strategy for all source files. Forced manual review for sensitive paths." 
            />
            <ProfileCard 
              name="balanced" 
              icon={ShieldCheck}
              color="text-blue-500"
              desc="Default posture. Auto-resolves trivial blocks but escalates on semantic structural changes." 
            />
            <ProfileCard 
              name="aggressive" 
              icon={Zap}
              color="text-green-500"
              desc="Maximum automation. Suitable for generated code, logs, or documentation where safety gates can be relaxed." 
            />
            <ProfileCard 
              name="auto" 
              icon={Shield}
              color="text-purple-500"
              desc="Dynamic resolution. Matches files against .gitresolve/policy.json for path-specific overrides." 
            />
          </div>
        </section>

        <section>
          <div className="mb-8">
            <h2 className="text-2xl font-bold text-white mb-4 tracking-tight">Configuration: policy.json</h2>
            <p className="text-[#a1a1aa] text-[16px] font-medium">
              Create a <code>.gitresolve/policy.json</code> at your repository root to define fine-grained rules.
            </p>
          </div>
          
          <TerminalWindow title=".gitresolve/policy.json">
            <div className="text-blue-500 font-mono text-[13px] whitespace-pre overflow-x-auto leading-relaxed">
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
          </TerminalWindow>
        </section>

        <section className="pb-16">
          <div className="mb-8">
            <h2 className="text-2xl font-bold text-white mb-4 tracking-tight">Resolution Logic</h2>
            <p className="text-[#a1a1aa] text-[16px] font-medium">How the engine determines the active profile.</p>
          </div>
          <div className="space-y-6">
            <div className="p-6 rounded-xl bg-black border border-white/[0.05] hover-card">
               <h4 className="text-white font-bold text-lg mb-3 tracking-tight">Path Overrides</h4>
               <p className="text-[15px] text-[#a1a1aa] font-medium leading-relaxed max-w-4xl">The engine uses a <strong>longest-path match</strong>. If a file is in <code>internal/auth/utils.go</code>, it receives the &quot;strict&quot; profile even if the rest of <code>internal/</code> is balanced.</p>
            </div>
            <div className="p-6 rounded-xl bg-black border border-white/[0.05] hover-card">
               <h4 className="text-white font-bold text-lg mb-3 tracking-tight">Previewing Policy</h4>
               <p className="text-[15px] text-[#a1a1aa] font-medium leading-relaxed mb-6 max-w-4xl">You can preview which profile is active for a file by running:</p>
               <div className="px-4 py-2 rounded-lg bg-[#050505] border border-white/[0.1] font-mono text-[13px] text-blue-500 inline-flex items-center gap-3">
                 <span className="text-white opacity-30 font-bold">$</span>
                 gitresolve resolve --policy-profile auto --dry-run path/to/file.go
               </div>
            </div>
          </div>
        </section>
      </div>
    </DocsShell>
  );
}

function ProfileCard({ name, color, desc, icon: Icon }) {
  const bgClass = color.replace('text', 'bg');
  return (
    <div className="p-6 rounded-xl bg-black border border-white/[0.05] hover-card flex flex-col gap-6">
      <div className="flex items-center justify-between">
        <div className={`w-10 h-10 rounded-lg ${bgClass}/10 flex items-center justify-center shadow-lg`}>
          <Icon className={`w-5 h-5 ${color}`} />
        </div>
        <span className={`text-[10px] font-extrabold uppercase tracking-[0.2em] px-2 py-1 rounded-md border ${color} ${color.replace('text', 'border')}/20 bg-black/50`}>
          Profile
        </span>
      </div>
      <div>
        <h3 className={`text-xl font-extrabold text-white mb-2 tracking-tight uppercase`}>{name}</h3>
        <p className="text-[14px] text-[#a1a1aa] font-medium leading-relaxed">{desc}</p>
      </div>
    </div>
  );
}
