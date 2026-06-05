import { expect, test } from '@playwright/test';

const adminEmail = process.env.E2E_ADMIN_EMAIL ?? 'admin@example.com';
const adminPassword = process.env.E2E_ADMIN_PASSWORD ?? 'AdminChangeMe-001!';
const salesPassword = 'SalesOplog-001!';

test.beforeEach(async ({ page }) => {
  await page.goto('/');
  await signIn(page, adminEmail, adminPassword);
  await expect(page.getByRole('heading', { name: '工作台' })).toBeVisible();
});

test('TEST-OPLOG-001/002/005 administrator sees read-only global operation logs', async ({ page }) => {
  const suffix = Date.now();
  await createSalesUser(page, `oplog-admin-${suffix}@example.com`, salesPassword, 'Oplog Admin Evidence');

  await page.getByRole('button', { name: '操作日志' }).click();
  await expect(page.getByRole('heading', { name: '操作日志' })).toBeVisible();

  const logTable = page.getByLabel('操作日志表');
  await expect(logTable).toContainText('EVT-USER-ADMIN-CHANGED');
  await expect(logTable).toContainText('新建用户');
  await expect(logTable).toContainText('usr_seed_admin');
  await expect(logTable).toContainText('用户');
  await expect(logTable).toContainText('成功');
  await expect(logTable).toContainText('变更前');
  await expect(logTable).toContainText('变更后');
  await expect(logTable.getByRole('button', { name: /编辑|删除|保存|edit|delete|save/i })).toHaveCount(0);
});

test('TEST-OPLOG-004 sales is denied global operation logs without leakage', async ({ page }) => {
  const suffix = Date.now();
  const salesEmail = `oplog-sales-${suffix}@example.com`;
  await createSalesUser(page, salesEmail, salesPassword, 'Oplog Sales');

  await page.getByRole('button', { name: '退出登录' }).click();
  await signIn(page, salesEmail, salesPassword);
  await expect(page.locator('.topbar').getByText('Oplog Sales')).toBeVisible();

  const denied = await page.evaluate(async () => {
    const response = await fetch('/api/operation-log', { credentials: 'include' });
    return { status: response.status, body: await response.text() };
  });
  expect(denied.status).toBe(403);
  expect(denied.body).not.toContain('EVT-USER-ADMIN-CHANGED');
  await expect(page.getByRole('button', { name: '操作日志' })).toHaveCount(0);
});

async function signIn(page: import('@playwright/test').Page, email: string, password: string) {
  await page.getByLabel('邮箱').fill(email);
  await page.getByLabel('密码').fill(password);
  await page.getByRole('button', { name: '登录' }).click();
}

async function createSalesUser(page: import('@playwright/test').Page, email: string, password: string, displayName: string) {
  await page.evaluate(async ({ email, password, displayName }) => {
    const response = await fetch('/admin/users', {
      method: 'POST',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, displayName, password, role: 'Sales' })
    });
    if (!response.ok) {
      throw new Error(`create sales user failed: ${response.status}`);
    }
  }, { email, password, displayName });
}
