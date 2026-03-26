export default function Architecture() {
  return (
    <div className="space-y-10">
      <div className="space-y-4">
        <h1 className="text-3xl font-bold text-[#ededed] tracking-tight">System Architecture</h1>
        <p className="text-[16px] text-[#888]">
          Scaling classic Git workflows computationally via strict AST modeling and transactional isolation domains.
        </p>
      </div>

      <div className="space-y-6">
        <div className="rounded-xl bg-[#000] p-8 my-10">
          <div className="flex flex-col gap-6 font-mono text-[13px] items-center text-[#ededed]">
            <div className="w-full flex justify-center text-white pb-4 font-bold">
              [ CLI COMMAND: gitresolve merge ]
            </div>
            
            <div className="flex flex-col sm:flex-row gap-4 items-center">
               <div className="bg-[#111] px-4 py-2 rounded">.git/ objects</div>
               <span className="text-[#666] hidden sm:block">→</span>
               <div className="text-[#666] sm:hidden">↓</div>
               <div className="bg-[#111] text-[#ededed] px-4 py-2 rounded">Extract Blobs</div>
               <span className="text-[#666] hidden sm:block">→</span>
               <div className="text-[#666] sm:hidden">↓</div>
               <div className="bg-[#111] px-4 py-2 rounded">Conflict Parser</div>
            </div>

            <div className="text-[#666]">↓</div>

            <div className="flex flex-col sm:flex-row gap-4 items-center text-center">
               <div className="bg-[#111] text-[#ededed] px-4 py-2 rounded">Tree-Sitter AST Analysis</div>
               <span className="text-[#666] hidden sm:block">→</span>
               <div className="text-[#666] sm:hidden">↓</div>
               <div className="bg-[#111] px-4 py-2 rounded flex flex-col gap-1 text-[11px] text-left text-[#888]">
                 <span className="text-[#ededed]">Sev: 1 - Format/Space</span>
                 <span className="text-[#ededed]">Sev: 3 - JSON Config</span>
                 <span className="text-[#ededed]">Sev: 9 - Auth Logic</span>
               </div>
            </div>

            <div className="text-[#666]">↓</div>

            <div className="flex gap-8 items-center text-center mt-4">
               <div className="flex flex-col gap-2">
                 <div className="text-[#ededed] font-bold text-[11px] uppercase tracking-wider">Deterministic</div>
                 <div className="bg-[#111] px-4 py-2 rounded">Atomic Disk Write</div>
               </div>
               
               <div className="flex flex-col gap-2">
                 <div className="text-[#888] font-bold text-[11px] uppercase tracking-wider">Blocked</div>
                 <div className="bg-[#111] px-4 py-2 rounded">Manual Escalation</div>
               </div>
            </div>
          </div>
        </div>
      </div>

      <div className="space-y-6">
        <h2 className="text-xl font-semibold text-white pb-2">AST Analysis (Not Text Checking)</h2>
        <p className="text-[#888] text-[15px] leading-relaxed">
          Traditional Git engines like Myers-diff parse visual edit distancing, causing terrifying semantic logic fractures. 
          The <span className="text-white font-medium">gitresolve</span> engine bypasses visual distances natively.
        </p>

        <div className="bg-[#0a0a0a] p-5 rounded-xl font-mono text-[13px] overflow-x-auto">
          <pre className="text-[#888]">
{`// internal/conflict/classifier.go

func Classify(c *Conflict) {
  // Automatically solves spaces mathematically
  if isWhitespaceOnly(c.OurLines, c.TheirLines) {
      c.Type = TypeWhitespace
      c.CanAutoResolve = true
  }

  // Defends critical pathways aggressively 
  if isSensitivePath(c.FilePath) { // "auth/**", "payments/**"
      c.Severity = SeverityCritical
      c.CanAutoResolve = false
  }
}`}
          </pre>
        </div>
      </div>
    </div>
  );
}
