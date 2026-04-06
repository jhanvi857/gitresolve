"use client";

import Link from "next/link";
import Image from "next/image";
import { motion, useScroll, useTransform } from "framer-motion";
import { useRef } from "react";

const FeatureCard = ({ title, description, icon, index }) => (
  <motion.div 
    initial={{ opacity: 0, y: 20 }}
    whileInView={{ opacity: 1, y: 0 }}
    viewport={{ once: true }}
    transition={{ delay: index * 0.1, duration: 0.5, ease: "easeOut" }}
    className="p-10 rounded-xl border border-[#222] bg-[#0a0a0a] hover:border-white hover:bg-[#0d0d0d] transition-all duration-500 group relative overflow-hidden"
  >
    <div className="absolute top-0 left-0 w-full h-1 bg-gradient-to-r from-transparent via-[#333] to-transparent opacity-0 group-hover:opacity-100 transition-opacity" />
    <div className="w-10 h-10 rounded bg-[#111] border border-[#222] flex items-center justify-center mb-6 text-white group-hover:border-white group-hover:scale-110 transition-all duration-300">
      {icon}
    </div>
    <h3 className="text-[19px] font-semibold text-white mb-3 tracking-tight">{title}</h3>
    <p className="text-gray-400 text-[15px] leading-relaxed">
      {description}
    </p>
  </motion.div>
);

