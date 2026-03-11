INSERT INTO webhook_keys (environment, public_key, secret)
VALUES
  ('staging', 'whk_staging_demo_01', 'whsec_staging_demo_secret_01'),
  ('production', 'whk_prod_demo_01', 'whsec_prod_demo_secret_01')
ON CONFLICT (public_key) DO NOTHING;

INSERT INTO workflows (id, name, description, enabled, trigger_event, conditions, actions, webhook_key)
VALUES
(
  '8f478d6f-aede-4cc5-96f8-d2f7b2eced2f',
  'Lead Qualification Sync',
  'Route high-intent leads to CRM and send Slack notification.',
  TRUE,
  'lead.created',
  '[{"field":"lead.score","operator":"gt","value":70},{"field":"lead.region","operator":"eq","value":"us"}]'::jsonb,
  '[{"name":"push-to-crm","method":"POST","url":"https://httpbin.org/post","headers":{"X-Target":"crm"},"timeout_ms":8000,"body":{"channel":"crm"}},
    {"name":"notify-sales","method":"POST","url":"https://httpbin.org/post","headers":{"X-Target":"slack"},"timeout_ms":8000,"body":{"channel":"sales-alerts"}}]'::jsonb,
  'whk_staging_demo_01'
),
(
  '86ef78d4-4e37-46d8-bea8-6ac86d413860',
  'Invoice Recovery Trigger',
  'Escalate failed payments after the second retry.',
  TRUE,
  'invoice.failed',
  '[{"field":"invoice.retry_attempt","operator":"gt","value":1}]'::jsonb,
  '[{"name":"notify-recovery-system","method":"POST","url":"https://httpbin.org/status/500","headers":{"X-Target":"billing"},"timeout_ms":5000,"body":{"priority":"high"}}]'::jsonb,
  'whk_prod_demo_01'
)
ON CONFLICT (id) DO NOTHING;
