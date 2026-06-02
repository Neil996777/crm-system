import { expect, test } from '@playwright/test';

const adminEmail = process.env.E2E_ADMIN_EMAIL ?? 'admin@example.com';
const adminPassword = process.env.E2E_ADMIN_PASSWORD ?? 'AdminChangeMe-001!';
const salesPassword = 'SalesOplog-001!';

test.beforeEach(async ({ page }) => {
  await page.goto('/');
  await signIn(page, adminEmail, adminPassword);
  await expect(page.getByRole('heading', { name: 'Work Overview' })).toBeVisible();
});

test('TEST-OPLOG-001/002/005 administrator sees read-only global operation logs', async ({ page }) => {
  const suffix = Date.now();
  await createSalesUser(page, `oplog-admin-${suffix}@example.com`, salesPassword, 'Oplog Admin Evidence');

  await page.getByRole('button', { name: 'Operation Logs' }).click();
  await expect(page.getByRole('heading', { name: 'Operation Logs' })).toBeVisible();

  const logTable = page.getByLabel('Operation log table');
  await expect(logTable).toContainText('EVT-USER-ADMIN-CHANGED');
  await expect(logTable).toContainText('create_user');
  await expect(logTable).toContainText('usr_seed_admin');
  await expect(logTable).toContainText('User');
  await expect(logTable).toContainText('success');
  await expect(logTable).toContainText('Before');
  await expect(logTable).toContainText('After');
  await expect(logTable.getByRole('button', { name: /edit|delete|save/i })).toHaveCount(0);
});

test('TEST-OPLOG-004 sales is denied global operation logs without leakage', async ({ page }) => {
  const suffix = Date.now();
  const salesEmail = `oplog-sales-${suffix}@example.com`;
  await createSalesUser(page, salesEmail, salesPassword, 'Oplog Sales');

  await page.getByRole('button', { name: 'Sign out' }).click();
  await signIn(page, salesEmail, salesPassword);
  await expect(page.locator('.topbar').getByText('Oplog Sales')).toBeVisible();

  const denied = await page.evaluate(async () => {
    const response = await fetch('/api/operation-log', { credentials: 'include' });
    return { status: response.status, body: await response.text() };
  });
  expect(denied.status).toBe(403);
  expect(denied.body).not.toContain('EVT-USER-ADMIN-CHANGED');
  await expect(page.getByRole('button', { name: 'Operation Logs' })).toHaveCount(0);
});

async function signIn(page: import('@playwright/test').Page, email: string, password: string) {
  await page.getByLabel('Email').fill(email);
  await page.getByLabel('Password').fill(password);
  await page.getByRole('button', { name: 'Sign in' }).click();
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
