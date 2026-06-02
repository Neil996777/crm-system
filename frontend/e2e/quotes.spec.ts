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

test('TEST-QUOTE-LIFECYCLE-002 validates quote required fields', async ({ page }) => {
  await page.getByRole('button', { name: 'Quotes' }).click();
  await expect(page.getByRole('heading', { name: 'Quotes' })).toBeVisible();

  await page.getByRole('button', { name: 'New quote' }).click();
  await page.getByLabel('Opportunity ID').fill(`opp_quote_missing_${Date.now()}`);
  await page.getByLabel('Customer ID').fill(`acct_quote_missing_${Date.now()}`);
  await page.getByLabel('Owner ID').fill('sales-1');
  await page.getByRole('button', { name: 'Save quote' }).click();

  await expect(page.getByRole('alert')).toContainText('The quote input is invalid.');
});

test('TEST-QUOTE-LIFECYCLE-002 shows expired quote warning and blocks contract link', async ({ page }) => {
  const opportunityId = `opp_quote_expire_${Date.now()}`;

  await page.getByRole('button', { name: 'Quotes' }).click();
  await page.getByRole('button', { name: 'New quote' }).click();
  await page.getByLabel('Opportunity ID').fill(opportunityId);
  await page.getByLabel('Customer ID').fill(`acct_${opportunityId}`);
  await page.getByLabel('Amount').fill('8800.00');
  await page.getByLabel('Validity end').fill('2027-10-31');
  await page.getByLabel('Owner ID').fill('sales-1');
  await page.getByRole('button', { name: 'Save quote' }).click();

  await page.getByRole('button', { name: opportunityId }).click();
  await page.getByRole('button', { name: 'Expire', exact: true }).click();

  await expect(page.getByText('Status: Expired')).toBeVisible();
  await expect(page.getByRole('alert')).toContainText('Expired quote cannot be linked to a contract.');
  await expect(page.getByText('Contract link blocked')).toBeVisible();
});

test('TEST-QUOTE-ACCEPT-001 creates sends and accepts a quote with contract link indicator', async ({ page }) => {
  const opportunityId = `opp_quote_accept_${Date.now()}`;

  await page.getByRole('button', { name: 'Quotes' }).click();
  await page.getByRole('button', { name: 'New quote' }).click();
  await page.getByLabel('Opportunity ID').fill(opportunityId);
  await page.getByLabel('Customer ID').fill(`acct_${opportunityId}`);
  await page.getByLabel('Amount').fill('12000.00');
  await page.getByLabel('Validity end').fill('2027-09-30');
  await page.getByLabel('Owner ID').fill('sales-1');
  await page.getByRole('button', { name: 'Save quote' }).click();

  await page.getByRole('button', { name: opportunityId }).click();
  await expect(page.getByText('Status: Draft')).toBeVisible();
  await page.getByRole('button', { name: 'Send' }).click();
  await expect(page.getByText('Status: Sent')).toBeVisible();
  await page.getByRole('button', { name: 'Accept', exact: true }).click();
  await expect(page.getByText('Status: Accepted')).toBeVisible();
  await expect(page.getByText('Contract link available')).toBeVisible();
  await expect(page.getByRole('button', { name: 'Create another quote for this opportunity' })).toHaveCount(0);
});
