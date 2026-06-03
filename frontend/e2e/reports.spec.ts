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

test('TEST-BASIC-REPORT-002 renders basic reports empty state through gateway', async ({ page }) => {
  await page.getByRole('button', { name: 'Reports' }).click();
  await expect(page.getByRole('heading', { name: 'Basic Sales Reports' })).toBeVisible();
  await expect(page.locator('section[aria-label="Basic report metrics"]')).toBeVisible();
  await expect(page.getByText('Leads by Status')).toBeVisible();
  await expect(page.getByText('Payments by Status')).toBeVisible();
});
