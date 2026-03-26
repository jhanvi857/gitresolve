import Link from "next/link";

export default function Home() {
  return (
    <div className="flex flex-col min-h-screen bg-black text-[#ededed] font-sans">
      <nav className="bg-black/80 sticky top-0 z-50 backdrop-blur-md">
        <div className="w-full max-w-[90rem] mx-auto px-6 h-14 flex items-center justify-between">
          <div className="font-semibold tracking-tight text-white flex items-center">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="mr-3"><polygon points="12 2 22 8.5 22 15.5 12 22 2 15.5 2 8.5 12 2"></polygon><line x1="12" y1="22" x2="12" y2="15.5"></line><polyline points="22 8.5 12 15.5 2 8.5"></polyline><polyline points="2 15.5 12 8.5 22 15.5"></polyline><line x1="12" y1="2" x2="12" y2="8.5"></line></svg>
            gitresolve
          </div>
          <div className="flex items-center gap-6 text-[14px] text-[#888] font-medium">
            <Link href="/docs/installation" className="hover:text-white transition">Documentation</Link>
            <a href="https://github.com/jhanvi857/gitresolve" className="hover:text-white transition">GitHub</a>
            <a href="#features" className="hover:text-white transition">Features</a>
          </div>
        </div>
      </nav>

      <main className="flex-1 w-full flex flex-col items-center">
        {/* Hero Section */}
        <section className="w-full max-w-5xl mx-auto px-6 pt-32 pb-24 flex flex-col items-center text-center">
          <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-[#111] text-[#888] text-[13px] font-medium mb-8 transition hover:bg-[#222]">
            <span>v1.0.0 Production Release</span>
            <span className="w-4 h-4 rounded-full bg-[#333] flex items-center justify-center text-[10px] text-white overflow-hidden">→</span>
          </div>
          
          <h1 className="text-[56px] sm:text-[72px] font-bold tracking-tighter text-white leading-tight mb-8">
            Deterministic Conflict Resolution
          </h1>
          
          <p className="max-w-[42rem] text-[20px] text-[#888] mb-12 leading-normal">
            A purely offline, highly deterministic Git conflict resolution engine scaling standard 20-year-old algorithms into modern Abstract Syntax Tree (AST) structural analysis. Built without LLMs.
          </p>

          <div className="flex flex-col sm:flex-row items-center gap-4">
            <Link href="/docs/installation" className="bg-white hover:bg-[#ededed] text-black font-semibold px-6 py-3 rounded-md transition text-base">
              Read Documentation
            </Link>
            <a href="#features" className="bg-[#111] hover:bg-[#222] text-white font-semibold px-6 py-3 rounded-md transition text-base">
              Explore Architecture
            </a>
          </div>
        </section>

        {/* Feature Grid / Core Infrastructure */}
        <section id="features" className="w-full bg-[#050505]">
          <div className="max-w-6xl mx-auto px-6 py-24">
            <div className="mb-16">
              <h2 className="text-3xl font-bold tracking-tight text-white mb-4">Core Infrastructure</h2>
              <p className="text-[16px] text-[#888] max-w-2xl">
                Git solves diff algorithms through visual line-breaking. gitresolve scales natively into the codebase structure, fixing whitespace, syntax logic, and imports implicitly without altering semantic behavior or breaking.
              </p>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              
              <div className="p-8 rounded-xl bg-[#111] transition duration-300">
                <h3 className="text-xl font-bold text-white mb-3 tracking-tight">Zero Large Language Models</h3>
                <p className="text-[#888] text-[15px] leading-relaxed">
                  Resolving internal code disputes via generative LLM APIs acts as a major enterprise leak. Our algorithm calculates AST mismatches entirely locally without probabilistic networks hallucinating over code.
                </p>
              </div>

              <div className="p-8 rounded-xl bg-[#111] transition duration-300">
                <h3 className="text-xl font-bold text-white mb-3 tracking-tight">POSIX Atomic Native Writes</h3>
                <p className="text-[#888] text-[15px] leading-relaxed">
                  We leverage <code className="bg-[#222] px-1 py-0.5 rounded text-[#ededed]">writeAtomic()</code> mappings to parse resolved buffers straight natively into `.tmp` disks before flushing via strict `os.Rename`. Safe from OS crashes or panics.
                </p>
              </div>

              <div className="p-8 rounded-xl bg-[#111] transition duration-300">
                <h3 className="text-xl font-bold text-white mb-3 tracking-tight">Structured File Syncing</h3>
                <p className="text-[#888] text-[15px] leading-relaxed">
                  When developers concurrently map non-overlapping JSON objects or modify top-level arrays, gitresolve implicitly deduplicates structural imports automatically.
                </p>
              </div>

            </div>
          </div>
        </section>

        {/* Engine Code Demo section */}
        <section className="w-full bg-[#000]">
          <div className="max-w-6xl mx-auto px-6 py-32 flex flex-col lg:flex-row items-center gap-16">
            <div className="lg:w-1/2">
              <h2 className="text-4xl font-bold tracking-tight text-white mb-6">Engine Bootup <br/> In Native Speed.</h2>
              <p className="text-[16px] text-[#888] leading-relaxed mb-6">
                Compile locally directly from source using standard go toolchains. The executable connects securely to your `.git/` tracking endpoints.
              </p>
              <div className="text-[14px] flex flex-col gap-3 font-mono text-[#a1a1aa] pl-4">
                <div><span>1.</span> <span className="text-white">gitresolve scan</span> — predict conflicts</div>
                <div><span>2.</span> <span className="text-white">gitresolve merge</span> — resolve deterministically</div>
                <div><span>3.</span> <span className="text-white">gitresolve undo</span> — rollback SQLite snapshot</div>
              </div>
            </div>
            
            <div className="lg:w-1/2 w-full">
              <div className="w-full rounded-xl bg-[#111] shadow-2xl overflow-hidden">
                <div className="bg-[#222] px-4 py-3 flex items-center gap-2">
                  <div className="w-2.5 h-2.5 rounded-full bg-[#333]"></div>
                  <div className="w-2.5 h-2.5 rounded-full bg-[#333]"></div>
                  <div className="w-2.5 h-2.5 rounded-full bg-[#333]"></div>
                </div>
                <div className="p-6 font-mono text-[13px] leading-extralight text-[#888]">
                  <div><span className="text-[#888]">$</span> <span className="text-white">gitresolve merge --dry-run</span></div>
                  <div className="mt-4">Engine Bootup: Initializing inside directory '.'</div>
                  <div>Scanning index. Found 2 unmerged conflicts...</div>
                  
                  <div className="mt-4">--- Processing pkg/auth/login.go ---</div>
                  <div className="text-[#ededed] bg-[#222] px-2 py-1 inline-block mt-1 rounded">! Escalated Severity 9 [Type: Logic]</div>
                  
                  <div className="mt-4">--- Processing web/config.yaml ---</div>
                  <div className="text-white">✓ Auto-resolved 100% conflicts securely.</div>
                  
                  <div className="mt-6 text-[#666]">1 conflict resolved atomically. 1 requires manual intervention.</div>
                </div>
              </div>
            </div>
          </div>
        </section>
      </main>

      <footer className="bg-[#050505] py-12">
        <div className="max-w-[90rem] mx-auto px-6 text-[14px] text-[#666] flex justify-between">
          <p>© 2026 gitresolve. Built purely offline.</p>
          <div className="flex gap-4">
            <Link href="/docs/installation" className="hover:text-white transition">Documentation</Link>
            <a href="https://github.com/jhanvi857/gitresolve" className="hover:text-white transition">GitHub</a>
          </div>
        </div>
      </footer>
    </div>
  );
}