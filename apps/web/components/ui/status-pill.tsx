import { clsx } from "clsx";

type Status = "pending" | "running" | "succeeded" | "failed";

export function StatusPill({ status }: { status: Status }) {
  return (
    <span
      className={clsx(
        "inline-flex items-center rounded-full px-3 py-1 text-xs font-medium uppercase tracking-wide",
        status === "succeeded" && "bg-emerald-400/15 text-emerald-300",
        status === "failed" && "bg-rose-400/15 text-rose-300",
        status === "running" && "bg-amber-300/15 text-amber-300",
        status === "pending" && "bg-slate-400/15 text-slate-300"
      )}
    >
      {status}
    </span>
  );
}
