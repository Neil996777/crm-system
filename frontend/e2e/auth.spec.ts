import { expect, test } from '@playwright/test';

const adminEmail = process.env.E2E_ADMIN_EMAIL ?? 'admin@example.com';
const adminPassword = process.env.E2E_ADMIN_PASSWORD ?? 'AdminChangeMe-001!';

test('TEST-AUTH-LOGIN-001/005 signs in through gateway and persists session', async ({ page }) => {
  await page.goto('/');
  await expect(page.getByRole('heading', { name: 'CRM System' })).toBeVisible();

  await page.getByLabel('Email').fill(adminEmail);
  await page.getByLabel('Password').fill(adminPassword);
  await page.getByRole('button', { name: 'Sign in' }).click();

  await expect(page.getByRole('heading', { name: 'Work Overview' })).toBeVisible();
  await expect(page.getByText('Administrator', { exact: true })).toBeVisible();
  await expect(page.getByRole('navigation').getByText('Admin: Users/Roles')).toBeVisible();
  await expect(page.getByRole('navigation').getByText('Operation Logs')).toBeVisible();

  await page.reload();
  await expect(page.getByRole('heading', { name: 'Work Overview' })).toBeVisible();
  await expect(page.getByText('Seed Administrator')).toBeVisible();
});

test('TEST-AUTH-LOGIN-002 shows one generic sign-in failure', async ({ page }) => {
  await page.goto('/');
  await page.getByLabel('Email').fill('missing@example.com');
  await page.getByLabel('Password').fill('wrong-password');
  await page.getByRole('button', { name: 'Sign in' }).click();

  await expect(page.getByRole('alert')).toHaveText('Authentication failed.');
});
