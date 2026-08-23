"use client";

import React from 'react';
import DocsShell from '@/components/DocsShell';
import TerminalWindow from '@/components/TerminalWindow';
import { Settings } from 'lucide-react';

export default function ConfigCommand() {
  return (
    <DocsShell 
      title="config" 
      subtitle="Manage global and repository-specific settings and policy profiles."
    >
      <div className="space-y-12">
        <section>
          <div className="flex items-center gap-3 mb-8">
            <div className="w-10 h-10 rounded-lg bg-blue-500/10 border border-blue-500/20 flex items-center justify-center">
              <Settings className="w-5 h-5 text-blue-500" />
            </div>
            <code className="text-2xl font-bold text-white bg-black px-3 py-1 rounded-lg border border-white/[0.05] tracking-tight">config</code>
          </div>
          
          <div className="docs-prose">
            <p className="text-[17px]">
              The <code>config</code> command allows you to view and modify your gitresolve environment. It manages the <code>.gitresolve/policy.json</code> and global user preferences.
            </p>

            <TerminalWindow title="bash">
              <div className="space-y-4 text-[13px]">
                <div className="flex gap-3">
                  <span className="text-blue-500 font-bold">$</span>
                  <span className="text-white font-bold">gitresolve config --list</span>
                </div>
                <div className="font-mono pt-2">
                  <div className="text-white">core.editor=vim</div>
                  <div className="text-white">policy.default=balanced</div>
                  <div className="text-white">audit.enabled=true</div>
                  <div className="text-white">ui.theme=terminal</div>
                </div>
              </div>
            </TerminalWindow>

            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 mt-8">
              <FlagItem flag="--list" desc="Display all current configuration variables and their sources." />
              <FlagItem flag="--get <key>" desc="Retrieve the value for a specific configuration key." />
              <FlagItem flag="--set <key> <value>" desc="Set a configuration key to a specific value." />
              <FlagItem flag="--global" desc="Perform the operation on the global config instead of the local one." />
              <FlagItem flag="--edit" desc="Open the configuration file in your default system editor." />
            </div>
          </div>
        </section>
      </div>
    </DocsShell>
  );
}

function FlagItem({ flag, desc }) {
  return (
    <div className="p-6 rounded-xl bg-black border border-white/[0.05] hover-card group">
       <code className="text-blue-500 font-bold text-[14px] mb-2 block group-hover:text-white transition-colors tracking-tight">{flag}</code>
       <p className="text-[14px] text-[#a1a1aa] font-medium leading-relaxed">{desc}</p>
    </div>
  );
}
