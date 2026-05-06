"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import Image from "next/image";

const navItems = [
  { name: "Getting Started", section: true },
  { name: "Installation & Setup", path: "/docs/installation" },
  { name: "CLI Reference", section: true },
  { name: "Scan (Predictive)", path: "/docs/commands#scan", indent: true },
  { name: "Status (Diagnostic)", path: "/docs/commands#status", indent: true },
  { name: "Merge (Auto-resolve)", path: "/docs/commands#merge", indent: true },
  { name: "Resolve (Interactive)", path: "/docs/commands#resolve", indent: true },
  { name: "Blame (Auditability)", path: "/docs/commands#blame", indent: true },
  { name: "Undo (History)", path: "/docs/commands#undo", indent: true },
  { name: "Observability", section: true },
  { name: "Stats & Metrics", path: "/docs/stats" },
  { name: "Decision History", path: "/docs/commands#blame" },
  { name: "Configuration", section: true },
  { name: "Policy Profiles", path: "/docs/policy" },
  { name: "Core Engine", section: true },
  { name: "System Architecture", path: "/docs/architecture" },
  { name: "Security & Privacy", path: "/docs/security" },
];

export default function DocsLayout({ children }) {
  const pathname = usePathname();

  return (
    <div className="min-h-screen bg-black text-[#ededed] font-sans selection:bg-[#3b82f6]">
      {/* Top Navbar */}
      <nav className="border-b border-[#222] bg-black sticky top-0 z-50">
        <div className="w-full px-6 md:px-12 h-16 flex items-center justify-between">
          <Link href="/" className="font-bold tracking-tight text-white flex items-center gap-3">
            <Image src="/logo.png" alt="gitresolve logo" width={24} height={24} className="rounded" />
            gitresolve
          </Link>
          <div className="flex text-sm text-gray-500 gap-8 font-medium">
            <Link href="/docs/installation" className="text-white hover:text-gray-300 transition-colors">Documentation</Link>
            <a href="https://github.com/jhanvi857/gitresolve" target="_blank" className="hover:text-white transition-colors">GitHub</a>
            {/* <Link href="/enterprise" className="hover:text-white transition-colors">Enterprise</Link> */}
          </div>
        </div>
      </nav>

      <div className="flex max-w-[90rem] mx-auto w-full">
        {/* Sidebar Nav */}
        <aside className="w-64 flex-shrink-0 hidden md:block pt-12 pr-6 border-r border-[#222]" style={{ height: "calc(100vh - 64px)", position: "sticky", top: "64px" }}>
          <div className="flex flex-col gap-0.5 w-full text-[14px] overflow-y-auto pb-10 custom-scrollbar">
            {navItems.map((item, i) => (
              item.section ? (
                <div key={i} className={`font-bold text-[#888] mt-8 mb-4 px-3 text-[11px] uppercase tracking-widest leading-none ${i === 0 ? 'mt-4' : ''}`}>
                  {item.name}
                </div>
              ) : (
                <Link
                  key={i}
                  href={item.path}
                  className={`px-3 py-1.5 rounded-md transition-all font-medium flex items-center gap-2 ${item.indent ? 'ml-4 text-[13px] opacity-80 hover:opacity-100' : ''} ${pathname === item.path ? 'bg-white/5 text-white border border-[#333]' : 'text-gray-500 hover:text-white hover:bg-white/5'}`}
                >
                  {item.indent && <span className="w-1.5 h-1.5 rounded-full bg-[#333] shrink-0" />}
                  {item.name}
                </Link>
              )
            ))}
          </div>
        </aside>

        {/* Content Area */}
        <main className="flex-1 px-6 md:px-16 py-12 max-w-5xl">
          {children}
        </main>
      </div>

    </div>
  );
}
