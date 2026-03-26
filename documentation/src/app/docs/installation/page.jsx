export default function Installation() {
  return (
    <div className="space-y-10 font-sans">
      <div className="space-y-4">
        <h1 className="text-3xl font-bold text-[#ededed] tracking-tight">Installation & Setup</h1>
        <p className="text-[16px] text-[#888]">
          gitresolve compiles securely into a single static binary. No bloated network requests or external dependencies. 
        </p>
      </div>

      <div className="space-y-4">
        <h2 className="text-[18px] font-semibold text-[#ededed] pb-2">Quick Start</h2>
        <p className="text-[#a1a1aa] text-[15px]">
          The fastest way to install gitresolve across your terminal paths is utilizing the standard Go version 1.20+ toolchain.
        </p>
        
        <div className="bg-[#111] p-4 rounded-lg font-mono text-[13px] text-[#ededed]">
          <span className="text-white">$ go install</span> github.com/jhanvi857/gitresolve@latest
        </div>
      </div>

      <div className="space-y-4">
        <h2 className="text-[18px] font-semibold text-[#ededed] pb-2">Workspace Ownership Configuration</h2>
        <p className="text-[#a1a1aa] text-[15px]">
          In order to initialize pre-push overlap warnings tied to massive engineering units, map specific directories to discrete ownership namespaces inside a <code className="text-[#ededed] bg-[#111] px-1.5 py-0.5 rounded text-[13px]">.gitresolve/owners.json</code> file.
        </p>

        <div className="bg-[#111] p-5 rounded-lg font-mono text-[13px] overflow-x-auto text-[#888]">
          <pre>
{`{
  "frontend_team": ["web/**", "components/**"],
  "platform_team": ["pkg/**", "internal/cloud/**"],
  "security_team": ["internal/secrets/**"]
}`}
          </pre>
        </div>
      </div>
    </div>
  );
}
