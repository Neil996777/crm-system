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

test('TEST-CUSTOMER-CRUD-002 validates required customer fields', async ({ page }) => {
  await page.getByRole('button', { name: 'Companies/Customers' }).click();
  await expect(page.getByRole('heading', { name: 'Companies/Customers' })).toBeVisible();

  await page.getByRole('button', { name: 'New customer' }).click();
  await page.getByLabel('Company name').fill('E2E Missing Status');
  await page.getByRole('button', { name: 'Save customer' }).click();

  await expect(page.getByRole('alert')).toContainText('The account input is invalid.');
});

test('TEST-CONTACT-LINK-003 creates two contacts visible in customer context', async ({ page }) => {
  const suffix = Date.now();
  const companyName = `E2E Account ${suffix}`;

  await page.getByRole('button', { name: 'Companies/Customers' }).click();
  await page.getByRole('button', { name: 'New customer' }).click();
  await page.getByLabel('Company name').fill(companyName);
  await page.getByLabel('Customer status').fill('Prospect');
  await page.getByLabel('Owner ID').fill('sales-1');
  await page.getByRole('button', { name: 'Save customer' }).click();

  await page.getByRole('button', { name: companyName }).click();
  await page.getByRole('button', { name: 'Add contact', exact: true }).click();
  await expect(page.getByLabel('Contact name')).toBeVisible();
  await page.getByLabel('Contact name').fill('Primary Buyer');
  await page.getByLabel('Email').fill(`buyer-${suffix}@example.com`);
  await page.getByRole('button', { name: 'Save contact' }).click();
  await expect(page.getByRole('table', { name: 'Contacts' }).getByText('Primary Buyer')).toBeVisible();
  await expect(page.getByLabel('Contact name')).toHaveCount(0);

  await page.getByRole('button', { name: 'Add contact', exact: true }).click();
  await expect(page.getByLabel('Contact name')).toBeVisible();
  await page.getByLabel('Contact name').fill('Technical Reviewer');
  await page.getByLabel('Email').fill(`technical-${suffix}@example.com`);
  await page.getByLabel('Role note').fill('Technical review');
  await page.getByRole('button', { name: 'Save contact' }).click();

  await expect(page.getByRole('table', { name: 'Contacts' }).getByText('Primary Buyer')).toBeVisible();
  await expect(page.getByRole('table', { name: 'Contacts' }).getByText('Technical Reviewer')).toBeVisible();
});
