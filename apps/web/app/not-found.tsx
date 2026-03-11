import Link from "next/link";

export default function NotFoundPage() {
  return (
    <main className="flex min-h-screen items-center justify-center bg-slate-950 px-6 text-center text-slate-200">
      <div>
        <h1 className="text-3xl font-semibold text-white">Page not found</h1>
        <p className="mt-2 text-sm text-slate-400">The requested resource does not exist.</p>
        <Link href="/dashboard" className="mt-6 inline-block rounded-lg border border-cyan-300/40 px-4 py-2 text-cyan-200">
          Back to dashboard
        </Link>
      </div>
    </main>
  );
}
