import { Panel } from "@/components/ui/panel";

type MetricCardProps = {
  label: string;
  value: string;
  helper: string;
};

export function MetricCard({ label, value, helper }: MetricCardProps) {
  return (
    <Panel>
      <p className="text-sm text-slate-400">{label}</p>
      <p className="mt-2 text-3xl font-semibold text-white">{value}</p>
      <p className="mt-2 text-xs text-slate-500">{helper}</p>
    </Panel>
  );
}
