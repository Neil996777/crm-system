import { expect, test, type Page } from '@playwright/test';

const adminEmail = process.env.E2E_ADMIN_EMAIL ?? 'admin@example.com';
const adminPassword = process.env.E2E_ADMIN_PASSWORD ?? 'AdminChangeMe-001!';
const dashboardPassword = 'Dashboard-001!';

test.beforeEach(async ({ page }) => {
  await page.goto('/');
  await page.getByLabel('邮箱').fill(adminEmail);
  await page.getByLabel('密码').fill(adminPassword);
  await page.getByRole('button', { name: '登录' }).click();
  await expect(page.getByRole('heading', { name: '工作台' })).toBeVisible();
});

test('TEST-TEAM-OVERVIEW-003 renders manager overview empty state through gateway', async ({ page }) => {
  const overview = await page.evaluate(async () => {
    const response = await fetch('/api/reports/team-overview', { credentials: 'include' });
    const body = await response.json();
    return { status: response.status, scope: body.data?.scope, hasMetrics: Boolean(body.data?.metrics) };
  });
  expect(overview.status).toBe(200);
  expect(overview.scope).toBe('all');
  expect(overview.hasMetrics).toBe(true);

  await page.getByLabel('主导航').getByRole('button', { name: '报表' }).click();
  await expect(page.getByRole('heading', { name: '团队报表' })).toBeVisible();
  await expect(page.getByRole('heading', { name: '经理团队总览' })).toHaveCount(0);
  await expect(page.locator('[data-uiux="reports-team"]')).toBeVisible();
  await expect(page.locator('.reportsKpiStrip .metricTile')).toHaveCount(9);
  await expect(page.getByLabel('管道分析')).toBeVisible();
  await expect(page.getByLabel('负责人分组')).toBeVisible();
});

test('TEST-UIUX-DASHBOARD-001 renders manager dashboard 8 cards and per-card focus stage', async ({ page }) => {
  await expect(page.locator('[data-uiux="dashboard"]')).toBeVisible();
  await expect(page.locator('section[aria-label="今日实时战报"]')).toBeVisible();
  await expect(page.locator('.dashboardKpis .metricTile')).toHaveCount(4);
  await expect(page.locator('[aria-label="管端工作台数据卡"] [data-dashboard-card]')).toHaveCount(8);
  await expect(page.getByLabel('管端工作台数据卡')).toContainText('团队销售漏斗');
  await expect(page.getByLabel('管端工作台数据卡')).toContainText('商机阶段构成');
  await expect(page.getByLabel('管端工作台数据卡')).toContainText('团队赢单金额趋势');
  await expect(page.getByLabel('管端工作台数据卡')).toContainText('销售业绩榜');
  await expect(page.getByLabel('管端工作台数据卡')).toContainText('团队待办与预警');
  await expect(page.getByLabel('管端工作台数据卡')).toContainText('团队回款到账');
  await expect(page.getByLabel('管端工作台数据卡')).toContainText('重点商机');
  await expect(page.getByLabel('管端工作台数据卡')).toContainText('团队最近活动');
  await expect(page.getByRole('button', { name: '专注模式' })).toHaveCount(0);

  const managerFunnelCard = page.locator('[data-dashboard-card="funnel"]');
  await expect(managerFunnelCard).toHaveAttribute('role', 'button');
  await expect(managerFunnelCard).toHaveAttribute('tabindex', '0');
  await managerFunnelCard.click();
  await expect(page.locator('[data-uiux="dashboard-focus"]')).toBeVisible();
  await expect(page.locator('.shell.focusMode')).toBeVisible();
  await expect(page.getByRole('heading', { name: '团队销售漏斗' })).toBeVisible();
  await expect(page.getByLabel('折叠卡片').locator('.sideCard')).toHaveCount(7);
  await expect(page.getByRole('button', { name: /商机阶段构成/ })).toBeVisible();
  await page.getByRole('button', { name: /商机阶段构成/ }).click();
  await expect(page.getByRole('heading', { name: '团队商机阶段构成' })).toBeVisible();
  await page.keyboard.press('Escape');
  await expect(page.locator('[data-uiux="dashboard"]')).toBeVisible();
  await expect(page.locator('.shell.focusMode')).toHaveCount(0);
  await page.locator('[data-dashboard-card="activity"]').click();
  await expect(page.locator('[data-uiux="dashboard-focus"]')).toBeVisible();
  await page.getByRole('button', { name: '返回' }).click();
  await expect(page.locator('[data-uiux="dashboard"]')).toBeVisible();
});

