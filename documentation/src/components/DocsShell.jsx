"use client";

import PageTransition from "./PageTransition";

export default function DocsShell({ title, subtitle, children }) {
  return (
    <PageTransition>
      <div className="mb-12">
        <div className="flex gap-2 mb-4">
          <span className="badge py-1 px-3 bg-white/[0.05] text-[#888] border border-white/[0.1] rounded-full text-xs font-bold uppercase tracking-wider">Documentation</span>
          <span className="badge py-1 px-3 bg-blue-500/10 text-blue-400 border border-blue-500/20 rounded-full text-xs font-bold uppercase tracking-wider">v1.4.0</span>
        </div>
        <h1 className="text-4xl md:text-5xl font-extrabold tracking-tight text-white mb-4 leading-tight">{title}</h1>
        {subtitle && <p className="text-[#a1a1aa] text-lg font-medium leading-relaxed max-w-2xl">{subtitle}</p>}
      </div>
      
      <div className="docs-prose">
        {children}
      </div>
    </PageTransition>
  );
}
