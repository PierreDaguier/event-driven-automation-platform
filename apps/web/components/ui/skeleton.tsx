export function Skeleton({ className }: { className: string }) {
  return <div className={`animate-pulse rounded-xl bg-slate-800/70 ${className}`} aria-hidden />;
}
