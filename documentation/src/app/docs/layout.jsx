"use client";

import { useState, useEffect } from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import Image from "next/image";
import { ChevronRight, Book, Shield, Cpu, Activity, Zap, BarChart3, Menu, X, GitMerge, History, RotateCcw, Database } from "lucide-react";
import Footer from "@/components/Footer";

const navItems = [
  { group: "Getting Started", items: [
    { name: "Installation", path: "/docs/installation", icon: Zap },
    { name: "Quick Start", path: "/docs/quick-start", icon: Activity },
    { name: "History Escalation", path: "/docs/history-escalation", icon: Shield, badge: "NEW" },
  ]},
  { group: "CLI Reference", items: [
    { name: "resolve", path: "/docs/commands/resolve", icon: Cpu },
    { name: "merge", path: "/docs/commands/merge", icon: GitMerge },
    { name: "scan", path: "/docs/commands/scan", icon: Shield },
    { name: "status", path: "/docs/commands/status", icon: Activity },
    { name: "blame", path: "/docs/commands/blame", icon: History },
    { name: "policy check", path: "/docs/commands/policy", icon: Shield },
    { name: "stats", path: "/docs/commands/stats", icon: BarChart3 },
    { name: "undo", path: "/docs/commands/undo", icon: RotateCcw },
    { name: "db repair", path: "/docs/commands/db", icon: Database },
  ]},
  { group: "Core Engine", items: [
    { name: "Architecture", path: "/docs/architecture", icon: Book },
    { name: "Security", path: "/docs/security", icon: Shield },
    { name: "Merge Logic & Flow", path: "/docs/merge-flow", icon: GitMerge },
    { name: "Policy Profiles", path: "/docs/policy", icon: Cpu },
  ]},
];

function SidebarContent({ pathname, currentHash, onItemClick }) {
  const hasExactMatch = navItems.some(group => 
    group.items.some(item => {
      const [itemPath, itemHashSuffix] = item.path.split('#');
      const itemHash = itemHashSuffix ? '#' + itemHashSuffix : '';
      return itemHash 
        ? pathname === itemPath && currentHash === itemHash
        : pathname === item.path && !currentHash;
    })
  );

  return (
    <div className="space-y-12">
      {navItems.map((group, i) => (
        <div key={i}>
          <h4 className="text-[11px] font-extrabold uppercase tracking-[0.3em] text-[#333] mb-6 px-6">
            {group.group}
          </h4>
          <div className="space-y-1.5">
            {group.items.map((item, j) => {
              const Icon = item.icon;
              const [itemPath, itemHashSuffix] = item.path.split('#');
              const itemHash = itemHashSuffix ? '#' + itemHashSuffix : '';
              const isActive = hasExactMatch
                ? (itemHash ? pathname === itemPath && currentHash === itemHash : pathname === item.path && !currentHash)
                : (pathname === itemPath);
              return (
                <Link
                  key={j}
                  href={item.path}
                  onClick={onItemClick}
                  className={`flex items-center gap-3 px-6 py-2.5 rounded-xl text-[14px] font-bold transition-all duration-200 group ${
                    isActive 
                      ? "bg-blue-500/10 text-blue-500 shadow-[inset_0_0_10px_rgba(0,112,243,0.05)]" 
                      : "text-[#a1a1aa] hover:text-white hover:bg-white/[0.03]"
                  }`}
                >
                  <Icon className={`w-4 h-4 transition-colors ${isActive ? "text-blue-500" : "text-[#333] group-hover:text-[#666]"}`} />
                  <span>{item.name}</span>
                  {item.badge && (
                    <span className="ml-auto px-2 py-0.5 rounded-full text-[9px] font-extrabold uppercase tracking-wider bg-amber-500/10 text-amber-400 border border-amber-500/20 shadow-[0_0_8px_rgba(245,158,11,0.2)]">
                      {item.badge}
                    </span>
                  )}
                </Link>
              );
            })}
          </div>
        </div>
      ))}
    </div>
  );
}

