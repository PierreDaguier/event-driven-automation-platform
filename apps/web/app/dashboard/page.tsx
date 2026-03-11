import { SuccessTrend } from "@/components/charts/success-trend";
import { MetricCard } from "@/components/ui/metric-card";
import { Panel } from "@/components/ui/panel";
import { StatusPill } from "@/components/ui/status-pill";
import { getLogs, getOverview } from "@/lib/api";

export default async function DashboardOverviewPage() {
  const [overview, logs] = await Promise.all([getOverview(), getLogs("", 8)]);

  return (
    <div className="space-y-6">
      <header>
        <h2 className="text-3xl font-semibold text-white">Operations Overview</h2>
        <p className="mt-2 text-sm text-slate-400">Monitor run throughput, reliability, and latency in one place.</p>
      </header>

      <section className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
        <MetricCard label="Runs (24h)" value={overview.total_runs.toLocaleString()} helper="All processed actions" />
        <MetricCard label="Success Rate" value={`${overview.success_rate.toFixed(1)}%`} helper="Successful outcomes / terminal runs" />
        <MetricCard label="Avg Latency" value={`${overview.avg_latency_ms.toFixed(0)} ms`} helper="HTTP action execution" />
        <MetricCard label="Failed Runs" value={overview.failed_runs.toString()} helper="Candidates for DLQ replay" />
      </section>

      <section className="grid gap-4 xl:grid-cols-3">
        <div className="xl:col-span-2">
          <Panel>
            <div className="mb-4 flex items-center justify-between">
              <h3 className="text-sm font-medium text-slate-200">Recent Executions</h3>
              <span className="text-xs text-slate-400">Last 8 actions</span>
            </div>
            {logs.length === 0 ? (
              <div className="rounded-xl border border-dashed border-slate-700 p-6 text-sm text-slate-400">
                No executions yet. Trigger a webhook to populate this timeline.
              </div>
            ) : (
              <ul className="space-y-3">
                {logs.map((log) => (
                  <li key={log.id} className="flex flex-col gap-2 rounded-xl border border-white/5 bg-slate-900/70 p-4 md:flex-row md:items-center md:justify-between">
                    <div>
                      <p className="text-sm text-slate-200">Run {log.run_id.slice(0, 8)}</p>
                      <p className="text-xs text-slate-500">{new Date(log.created_at).toLocaleString()}</p>
                    </div>
                    <div className="flex items-center gap-3">
                      <span className="text-xs text-slate-400">{log.latency_ms} ms</span>
                      <StatusPill status={log.status} />
                    </div>
                  </li>
                ))}
              </ul>
            )}
          </Panel>
        </div>

        <SuccessTrend successRate={overview.success_rate} />
      </section>
    </div>
  );
}
