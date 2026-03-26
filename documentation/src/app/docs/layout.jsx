import Link from "next/link";

export default function DocsLayout({ children }) {
  const routes = [
    { name: "Getting Started", section: true },
    { name: "Installation & Setup", path: "/docs/installation" },
    { name: "Commands", section: true },
    { name: "CLI Reference", path: "/docs/commands" },
    { name: "Core Engine", section: true },
    { name: "System Architecture", path: "/docs/architecture" },
    { name: "Security & Privacy", path: "/docs/security" },
  ];

  return (
    <div className="min-h-screen bg-black text-[#ededed] font-sans">
      
      {/* Top Navbar */}
      <nav className="bg-black/80 sticky top-0 z-50 backdrop-blur-md">
        <div className="w-full px-6 md:px-12 h-14 flex items-center justify-between">
          <Link href="/" className="font-semibold text-sm tracking-tight text-white flex items-center">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="mr-3"><polygon points="12 2 22 8.5 22 15.5 12 22 2 15.5 2 8.5 12 2"></polygon><line x1="12" y1="22" x2="12" y2="15.5"></line><polyline points="22 8.5 12 15.5 2 8.5"></polyline><polyline points="2 15.5 12 8.5 22 15.5"></polyline><line x1="12" y1="2" x2="12" y2="8.5"></line></svg>
            gitresolve
          </Link>
          <div className="flex text-sm text-[#888] gap-6 font-medium">
             <Link href="/docs/installation" className="text-white transition">Documentation</Link>
             <a href="https://github.com/jhanvi857/gitresolve" target="_blank" className="hover:text-white transition">GitHub</a>
             <Link href="/" className="hover:text-white transition">Enterprise</Link>
          </div>
        </div>
      </nav>

      <div className="flex max-w-[90rem] mx-auto w-full">
        {/* Sidebar Nav */}
        <aside className="w-64 flex-shrink-0 hidden md:block pt-8 pr-6" style={{ height: "calc(100vh - 56px)", position: "sticky", top: "56px" }}>
          <div className="flex flex-col gap-[2px] w-full text-[14px]">
            {routes.map((route, i) => (
              route.section ? (
                <div key={i} className="font-semibold text-white mt-6 mb-2 tracking-tight px-3">{route.name}</div>
              ) : (
                <Link 
                  key={i} 
                  href={route.path} 
                  className="text-[#a1a1aa] font-medium px-3 py-1.5 rounded-md hover:text-white hover:bg-[#111] transition-colors block"
                >
                  {route.name}
                </Link>
              )
            ))}
          </div>
        </aside>

        {/* Content Area */}
        <main className="flex-1 px-6 md:px-16 py-12 max-w-4xl text-[15px] leading-relaxed">
          {children}
        </main>
      </div>

    </div>
  );
}
