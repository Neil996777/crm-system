import { expect, test } from '@playwright/test';

const adminEmail = process.env.E2E_ADMIN_EMAIL ?? 'admin@example.com';
const adminPassword = process.env.E2E_ADMIN_PASSWORD ?? 'AdminChangeMe-001!';

test.beforeEach(async ({ page }) => {
  await page.goto('/');
  await page.getByLabel('邮箱').fill(adminEmail);
  await page.getByLabel('密码').fill(adminPassword);
  await page.getByRole('button', { name: '登录' }).click();
  await expect(page.getByRole('heading', { name: '工作台' })).toBeVisible();
});

test('TEST-ARCHIVE-001/003/004 TEST-INV-ARCHIVEBLOCK-001 and TEST-ABUSE-ARCHIVED-001 archives only after real open task obligation is resolved', async ({ page }) => {
  const suffix = Date.now();
  const companyName = `Archive E2E ${suffix}`;
  const account = await createAccount(page, companyName);
  const task = await createOpenTask(page, account.id, `Archive blocker ${suffix}`);

  await page.getByRole('button', { name: '公司/客户' }).click();
  await page.getByLabel('搜索').fill(companyName);
  await page.getByRole('button', { name: '搜索' }).click();
  await page.getByRole('button', { name: companyName }).click();

  await page.getByRole('button', { name: '归档', exact: true }).click();
  await expect(page.getByRole('alert')).toContainText('仍有未完成事项');
  await expect(page.getByText(`Archive blocker ${suffix}`)).toBeVisible();

  await completeTask(page, task.id, task.version);
  await page.getByRole('button', { name: '归档', exact: true }).click();
  await page.getByRole('button', { name: '确认归档' }).click();
  await expect(page.getByRole('button', { name: companyName })).toHaveCount(0);

  await page.getByLabel('包含已归档').check();
  await page.getByRole('button', { name: '搜索' }).click();
  await expect(page.getByRole('button', { name: companyName })).toBeVisible();
});

async function createAccount(page: import('@playwright/test').Page, companyName: string) {
  return page.evaluate(async (companyName) => {
    const response = await fetch('/api/accounts', {
      method: 'POST',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ companyName, customerStatus: 'Active', ownerId: 'sales-1' })
    });
    const body = await response.json();
    if (!response.ok) throw new Error(JSON.stringify(body));
    return body.data as { id: string; version: number };
  }, companyName);
}

async function createOpenTask(page: import('@playwright/test').Page, accountId: string, title: string) {
  return page.evaluate(async ({ accountId, title }) => {
    const response = await fetch('/api/tasks', {
      method: 'POST',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        relatedType: 'Customer',
        relatedId: accountId,
        title,
        dueDate: '2026-06-02',
        ownerId: 'sales-1'
      })
    });
    const body = await response.json();
    if (!response.ok) throw new Error(JSON.stringify(body));
    return body.data as { id: string; version: number };
  }, { accountId, title });
}

async function completeTask(page: import('@playwright/test').Page, taskId: string, expectedVersion: number) {
  await page.evaluate(async ({ taskId, expectedVersion }) => {
    const response = await fetch(`/api/tasks/${taskId}/status`, {
      method: 'POST',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ toStatus: 'Completed', expectedVersion })
    });
    const body = await response.json().catch(() => ({}));
    if (!response.ok) throw new Error(JSON.stringify(body));
  }, { taskId, expectedVersion });
}
