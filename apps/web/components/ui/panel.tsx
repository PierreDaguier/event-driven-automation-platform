import { clsx } from "clsx";
import { PropsWithChildren } from "react";

type PanelProps = PropsWithChildren<{
  className?: string;
}>;

export function Panel({ className, children }: PanelProps) {
  return (
    <section
      className={clsx(
        "rounded-2xl border border-white/10 bg-slate-950/70 p-5 shadow-[0_0_0_1px_rgba(245,158,11,0.08),0_16px_48px_rgba(15,23,42,0.5)] backdrop-blur",
        className
      )}
    >
      {children}
    </section>
  );
}
