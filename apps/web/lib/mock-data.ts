import type { ExecutionLog, Overview, SettingsResponse, Workflow } from "@/lib/types";

export const mockOverview: Overview = {
  total_runs: 248,
  success_rate: 96.4,
  avg_latency_ms: 184,
  failed_runs: 9,
  pending_runs: 4
};

export const mockWorkflows: Workflow[] = [
  {
    id: "8f478d6f-aede-4cc5-96f8-d2f7b2eced2f",
    name: "Lead Qualification Sync",
    description: "Route high-intent leads to CRM + sales alerting.",
    enabled: true,
    trigger: "lead.created",
    webhook_key: "whk_staging_demo_01",
    conditions: [
      { field: "lead.score", operator: "gt", value: 70 },
      { field: "lead.region", operator: "eq", value: "us" }
    ],
    actions: [
      { name: "push-to-crm", method: "POST", url: "https://httpbin.org/post" },
      { name: "notify-sales", method: "POST", url: "https://httpbin.org/post" }
    ],
    created_at: "2026-03-10T12:00:00Z",
    updated_at: "2026-03-10T12:00:00Z"
  },
  {
    id: "86ef78d4-4e37-46d8-bea8-6ac86d413860",
    name: "Invoice Recovery Trigger",
    description: "Escalate failed payments after second retry.",
    enabled: true,
    trigger: "invoice.failed",
    webhook_key: "whk_prod_demo_01",
    conditions: [{ field: "invoice.retry_attempt", operator: "gt", value: 1 }],
    actions: [{ name: "notify-recovery-system", method: "POST", url: "https://httpbin.org/status/500" }],
    created_at: "2026-03-10T12:00:00Z",
    updated_at: "2026-03-10T12:00:00Z"
  }
];

export const mockLogs: ExecutionLog[] = [
  {
    id: 1004,
    run_id: "beef9a17-b0f6-4f95-b7d0-04510ce4b548",
    workflow_id: mockWorkflows[0].id,
    status: "succeeded",
    latency_ms: 139,
    request_preview: { payload: { lead: { score: 84, email: "***redacted***" } } },
    response_preview: "200 OK",
    error: "",
    created_at: "2026-03-11T01:20:00Z"
  },
  {
    id: 1003,
    run_id: "75f024be-70d7-4401-99f1-4f93e2e4f4b9",
    workflow_id: mockWorkflows[1].id,
    status: "failed",
    latency_ms: 5034,
    request_preview: { payload: { invoice: { id: "INV-2991" } } },
    response_preview: "500 Internal Server Error",
    error: "downstream status 500",
    created_at: "2026-03-11T01:10:00Z"
  }
];

export const mockSettings: SettingsResponse = {
  environment: "staging",
  webhook_keys: [
    {
      id: 1,
      environment: "staging",
      public_key: "whk_staging_demo_01",
      secret: "whse*****************t_01",
      created_at: "2026-03-10T10:00:00Z"
    },
    {
      id: 2,
      environment: "production",
      public_key: "whk_prod_demo_01",
      secret: "whse*****************t_01",
      created_at: "2026-03-10T10:05:00Z"
    }
  ]
};
