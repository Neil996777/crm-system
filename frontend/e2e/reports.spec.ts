import { expect, test } from '@playwright/test';

const adminEmail = process.env.E2E_ADMIN_EMAIL ?? 'admin@example.com';
const adminPassword = process.env.E2E_ADMIN_PASSWORD ?? 'AdminChangeMe-001!';
const reportPassword = 'Reports-001!';

test.beforeEach(async ({ page }) => {
  await page.goto('/');
  await page.getByLabel('邮箱').fill(adminEmail);
  await page.getByLabel('密码').fill(adminPassword);
  await page.getByRole('button', { name: '登录' }).click();
  await expect(page.getByRole('heading', { name: '工作台' })).toBeVisible();
});

test('TEST-BASIC-REPORT-002 renders basic reports empty state through gateway', async ({ page }) => {
  await page.getByLabel('主导航').getByRole('button', { name: '报表' }).click();
  await expect(page.getByRole('heading', { name: '团队报表' })).toBeVisible();
  await expect(page.getByRole('heading', { name: '经理团队总览' })).toHaveCount(0);
  await expect(page.locator('[data-uiux="reports-team"]')).toBeVisible();
  await expect(page.locator('.reportsKpiStrip .metricTile')).toHaveCount(9);
  await expect(page.locator('.reportsKpiStrip .currencyMetric')).toHaveCount(4);
  const overflowingKpis = await page.locator('.reportsKpiStrip .metricTile').evaluateAll((cards) => (
    cards.filter((card) => card.scrollWidth > card.clientWidth + 1).length
  ));
  expect(overflowingKpis).toBe(0);
  await expect(page.getByLabel('管道分析')).toBeVisible();
  await expect(page.getByLabel('管道分析').locator('.pipelineViz')).toBeVisible();
  await expect(page.getByLabel('负责人分组')).toBeVisible();
  await expect(page.getByLabel('负责人分组').locator('.dataTable')).toBeVisible();
  await expect(page.getByLabel('状态阶段分解').locator('.breakdownCard')).toHaveCount(5);
});

test('TEST-UIUX-A4-REPORT-001 sales cannot see reports and direct API access is denied without leakage', async ({ page }) => {
  const suffix = Date.now();
  const sales = await createUser(page, `report-sales-${suffix}@example.com`, `Report Sales ${suffix}`, 'Sales');

  await page.getByRole('button', { name: '退出登录' }).click();
  await signIn(page, sales.email, reportPassword);
  await expect(page.locator('.topbar').getByText(sales.displayName)).toBeVisible();
  await expect(page.getByRole('button', { name: '报表' })).toHaveCount(0);

  const teamDenied = await fetchText(page, '/api/reports/team-overview');
  expect(teamDenied.status).toBe(403);
  expect(teamDenied.body).not.toContain('leadCount');
  expect(teamDenied.body).not.toContain('contractAmount');

  const basicDenied = await fetchText(page, '/api/reports/sales-overview');
  expect(basicDenied.status).toBe(403);
  expect(basicDenied.body).not.toContain('leadCount');
  expect(basicDenied.body).not.toContain('contractAmount');
});

test('TEST-UIUX-A4-REPORT-002 manager uses team scope and administrator uses all scope', async ({ page }) => {
  const suffix = Date.now();
  const manager = await createUser(page, `report-manager-${suffix}@example.com`, `Report Manager ${suffix}`, 'Sales Manager');

  const adminBasic = await fetchJson(page, '/api/reports/sales-overview');
  expect(adminBasic.status).toBe(200);
  expect(adminBasic.body.data.scope).toBe('all');

  await page.getByRole('button', { name: '退出登录' }).click();
  await signIn(page, manager.email, reportPassword);
  await expect(page.locator('.topbar').getByText(manager.displayName)).toBeVisible();
  await expect(page.getByRole('button', { name: '报表' }).first()).toBeVisible();

  const managerOverview = await fetchJson(page, '/api/reports/team-overview');
  expect(managerOverview.status).toBe(200);
  expect(managerOverview.body.data.scope).toBe('team');
  const managerBasic = await fetchJson(page, '/api/reports/sales-overview');
  expect(managerBasic.status).toBe(200);
  expect(managerBasic.body.data.scope).toBe('team');
});

async function signIn(page: import('@playwright/test').Page, email: string, password: string) {
  await page.getByLabel('邮箱').fill(email);
  await page.getByLabel('密码').fill(password);
  await page.getByRole('button', { name: '登录' }).click();
}

async function createUser(page: import('@playwright/test').Page, email: string, displayName: string, role: 'Sales Manager' | 'Sales') {
  return page.evaluate(async ({ email, displayName, password, role }) => {
    const response = await fetch('/admin/users', {
      method: 'POST',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, displayName, password, role })
    });
    const body = await response.json();
    if (!response.ok) throw new Error(JSON.stringify(body));
    return body.user as { id: string; email: string; displayName: string; role: string; status: string };
  }, { email, displayName, password: reportPassword, role });
}

async function fetchText(page: import('@playwright/test').Page, path: string) {
  return page.evaluate(async (path) => {
    const response = await fetch(path, { credentials: 'include' });
    return { status: response.status, body: await response.text() };
  }, path);
}

async function fetchJson(page: import('@playwright/test').Page, path: string) {
  return page.evaluate(async (path) => {
    const response = await fetch(path, { credentials: 'include' });
    return { status: response.status, body: await response.json() };
  }, path);
}
