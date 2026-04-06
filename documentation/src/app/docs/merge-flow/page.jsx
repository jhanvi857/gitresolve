"use client";

import React from 'react';
import DocsShell from '@/components/DocsShell';

export default function MergeFlow() {
  return (
    <DocsShell 
      title="Merge Flow" 
      subtitle="How gitresolve handles different types of conflicts across multi-branch merges."
    >
      <div className="space-y-12">
        <section>
          <h2 className="text-xl font-semibold text-white mb-6">Automated Triage</h2>
          <p className="text-gray-400 mb-8 leading-relaxed">
            Conflicts are triaged based on confidence. If the tool is 100% certain (e.g., whitespace), it auto-resolves. If it's 50% certain, it prompts for human feedback.
          </p>
          <div className="space-y-4">
             <div className="p-4 rounded border border-[#222] bg-[#0a0a0a]">
                <h3 className="text-blue-400 text-sm font-semibold mb-2 uppercase tracking-wide">Category: Trivial</h3>
                <p className="text-sm text-gray-500">Whitespace shifts, identical line changes, or purely comment updates. <strong>Auto-Resolution: Yes.</strong></p>
             </div>
             <div className="p-4 rounded border border-[#222] bg-[#0a0a0a]">
                <h3 className="text-green-400 text-sm font-semibold mb-2 uppercase tracking-wide">Category: Structured</h3>
                <p className="text-sm text-gray-500">JSON, YAML, or TOML edits that don't overlap on the same key. <strong>Auto-Resolution: Yes (Semantic).</strong></p>
             </div>
             <div className="p-4 rounded border border-[#222] bg-[#0a0a0a]">
                <h3 className="text-yellow-400 text-sm font-semibold mb-2 uppercase tracking-wide">Category: Imports</h3>
                <p className="text-sm text-gray-500">Additions of new packages or libraries across branches. <strong>Auto-Resolution: Yes (Deduplicated).</strong></p>
             </div>
             <div className="p-4 rounded border border-red-900/50 bg-[#0a0a0a]">
                <h3 className="text-red-400 text-sm font-semibold mb-2 uppercase tracking-wide">Category: Logic</h3>
                <p className="text-sm text-gray-500">Function signature changes, deletion of used logic, or conflicting business rules. <strong>Auto-Resolution: No.</strong></p>
             </div>
          </div>
        </section>

        {/* IMAGE PLACEHOLDER 2 */}
        <section className="pt-10">
          <div className="w-full aspect-video rounded-xl border-2 border-dashed border-[#222] bg-[#050505] flex flex-col items-center justify-center text-center p-10">
            <div className="w-16 h-16 rounded-full bg-[#111] flex items-center justify-center mb-4 text-gray-600">
               <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round"><path d="M7 11V7a5 5 0 0 1 10 0v4"></path><rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect></svg>
            </div>
            <p className="text-white font-medium mb-1">Merge Workflow Placeholder</p>
            <p className="text-sm text-gray-500 max-w-xs">Detailed illustration of a three-way merge resolution process for overlapping YAML blocks.</p>
          </div>
        </section>
      </div>
    </DocsShell>
  );
}
