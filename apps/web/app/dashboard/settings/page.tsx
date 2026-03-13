import { Panel } from "@/components/ui/panel";
import { getSettings } from "@/lib/api";

export default async function SettingsPage() {
  const settings = await getSettings();

  return (
    <div className="space-y-6">
      <header>
        <h2 className="text-3xl font-semibold text-white">Settings</h2>
        <p className="mt-2 text-sm text-slate-400">Webhook keys and environment configuration.</p>
      </header>

      <Panel>
        <p className="text-xs uppercase tracking-wide text-amber-300">Environment</p>
        <p className="mt-2 text-lg text-white">{settings.environment}</p>
      </Panel>

      <Panel>
        <h3 className="text-sm uppercase tracking-wide text-amber-300">Webhook Keys</h3>
        <div className="mt-4 overflow-x-auto">
          <table className="min-w-full text-left text-sm">
            <thead className="text-xs uppercase tracking-wide text-slate-500">
              <tr>
                <th className="pb-3">Environment</th>
                <th className="pb-3">Public Key</th>
                <th className="pb-3">Secret</th>
                <th className="pb-3">Created</th>
              </tr>
            </thead>
            <tbody className="text-slate-200">
              {settings.webhook_keys.map((key) => (
                <tr key={key.id} className="border-t border-white/10">
                  <td className="py-3">{key.environment}</td>
                  <td className="py-3 code">{key.public_key}</td>
                  <td className="py-3 code">{key.secret}</td>
                  <td className="py-3">{new Date(key.created_at).toLocaleString()}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </Panel>
    </div>
  );
}
