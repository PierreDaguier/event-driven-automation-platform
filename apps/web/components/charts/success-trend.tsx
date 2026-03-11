import { Panel } from "@/components/ui/panel";

export function SuccessTrend({ successRate }: { successRate: number }) {
  return (
    <Panel>
      <div className="mb-4 flex items-center justify-between">
        <h3 className="text-sm font-medium text-slate-200">Success Rate Trend</h3>
        <span className="text-xs text-slate-400">Last 24h</span>
      </div>
      <div className="space-y-3">
        {[72, 76, 81, 88, 93, successRate].map((v, i) => (
          <div key={i} className="flex items-center gap-2">
            <div className="h-2 flex-1 rounded-full bg-slate-800">
              <div
                className="h-2 rounded-full bg-gradient-to-r from-cyan-400 to-lime-300"
                style={{ width: `${v}%` }}
              />
            </div>
            <span className="w-12 text-right text-xs text-slate-300">{v}%</span>
          </div>
        ))}
      </div>
    </Panel>
  );
}