export default function Home() {
  const targetRef = useRef(null);
  const { scrollYProgress } = useScroll({
    target: targetRef,
    offset: ["start end", "end start"],
  });

  const backgroundY = useTransform(scrollYProgress, [0, 1], ["0%", "20%"]);

  return (
    <div className="flex flex-col min-h-screen bg-black text-[#ededed] font-sans selection:bg-[#3b82f6]">
      <nav className="bg-black/80 backdrop-blur-md sticky top-0 z-50 border-b border-[#222]">
        <div className="w-full max-w-7xl mx-auto px-6 h-16 flex items-center justify-between">
          <Link href="/" className="font-bold tracking-tight text-[18px] text-white flex items-center group gap-3">
             <div className="p-1 rounded bg-[#111] border border-[#333] group-hover:border-white transition-colors">
               <Image src="/logo.png" alt="gitresolve logo" width={20} height={20} className="rounded-sm" />
             </div>
            gitresolve
          </Link>
          <div className="hidden md:flex items-center gap-10 text-[14px] text-gray-500 font-medium">
            <Link href="/docs/installation" className="hover:text-white transition-colors">Documentation</Link>
            <Link href="/docs/architecture" className="hover:text-white transition-colors">Architecture</Link>
            <a href="https://github.com/jhanvi857/gitresolve" target="_blank" className="hover:text-white transition-colors">GitHub</a>
          </div>
          <a href="https://github.com/jhanvi857/gitresolve" target="_blank" className="bg-white text-black px-4 py-1.5 rounded-md text-[13px] font-semibold hover:bg-gray-200 transition-all active:scale-95 shadow-[0_0_15px_rgba(255,255,255,0.1)]">
            View Source Code
          </a>
        </div>
      </nav>

      <main className="flex-1 w-full flex flex-col items-center">
        {/* HERO SECTION */}
        <section ref={targetRef} className="relative min-h-[calc(100vh-64px)] w-full overflow-hidden flex flex-col items-center justify-center text-center">
          <motion.div 
             style={{ y: backgroundY }}
             className="absolute inset-0 z-0 opacity-20 pointer-events-none"
          >
             <div className="absolute top-1/4 left-1/2 -translate-x-1/2 w-[800px] h-[500px] bg-blue-600/20 blur-[120px] rounded-full" />
             <div className="absolute bottom-1/4 left-1/3 -translate-x-1/2 w-[600px] h-[400px] bg-purple-600/10 blur-[100px] rounded-full" />
          </motion.div>

          <div className="relative z-10 max-w-5xl mx-auto px-6 flex flex-col items-center">
            <motion.div
              initial={{ opacity: 0, scale: 0.95 }}
              animate={{ opacity: 1, scale: 1 }}
              transition={{ duration: 0.8, ease: [0.23, 1, 0.32, 1] }}
            >
              <h1 className="text-4xl md:text-6xl font-bold tracking-tighter text-white leading-[1.1] mb-8 max-w-2xl px-2">
                Deterministic <span className="text-transparent bg-clip-text bg-gradient-to-r from-white via-gray-400 to-gray-600">Conflict Resolution.</span>
              </h1>
            </motion.div>

            <motion.p 
              initial={{ opacity: 0, y: 10 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.6, delay: 0.2 }}
              className="max-w-xl text-lg md:text-xl text-gray-400 mb-12 leading-relaxed"
            >
              A privacy-first, purely offline engine using AST structural analysis to solve Git conflicts with mathematical precision. No LLMs, no hallucinations.
            </motion.p>

            <motion.div 
              initial={{ opacity: 0, y: 10 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.6, delay: 0.4 }}
              className="flex flex-col sm:flex-row items-center justify-center gap-4"
            >
              <Link href="/docs/installation" className="group relative bg-white text-black font-semibold px-10 py-3 rounded hover:bg-gray-200 active:scale-95 transition-all text-[16px]">
                <span className="relative z-10">Get Started</span>
                <div className="absolute inset-0 bg-white blur-md opacity-0 group-hover:opacity-20 transition-opacity" />
              </Link>
              <a href="https://github.com/jhanvi857/gitresolve" target="_blank" className="bg-black hover:bg-[#111] text-white font-semibold px-10 py-3 rounded transition-all text-[16px] border border-[#222] hover:border-[#444]">
                Star on GitHub
              </a>
            </motion.div>
          </div>

          <motion.div 
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 1, duration: 1 }}
            className="absolute bottom-10 left-1/2 -translate-x-1/2 flex flex-col items-center gap-4"
          >
             <span className="text-[10px] uppercase tracking-[0.3em] font-bold text-gray-600">Scroll to explore</span>
             <div className="w-px h-12 bg-gradient-to-b from-[#222] to-transparent" />
          </motion.div>
        </section>

        {/* FEATURES GRID */}
        <section className="w-full border-t border-[#222] bg-black">
          <div className="max-w-7xl mx-auto px-6 py-32">
            <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
              <FeatureCard
                index={0}
                title="Zero Dependencies"
                description="Resolving conflicts via remote APIs is a security nightmare. Our engine runs 100% locally using tree-sitter based AST logic."
                icon={<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect><path d="M7 11V7a5 5 0 0 1 10 0v4"></path></svg>}
              />
              <FeatureCard
                index={1}
                title="AST Intelligence"
                description="Detects semantic overlaps and structural mismatches. Perfectly handles duplicate imports and non-conflicting schema shifts."
                icon={<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round"><polyline points="16 18 22 12 16 6"></polyline><polyline points="8 6 2 12 8 18"></polyline></svg>}
              />
              <FeatureCard
                index={2}
                title="Atomic Reliability"
                description="Leverages POSIX atomic writes and a local SQLite session log. Revert any state change with a single command."
                icon={<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M23 4v6h-6"></path><path d="M20.49 15a9 9 0 1 1-2.12-9.36L23 10"></path></svg>}
              />
            </div>
          </div>
        </section>

        {/* TERMINAL DEMO */}
        <section className="w-full border-t border-[#222] flex flex-col items-center justify-center py-32">
          <div className="max-w-4xl mx-auto w-full px-6 text-center">
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              className="mb-16"
            >
              <h2 className="text-3xl font-bold text-white mb-4 tracking-tight">Blasing Fast CLI.</h2>
              <p className="text-[16px] text-gray-500">Native performance, zero-dependency binary execution.</p>
            </motion.div>

            <motion.div 
               initial={{ opacity: 0, scale: 0.98 }}
               whileInView={{ opacity: 1, scale: 1 }}
               viewport={{ once: true }}
               transition={{ duration: 0.8 }}
               className="code-window shadow-[0_30px_60px_-15px_rgba(0,0,0,0.5)] text-left font-mono max-w-2xl mx-auto backdrop-blur-sm border-[#333]"
            >
              <div className="code-header border-b border-[#222] bg-[#0c0c0c]/80">
                <div className="flex gap-2 mr-4">
                  <div className="w-3 h-3 rounded-full bg-[#333]" />
                  <div className="w-3 h-3 rounded-full bg-[#333]" />
                  <div className="w-3 h-3 rounded-full bg-[#333]" />
                </div>
                <span className="text-[10px] text-gray-500 uppercase tracking-widest font-bold">terminal — gitresolve</span>
              </div>
              <div className="code-content whitespace-pre overflow-x-auto min-h-[240px] p-6 text-[13px] leading-relaxed">
                <span className="text-blue-500">$</span> <span className="text-white">gitresolve resolve</span><br />
                <span className="text-gray-600 block mt-3">Analyzing dependency graph...</span>
                <span className="text-white block">Found 3 unmerged files.</span>
                <br />
                <div className="space-y-1">
                  <div className="flex gap-4">
                    <span className="text-gray-600 shrink-0">PARSE</span>
                    <span className="text-white">internal/net/socket.go</span>
                  </div>
                  <div className="flex gap-4">
                    <span className="text-gray-600 shrink-0">INDEX</span>
                    <span className="text-green-500">Auto-resolved: Duplicate Imports</span>
                  </div>
                </div>
                <br />
                <div className="space-y-1">
                  <div className="flex gap-4">
                    <span className="text-gray-600 shrink-0">PARSE</span>
                    <span className="text-white">pkg/auth/login_flow.go</span>
                  </div>
                  <div className="flex gap-4">
                    <span className="text-gray-600 shrink-0">ERROR</span>
                    <span className="text-red-500">Critical: Logic Overlap detected in Auth Handler</span>
                  </div>
                </div>
                <br />
                <span className="text-gray-500 mt-4 block underline underline-offset-4 pointer-events-none">Scan complete. 1 file resolved, 2 require manual review.</span>
              </div>
            </motion.div>
          </div>
        </section>

        {/* CTA Page SECTION */}
        <section className="relative w-full border-t border-[#222] bg-[#050505] py-40 overflow-hidden text-center flex flex-col items-center">
          <div className="absolute top-0 left-1/2 -translate-x-1/2 w-[1000px] h-[400px] bg-blue-900/10 blur-[150px] rounded-full pointer-events-none" />
          
          <motion.div 
            initial={{ opacity: 0, y: 30 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            className="relative z-10 max-w-3xl mx-auto px-6"
          >
            <div className="w-16 h-16 rounded-xl bg-black border border-[#222] flex items-center justify-center mb-10 mx-auto shadow-xl group hover:border-white transition-colors duration-500">
              <Image src="/logo.png" alt="gitresolve logo" width={32} height={32} className="grayscale group-hover:grayscale-0 transition-all duration-500" />
            </div>
            
            <h2 className="text-4xl md:text-5xl font-bold text-white mb-8 tracking-tighter">
              Stop fighting <span className="underline decoration-blue-500/30 decoration-8 underline-offset-[-2px]">Git markers.</span>
            </h2>
            
            <p className="text-lg text-gray-400 mb-12 leading-relaxed max-w-xl mx-auto">
              Experience the power of structural conflict resolution.
              Purely offline, extremely predictable, built for performance.
            </p>
            
            <motion.div 
              whileHover={{ scale: 1.02 }}
              whileTap={{ scale: 0.98 }}
              className="flex flex-col sm:flex-row items-center justify-center gap-6"
            >
              <Link href="/docs/installation" className="relative group bg-white text-black font-bold px-12 py-4 rounded-md text-[17px] shadow-[0_0_30px_rgba(255,255,255,0.1)] transition-all">
                Go to Documentation
                <div className="absolute inset-0 bg-white blur-xl opacity-0 group-hover:opacity-20 transition-opacity" />
              </Link>
            </motion.div>
          </motion.div>
        </section>
      </main>

      <footer className="w-full py-20 px-6 border-t border-[#222] bg-black">
        <div className="max-w-7xl mx-auto grid grid-cols-1 md:grid-cols-4 gap-16 text-center md:text-left">
          <div className="md:col-span-2">
            <Link href="/" className="font-bold tracking-tighter text-[20px] text-white flex items-center justify-center md:justify-start mb-8 transition-opacity hover:opacity-80">
              <Image src="/logo.png" alt="gitresolve" width={24} height={24} className="mr-4 rounded-sm" />
              gitresolve
            </Link>
            <p className="text-gray-500 max-w-sm leading-relaxed text-[15px] mx-auto md:mx-0">
              Computational conflict resolution for high-performance engineering systems. Determinism by design.
            </p>
          </div>
          <div>
            <h4 className="text-white font-bold mb-8 text-[11px] uppercase tracking-[0.2em] px-2 text-gray-400">Documentation</h4>
            <div className="flex flex-col gap-4 text-gray-500 font-medium text-[15px]">
              <Link href="/docs/installation" className="hover:text-white transition-all transform hover:translate-x-1">Get Started</Link>
              <Link href="/docs/architecture" className="hover:text-white transition-all transform hover:translate-x-1">Architecture</Link>
              <Link href="/docs/merge-flow" className="hover:text-white transition-all transform hover:translate-x-1">Merge Flow</Link>
              <Link href="/docs/commands" className="hover:text-white transition-all transform hover:translate-x-1">Commands</Link>
            </div>
          </div>
          <div>
            <h4 className="text-white font-bold mb-8 text-[11px] uppercase tracking-[0.2em] px-2 text-gray-400">Resources</h4>
            <div className="flex flex-col gap-4 text-gray-500 font-medium text-[15px]">
              <a href="https://github.com/jhanvi857/gitresolve" target="_blank" className="hover:text-white transition-all transform hover:translate-x-1">Source Code</a>
              <a href="#" className="hover:text-white transition-all transform hover:translate-x-1">Discord Support</a>
              <a href="#" className="hover:text-white transition-all transform hover:translate-x-1">Changelog</a>
            </div>
          </div>
        </div>
        <div className="max-w-7xl mx-auto mt-24 pt-10 border-t border-[#111] flex flex-col md:flex-row justify-between items-center gap-8 text-gray-600 text-[10px] font-bold uppercase tracking-[0.2em]">
          <p>© 2026 gitresolve. Purely offline resolution.</p>
          <div className="flex gap-12">
            <a href="https://github.com/jhanvi857/gitresolve" target="_blank" className="hover:text-white transition-colors">GitHub</a>
            <a href="#" className="hover:text-white transition-colors">Twitter</a>
            <a href="#" className="hover:text-white transition-colors">LinkedIn</a>
          </div>
        </div>
      </footer>
    </div>
  );
}