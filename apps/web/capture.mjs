import { chromium } from 'playwright-core';
import path from 'node:path';

const pages = [
  ['01-landing.png', 'http://localhost:3000/'],
  ['02-dashboard-overview.png', 'http://localhost:3000/dashboard'],
  ['03-workflows.png', 'http://localhost:3000/dashboard/workflows'],
  ['04-workflow-detail.png', 'http://localhost:3000/dashboard/workflows/8f478d6f-aede-4cc5-96f8-d2f7b2eced2f'],
  ['05-logs.png', 'http://localhost:3000/dashboard/logs'],
  ['06-settings.png', 'http://localhost:3000/dashboard/settings']
];

const browser = await chromium.launch({
  executablePath: '/usr/bin/chromium',
  headless: true,
  args: ['--no-sandbox', '--disable-gpu']
});

const context = await browser.newContext({ viewport: { width: 1440, height: 900 } });

for (const [name, url] of pages) {
  const page = await context.newPage();
  await page.goto(url, { waitUntil: 'networkidle', timeout: 45000 });
  await page.screenshot({ path: path.join('..', '..', 'docs', 'screenshots', name), fullPage: true });
  await page.close();
}

await browser.close();
