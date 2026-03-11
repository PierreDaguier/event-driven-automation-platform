import { chromium } from 'playwright-core';
import path from 'node:path';

const pages = [
  ['07-mobile-landing.png', 'http://localhost:3000/'],
  ['08-mobile-dashboard.png', 'http://localhost:3000/dashboard']
];

const browser = await chromium.launch({
  executablePath: '/usr/bin/chromium',
  headless: true,
  args: ['--no-sandbox', '--disable-gpu']
});

const context = await browser.newContext({ viewport: { width: 390, height: 844 } });
for (const [name, url] of pages) {
  const page = await context.newPage();
  await page.goto(url, { waitUntil: 'networkidle', timeout: 45000 });
  await page.screenshot({ path: path.join('..', '..', 'docs', 'screenshots', name), fullPage: true });
  await page.close();
}
await browser.close();
