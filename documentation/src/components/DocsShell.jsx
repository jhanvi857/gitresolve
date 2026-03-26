"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";

const nav = [
  { href: "/", label: "Overview" },
  { href: "/architecture", label: "Architecture" },
  { href: "/merge-flow", label: "Merge Flow" },
  { href: "/commands", label: "Commands" },
  { href: "/operations", label: "Operations" },
];

export default function DocsShell({ title, subtitle, children }) {
  const pathname = usePathname();

  const isActive = (href) => pathname === href;

  return (
    <div className="site-wrap">
      <header className="topbar">
        <div className="topbar-inner">
          <Link href="/" className="brand">
            <span className="brand-mark">{"//"}</span>
            <span>gitresolve docs</span>
          </Link>
          <div className="top-links-wrap">
            <nav className="top-links" aria-label="Primary">
              {nav.map((item) => (
                <Link
                  key={item.href}
                  href={item.href}
                  className={isActive(item.href) ? "nav-link is-active" : "nav-link"}
                  aria-current={isActive(item.href) ? "page" : undefined}
                >
                  {item.label}
                </Link>
              ))}
            </nav>
            <Link href="/get-started" className="top-cta">
              Get Started
            </Link>
          </div>
        </div>
      </header>

      <main className="layout-grid">
        <aside className="sidebar" aria-label="Documentation sections">
          <p className="sidebar-kicker">Docs Navigation</p>
          {nav.map((item) => (
            <Link
              key={item.href}
              href={item.href}
              className={isActive(item.href) ? "sidebar-link is-active" : "sidebar-link"}
              aria-current={isActive(item.href) ? "page" : undefined}
            >
              {item.label}
            </Link>
          ))}
          <Link href="/get-started" className="sidebar-link sidebar-link-cta">
            Get Started
          </Link>
        </aside>

        <article className="content doc-page">
          <header className="page-header">
            <div className="page-meta-row">
              <span className="page-meta-chip">Documentation</span>
              <span className="page-meta-chip muted">Core: gitresolve</span>
            </div>
            <h1>{title}</h1>
            {subtitle ? <p>{subtitle}</p> : null}
          </header>
          {children}
        </article>
      </main>
    </div>
  );
}