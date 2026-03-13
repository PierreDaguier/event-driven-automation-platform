import { notFound } from "next/navigation";
import { Panel } from "@/components/ui/panel";
import { getWorkflowById } from "@/lib/api";

export default async function WorkflowDetailPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  const workflow = await getWorkflowById(id);

  if (!workflow) {
    notFound();
  }

  return (
    <div className="space-y-6">
      <header>
        <h2 className="text-3xl font-semibold text-white">{workflow.name}</h2>
        <p className="mt-2 text-sm text-slate-400">{workflow.description}</p>
      </header>

      <div className="grid gap-4 lg:grid-cols-2">
        <Panel>
          <h3 className="text-sm uppercase tracking-wide text-amber-300">Trigger & Conditions</h3>
          <p className="mt-3 text-sm text-slate-300">Trigger event: <span className="code text-amber-200">{workflow.trigger}</span></p>
          <ul className="mt-4 space-y-3">
            {workflow.conditions.map((condition) => (
              <li key={`${condition.field}-${condition.operator}`} className="rounded-xl border border-white/10 bg-slate-900/70 p-3 text-sm text-slate-200">
                <span className="code text-amber-200">{condition.field}</span> {condition.operator} <span className="code text-lime-200">{String(condition.value)}</span>
              </li>
            ))}
          </ul>
        </Panel>

        <Panel>
          <h3 className="text-sm uppercase tracking-wide text-amber-300">Actions</h3>
          <ul className="mt-4 space-y-3">
            {workflow.actions.map((action) => (
              <li key={action.name} className="rounded-xl border border-white/10 bg-slate-900/70 p-3 text-sm text-slate-200">
                <p className="font-medium text-white">{action.name}</p>
                <p className="mt-1 text-xs text-slate-400">
                  {action.method} <span className="code text-amber-200">{action.url}</span>
                </p>
              </li>
            ))}
          </ul>
        </Panel>
      </div>
    </div>
  );
}