test('TEST-UIUX-A6-001 dashboard remains stable on narrow desktop viewport', async ({ page }) => {
  await page.setViewportSize({ width: 1440, height: 900 });
  await expect(page.locator('[data-uiux="dashboard"]')).toBeVisible();
  await expectDashboardCardsNotClipped(page);

  await page.setViewportSize({ width: 1680, height: 900 });
  await expect(page.locator('[data-uiux="dashboard"]')).toBeVisible();
  await expectDashboardCardsNotClipped(page);

  await page.setViewportSize({ width: 900, height: 720 });
  await expect(page.locator('[data-uiux="dashboard"]')).toBeVisible();
  await expect(page.locator('section[aria-label="今日实时战报"]')).toBeVisible();
  await expect(page.locator('[data-dashboard-card]')).toHaveCount(8);
  const overflow = await page.evaluate(() => document.documentElement.scrollWidth - document.documentElement.clientWidth);
  expect(overflow).toBeLessThanOrEqual(1);
});

test('TEST-UIUX-DASHBOARD-002 renders sales personal dashboard variant without manager cards', async ({ page }) => {
  const suffix = Date.now();
  const sales = await createUser(page, `dashboard-sales-${suffix}@example.com`, `Dashboard Sales ${suffix}`, 'Sales');

  await page.getByRole('button', { name: '退出登录' }).click();
  await signIn(page, sales.email, dashboardPassword);
  await expect(page.getByRole('heading', { name: '我的工作台' })).toBeVisible();
  await expect(page.locator('[data-uiux="dashboard"]')).toBeVisible();
  await expect(page.locator('.dashboardKpis .metricTile')).toHaveCount(4);
  await expect(page.locator('[aria-label="销售工作台数据卡"] [data-dashboard-card]')).toHaveCount(6);
  await expect(page.getByLabel('销售工作台数据卡')).toContainText('我的销售漏斗');
  await expect(page.getByLabel('销售工作台数据卡')).toContainText('我的待办与预警');
  await expect(page.getByLabel('销售工作台数据卡')).toContainText('我的商机阶段构成');
  await expect(page.getByLabel('销售工作台数据卡')).toContainText('我的赢单金额趋势');
  await expect(page.getByLabel('销售工作台数据卡')).toContainText('我的回款到账');
  await expect(page.getByLabel('销售工作台数据卡')).toContainText('我的最近活动');
  await expect(page.getByLabel('销售工作台数据卡')).not.toContainText('销售业绩榜');

  await page.locator('[data-dashboard-card="funnel"]').click();
  await expect(page.locator('[data-uiux="dashboard-focus"]')).toBeVisible();
  await expect(page.getByRole('heading', { name: '我的销售漏斗' })).toBeVisible();
  await expect(page.getByLabel('折叠卡片').locator('.sideCard')).toHaveCount(5);
  await page.keyboard.press('Escape');
  await expect(page.locator('[data-uiux="dashboard"]')).toBeVisible();
});

test('TEST-UIUX-A7-001 card focus respects reduced-motion mode and still snaps between states', async ({ page }) => {
  await page.emulateMedia({ reducedMotion: 'reduce' });
  await expect(page.locator('[data-uiux="dashboard"]')).toBeVisible();
  const paymentsCard = page.locator('[data-dashboard-card="payments"]');
  await paymentsCard.focus();
  await expect(paymentsCard).toBeFocused();
  await page.keyboard.press(' ');
  await expect(page.locator('[data-uiux="dashboard-focus"]')).toBeVisible();
  await expect(page.locator('.shell.focusMode')).toBeVisible();
  await expect(page.getByLabel('折叠卡片').locator('.sideCard')).toHaveCount(7);
  await page.keyboard.press('Escape');
  await expect(page.locator('[data-uiux="dashboard"]')).toBeVisible();
});

test('TEST-UIUX-A5-001 main navigation is keyboard reachable', async ({ page }) => {
  const leadsNav = page.getByLabel('主导航').getByRole('button', { name: '线索' });
  await leadsNav.focus();
  await expect(leadsNav).toBeFocused();
  await page.keyboard.press('Enter');
  await expect(page.getByRole('heading', { name: '线索', exact: true })).toBeVisible();
});

async function signIn(page: import('@playwright/test').Page, email: string, password: string) {
  await page.getByLabel('邮箱').fill(email);
  await page.getByLabel('密码').fill(password);
  await page.getByRole('button', { name: '登录' }).click();
}

async function expectDashboardCardsNotClipped(page: Page) {
  const clippedCards = await page.locator('[data-dashboard-card]').evaluateAll((cards) => (
    cards
      .map((card) => {
        const element = card as HTMLElement;
        return {
          key: element.dataset.dashboardCard ?? '',
          verticalOverflow: element.scrollHeight - element.clientHeight,
          horizontalOverflow: element.scrollWidth - element.clientWidth
        };
      })
      .filter((card) => card.verticalOverflow > 1 || card.horizontalOverflow > 1)
  ));
  expect(clippedCards).toEqual([]);
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
  }, { email, displayName, password: dashboardPassword, role });
}
