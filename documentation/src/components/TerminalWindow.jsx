"use client";

import React from 'react';

export default function TerminalWindow({ children, title }) {
  return (
    <div className="terminal-window my-8">
      <div className="terminal-header">
        <div className="terminal-dot bg-[#ff5f56] shadow-[0_0_8px_rgba(255,95,86,0.3)]" />
        <div className="terminal-dot bg-[#ffbd2e] shadow-[0_0_8px_rgba(255,189,46,0.3)]" />
        <div className="terminal-dot bg-[#27c93f] shadow-[0_0_8px_rgba(39,201,63,0.3)]" />
        {title && <span className="ml-6 text-[12px] font-bold text-gray-500 uppercase tracking-[0.2em] font-mono">{title}</span>}
      </div>
      <div className="terminal-content bg-[#000]">
        {children}
      </div>
    </div>
  );
}
