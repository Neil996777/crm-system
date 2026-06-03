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

test('TEST-CSV-IMPORT-001/002 imports valid CSV rows and shows row errors', async ({ page }) => {
  const suffix = Date.now();
  const companyName = `Import E2E ${suffix}`;
  const csv = `companyName,leadName,source,ownerId\n${companyName},Imported Lead,Website,sales-1\nBroken Lead,, ,sales-1\n`;

  await page.getByRole('button', { name: 'Import/Export' }).click();
  await expect(page.getByRole('heading', { name: 'Import/Export' })).toBeVisible();
  await page.getByLabel('Object type').selectOption('lead');
  await page.getByLabel('CSV file').setInputFiles({
    name: 'leads.csv',
    mimeType: 'text/csv',
    buffer: Buffer.from(csv)
  });
  await page.getByRole('button', { name: 'Start import' }).click();

  await expect(page.getByText('Imported 1 of 2 rows')).toBeVisible();
  await expect(page.getByText('Row 3')).toBeVisible();

  await page.getByRole('button', { name: 'Leads' }).click();
  await expect(page.getByText(companyName)).toBeVisible();
});
