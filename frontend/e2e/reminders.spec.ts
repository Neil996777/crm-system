import { expect, test } from '@playwright/test';

const adminEmail = process.env.E2E_ADMIN_EMAIL ?? 'admin@example.com';
const adminPassword = process.env.E2E_ADMIN_PASSWORD ?? 'AdminChangeMe-001!';

type Quote = { id: string; opportunityId: string; customerId: string; amount: string; version: number };
type Contract = { id: string; version: number };

test.beforeEach(async ({ page }) => {
  await page.goto('/');
  await page.getByLabel('邮箱').fill(adminEmail);
  await page.getByLabel('密码').fill(adminPassword);
  await page.getByRole('button', { name: '登录' }).click();
  await expect(page.getByRole('heading', { name: '工作台' })).toBeVisible();
});

test('TEST-REMINDER-001/002/003 groups task contract and payment reminders by business date', async ({ page }) => {
  const key = `reminder_${Date.now()}`;
  await createReminderData(page, key);

  await page.getByRole('button', { name: '提醒中心' }).click();
  await expect(page.getByRole('heading', { name: '提醒中心' })).toBeVisible();
  await page.getByLabel('业务日期').fill('2026-06-02');
  await page.getByRole('button', { name: '刷新提醒', exact: true }).click();

  await expect(page.getByRole('heading', { name: '任务提醒' })).toBeVisible();
  await expect(page.getByText(`Reminder task ${key}`)).toBeVisible();
  await expect(page.getByRole('heading', { name: '合同提醒' })).toBeVisible();
  await expect(page.getByText(`opp_contract_${key}`)).toBeVisible();
  await expect(page.getByRole('heading', { name: '回款提醒' })).toBeVisible();
  await expect(page.getByText(`opp_payment_${key}`)).toBeVisible();
});

test('TEST-REMINDER-004 suppresses completed task reminders after refresh', async ({ page }) => {
  const task = await createTask(page, `suppress_${Date.now()}`);

  await page.getByRole('button', { name: '提醒中心' }).click();
  await page.getByLabel('业务日期').fill('2026-06-02');
  await page.getByRole('button', { name: '刷新提醒', exact: true }).click();
  await expect(page.getByText(task.title)).toBeVisible();

  await page.evaluate(async ({ id, version }) => {
    const response = await fetch(`/api/tasks/${id}/status`, {
      method: 'POST',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ toStatus: 'Completed', expectedVersion: version })
    });
    if (!response.ok) throw new Error('complete task failed');
  }, { id: task.id, version: task.version });

  await page.getByRole('button', { name: '刷新提醒', exact: true }).click();
  await expect(page.getByText(task.title)).toHaveCount(0);
});

async function createReminderData(page: import('@playwright/test').Page, key: string) {
  await createTask(page, key);
  await createPendingContract(page, `opp_contract_${key}`, '2026-01-01');
  const paymentContract = await createPendingContract(page, `opp_payment_${key}`, '2027-01-01');
  await request(page, `/api/contracts/${paymentContract.id}/payment-plans`, {
    method: 'POST',
    body: JSON.stringify({ dueAmount: '10000.00', dueDate: '2026-01-01', currency: 'CNY' })
  });
}

async function createTask(page: import('@playwright/test').Page, key: string) {
  return request<{ id: string; title: string; version: number }>(page, '/api/tasks', {
    method: 'POST',
    body: JSON.stringify({
      relatedType: 'Opportunity',
      relatedId: `opp_task_${key}`,
      title: `Reminder task ${key}`,
      dueDate: '2026-01-01',
      ownerId: 'sales-1'
    })
  });
}

async function createPendingContract(page: import('@playwright/test').Page, opportunityId: string, expectedSignedDate: string) {
  const quote = await request<Quote>(page, '/api/quotes', {
    method: 'POST',
    body: JSON.stringify({
      opportunityId,
      customerId: `acct_${opportunityId}`,
      amount: '10000.00',
      status: 'Draft',
      validityEnd: '2027-12-31',
      ownerId: 'sales-1'
    })
  });
  const sent = await request<Quote>(page, `/api/quotes/${quote.id}/status`, {
    method: 'POST',
    body: JSON.stringify({ expectedVersion: quote.version, toStatus: 'Sent' })
  });
  const accepted = await request<Quote>(page, `/api/quotes/${quote.id}/status`, {
    method: 'POST',
    body: JSON.stringify({ expectedVersion: sent.version, toStatus: 'Accepted' })
  });
  return request<Contract>(page, '/api/contracts', {
    method: 'POST',
    body: JSON.stringify({
      quoteId: accepted.id,
      opportunityId,
      customerId: accepted.customerId,
      amount: accepted.amount,
      status: 'Pending Signature',
      contractNote: 'Reminder E2E contract note',
      expectedSignedDate,
      ownerId: 'sales-1'
    })
  });
}

async function request<T>(page: import('@playwright/test').Page, path: string, init: RequestInit): Promise<T> {
  return page.evaluate(async ({ path, init }) => {
    const response = await fetch(path, {
      ...init,
      credentials: 'include',
      headers: { 'Content-Type': 'application/json', ...(init.headers ?? {}) }
    });
    const body = await response.json();
    if (!response.ok) throw new Error(body.error?.safeMessage ?? 'Request failed.');
    return body.data as T;
  }, { path, init });
}
