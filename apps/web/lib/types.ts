export type Overview = {
  total_runs: number;
  success_rate: number;
  avg_latency_ms: number;
  failed_runs: number;
  pending_runs: number;
};

export type Condition = {
  field: string;
  operator: string;
  value: string | number;
};

export type WorkflowAction = {
  name: string;
  method: string;
  url: string;
};

export type Workflow = {
  id: string;
  name: string;
  description: string;
  enabled: boolean;
  trigger: string;
  conditions: Condition[];
  actions: WorkflowAction[];
  webhook_key: string;
  created_at: string;
  updated_at: string;
};

export type ExecutionLog = {
  id: number;
  run_id: string;
  workflow_id: string;
  status: "pending" | "running" | "succeeded" | "failed";
  latency_ms: number;
  request_preview: Record<string, unknown>;
  response_preview: string;
  error: string;
  created_at: string;
};

export type SettingsResponse = {
  environment: string;
  webhook_keys: Array<{
    id: number;
    environment: string;
    public_key: string;
    secret: string;
    created_at: string;
  }>;
};
