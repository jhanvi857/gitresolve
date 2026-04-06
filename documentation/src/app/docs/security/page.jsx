"use client";

import React from 'react';
import DocsShell from '@/components/DocsShell';

export default function Security() {
  return (
    <DocsShell 
      title="Security & Privacy" 
      subtitle="How we protect your codebase by staying purely offline."
    >
      <div className="space-y-12">
        <section>
          <h2 className="text-xl font-semibold text-white mb-4">Zero Data Leakage</h2>
          <p className="text-gray-400 leading-relaxed mb-6">
            In an era of AI-driven tools, your source code is often treated as training data. gitresolve takes a different path. We believe your code is your most valuable asset.
          </p>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
            <div className="p-6 rounded-lg bg-[#0a0a0a] border border-[#222]">
              <div className="w-8 h-8 rounded bg-blue-500/10 flex items-center justify-center mb-4">
                <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#3b82f6" strokeWidth="2"><rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect><path d="M7 11V7a5 5 0 0 1 10 0v4"></path></svg>
              </div>
              <h3 className="text-white font-medium mb-2">100% Offline</h3>
              <p className="text-sm text-gray-500">The tool contains no networking code. It cannot send your code to a remote server because it doesn't know how to talk to the internet.</p>
            </div>
            <div className="p-6 rounded-lg bg-[#0a0a0a] border border-[#222]">
              <div className="w-8 h-8 rounded bg-green-500/10 flex items-center justify-center mb-4">
                <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#10b981" strokeWidth="2"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"></path></svg>
              </div>
              <h3 className="text-white font-medium mb-2">AST Parsing</h3>
              <p className="text-sm text-gray-500">We analyze the structure of your code using local parsers. No heuristics that "guess" your logic; only mathematical structural matching.</p>
            </div>
          </div>
        </section>

        <section>
          <h2 className="text-xl font-semibold text-white mb-4">Integrity Verification</h2>
          <p className="text-gray-400 leading-relaxed mb-6">
            gitresolve uses a two-stage verification process to ensure that resolutions are not just "done" but "correct".
          </p>
          <ul className="space-y-4">
            <li className="flex gap-4">
               <span className="text-blue-500 font-mono">01.</span>
               <div>
                 <h4 className="text-white font-medium">No-Marker Guarantee</h4>
                 <p className="text-sm text-gray-500">The final write operation is blocked if any byte sequence matching a Git conflict marker is detected in the output stream.</p>
               </div>
            </li>
            <li className="flex gap-4">
               <span className="text-blue-500 font-mono">02.</span>
               <div>
                 <h4 className="text-white font-medium">Syntax Validation</h4>
                 <p className="text-sm text-gray-500">The merged code is passed through a language-specific syntax checker. If the merge creates invalid syntax, the operation is rolled back.</p>
               </div>
            </li>
          </ul>
        </section>
      </div>
    </DocsShell>
  );
}
