"use client";

import Link from "next/link";
import Image from "next/image";
import { motion } from "framer-motion";
import { Terminal, Shield, Cpu, Activity, Zap, ChevronRight } from "lucide-react";
import TerminalWindow from "@/components/TerminalWindow";
import Footer from "@/components/Footer";

const FeatureCard = ({ title, description, icon: Icon, index }) => (
  <motion.div
    initial={{ opacity: 0, y: 20 }}
    whileInView={{ opacity: 1, y: 0 }}
    viewport={{ once: true }}
    transition={{ delay: index * 0.1, duration: 0.5 }}
    className="p-8 rounded-xl hover-card bg-black group"
  >
    <div className="w-10 h-10 rounded-lg bg-[#111] border border-[#222] flex items-center justify-center mb-6 group-hover:border-blue-500/50 transition-colors">
      <Icon className="w-5 h-5 text-white group-hover:text-blue-500 transition-colors" />
    </div>
    <h3 className="text-xl font-bold text-white mb-3 tracking-tight">{title}</h3>
    <p className="text-[#a1a1aa] text-[15px] leading-relaxed">
      {description}
    </p>
  </motion.div>
);

export default function Home() {
  return (
    <div className="flex flex-col min-h-screen bg-black text-white selection:bg-blue-500/30 font-sans">
      {/* Grid Background */}
      <div className="fixed inset-0 grid-bg opacity-10 pointer-events-none z-0" />

      {/* Radial Gradient Overlay */}
      <div className="fixed inset-0 bg-[radial-gradient(circle_at_50%_-20%,rgba(0,112,243,0.05),transparent_70%)] pointer-events-none z-0" />

      <nav className="sticky top-0 z-50 border-b border-white/[0.05] bg-black/80 backdrop-blur-xl">
        <div className="max-w-7xl mx-auto px-8 h-16 flex items-center justify-between">
          <Link href="/" className="flex items-center gap-3 group">
            <div className="p-1.5 rounded-lg bg-black border border-white/[0.1] group-hover:border-blue-500 transition-all">
              <Image src="/logo.png" alt="logo" width={20} height={20} className="opacity-90" />
            </div>
            <span className="font-extrabold tracking-tighter text-xl">gitresolve</span>
          </Link>
          <div className="hidden md:flex items-center gap-10 text-[14px] font-bold text-[#555]">
            <Link href="/docs/installation" className="hover:text-white transition-colors">Docs</Link>
            <Link href="/docs/architecture" className="hover:text-white transition-colors">Architecture</Link>
            <a href="https://github.com/jhanvi857/gitresolve" target="_blank" className="hover:text-white transition-colors">GitHub</a>
          </div>
          <Link href="/docs/installation" className="bg-white text-black px-6 py-2 rounded-lg text-[14px] font-bold hover:bg-[#e1e1e1] transition-all">
            Read Docs
          </Link>
        </div>
      </nav>

      <main className="relative z-10">
        {/* Hero Section */}
        <section className="pt-24 md:pt-32 pb-20 px-8">
          <div className="max-w-7xl mx-auto grid grid-cols-1 lg:grid-cols-2 gap-16 items-center">
            <motion.div
              initial={{ opacity: 0, x: -20 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ duration: 0.6 }}
              className="text-left"
            >
              <div className="inline-flex items-center gap-3 px-4 py-2 rounded-full bg-[#111] border border-white/[0.05] mb-8">
                <div className="w-2 h-2 rounded-full bg-blue-500 shadow-[0_0_8px_rgba(0,112,243,0.8)]" />
                <span className="text-[12px] font-bold text-[#888] tracking-wide">
                  Deterministic Conflict Resolution <span className="mx-2 text-[#333]">|</span> v1.4.0 Stable
                </span>
              </div>
              
              <h1 className="text-5xl md:text-6xl font-extrabold tracking-tight mb-8 leading-[1.1]">
                GitResolve: A Precision <br />
                <span className="text-white opacity-90">Conflict Engine.</span>
              </h1>
              
              <p className="text-[#a1a1aa] text-lg md:text-xl max-w-xl mb-10 leading-relaxed font-medium">
                GitResolve is a privacy-first, purely offline engine using AST structural analysis to solve Git conflicts with mathematical precision and zero data leakage.
              </p>
              
              <div className="flex flex-col sm:flex-row items-center gap-4">
                <Link href="/docs/architecture" className="w-full sm:w-auto bg-white text-black px-8 py-3.5 rounded-lg font-bold hover:bg-[#e1e1e1] transition-all text-[15px]">
                  Read Architecture Docs
                </Link>
                <a href="https://github.com/jhanvi857/gitresolve" target="_blank" className="w-full sm:w-auto bg-black text-white px-8 py-3.5 rounded-lg font-bold border border-white/[0.1] hover:bg-[#111] transition-all flex items-center justify-center gap-3 text-[15px]">
                  View Source Code
                </a>
              </div>
            </motion.div>

            <motion.div
              initial={{ opacity: 0, x: 20 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ duration: 0.6, delay: 0.2 }}
              className="relative"
            >
              <TerminalWindow title="gitresolve_server – bash">
                <div className="space-y-1 font-mono text-[13px] leading-relaxed">
                  <div className="flex gap-3 mb-4">
                    <span className="text-green-500 font-bold">$</span>
                    <span className="text-white">gitresolve install -g @jhanvi857/gitresolve</span>
                  </div>
                  <div className="text-[#888]">[INFO] gitresolve: Linked global binary successfully.</div>
                  <div className="flex gap-3 mt-4">
                    <span className="text-green-500 font-bold">$</span>
                    <span className="text-white">gitresolve resolve</span>
                  </div>
                  <div className="text-[#888] mt-2">[INFO] Watcher: Monitoring internal/auth/login.go...</div>
                  <div className="text-[#888]">[INFO] Bootstrap: Starting GitResolve v1.4.0</div>
                  <div className="text-blue-500">[READY] Conflict analysis live with AST validation.</div>
                  
                  <div className="flex gap-3 mt-6">
                    <span className="text-green-500 font-bold">$</span>
                    <span className="text-white">gitresolve status --json</span>
                  </div>
                  <div className="text-[#888]">{`{ "conflicts": 0, "confidence_score": 99.8 }`}</div>
                </div>
              </TerminalWindow>
            </motion.div>
          </div>
        </section>

        {/* Features Grid */}
        <section className="py-20 px-8">
          <div className="max-w-7xl mx-auto">
            <div className="mb-16">
              <h2 className="text-3xl font-extrabold tracking-tight mb-4">Built for scale.</h2>
              <p className="text-[#a1a1aa] text-lg font-medium">Every feature designed for production reliability.</p>
            </div>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
              <FeatureCard
                index={0}
                title="AST Classification"
                icon={Cpu}
                description="Go-tree-sitter integration maps conflicts to syntax trees. Detect function signature changes, not just text diffs."
              />
              <FeatureCard
                index={1}
                title="CWE-22 Hardened"
                icon={Shield}
                description="Mandatory os.Root sandboxing. Every file operation is cryptographically verified within repository boundaries."
              />
              <FeatureCard
                index={2}
                title="Deep Merging"
                icon={Terminal}
                description="Recursive map merges for JSON, YAML, and TOML. Handle complex configuration conflicts with native parsers."
              />
            </div>
          </div>
        </section>

        {/* CTA Section */}
        <section className="py-32 px-8 border-t border-white/[0.05]">
          <div className="max-w-4xl mx-auto text-center">
            <h2 className="text-4xl md:text-6xl font-extrabold tracking-tight mb-8">
              Ready to automate?
            </h2>
            <p className="text-[#a1a1aa] text-lg md:text-xl mb-12 font-medium">
              Join teams resolving conflicts with mathematical certainty. <br className="hidden md:block" />
              Purely offline. Completely deterministic.
            </p>
            <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
              <Link href="/docs/installation" className="w-full sm:w-auto bg-white text-black px-10 py-3.5 rounded-lg font-bold hover:bg-[#e1e1e1] transition-all">
                Read Documentation
              </Link>
              <a href="https://github.com/jhanvi857/gitresolve" target="_blank" className="w-full sm:w-auto px-10 py-3.5 rounded-lg font-bold border border-white/[0.1] hover:bg-white/5 transition-all">
                Star on GitHub
              </a>
            </div>
          </div>
        </section>
      </main>

      <Footer />
    </div>
  );
}