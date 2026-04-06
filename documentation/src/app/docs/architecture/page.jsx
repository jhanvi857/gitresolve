"use client";

import React from 'react';
import DocsShell from '@/components/DocsShell';

export default function Architecture() {
  return (
    <DocsShell
      title="Architecture"
      subtitle="How gitresolve handles conflicts without a centralized brain."
    >
      <div className="space-y-12">
        <section>
          <h2 className="text-xl font-semibold text-white mb-4">Core Principles</h2>
          <p className="text-gray-400 leading-relaxed mb-6">
            gitresolve is built on the principle of <strong>Deterministic Resolution</strong>. Unlike AI-based tools that might provide different answers for the same input, gitresolve uses a fixed engine with zero side effects.
          </p>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div className="p-6 rounded-lg bg-[#0a0a0a] border border-[#222]">
              <h3 className="text-white font-medium mb-2">AST Analysis</h3>
              <p className="text-sm text-gray-500">Instead of lines of text, we see structures. This allows us to safely merge import blocks and detect function signature changes.</p>
            </div>
            <div className="p-6 rounded-lg bg-[#0a0a0a] border border-[#222]">
              <h3 className="text-white font-medium mb-2">Safety First</h3>
              <p className="text-sm text-gray-500">Every resolution is verified against the language's grammar. If it doesn't parse, it doesn't merge.</p>
            </div>
          </div>
        </section>

        <section>
          <h2 className="text-xl font-semibold text-white mb-4">The Pipeline</h2>
          <div className="relative pl-8 border-l border-[#222] space-y-8">
            <div className="relative">
              <div className="absolute -left-[37px] top-1 w-4 h-4 rounded-full bg-blue-500 shadow-[0_0_10px_rgba(59,130,246,0.5)]"></div>
              <h4 className="text-white font-medium mb-1">1. Block Identification</h4>
              <p className="text-sm text-gray-500">Regex-based scanning identifying {"<<<<<<<"}, {"======="}, and {">>>>>>>"} markers.</p>
            </div>
            <div className="relative">
              <div className="absolute -left-[37px] top-1 w-4 h-4 rounded-full bg-[#333]"></div>
              <h4 className="text-white font-medium mb-1">2. Heuristic Classification</h4>
              <p className="text-sm text-gray-500">Categorizing blocks into: Whitespace, Identical, Imports, Structured, or Logic.</p>
            </div>
            <div className="relative">
              <div className="absolute -left-[37px] top-1 w-4 h-4 rounded-full bg-[#333]"></div>
              <h4 className="text-white font-medium mb-1">3. Semantic Resolution</h4>
              <p className="text-sm text-gray-500">Applying language-specific rules (e.g., Go AST import deduplication).</p>
            </div>
            <div className="relative">
              <div className="absolute -left-[37px] top-1 w-4 h-4 rounded-full bg-[#333]"></div>
              <h4 className="text-white font-medium mb-1">4. Verification Gate</h4>
              <p className="text-sm text-gray-500">Running the language parser (e.g., `go/parser`) on the resulting code snippet.</p>
            </div>
          </div>
        </section>

        {/* IMAGE PLACEHOLDER 1 */}
        <section className="pt-10">
          <div className="w-full aspect-video rounded-xl border-2 border-dashed border-[#222] bg-[#050505] flex flex-col items-center justify-center text-center p-10">
            <div className="w-16 h-16 rounded-full bg-[#111] flex items-center justify-center mb-4 text-gray-600">
              <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round"><rect x="3" y="3" width="18" height="18" rx="2" ry="2"></rect><line x1="9" y1="3" x2="9" y2="21"></line></svg>
            </div>
            <p className="text-white font-medium mb-1">Architecture Visualization Placeholder</p>
            <p className="text-sm text-gray-500 max-w-xs">Detailed diagram of the deterministic engine resolving a complex Go import conflict across three branches.</p>
          </div>
        </section>
      </div>
    </DocsShell>
  );
}
