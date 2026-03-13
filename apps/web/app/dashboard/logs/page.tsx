import { Panel } from "@/components/ui/panel";
import { StatusPill } from "@/components/ui/status-pill";
import { getLogs } from "@/lib/api";

export default async function LogsPage() {
  const logs = await getLogs();

  return (
    <div className="space-y-6">
      <header>
        <h2 className="text-3xl font-semibold text-white">Execution Timeline</h2>
        <p className="mt-2 text-sm text-slate-400">Complete audit trail of action runs with latency and redacted payload preview.</p>
      </header>

      <Panel>
        {logs.length === 0 ? (
          <div className="rounded-xl border border-dashed border-slate-700 p-8 text-center text-sm text-slate-400">
            No logs available. Trigger a workflow to see timeline entries.
          </div>
        ) : (
          <ol className="space-y-4">
            {logs.map((log) => (
              <li key={log.id} className="rounded-xl border border-white/10 bg-slate-900/70 p-4">
                <div className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
                  <div>
                    <p className="text-sm text-white">Run {log.run_id.slice(0, 8)} • Workflow {log.workflow_id.slice(0, 8)}</p>
                    <p className="mt-1 text-xs text-slate-500">{new Date(log.created_at).toLocaleString()}</p>
                  </div>
                  <div className="flex items-center gap-3">
                    <span className="text-xs text-slate-400">{log.latency_ms} ms</span>
                    <StatusPill status={log.status} />
                  </div>
                </div>
                <details className="mt-3 rounded-lg border border-white/10 bg-slate-950/80 p-3">
                  <summary className="cursor-pointer text-xs text-amber-300">Payload preview (redacted)</summary>
                  <pre className="mt-2 overflow-x-auto text-xs text-slate-300">{JSON.stringify(log.request_preview, null, 2)}</pre>
                </details>
                {log.error && <p className="mt-3 text-xs text-rose-300">Error: {log.error}</p>}
              </li>
            ))}
          </ol>
        )}
      </Panel>
    </div>
  );
}
