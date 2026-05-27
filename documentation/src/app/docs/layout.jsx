"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import Image from "next/image";
import { ChevronRight, Book, Shield, Cpu, Activity, Zap, BarChart3 } from "lucide-react";
import Footer from "@/components/Footer";

const navItems = [
  { group: "Getting Started", items: [
    { name: "Installation", path: "/docs/installation", icon: Zap },
    { name: "Quick Start", path: "/docs/installation#quick-start", icon: Activity },
  ]},
  { group: "CLI Reference", items: [
    { name: "init", path: "/docs/commands/init", icon: Zap },
    { name: "scan", path: "/docs/commands/scan", icon: Shield },
    { name: "status", path: "/docs/commands/status", icon: Activity },
    { name: "resolve", path: "/docs/commands/resolve", icon: Cpu },
    { name: "config", path: "/docs/commands/config", icon: Book },
    { name: "stats", path: "/docs/commands/stats", icon: BarChart3 },
    { name: "audit", path: "/docs/commands/audit", icon: Shield },
  ]},
  { group: "Core Engine", items: [
    { name: "Architecture", path: "/docs/architecture", icon: Book },
    { name: "Security", path: "/docs/security", icon: Shield },
    { name: "Policy Profiles", path: "/docs/policy", icon: Cpu },
  ]},
];

export default function DocsLayout({ children }) {
  const pathname = usePathname();

  return (
    <div className="min-h-screen bg-black text-white selection:bg-blue-500/30">
      <div className="fixed inset-0 grid-bg opacity-10 pointer-events-none" />
      
      {/* Top Navbar */}
      <nav className="sticky top-0 z-50 border-b border-white/[0.05] bg-black/80 backdrop-blur-xl">
        <div className="w-full px-8 h-16 flex items-center justify-between">
          <div className="flex items-center gap-12">
            <Link href="/" className="flex items-center gap-3 group">
              <div className="p-1.5 rounded-lg bg-black border border-white/[0.1] group-hover:border-blue-500 transition-all duration-300">
                <Image src="/logo.png" alt="logo" width={22} height={22} className="opacity-90" />
              </div>
              <span className="font-extrabold tracking-tighter text-xl">gitresolve</span>
            </Link>
            <div className="hidden md:flex items-center gap-8 text-[14px] font-bold">
              <Link href="/docs/installation" className={`pb-1 border-b-2 transition-all ${pathname === "/docs/installation" ? "text-white border-blue-500" : "text-[#555] border-transparent hover:text-white"}`}>Docs</Link>
              <Link href="/docs/architecture" className={`transition-all ${pathname === "/docs/architecture" ? "text-white" : "text-[#555] hover:text-white"}`}>Architecture</Link>
              <a href="https://github.com/jhanvi857/gitresolve" target="_blank" className="text-[#555] hover:text-white transition-colors">Github</a>
            </div>
          </div>
          <div className="flex items-center gap-6">
             <Link href="/docs/installation" className="bg-white text-black px-5 py-2 rounded-full text-[14px] font-bold hover:bg-[#e1e1e1] active:scale-95 transition-all flex items-center gap-2 group">
              Get Started
              <ChevronRight className="w-4 h-4 group-hover:translate-x-0.5 transition-transform" />
            </Link>
          </div>
        </div>
      </nav>

      <div className="flex max-w-[1500px] mx-auto">
        {/* Sidebar */}
        <aside className="w-72 flex-shrink-0 hidden md:block pt-16 pr-10 border-r border-white/[0.05] sticky top-16 h-[calc(100vh-64px)] overflow-y-auto">
          <div className="space-y-12">
            {navItems.map((group, i) => (
              <div key={i}>
                <h4 className="text-[11px] font-extrabold uppercase tracking-[0.3em] text-[#333] mb-6 px-6">
                  {group.group}
                </h4>
                <div className="space-y-1.5">
                  {group.items.map((item, j) => {
                    const Icon = item.icon;
                    const isActive = pathname === item.path || (item.path.includes('#') && pathname === item.path.split('#')[0]);
                    return (
                      <Link
                        key={j}
                        href={item.path}
                        className={`flex items-center gap-3 px-6 py-2.5 rounded-xl text-[14px] font-bold transition-all duration-200 group ${
                          isActive 
                            ? "bg-blue-500/10 text-blue-500 shadow-[inset_0_0_10px_rgba(0,112,243,0.05)]" 
                            : "text-[#a1a1aa] hover:text-white hover:bg-white/[0.03]"
                        }`}
                      >
                        <Icon className={`w-4 h-4 transition-colors ${isActive ? "text-blue-500" : "text-[#333] group-hover:text-[#666]"}`} />
                        {item.name}
                      </Link>
                    );
                  })}
                </div>
              </div>
            ))}
          </div>
        </aside>

        {/* Content */}
        <main className="flex-1 min-w-0 flex flex-col">
          <div className="flex-1 px-8 md:px-20 py-16">
            <div className="max-w-4xl mx-auto">
              {children}
            </div>
          </div>
          <Footer />
        </main>
      </div>
    </div>
  );
}
