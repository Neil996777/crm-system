import { expect, test } from '@playwright/test';

const adminEmail = process.env.E2E_ADMIN_EMAIL ?? 'admin@example.com';
const adminPassword = process.env.E2E_ADMIN_PASSWORD ?? 'AdminChangeMe-001!';
const salesPassword = 'SalesHistory-001!';

test.beforeEach(async ({ page }) => {
  await page.goto('/');
  await signIn(page, adminEmail, adminPassword);
  await expect(page.getByRole('heading', { name: 'Work Overview' })).toBeVisible();
});

test('TEST-HISTORY-001 and TEST-HISTORY-004 shows read-only record-local history after a business mutation', async ({ page }) => {
  const companyName = `E2E History ${Date.now()}`;

  await page.getByRole('button', { name: 'Leads' }).click();
  await page.getByRole('button', { name: 'New lead' }).click();
  await page.getByLabel('Company name').fill(companyName);
  await page.getByLabel('Source').fill('Website');
  await page.getByLabel('Owner ID').fill('sales-1');
  await page.getByRole('button', { name: 'Save lead' }).click();

  await page.getByRole('button', { name: companyName }).click();
  await page.getByRole('button', { name: 'Qualify valid' }).click();
  await expect(page.getByLabel('Lead detail').getByText('Valid', { exact: true })).toBeVisible();

  const history = page.getByLabel('Record history');
  await expect(history.getByRole('heading', { name: 'History' })).toBeVisible();
  await expect(history).toContainText('EVT-LEAD-QUALIFIED');
  await expect(history).toContainText('Lead qualified as Valid');
  await expect(history).toContainText('Lead');
  await expect(history.getByText(/Actor: usr_seed_admin/)).toBeVisible();
  await expect(history.getByText(/Resource: Lead/)).toBeVisible();
  await expect(history.getByText(/Occurred:/)).toBeVisible();
  await expect(history.getByText(/Before:/)).toBeVisible();
  await expect(history.getByText(/After:/)).toBeVisible();
  await expect(history.getByRole('button', { name: /edit|delete|save/i })).toHaveCount(0);
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

  await page.getByRole('button', { name: 'Leads' }).click();
  await page.getByRole('button', { name: 'New lead' }).click();
  await page.getByLabel('Company name').fill(companyName);
  await page.getByLabel('Source').fill('Referral');
  await page.getByLabel('Owner ID').fill('sales-1');
  await page.getByRole('button', { name: 'Save lead' }).click();
  await page.getByRole('button', { name: companyName }).click();
  const leadId = await selectedLeadId(page);
  await page.getByRole('button', { name: 'Qualify valid' }).click();
  await expect(page.getByLabel('Lead detail').getByText('Valid', { exact: true })).toBeVisible();

  await page.getByRole('button', { name: 'Sign out' }).click();
  await signIn(page, salesEmail, salesPassword);
  await expect(page.locator('.topbar').getByText('History Sales')).toBeVisible();
  await expect(page.locator('.topbar').getByText('Sales', { exact: true })).toBeVisible();

  const denied = await page.evaluate(async (id) => {
    const response = await fetch(`/api/leads/${id}/history`, { credentials: 'include' });
    return { status: response.status, body: await response.text() };
  }, leadId);
  expect(denied.status).toBe(403);
  expect(denied.body).not.toContain('Lead qualified as Valid');
});

async function signIn(page: import('@playwright/test').Page, email: string, password: string) {
  await page.getByLabel('Email').fill(email);
  await page.getByLabel('Password').fill(password);
  await page.getByRole('button', { name: 'Sign in' }).click();
}

async function selectedLeadId(page: import('@playwright/test').Page) {
  const id = await page.locator('[data-record-id]').getAttribute('data-record-id');
  if (!id) {
    throw new Error('selected lead id not found');
  }
  return id;
}
