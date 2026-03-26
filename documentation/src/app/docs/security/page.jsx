export default function Security() {
  return (
    <div className="space-y-10">
      <div className="space-y-4">
        <h1 className="text-3xl font-bold text-white tracking-tight">Security & Privacy Standard</h1>
        <p className="text-[16px] text-[#888]">
          The ultimate defense infrastructure built directly for internal enterprise engineering environments.
        </p>
      </div>

      {/* No LLM Banner */}
      <div className="p-8 rounded-xl my-8">
        <h2 className="text-xl font-bold text-white mb-3 flex items-center gap-2 tracking-tight">
          <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 20 20"><path fillRule="evenodd" d="M10 1.944A11.954 11.954 0 012.166 5C2.056 5.642 2 6.319 2 7c0 5.225 3.34 9.67 8 11.317C14.66 16.67 18 12.225 18 7c0-.682-.057-1.358-.166-2.001A11.954 11.954 0 0110 1.944zM11 14a1 1 0 11-2 0 1 1 0 012 0zm0-7a1 1 0 10-2 0v3a1 1 0 102 0V7z" clipRule="evenodd"></path></svg>
          Zero LLM Integrations
        </h2>
        <p className="text-[#a1a1aa] leading-relaxed text-[15px]">
          <strong className="text-white font-medium">gitresolve strictly does not use LLMs (Large Language Models) or probabilistic AI networks to resolve your source code.</strong> 
          Generating intelligent code resolution natively requires zero interaction with external generative APIs.<br /><br />
          Sending proprietary intellectual property to clouded APIs is a massive compliance vulnerability. Code hallucinates. 
          Compliance breaks. Security teams aggressively reject it.<br /><br />
          Instead, gitresolve applies strict, verifiable mathematical trees to parse conflicts locally. <strong className="text-white font-medium">100% Privacy. 100% Accuracy. 0 API Calls.</strong>
        </p>
      </div>

      <div className="space-y-8">
        <div className="space-y-4">
          <h3 className="text-[18px] font-semibold text-white pb-2 tracking-tight">POSIX Atomic Disk Defense</h3>
          <p className="text-[#888] text-[15px]">Writing directly to file handles during automated merges corrupts states if processes hang or crash natively.</p>
          <div className="flex gap-4">
            <div className="w-1 bg-[#fff] rounded"></div>
            <p className="text-[#ededed] bg-[#111] px-4 py-3 rounded-md font-mono text-[13px]">
               func writeAtomic() → fsync() → os.Rename()
            </p>
          </div>
          <p className="text-[14px] text-[#666]">Every file successfully auto-resolved is initially written to a <code className="text-[#ededed] bg-[#111] px-1 py-0.5 rounded">.gitresolve-tmp</code> chunk before atomic system renames perfectly swap pointers. Power loss? No data loss.</p>
        </div>

        <div className="space-y-4">
          <h3 className="text-[18px] font-semibold text-white pb-2 tracking-tight">Replay Snapshots (Undo State)</h3>
          <p className="text-[#888] text-[15px] leading-relaxed">
            Every distinctive gitresolve command executes a snapshot signature to a localized tracking layer. Executing <span className="font-mono text-[#ededed] bg-[#111] px-1.5 py-0.5 rounded text-[13px]">gitresolve undo</span> forces a rigorous rollback protocol using the original <code className="text-[#ededed] bg-[#111] px-1.5 py-0.5 rounded">.gitresolve-orig</code> file backups to perfectly restore original unmerged file pointers flawlessly.
          </p>
        </div>
      </div>
    </div>
  );
}
