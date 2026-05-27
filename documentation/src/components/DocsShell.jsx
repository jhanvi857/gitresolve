"use client";

import PageTransition from "./PageTransition";

export default function DocsShell({ title, subtitle, children }) {
  return (
    <PageTransition>
      <div className="mb-12">
        <div className="flex gap-2 mb-4">
          <span className="badge py-1 px-3">Documentation</span>
          <span className="badge py-1 px-3 bg-green-500/10 text-green-500 border-green-500/20">v1.4.0</span>
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
