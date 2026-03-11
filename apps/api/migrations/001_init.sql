CREATE TABLE IF NOT EXISTS webhook_keys (
    id SERIAL PRIMARY KEY,
    environment TEXT NOT NULL,
    public_key TEXT NOT NULL UNIQUE,
    secret TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS workflows (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    trigger_event TEXT NOT NULL,
    conditions JSONB NOT NULL DEFAULT '[]'::jsonb,
    actions JSONB NOT NULL DEFAULT '[]'::jsonb,
    webhook_key TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS events (
    id UUID PRIMARY KEY,
    source TEXT NOT NULL,
    idempotency_key TEXT NOT NULL UNIQUE,
    event_type TEXT NOT NULL,
    payload JSONB NOT NULL,
    received_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS runs (
    id UUID PRIMARY KEY,
    workflow_id UUID NOT NULL REFERENCES workflows(id) ON DELETE CASCADE,
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    action_name TEXT NOT NULL,
    status TEXT NOT NULL,
    retries INTEGER NOT NULL DEFAULT 0,
    latency_ms INTEGER NOT NULL DEFAULT 0,
    error_msg TEXT NOT NULL DEFAULT '',
    request_preview JSONB NOT NULL DEFAULT '{}'::jsonb,
    response_preview TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS execution_logs (
    id BIGSERIAL PRIMARY KEY,
    run_id UUID NOT NULL REFERENCES runs(id) ON DELETE CASCADE,
    workflow_id UUID NOT NULL REFERENCES workflows(id) ON DELETE CASCADE,
    status TEXT NOT NULL,
    latency_ms INTEGER NOT NULL,
    request_preview JSONB NOT NULL,
    response_preview TEXT NOT NULL,
    error_msg TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_workflows_trigger ON workflows(trigger_event);
CREATE INDEX IF NOT EXISTS idx_runs_workflow ON runs(workflow_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_execution_logs_workflow ON execution_logs(workflow_id, created_at DESC);
