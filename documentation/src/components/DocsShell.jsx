"use client";

import PageTransition from "./PageTransition";

export default function DocsShell({ title, subtitle, children }) {
  return (
    <PageTransition>
      <header className="page-header">
        <div className="page-meta-row flex gap-3 mb-6">
          <span className="page-meta-chip text-[10px] font-bold uppercase tracking-widest px-2 py-1 bg-[#111] border border-[#222] rounded text-[#888]">Docs v1.0.2</span>
          <span className="page-meta-chip text-[10px] font-bold uppercase tracking-widest px-2 py-1 bg-[#111] border border-[#222] rounded text-blue-500/50">Determinism Guaranteed</span>
        </div>
        <h1 className="text-4xl md:text-5xl font-bold tracking-tighter text-white mb-6 underline decoration-blue-500/10 decoration-8 underline-offset-[-2px]">{title}</h1>
        {subtitle ? <p className="text-gray-400 text-lg md:text-xl leading-relaxed max-w-2xl">{subtitle}</p> : null}
      </header>
      
      <div className="mt-16 prose-layout font-inter text-[#999] leading-relaxed text-[16px]">
        {children}
      </div>
    </PageTransition>
  );
}
