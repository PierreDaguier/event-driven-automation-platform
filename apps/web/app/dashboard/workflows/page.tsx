import Link from "next/link";
import { Panel } from "@/components/ui/panel";
import { getWorkflows } from "@/lib/api";

export default async function WorkflowsPage() {
  const workflows = await getWorkflows();

  return (
    <div className="space-y-6">
      <header>
        <h2 className="text-3xl font-semibold text-white">Workflows</h2>
        <p className="mt-2 text-sm text-slate-400">Trigger rules and downstream automation actions.</p>
      </header>

      {workflows.length === 0 ? (
        <Panel>
          <p className="text-sm text-slate-400">No workflow configured yet. Use seed data or create one via API.</p>
        </Panel>
      ) : (
        <div className="grid gap-4">
          {workflows.map((workflow) => (
            <Panel key={workflow.id} className="transition hover:border-cyan-300/40">
              <div className="flex flex-col gap-3 md:flex-row md:items-start md:justify-between">
                <div>
                  <h3 className="text-lg font-medium text-white">{workflow.name}</h3>
                  <p className="mt-1 text-sm text-slate-400">{workflow.description}</p>
                  <p className="mt-3 text-xs uppercase tracking-wide text-cyan-300">Trigger: {workflow.trigger}</p>
                </div>
                <div className="flex items-center gap-3 text-sm">
                  <span className={workflow.enabled ? "text-emerald-300" : "text-rose-300"}>{workflow.enabled ? "Active" : "Paused"}</span>
                  <Link
                    href={`/dashboard/workflows/${workflow.id}`}
                    className="rounded-lg border border-white/20 px-3 py-2 text-slate-100 transition hover:border-cyan-300 hover:text-cyan-200"
                  >
                    View Details
                  </Link>
                </div>
              </div>
            </Panel>
          ))}
        </div>
      )}
    </div>
  );
}
