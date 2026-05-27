"use client";

import React from 'react';
import Link from 'next/link';
import Image from 'next/image';

// Inline SVGs for social icons to avoid lucide version mismatches
const GithubIcon = ({ className }) => (
  <svg className={className} width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
    <path d="M9 19c-5 1.5-5-2.5-7-3m14 6v-3.87a3.37 3.37 0 0 0-.94-2.61c3.14-.35 6.44-1.54 6.44-7A5.44 5.44 0 0 0 20 4.77 5.07 5.07 0 0 0 19.91 1S18.73.65 16 2.48a13.38 13.38 0 0 0-7 0C6.27.65 5.09 1 5.09 1A5.07 5.07 0 0 0 5 4.77a5.44 5.44 0 0 0-1.5 3.78c0 5.42 3.3 6.61 6.44 7A3.37 3.37 0 0 0 9 18.13V22"></path>
  </svg>
);

const TwitterIcon = ({ className }) => (
  <svg className={className} width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
    <path d="M22 4s-.7 2.1-2 3.4c1.6 10-9.4 17.3-18 11.6 2.2.1 4.4-.6 6-2C3 15.5.5 9.6 3 5c2.2 2.6 5.6 4.1 9 4-.9-4.2 4-6.6 7-3.8 1.1 0 3-1.2 3-1.2z"></path>
  </svg>
);

const DiscordIcon = ({ className }) => (
  <svg className={className} width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
    <circle cx="9" cy="12" r="1"></circle>
    <circle cx="15" cy="12" r="1"></circle>
    <path d="M7.5 7.1A9 9 0 0 0 2 12a10 10 0 0 0 2 6l1.5-1.5a10 10 0 0 1 1-6.5M16.5 7.1A9 9 0 0 1 22 12a10 10 0 0 1-2 6l-1.5-1.5a10 10 0 0 0-1-6.5"></path>
    <path d="M7.5 7.1C9 5.5 11 5 12 5s3 .5 4.5 2.1"></path>
    <path d="M7.5 18.5c1.5 1.5 3 2 4.5 2s3-.5 4.5-2.1"></path>
  </svg>
);

export default function Footer() {
  return (
    <footer className="border-t border-white/[0.05] bg-black py-16 px-8 relative overflow-hidden">
      <div className="absolute top-0 left-1/2 -translate-x-1/2 w-full h-px bg-gradient-to-r from-transparent via-blue-500/20 to-transparent" />
      
      <div className="max-w-7xl mx-auto">
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-12 mb-16">
          <div className="space-y-6">
            <Link href="/" className="flex items-center gap-3 group">
              <div className="p-1.5 rounded-lg bg-black border border-white/[0.1] group-hover:border-blue-500 transition-all">
                <Image src="/logo.png" alt="logo" width={20} height={20} className="opacity-90" />
              </div>
              <span className="font-extrabold tracking-tighter text-xl text-white">gitresolve</span>
            </Link>
            <p className="text-[#555] text-[14px] leading-relaxed max-w-xs font-medium">
              A high-precision, purely offline deterministic conflict resolution engine for modern engineering teams.
            </p>
            <div className="flex gap-4">
              <SocialLink href="https://github.com/jhanvi857/gitresolve" icon={GithubIcon} />
              <SocialLink href="#" icon={TwitterIcon} />
              <SocialLink href="#" icon={DiscordIcon} />
            </div>
          </div>

          <FooterGroup title="Platform">
            <FooterLink href="/docs/commands/init">init</FooterLink>
            <FooterLink href="/docs/commands/scan">scan</FooterLink>
            <FooterLink href="/docs/commands/resolve">resolve</FooterLink>
            <FooterLink href="/docs/commands/status">status</FooterLink>
          </FooterGroup>

          <FooterGroup title="Resources">
            <FooterLink href="/docs/installation">Installation</FooterLink>
            <FooterLink href="/docs/architecture">Architecture</FooterLink>
            <FooterLink href="/docs/security">Security</FooterLink>
            <FooterLink href="/docs/policy">Policy Profiles</FooterLink>
          </FooterGroup>

          <FooterGroup title="Community">
            <FooterLink href="https://github.com/jhanvi857/gitresolve/issues">GitHub Issues</FooterLink>
            <FooterLink href="/docs/security#report">Report Vulnerability</FooterLink>
            <FooterLink href="#">Discord Server</FooterLink>
            <FooterLink href="#" external>Twitter / X</FooterLink>
          </FooterGroup>
        </div>

        <div className="pt-8 border-t border-white/[0.05] flex flex-col md:flex-row justify-between items-center gap-6">
          <p className="text-[#333] text-[11px] font-extrabold uppercase tracking-[0.2em]">
            © 2026 GitResolve Engine. Built for Determinism.
          </p>
          <div className="flex gap-8 text-[11px] font-extrabold uppercase tracking-[0.2em] text-[#333]">
             <Link href="#" className="hover:text-[#555] transition-colors">Privacy Policy</Link>
             <Link href="#" className="hover:text-[#555] transition-colors">Terms of Service</Link>
             <Link href="#" className="hover:text-[#555] transition-colors">Security.md</Link>
          </div>
        </div>
      </div>
    </footer>
  );
}

function FooterGroup({ title, children }) {
  return (
    <div className="space-y-6">
      <h4 className="text-[11px] font-extrabold uppercase tracking-[0.3em] text-white/40">{title}</h4>
      <div className="flex flex-col gap-3">
        {children}
      </div>
    </div>
  );
}

function FooterLink({ href, children, external }) {
  const isExternal = external || href.startsWith('http');
  const Component = isExternal ? 'a' : Link;
  const props = isExternal ? { target: '_blank', rel: 'noopener noreferrer' } : {};

  return (
    <Component 
      href={href} 
      {...props}
      className="text-[#555] hover:text-white transition-colors text-[14px] font-bold flex items-center gap-2 group"
    >
      {children}
      {isExternal && (
        <svg className="w-3 h-3 opacity-0 group-hover:opacity-100 transition-opacity" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
          <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"></path>
          <polyline points="15 3 21 3 21 9"></polyline>
          <line x1="10" y1="14" x2="21" y2="3"></line>
        </svg>
      )}
    </Component>
  );
}

function SocialLink({ href, icon: Icon }) {
  return (
    <a 
      href={href} 
      target="_blank" 
      className="p-2 rounded-lg bg-black border border-white/[0.05] text-[#333] hover:text-white hover:border-white/10 transition-all"
    >
      <Icon className="w-4 h-4" />
    </a>
  );
}
