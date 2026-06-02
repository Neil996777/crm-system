import { expect, test } from '@playwright/test';

const adminEmail = process.env.E2E_ADMIN_EMAIL ?? 'admin@example.com';
const adminPassword = process.env.E2E_ADMIN_PASSWORD ?? 'AdminChangeMe-001!';

test.beforeEach(async ({ page }) => {
  await page.goto('/');
  await page.getByLabel('Email').fill(adminEmail);
  await page.getByLabel('Password').fill(adminPassword);
  await page.getByRole('button', { name: 'Sign in' }).click();
  await expect(page.getByRole('heading', { name: 'Work Overview' })).toBeVisible();
});

test('TEST-OPP-STAGE-002 shows backend-backed blocked transition alert', async ({ page }) => {
  const title = `E2E Blocked Stage ${Date.now()}`;
  await createOpportunity(page, title);
  await page.getByRole('button', { name: title }).click();

  await page.getByRole('button', { name: 'Quote', exact: true }).click();
  await expect(page.getByRole('alert')).toContainText('The requested stage transition is not allowed.');
  await expect(page.getByText('Current stage: New Opportunity')).toBeVisible();
});

test('TEST-OPP-CLOSE-002 blocks Won until related contract is Signed', async ({ page }) => {
  const title = `E2E Early Won ${Date.now()}`;
  await createOpportunity(page, title);
  await page.getByRole('button', { name: title }).click();

  await page.getByRole('button', { name: 'Close Won' }).click();
  await page.getByLabel('Contract ID').fill('contract_missing_e2e');
  await page.getByLabel('Close date').fill('2027-07-01');
  await page.getByRole('button', { name: 'Confirm Won' }).click();

  await expect(page.getByRole('alert')).toContainText('Won requires a Signed related contract.');
  await expect(page.getByText('Current stage: New Opportunity')).toBeVisible();
});

test('TEST-OPP-CLOSE-003 closes Lost with reason and terminal detail is read-only', async ({ page }) => {
  const title = `E2E Lost ${Date.now()}`;
  await createOpportunity(page, title);
  await page.getByRole('button', { name: title }).click();

  await page.getByRole('button', { name: 'Close Lost' }).click();
  await page.getByLabel('Close date').fill('2027-07-02');
  await page.getByLabel('Lost reason').selectOption('PRICE');
  await page.getByLabel('Reason detail').fill('Competitor pricing');
  await page.getByRole('button', { name: 'Confirm Lost' }).click();

  await expect(page.getByText('Current stage: Lost')).toBeVisible();
  await expect(page.getByText('Terminal record')).toBeVisible();
  await expect(page.getByRole('button', { name: 'Close Lost' })).toBeDisabled();
});

async function createOpportunity(page: import('@playwright/test').Page, title: string) {
  await page.getByRole('button', { name: 'Opportunities' }).click();
  await expect(page.getByRole('heading', { name: 'Opportunities' })).toBeVisible();
  await page.getByRole('button', { name: 'New opportunity', exact: true }).click();
  await page.getByLabel('Title').fill(title);
  await page.getByLabel('Customer ID').fill(`acct_${Date.now()}`);
  await page.getByLabel('Owner ID').fill('sales-1');
  await page.getByLabel('Expected amount').fill('10000.00');
  await page.getByLabel('Expected close date').fill('2027-06-30');
  await page.getByRole('button', { name: 'Save opportunity' }).click();
  await expect(page.getByRole('button', { name: title })).toBeVisible();
}
