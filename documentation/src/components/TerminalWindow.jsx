"use client";

import React, { useState, useRef } from 'react';
import { Copy, Check } from 'lucide-react';

export default function TerminalWindow({ children, title }) {
  const [copied, setCopied] = useState(false);
  const contentRef = useRef(null);

  const handleCopy = async () => {
    if (contentRef.current) {
      const text = contentRef.current.innerText || '';
      const lines = text.split('\n');
      const cmdLine = lines.find(l => l.trim().startsWith('$'));
      const textToCopy = cmdLine ? cmdLine.replace(/^\$\s*/, '') : text;
      
      if (textToCopy) {
        try {
          await navigator.clipboard.writeText(textToCopy.trim());
          setCopied(true);
          setTimeout(() => setCopied(false), 2000);
        } catch (err) {
          console.error('Failed to copy: ', err);
        }
      }
    }
  };

  return (
    <div className="terminal-window my-8">
      <div className="terminal-header flex items-center justify-between">
        <div className="flex items-center">
          <div className="terminal-dot bg-[#ff5f56] shadow-[0_0_8px_rgba(255,95,86,0.3)]" />
          <div className="terminal-dot bg-[#ffbd2e] shadow-[0_0_8px_rgba(255,189,46,0.3)]" />
          <div className="terminal-dot bg-[#27c93f] shadow-[0_0_8px_rgba(39,201,63,0.3)]" />
          {title && <span className="ml-6 text-[12px] font-bold text-gray-500 uppercase tracking-[0.2em] font-mono">{title}</span>}
        </div>

        <button
          type="button"
          onClick={handleCopy}
          className="flex items-center gap-1.5 px-2 py-0.5 rounded text-[11px] font-medium text-[#888] hover:text-white bg-white/[0.04] hover:bg-white/[0.08] border border-white/[0.06] transition-all cursor-pointer"
          title="Copy command"
        >
          {copied ? (
            <>
              <Check className="w-3 h-3 text-emerald-400" />
              <span className="text-emerald-400">Copied!</span>
            </>
          ) : (
            <>
              <Copy className="w-3 h-3" />
              <span>Copy</span>
            </>
          )}
        </button>
      </div>
      <div ref={contentRef} className="terminal-content bg-[#000]">
        {children}
      </div>
    </div>
  );
}
