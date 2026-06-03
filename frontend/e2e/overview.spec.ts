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

test('TEST-TEAM-OVERVIEW-003 renders manager overview empty state through gateway', async ({ page }) => {
  await page.getByRole('button', { name: 'Reports' }).click();
  await expect(page.getByRole('heading', { name: 'Manager Team Overview' })).toBeVisible();
  await expect(page.locator('section[aria-label="Team metrics"]')).toBeVisible();
  await expect(page.getByText('Pipeline Status')).toBeVisible();
});
