"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { motion } from "motion/react";

const navItems = [
  { href: "/dashboard", label: "Overview" },
  { href: "/dashboard/workflows", label: "Workflows" },
  { href: "/dashboard/logs", label: "Execution Logs" },
  { href: "/dashboard/settings", label: "Settings" }
];

export function AppShell({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();

  return (
    <div className="min-h-screen bg-[radial-gradient(circle_at_top_right,#22d3ee22,transparent_45%),radial-gradient(circle_at_bottom_left,#a3e63522,transparent_45%),#030712] text-slate-100">
      <div className="mx-auto flex w-full max-w-7xl flex-col gap-8 px-4 py-8 md:px-6 lg:flex-row lg:px-8">
        <aside className="w-full rounded-2xl border border-cyan-300/20 bg-slate-900/70 p-4 backdrop-blur lg:w-64 lg:sticky lg:top-6 lg:h-fit">
          <p className="text-xs uppercase tracking-[0.3em] text-cyan-300">Automation OS</p>
          <h1 className="mt-2 text-xl font-semibold">Event Platform</h1>
          <nav className="mt-6 space-y-2" aria-label="Primary">
            {navItems.map((item) => {
              const active = pathname === item.href || pathname.startsWith(`${item.href}/`);
              return (
                <Link
                  key={item.href}
                  href={item.href}
                  className={`block rounded-xl px-3 py-2 text-sm transition ${
                    active ? "bg-cyan-400/20 text-cyan-200" : "text-slate-300 hover:bg-white/5 hover:text-white"
                  }`}
                >
                  {item.label}
                </Link>
              );
            })}
          </nav>
          <p className="mt-8 text-xs text-slate-400">Production-grade demo / case study technique.</p>
        </aside>

        <motion.main
          initial={{ opacity: 0, y: 8 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.3 }}
          className="flex-1"
        >
          {children}
        </motion.main>
      </div>
    </div>
  );
}
