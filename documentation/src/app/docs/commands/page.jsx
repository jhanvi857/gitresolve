export default function Commands() {
  return (
    <div className="space-y-10 font-sans">
      <div className="space-y-4">
        <h1 className="text-3xl font-bold text-white tracking-tight">CLI Command Reference</h1>
        <p className="text-[16px] text-[#888]">
          The primary executable endpoints used for analyzing, predicting, and merging states.
        </p>
      </div>

      <div className="space-y-8">
        
        {/* Merge */}
        <div className="p-6 rounded-xl bg-[#0a0a0a]">
          <h3 className="text-lg font-bold text-white mb-2 font-mono tracking-tight">gitresolve merge</h3>
          <p className="text-[#a1a1aa] text-[15px] mb-4">Engages the core algorithmic engine to triage conflicted index blocks internally.</p>
          <div className="bg-[#111] p-4 rounded font-mono text-[13px] text-[#888]">
            <span className="text-[#ededed]">$ gitresolve merge --dry-run</span><br/><br/>
            Engine Bootup: Initializing in dry-run isolation...<br/>
            Scanning index. Found 3 unmerged conflicts...<br/>
          </div>
        </div>

        {/* Scan */}
        <div className="p-6 rounded-xl bg-[#0a0a0a]">
          <h3 className="text-lg font-bold text-white mb-2 font-mono tracking-tight">gitresolve scan --target main</h3>
          <p className="text-[#a1a1aa] text-[15px]">Predicts and blocks pushing into production paths where conflict-overlap breaks exist without triggering git staging APIs.</p>
        </div>

        {/* Status */}
        <div className="p-6 rounded-xl bg-[#0a0a0a]">
          <h3 className="text-lg font-bold text-white mb-2 font-mono tracking-tight">gitresolve status</h3>
          <p className="text-[#a1a1aa] text-[15px]">Prints unmerged severity scores explicitly. Calculates logical overlaps and mathematically issues Severity (1-10) scoring rules to files currently staging.</p>
        </div>

        {/* Undo */}
        <div className="p-6 rounded-xl bg-[#0a0a0a]">
          <h3 className="text-lg font-bold text-white mb-2 font-mono tracking-tight">gitresolve undo --steps 1</h3>
          <p className="text-[#a1a1aa] text-[15px]">Safely inverts the SQLite auditing log tracking mapping back to original .gitresolve-orig structures natively without data fragmentation.</p>
        </div>

      </div>
    </div>
  );
}
