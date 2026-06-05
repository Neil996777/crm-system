import { expect, test } from '@playwright/test';

const adminEmail = process.env.E2E_ADMIN_EMAIL ?? 'admin@example.com';
const adminPassword = process.env.E2E_ADMIN_PASSWORD ?? 'AdminChangeMe-001!';
const salesPassword = 'SalesHistory-001!';

test.beforeEach(async ({ page }) => {
  await page.goto('/');
  await signIn(page, adminEmail, adminPassword);
  await expect(page.getByRole('heading', { name: '工作台' })).toBeVisible();
});

test('TEST-HISTORY-001 and TEST-HISTORY-004 shows read-only record-local history after a business mutation', async ({ page }) => {
  const companyName = `E2E History ${Date.now()}`;

  await page.getByRole('button', { name: '线索', exact: true }).click();
  await page.getByRole('button', { name: '新建线索' }).click();
  await page.getByLabel('公司名称').fill(companyName);
  await page.getByLabel('来源').fill('Website');
  await page.getByLabel('负责人 ID').fill('sales-1');
  await page.getByRole('button', { name: '保存线索' }).click();

  await page.getByRole('button', { name: companyName }).click();
  await page.getByRole('button', { name: '标记有效' }).click();
  await expect(page.getByLabel('线索详情').getByText('有效', { exact: true })).toBeVisible();
  const leadId = await selectedLeadId(page);
  await expect.poll(async () => page.evaluate(async (id) => {
    const response = await fetch(`/api/leads/${id}/history`, { credentials: 'include' });
    const body = await response.json();
    return (body.data?.events ?? []).map((event: { eventId: string }) => event.eventId).join(',');
  }, leadId)).toContain('EVT-LEAD-QUALIFIED');
  await reopenLead(page, companyName);

  const history = page.getByLabel('记录历史');
  await expect(history.getByRole('heading', { name: '历史' })).toBeVisible();
  await expect(history).toContainText('EVT-LEAD-QUALIFIED');
  await expect(history).toContainText('Lead qualified');
  const qualifiedEvent = history.locator('.timelineItem', { hasText: 'EVT-LEAD-QUALIFIED' });
  await expect(qualifiedEvent).toContainText('线索');
  await expect(qualifiedEvent.getByText(/操作者：usr_seed_admin/)).toBeVisible();
  await expect(qualifiedEvent.getByText(/资源：线索/)).toBeVisible();
  await expect(qualifiedEvent.getByText(/发生时间：/)).toBeVisible();
  await expect(qualifiedEvent.getByText(/变更前：/)).toBeVisible();
  await expect(qualifiedEvent.getByText(/变更后：/)).toBeVisible();
  await expect(qualifiedEvent.getByRole('button', { name: /编辑|删除|保存|edit|delete|save/i })).toHaveCount(0);
});

test('TEST-HISTORY-003 denies non-owned record-local history without leaking events', async ({ page }) => {
  const suffix = Date.now();
  const companyName = `E2E History Denied ${suffix}`;
  const salesEmail = `history-sales-${suffix}@example.com`;

  await page.evaluate(async ({ email, password }) => {
    const response = await fetch('/admin/users', {
      method: 'POST',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        email,
        displayName: 'History Sales',
        password,
        role: 'Sales'
      })
    });
    if (!response.ok) {
      throw new Error(`create sales user failed: ${response.status}`);
    }
  }, { email: salesEmail, password: salesPassword });

  await page.getByRole('button', { name: '线索', exact: true }).click();
  await page.getByRole('button', { name: '新建线索' }).click();
  await page.getByLabel('公司名称').fill(companyName);
  await page.getByLabel('来源').fill('Referral');
  await page.getByLabel('负责人 ID').fill('sales-1');
  await page.getByRole('button', { name: '保存线索' }).click();
  await page.getByRole('button', { name: companyName }).click();
  const leadId = await selectedLeadId(page);
  await page.getByRole('button', { name: '标记有效' }).click();
  await expect(page.getByLabel('线索详情').getByText('有效', { exact: true })).toBeVisible();

  await page.getByRole('button', { name: '退出登录' }).click();
  await signIn(page, salesEmail, salesPassword);
  await expect(page.locator('.topbar').getByText('History Sales')).toBeVisible();
  await expect(page.locator('.topbar').getByText('销售', { exact: true })).toBeVisible();

  const denied = await page.evaluate(async (id) => {
    const response = await fetch(`/api/leads/${id}/history`, { credentials: 'include' });
    return { status: response.status, body: await response.text() };
  }, leadId);
  expect(denied.status).toBe(404);
  expect(denied.body).not.toContain('Lead qualified as Valid');
});

async function signIn(page: import('@playwright/test').Page, email: string, password: string) {
  await page.getByLabel('邮箱').fill(email);
  await page.getByLabel('密码').fill(password);
  await page.getByRole('button', { name: '登录' }).click();
}

async function selectedLeadId(page: import('@playwright/test').Page) {
  const id = await page.locator('[data-record-id]').getAttribute('data-record-id');
  if (!id) {
    throw new Error('selected lead id not found');
  }
  return id;
}

async function reopenLead(page: import('@playwright/test').Page, companyName: string) {
  await page.reload();
  await expect(page.getByRole('heading', { name: '工作台' })).toBeVisible();
  await page.getByRole('button', { name: '线索', exact: true }).click();
  await page.getByLabel('搜索').fill(companyName);
  await page.getByRole('button', { name: '搜索', exact: true }).click();
  await page.getByRole('button', { name: companyName }).click();
}
