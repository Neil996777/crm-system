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

test('TEST-LEAD-CREATE-002 validates create lead required fields', async ({ page }) => {
  await page.getByRole('button', { name: 'Leads' }).click();
  await expect(page.getByRole('heading', { name: 'Leads' })).toBeVisible();

  await page.getByRole('button', { name: 'New lead' }).click();
  await page.getByLabel('Company name').fill('E2E Missing Source');
  await page.getByRole('button', { name: 'Save lead' }).click();

  await expect(page.getByRole('alert')).toContainText('The lead input is invalid.');
});

test('TEST-LEAD-QUALIFY-003 converts a valid lead through the UI', async ({ page }) => {
  const companyName = `E2E Convert ${Date.now()}`;

  await page.getByRole('button', { name: 'Leads' }).click();
  await page.getByRole('button', { name: 'New lead' }).click();
  await page.getByLabel('Company name').fill(companyName);
  await page.getByLabel('Source').fill('Website');
  await page.getByLabel('Owner ID').fill('sales-1');
  await page.getByRole('button', { name: 'Save lead' }).click();

  await page.getByRole('button', { name: companyName }).click();
  await page.getByRole('button', { name: 'Qualify valid' }).click();
  await expect(page.getByLabel('Lead detail').getByText('Valid', { exact: true })).toBeVisible();

  await page.getByRole('button', { name: 'Convert lead' }).click();
  await page.getByLabel('Expected amount').fill('99000.00');
  await page.getByLabel('Expected close date').fill('2026-12-15');
  await page.getByRole('button', { name: 'Convert', exact: true }).click();

  await expect(page.getByLabel('Lead detail').getByText('Converted To Opportunity', { exact: true })).toBeVisible();
  await expect(page.getByLabel('Lead detail').getByText(/Opportunity: opp_/)).toBeVisible();
});

test('TEST-LEAD-QUALIFY-004 shows Unassigned qualification as unavailable and backend-denied', async ({ page }) => {
  const companyName = `E2E Unassigned ${Date.now()}`;

  await page.getByRole('button', { name: 'Leads' }).click();
  await page.getByRole('button', { name: 'New lead' }).click();
  await page.getByLabel('Company name').fill(companyName);
  await page.getByLabel('Source').fill('Referral');
  await page.getByRole('button', { name: 'Save lead' }).click();

  await page.getByRole('button', { name: companyName }).click();
  await expect(page.getByRole('button', { name: 'Qualify valid' })).toBeDisabled();
  await expect(page.getByText('Unassigned leads cannot be qualified or converted.')).toBeVisible();
});
