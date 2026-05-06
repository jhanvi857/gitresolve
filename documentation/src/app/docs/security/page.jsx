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
              <p className="text-sm text-gray-500">The tool contains no networking code. It cannot send your code to a remote server because it doesn&apos;t know how to talk to the internet.</p>
            </div>
            <div className="p-6 rounded-lg bg-[#0a0a0a] border border-[#222]">
              <div className="w-8 h-8 rounded bg-green-500/10 flex items-center justify-center mb-4">
                <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#10b981" strokeWidth="2"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"></path></svg>
              </div>
              <h3 className="text-white font-medium mb-2">CWE-22 Sandboxing</h3>
              <p className="text-sm text-gray-500">Every file operation uses <code className="text-gray-400">os.Root</code> to ensure no read/write can escape the repository root. Path traversal is mathematically impossible.</p>
            </div>
            <div className="p-6 rounded-lg bg-[#0a0a0a] border border-[#222]">
              <div className="w-8 h-8 rounded bg-orange-500/10 flex items-center justify-center mb-4">
                <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#f97316" strokeWidth="2"><path d="M13 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V9z"></path><polyline points="13 2 13 9 20 9"></polyline></svg>
              </div>
              <h3 className="text-white font-medium mb-2">DoS Protection</h3>
              <p className="text-sm text-gray-500">A mandatory 10MB gate prevents memory exhaustion. Maliciously oversized conflict files are skipped and escalated to manual review.</p>
            </div>
            <div className="p-6 rounded-lg bg-[#0a0a0a] border border-[#222]">
              <div className="w-8 h-8 rounded bg-red-500/10 flex items-center justify-center mb-4">
                <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#ef4444" strokeWidth="2"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path></svg>
              </div>
              <h3 className="text-white font-medium mb-2">Advisory Locking</h3>
              <p className="text-sm text-gray-500">Uses native <code className="text-gray-400">flock(2)</code> and <code className="text-gray-400">LockFileEx</code>. Safe from PID-reuse attacks and race conditions in concurrent CI pipes.</p>
            </div>
          </div>
        </section>

         <section>
           <h2 className="text-xl font-semibold text-white mb-4">Integrity & Privacy</h2>
            <p className="text-gray-400 leading-relaxed mb-6">
              gitresolve uses a multi-stage verification and privacy process to ensure that resolutions are not just &quot;done&quot; but &quot;correct&quot; and &quot;private&quot;.
            </p>
           <ul className="space-y-4">
             <li className="flex gap-4">
                <span className="text-blue-500 font-mono">01.</span>
                <div>
                  <h4 className="text-white font-medium">PII Privacy (Hashing)</h4>
                  <p className="text-sm text-gray-500">Sensitive file content or conflict blocks are never stored in plain text in debug logs. We use 12-char SHA-256 hashes for event correlation.</p>
                </div>
             </li>
             <li className="flex gap-4">
                <span className="text-blue-500 font-mono">02.</span>
                <div>
                  <h4 className="text-white font-medium">Syntax Validation</h4>
                  <p className="text-sm text-gray-500">The merged code is passed through a language-specific syntax checker. If the merge creates invalid syntax, the operation is rolled back.</p>
                </div>
             </li>
             <li className="flex gap-4">
                <span className="text-blue-500 font-mono">03.</span>
                <div>
                  <h4 className="text-white font-medium">Supply Chain Security</h4>
                  <p className="text-sm text-gray-500">All releases are signed via Cosign (OIDC) and include a CycloneDX SBOM. Binaries are verifiable against the public Rekor transparency log.</p>
                </div>
             </li>
           </ul>
         </section>

        <section>
          <h2 className="text-xl font-semibold text-white mb-4">Vulnerability Disclosure</h2>
          <p className="text-gray-400 leading-relaxed">
            Security issues should be reported privately according to the process in <code className="text-gray-300">SECURITY.md</code> at the repository root. Avoid opening public issues for unpatched vulnerabilities.
          </p>
        </section>
      </div>
    </DocsShell>
  );
}
