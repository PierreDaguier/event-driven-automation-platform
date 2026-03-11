"use client";

import Link from "next/link";
import { motion } from "motion/react";

const valueProps = [
  "Secure webhook ingestion with HMAC + idempotency",
  "Retry/backoff queue processing with dead-letter visibility",
  "Business dashboard with success rate, latency, and execution timeline"
];

export default function LandingPage() {
  return (
    <main className="relative min-h-screen overflow-hidden bg-[radial-gradient(circle_at_10%_20%,#22d3ee33,transparent_40%),radial-gradient(circle_at_90%_80%,#a3e63533,transparent_35%),#030712]">
      <div className="mx-auto flex min-h-screen max-w-6xl flex-col justify-center px-6 py-16">
        <motion.p
          initial={{ opacity: 0, y: 8 }}
          animate={{ opacity: 1, y: 0 }}
          className="mb-4 text-xs uppercase tracking-[0.35em] text-cyan-300"
        >
          Event-Driven Automation Platform
        </motion.p>
        <motion.h1
          initial={{ opacity: 0, y: 12 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.05 }}
          className="max-w-4xl text-4xl font-semibold leading-tight text-white md:text-6xl"
        >
          Automate critical B2B workflows from webhook to outcome, with full execution traceability.
        </motion.h1>
        <motion.p
          initial={{ opacity: 0, y: 12 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.1 }}
          className="mt-6 max-w-2xl text-lg text-slate-300"
        >
          Production-grade demo / technical case study designed for client-facing conversations around reliability, latency, and operational visibility.
        </motion.p>

        <motion.div
          initial={{ opacity: 0, y: 12 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.2 }}
          className="mt-10 grid gap-4 md:grid-cols-3"
        >
          {valueProps.map((item) => (
            <article key={item} className="rounded-2xl border border-white/10 bg-slate-900/60 p-4 text-sm text-slate-200 backdrop-blur">
              {item}
            </article>
          ))}
        </motion.div>

        <motion.div
          initial={{ opacity: 0, y: 12 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.3 }}
          className="mt-10 flex flex-wrap gap-4"
        >
          <Link
            href="/dashboard"
            className="rounded-xl bg-cyan-400 px-5 py-3 font-medium text-slate-950 transition hover:bg-cyan-300 focus:outline-none focus-visible:ring-2 focus-visible:ring-cyan-200"
          >
            Open Live Dashboard
          </Link>
          <Link
            href="https://github.com"
            className="rounded-xl border border-white/20 px-5 py-3 font-medium text-slate-100 transition hover:border-cyan-300 hover:text-cyan-200"
          >
            View Engineering Artifacts
          </Link>
        </motion.div>
      </div>
    </main>
  );
}
