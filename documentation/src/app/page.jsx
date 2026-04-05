"use client";

import Link from "next/link";
import Image from "next/image";

const FeatureCard = ({ title, description, icon }) => (
  <div className="p-10 rounded-xl border border-[#222] bg-[#0a0a0a] hover:border-white transition-colors duration-300">
    <div className="w-9 h-9 rounded bg-[#111] border border-[#222] flex items-center justify-center mb-6 text-white group-hover:bg-white group-hover:text-black transition-colors scale-90">
      {icon}
    </div>
    <h3 className="text-[17px] font-semibold text-white mb-3 tracking-tight">{title}</h3>
    <p className="text-gray-400 text-[14px] leading-relaxed">
      {description}
    </p>
  </div>
);

export default function Home() {
  return (
    <div className="flex flex-col min-h-screen bg-black text-[#ededed] font-sans selection:bg-[#3b82f6]">
      <nav className="bg-black sticky top-0 z-50 border-b border-[#222]">
        <div className="w-full max-w-7xl mx-auto px-6 h-16 flex items-center justify-between">
          <Link href="/" className="font-bold tracking-tight text-[18px] text-white flex items-center group gap-3">
            <Image src="/logo.png" alt="gitresolve logo" width={24} height={24} className="rounded" />
            gitresolve
          </Link>
          <div className="hidden md:flex items-center gap-10 text-[14px] text-gray-500 font-medium">
            <Link href="/docs/installation" className="hover:text-white transition-colors">Documentation</Link>
            <a href="https://github.com/jhanvi857/gitresolve" target="_blank" className="hover:text-white transition-colors">GitHub</a>
            {/* <Link href="/enterprise" className="hover:text-white transition-colors">Enterprise</Link> */}
          </div>
          <a href="https://github.com/jhanvi857/gitresolve" target="_blank" className="bg-white text-black px-4 py-1.5 rounded-md text-[13px] font-semibold hover:bg-gray-200 transition-all active:scale-95">
            View Source Code
          </a>
        </div>
      </nav>

      <main className="flex-1 w-full flex flex-col items-center">
        <section className="min-h-[calc(100vh-64px)] w-full max-w-5xl mx-auto px-6 flex flex-col items-center justify-center text-center">
          {/* <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full border border-[#222] bg-[#0a0a0a] text-gray-500 text-[11px] font-bold mb-10 cursor-default tracking-widest uppercase">
            v1.2 Release: Deterministic Engine
          </div> */}
          <h1 className="text-3xl md:text-5xl font-semibold tracking-tight text-white leading-[1.2] mb-8 max-w-2xl px-2">
            Deterministic Conflict Resolution.
          </h1>
          {/* <Image src="/logo.png" alt="gitresolve logo" width={320} height={300} className="mb-8" /> */}
          <p className="max-w-xl text-base md:text-lg text-gray-400 mb-12 leading-relaxed">
            A privacy-first, purely offline engine using AST structural analysis to solve Git conflicts with mathematical precision. No LLMs, no hallucinations.
          </p>

          <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
            <Link href="/docs/installation" className="bg-white text-black font-semibold px-8 py-2.5 rounded hover:bg-gray-200 active:scale-95 transition-all text-[15px] border-2 border-white">
              Get Started
            </Link>
            <a href="https://github.com/jhanvi857/gitresolve" target="_blank" className="bg-black hover:bg-[#111] text-white font-semibold px-8 py-2.5 rounded transition-all text-[15px] border-2 border-[#222]">
              Star on GitHub
            </a>
          </div>
        </section>

        <section className="min-h-screen w-full border-t border-[#222] bg-black flex flex-col items-center justify-center">
          <div className="max-w-7xl mx-auto px-6 py-24">
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-12">
              <FeatureCard
                title="Zero Dependencies"
                description="Resolving conflicts via remote APIs is a security nightmare. Our engine runs 100% locally using tree-sitter based AST logic."
                icon={<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect><path d="M7 11V7a5 5 0 0 1 10 0v4"></path></svg>}
              />
              <FeatureCard
                title="AST Structural Intelligence"
                description="Detects semantic overlaps and structural mismatches. Perfectly handles duplicate imports and non-conflicting schema shifts."
                icon={<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round"><polyline points="16 18 22 12 16 6"></polyline><polyline points="8 6 2 12 8 18"></polyline></svg>}
              />
              <FeatureCard
                title="Atomic Reliability"
                description="Leverages POSIX atomic writes and a local SQLite session log. Revert any state change with a single command."
                icon={<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M23 4v6h-6"></path><path d="M20.49 15a9 9 0 1 1-2.12-9.36L23 10"></path></svg>}
              />
            </div>
          </div>
        </section>

        <section className="min-h-screen w-full border-t border-[#222] flex flex-col items-center justify-center py-24">
          <div className="max-w-4xl mx-auto w-full px-6 text-center">
            <div className="mb-12">
              <h2 className="text-2xl font-semibold text-white mb-4">Blazing Fast Go CLI.</h2>
              <p className="text-[15px] text-gray-500">Native performance, zero-dependency binary execution.</p>
            </div>

            <div className="code-window shadow-2xl text-left font-mono max-w-2xl mx-auto">
              <div className="code-header">
                <div className="dot dot-red" />
                <div className="dot dot-yellow" />
                <div className="dot dot-green" />
                <span className="ml-4 text-[9px] text-gray-500 uppercase tracking-[0.2em]">terminal — demo</span>
              </div>
              <div className="code-content whitespace-pre overflow-x-auto min-h-[220px]">
                <span className="token-output">$ </span><span className="token-cmd">gitresolve</span> <span className="token-arg">merge</span> <span className="token-flag">--dry-run</span><br />
                <span className="token-comment mt-2 block"># Scanning current conflicts...</span>
                <span className="token-output block">Found 2 unmerged files.</span>
                <br />
                <span className="token-output block">--- Processing pkg/auth/login.go ---</span>
                <span className="token-severity-high block">Severity 9: Complex Logic Overlap detected.</span>
                <span className="token-comment block"># Escalating to manual-resolve.</span>
                <br />
                <span className="token-output block">--- Processing web/config.yaml ---</span>
                <span className="token-severity-low block">Successfully auto-resolved 100% of blocks.</span>
                <br />
                <span className="token-output mt-4 block">1 file resolved atomically.</span>
              </div>
            </div>
          </div>
        </section>

        <section className="min-h-screen w-full border-t border-[#222] bg-[#050505] flex flex-col items-center justify-center">
          <div className="max-w-4xl mx-auto px-6 text-center">
            <Image src="/logo.png" alt="gitresolve logo" width={60} height={60} className="mb-8 mx-auto" />
            <h2 className="text-2xl md:text-3xl font-semibold text-white mb-6 tracking-tight">
              Stop fighting Git markers.
            </h2>
            <p className="text-base text-gray-400 mb-10 leading-relaxed max-w-xl mx-auto px-4">
              Experience the power of structural conflict resolution.
              Purely offline, extremely predictable.
            </p>
            <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
              <Link href="/docs/installation" className="bg-white text-black font-semibold px-8 py-2.5 rounded text-[15px] hover:bg-gray-200 transition-all active:scale-95">
                Go to Documentation
              </Link>
            </div>
          </div>
        </section>
      </main>

      <footer className="w-full py-16 px-6 border-t border-[#222] bg-black">
        <div className="max-w-7xl mx-auto grid grid-cols-1 md:grid-cols-4 gap-12 text-center md:text-left">
          <div className="md:col-span-2">
            <Link href="/" className="font-bold tracking-tight text-[18px] text-white flex items-center justify-center md:justify-start mb-6">
              <Image src="/logo.png" alt="gitresolve" width={20} height={20} className="mr-3" />
              gitresolve
            </Link>
            <p className="text-gray-500 max-w-sm leading-relaxed text-[15px] mx-auto md:mx-0">
              Computational conflict resolution for high-performance engineering systems.
            </p>
          </div>
          <div>
            <h4 className="text-white font-semibold mb-6 text-xs uppercase tracking-widest px-2">Resources</h4>
            <div className="flex flex-col gap-3 text-gray-500 font-medium text-[14px]">
              <Link href="/docs/installation" className="hover:text-white transition-colors">Documentation</Link>
              <Link href="/docs/architecture" className="hover:text-white transition-colors">Architecture</Link>
              <Link href="/docs/security" className="hover:text-white transition-colors">Security</Link>
              <Link href="/docs/commands" className="hover:text-white transition-colors">CLI reference</Link>
            </div>
          </div>
          {/* <div>
            <h4 className="text-white font-semibold mb-6 text-xs uppercase tracking-widest px-2">Legal</h4>
            <div className="flex flex-col gap-3 text-gray-500 font-medium tracking-tight text-[14px]">
              <a href="#" className="hover:text-white transition-colors">Privacy Policy</a>
              <a href="#" className="hover:text-white transition-colors">Terms of Service</a>
              <a href="#" className="hover:text-white transition-colors">Enterprise</a>
            </div>
          </div> */}
        </div>
        <div className="max-w-7xl mx-auto mt-16 pt-8 border-t border-[#222] flex flex-col md:flex-row justify-between items-center gap-6 text-gray-600 text-[11px] font-medium uppercase tracking-[0.05em]">
          <p>© 2026 gitresolve. Purely offline resolution.</p>
          <div className="flex gap-10">
            <a href="https://github.com/jhanvi857/gitresolve" target="_blank" className="hover:text-white transition-colors">GitHub</a>
            <a href="#" className="hover:text-white transition-colors">Discord</a>
            <a href="#" className="hover:text-white transition-colors">Twitter</a>
          </div>
        </div>
      </footer>
    </div>
  );
}