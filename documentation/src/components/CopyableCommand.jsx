"use client";

import React, { useState } from 'react';
import { Copy, Check, Terminal } from 'lucide-react';

export default function CopyableCommand({ command, label = "Usage" }) {
  const [copied, setCopied] = useState(false);

  const handleCopy = async () => {
    try {
      await navigator.clipboard.writeText(command.trim());
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch (err) {
      console.error('Failed to copy: ', err);
    }
  };

  return (
    <div className="my-6 p-4 rounded-xl bg-gradient-to-r from-[#0d0d0e] to-[#080808] border border-white/[0.08] flex items-center justify-between gap-4 group hover:border-blue-500/30 transition-all shadow-lg">
      <div className="flex items-center gap-3 overflow-x-auto min-w-0">
        <div className="w-8 h-8 rounded-lg bg-blue-500/10 border border-blue-500/20 flex items-center justify-center shrink-0">
          <Terminal className="w-4 h-4 text-blue-400" />
        </div>
        <div className="min-w-0">
          {label && (
            <div className="text-[10px] uppercase tracking-widest text-[#666] font-extrabold mb-0.5">
              {label}
            </div>
          )}
          <code className="font-mono text-[14px] text-white font-bold tracking-tight select-all">
            {command}
          </code>
        </div>
      </div>

      <button
        type="button"
        onClick={handleCopy}
        className="shrink-0 flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-[12px] font-bold text-[#a1a1aa] hover:text-white bg-white/[0.05] hover:bg-white/[0.1] border border-white/[0.08] transition-all"
        title="Copy command"
      >
        {copied ? (
          <>
            <Check className="w-3.5 h-3.5 text-emerald-400" />
            <span className="text-emerald-400 font-semibold">Copied</span>
          </>
        ) : (
          <>
            <Copy className="w-3.5 h-3.5" />
            <span>Copy</span>
          </>
        )}
      </button>
    </div>
  );
}