export default function DocsLayout({ children }) {
  const pathname = usePathname();
  const [isSidebarOpen, setIsSidebarOpen] = useState(false);
  const [currentHash, setCurrentHash] = useState("");

  useEffect(() => {
    const handleHashChange = () => {
      setCurrentHash(window.location.hash);
    };
    const timer = setTimeout(handleHashChange, 0);
    window.addEventListener("hashchange", handleHashChange);
    return () => {
      clearTimeout(timer);
      window.removeEventListener("hashchange", handleHashChange);
    };
  }, [pathname]);

  // Lock body scroll when mobile sidebar is open
  useEffect(() => {
    if (isSidebarOpen) {
      document.body.style.overflow = "hidden";
    } else {
      document.body.style.overflow = "";
    }
    return () => {
      document.body.style.overflow = "";
    };
  }, [isSidebarOpen]);

  const flatItems = navItems.flatMap(group => group.items);
  
  // Find current index
  let currentIndex = flatItems.findIndex(item => {
    const [itemPath, itemHashSuffix] = item.path.split('#');
    const itemHash = itemHashSuffix ? '#' + itemHashSuffix : '';
    return itemHash 
      ? pathname === itemPath && currentHash === itemHash
      : pathname === item.path && !currentHash;
  });

  // Fallback to path only
  if (currentIndex === -1) {
    currentIndex = flatItems.findIndex(item => item.path.split('#')[0] === pathname);
  }

  const prevTopic = currentIndex > 0 ? flatItems[currentIndex - 1] : null;
  const nextTopic = currentIndex !== -1 && currentIndex < flatItems.length - 1 ? flatItems[currentIndex + 1] : null;

  return (
    <div className="min-h-screen bg-black text-white selection:bg-blue-500/30">
      <div className="fixed inset-0 grid-bg opacity-10 pointer-events-none" />
      
      {/* Top Navbar */}
      <nav className="sticky top-0 z-40 border-b border-white/[0.05] bg-black/80 backdrop-blur-xl">
        <div className="w-full px-8 h-16 flex items-center justify-between">
          <div className="flex items-center gap-4 md:gap-12">
            <button 
              onClick={() => setIsSidebarOpen(true)}
              className="p-2 rounded-lg bg-black border border-white/[0.1] hover:border-blue-500/50 transition-all md:hidden cursor-pointer"
              aria-label="Toggle Sidebar"
            >
              <Menu className="w-5 h-5 text-white" />
            </button>
            
            <Link href="/" className="flex items-center gap-3 group">
              <div className="w-8 h-8 rounded-lg bg-[#111] border border-white/10 group-hover:border-blue-500/50 transition-all flex items-center justify-center">
                <GitMerge className="w-4 h-4 text-blue-400" />
              </div>
              <span className="font-extrabold tracking-tighter text-xl">gitresolve</span>
            </Link>
            
            <div className="hidden md:flex items-center gap-8 text-[14px] font-bold">
              <Link href="/docs/installation" className={`pb-1 border-b-2 transition-all ${pathname.startsWith("/docs") ? "text-white border-blue-500" : "text-[#555] border-transparent hover:text-white"}`}>Docs</Link>
              <Link href="/docs/history-escalation" className={`transition-all ${pathname === "/docs/history-escalation" ? "text-white" : "text-[#555] hover:text-white"}`}>History Escalation</Link>
              <Link href="/docs/architecture" className={`transition-all ${pathname === "/docs/architecture" ? "text-white" : "text-[#555] hover:text-white"}`}>Architecture</Link>
              <a href="https://github.com/jhanvi857/gitresolve" target="_blank" rel="noreferrer" className="text-[#555] hover:text-white transition-colors">GitHub</a>
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

      {/* Mobile Sidebar Drawer Overlay */}
      {isSidebarOpen && (
        <div 
          className="fixed inset-0 bg-black/60 backdrop-blur-sm z-50 md:hidden transition-opacity duration-300"
          onClick={() => setIsSidebarOpen(false)}
        />
      )}

      {/* Mobile Sidebar Drawer */}
      <aside 
        className={`fixed inset-y-0 left-0 w-80 max-w-[85vw] bg-black border-r border-white/[0.05] z-50 p-6 overflow-y-auto transform transition-transform duration-300 md:hidden ${
          isSidebarOpen ? "translate-x-0" : "-translate-x-full"
        }`}
      >
        <div className="flex items-center justify-between mb-8 pb-4 border-b border-white/[0.05]">
          <Link href="/" className="flex items-center gap-3" onClick={() => setIsSidebarOpen(false)}>
            <div className="w-8 h-8 rounded-lg bg-[#111] border border-white/10 flex items-center justify-center">
              <GitMerge className="w-4 h-4 text-blue-400" />
            </div>
            <span className="font-extrabold tracking-tighter text-lg">gitresolve</span>
          </Link>
          <button 
            onClick={() => setIsSidebarOpen(false)}
            className="p-2 rounded-lg bg-black border border-white/[0.1] hover:border-blue-500/50 transition-all cursor-pointer"
            aria-label="Close Sidebar"
          >
            <X className="w-5 h-5 text-white" />
          </button>
        </div>
        
        <SidebarContent 
          pathname={pathname} 
          currentHash={currentHash} 
          onItemClick={() => setIsSidebarOpen(false)} 
        />
      </aside>

      <div className="flex max-w-[1500px] mx-auto">
        {/* Desktop Sidebar */}
        <aside className="w-72 flex-shrink-0 hidden md:block pt-16 pr-10 border-r border-white/[0.05] sticky top-16 h-[calc(100vh-64px)] overflow-y-auto">
          <SidebarContent pathname={pathname} currentHash={currentHash} />
        </aside>

        {/* Content */}
        <main className="flex-1 min-w-0 flex flex-col">
          <div className="flex-1 px-6 md:px-20 py-16">
            <div className="max-w-4xl mx-auto">
              {children}

              {/* Previous & Next Navigation */}
              {(prevTopic || nextTopic) && (
                <div className="mt-16 pt-8 border-t border-white/[0.05] grid grid-cols-1 sm:grid-cols-2 gap-4">
                  {prevTopic ? (
                    <Link 
                      href={prevTopic.path}
                      className="p-5 rounded-xl border border-white/[0.05] hover:border-blue-500/30 bg-black/40 hover:bg-blue-500/[0.02] transition-all group flex flex-col text-left active:scale-[0.99]"
                    >
                      <span className="text-[11px] font-extrabold uppercase tracking-[0.2em] text-[#333] group-hover:text-blue-500 transition-colors mb-2">Previous Topic</span>
                      <span className="text-[15px] font-bold text-[#a1a1aa] group-hover:text-white transition-colors flex items-center gap-2">
                        <ChevronRight className="w-4 h-4 rotate-180 text-[#333] group-hover:text-blue-500 transition-colors" />
                        {prevTopic.name}
                      </span>
                    </Link>
                  ) : (
                    <div />
                  )}

                  {nextTopic ? (
                    <Link 
                      href={nextTopic.path}
                      className="p-5 rounded-xl border border-white/[0.05] hover:border-blue-500/30 bg-black/40 hover:bg-blue-500/[0.02] transition-all group flex flex-col text-right items-end active:scale-[0.99]"
                    >
                      <span className="text-[11px] font-extrabold uppercase tracking-[0.2em] text-[#333] group-hover:text-blue-500 transition-colors mb-2">Next Topic</span>
                      <span className="text-[15px] font-bold text-[#a1a1aa] group-hover:text-white transition-colors flex items-center gap-2">
                        {nextTopic.name}
                        <ChevronRight className="w-4 h-4 text-[#333] group-hover:text-blue-500 transition-colors" />
                      </span>
                    </Link>
                  ) : (
                    <div />
                  )}
                </div>
              )}
            </div>
          </div>
          <Footer />
        </main>
      </div>
    </div>
  );
}
