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

test('TEST-CSV-EXPORT-001 exports active authorized records after confirmation', async ({ page }) => {
  await page.getByRole('button', { name: 'Import/Export' }).click();
  await expect(page.getByRole('heading', { name: 'Import/Export' })).toBeVisible();
  await page.getByRole('checkbox', { name: 'Confirm export scope and audit log' }).check();
  await page.getByRole('button', { name: 'Start export' }).click();

  await expect(page.getByText(/Exported \d+ lead rows/)).toBeVisible();
  await expect(page.locator('.exportResult').getByText('Archived excluded')).toBeVisible();
});
