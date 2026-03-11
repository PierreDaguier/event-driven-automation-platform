import { mockLogs, mockOverview, mockSettings, mockWorkflows } from "@/lib/mock-data";
import type { ExecutionLog, Overview, SettingsResponse, Workflow } from "@/lib/types";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";

async function fetchJSON<T>(path: string, fallback: T): Promise<T> {
  try {
    const response = await fetch(`${API_BASE_URL}${path}`, {
      cache: "no-store"
    });
    if (!response.ok) {
      return fallback;
    }
    return (await response.json()) as T;
  } catch {
    return fallback;
  }
}

export function getOverview(): Promise<Overview> {
  return fetchJSON<Overview>("/api/v1/overview", mockOverview);
}

export function getWorkflows(): Promise<Workflow[]> {
  return fetchJSON<Workflow[]>("/api/v1/workflows", mockWorkflows);
}

export async function getWorkflowById(id: string): Promise<Workflow | null> {
  const workflows = await getWorkflows();
  const direct = workflows.find((workflow) => workflow.id === id);
  if (direct) {
    return direct;
  }
  return fetchJSON<Workflow | null>(`/api/v1/workflows/${id}`, null);
}

export function getLogs(workflowId = "", limit = 100): Promise<ExecutionLog[]> {
  const query = workflowId ? `?workflow_id=${workflowId}&limit=${limit}` : `?limit=${limit}`;
  return fetchJSON<ExecutionLog[]>(`/api/v1/logs${query}`, mockLogs);
}

export function getSettings(): Promise<SettingsResponse> {
  return fetchJSON<SettingsResponse>("/api/v1/settings", mockSettings);
}
