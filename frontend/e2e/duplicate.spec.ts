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

test('TEST-DUPLICATE-WARN-001/005 account duplicate warning proceeds without merge', async ({ page }) => {
  const companyName = `E2E Duplicate Account ${Date.now()}`;

  await page.getByRole('button', { name: 'Companies/Customers' }).click();
  await page.getByRole('button', { name: 'New customer' }).click();
  await page.getByLabel('Company name').fill(companyName);
  await page.getByLabel('Customer status').fill('Prospect');
  await page.getByLabel('Owner ID').fill('sales-1');
  await page.getByRole('button', { name: 'Save customer' }).click();
  await expect(page.getByRole('button', { name: companyName })).toBeVisible();

  await page.getByRole('button', { name: 'New customer' }).click();
  await page.getByLabel('Company name').fill(`  ${companyName.toUpperCase()}  `);
  await page.getByLabel('Customer status').fill('Prospect');
  await page.getByLabel('Owner ID').fill('sales-1');
  await page.getByRole('button', { name: 'Save customer' }).click();

  await expect(page.getByRole('alert')).toContainText('Possible duplicate');
  await page.getByRole('button', { name: 'Create anyway' }).click();
  await expect(page.getByRole('button', { name: companyName })).toHaveCount(2);
});

test('TEST-DUPLICATE-WARN-004 lead duplicate warning and unique no-warning path', async ({ page }) => {
  const suffix = Date.now();
  const companyName = `E2E Duplicate Lead ${suffix}`;
  const uniqueCompany = `E2E Unique Lead ${suffix}`;
  const email = `leaddup-${suffix}@example.com`;
  const phone = `139${String(suffix).slice(-8)}`;

  await page.getByRole('button', { name: 'Leads' }).click();
  await page.getByRole('button', { name: 'New lead' }).click();
  await page.getByLabel('Company name').fill(companyName);
  await page.getByLabel('Source').fill('Website');
  await page.getByLabel('Owner ID').fill('sales-1');
  await page.getByLabel('Email').fill(email.toUpperCase());
  await page.getByLabel('Phone').fill(`+86 ${phone}`);
  await page.getByRole('button', { name: 'Save lead' }).click();
  await expect(page.getByRole('button', { name: companyName })).toBeVisible();

  await page.getByRole('button', { name: 'New lead' }).click();
  await page.getByLabel('Company name').fill(` ${companyName.toLowerCase()} `);
  await page.getByLabel('Source').fill('Website');
  await page.getByLabel('Owner ID').fill('sales-1');
  await page.getByLabel('Email').fill(email);
  await page.getByLabel('Phone').fill(phone);
  await page.getByRole('button', { name: 'Save lead' }).click();

  await expect(page.getByRole('alert')).toContainText('Possible duplicate');
  await page.getByRole('button', { name: 'Create anyway' }).click();
  await expect(page.getByRole('button', { name: companyName })).toHaveCount(2);

  await page.getByRole('button', { name: 'New lead' }).click();
  await page.getByLabel('Company name').fill(uniqueCompany);
  await page.getByLabel('Source').fill('Referral');
  await page.getByLabel('Owner ID').fill('sales-1');
  await page.getByRole('button', { name: 'Save lead' }).click();
  await expect(page.getByRole('alert')).toHaveCount(0);
  await expect(page.getByRole('button', { name: uniqueCompany })).toBeVisible();
});
