import { expect, test } from '@playwright/test';

const adminEmail = process.env.E2E_ADMIN_EMAIL ?? 'admin@example.com';
const adminPassword = process.env.E2E_ADMIN_PASSWORD ?? 'AdminChangeMe-001!';
const password = 'UserAdmin-001!';

test.beforeEach(async ({ page }) => {
  await page.goto('/');
  await signIn(page, adminEmail, adminPassword);
  await expect(page.getByRole('heading', { name: 'Work Overview' })).toBeVisible();
});

test('TEST-USER-ADMIN-001 creates user and changes role/status with confirmation', async ({ page }) => {
  const suffix = Date.now();
  const email = `user-admin-${suffix}@example.com`;
  const displayName = `User Admin Evidence ${suffix}`;

  await page.getByRole('button', { name: 'Admin: Users/Roles' }).click();
  await expect(page.getByRole('heading', { name: 'User Management' })).toBeVisible();

  await page.getByLabel('Email').fill(email);
  await page.getByLabel('Display name').fill(displayName);
  await page.getByLabel('Password').fill(password);
  await page.getByLabel('Role').selectOption('Sales');
  await page.getByRole('button', { name: 'Create user' }).click();
  await expect(page.getByRole('row', { name: displayName })).toContainText('Sales');

  await page.getByRole('button', { name: `Edit ${displayName}` }).click();
  await page.getByLabel('New role').selectOption('Sales Manager');
  await page.getByRole('button', { name: 'Review role/status change' }).click();
  await expect(page.getByRole('dialog')).toContainText('Old role: Sales');
  await expect(page.getByRole('dialog')).toContainText('New role: Sales Manager');
  await expect(page.getByRole('dialog')).toContainText('Access impact');
  await expect(page.getByRole('dialog')).toContainText('Operation log');
  await page.getByRole('button', { name: 'Confirm change' }).click();
  await expect(page.getByRole('row', { name: displayName })).toContainText('Sales Manager');

  await page.getByRole('button', { name: `Edit ${displayName}` }).click();
  await page.getByLabel('New status').selectOption('Disabled');
  await page.getByRole('button', { name: 'Review role/status change' }).click();
  await expect(page.getByRole('dialog')).toContainText('New status: Disabled');
  await page.getByRole('button', { name: 'Confirm change' }).click();
  await expect(page.getByRole('row', { name: displayName })).toContainText('Disabled');
});

test('TEST-INV-LASTADMIN-001 blocks disabling or downgrading the last active Administrator', async ({ page }) => {
  await page.getByRole('button', { name: 'Admin: Users/Roles' }).click();
  await page.getByRole('button', { name: 'Edit Seed Administrator' }).click();
  await page.getByLabel('New role').selectOption('Sales');
  await expect(page.getByText('Last active Administrator change is blocked.')).toBeVisible();
  await expect(page.getByRole('button', { name: 'Review role/status change' })).toBeDisabled();
});

test('TEST-PERM-USERADMIN-002/003 sales is denied user administration', async ({ page }) => {
  const suffix = Date.now();
  const email = `user-admin-sales-${suffix}@example.com`;
  await createSalesUser(page, email);

  await page.getByRole('button', { name: 'Sign out' }).click();
  await signIn(page, email, password);
  await expect(page.locator('.topbar').getByText('Sales Denied User Admin')).toBeVisible();
  await expect(page.getByRole('button', { name: 'Admin: Users/Roles' })).toHaveCount(0);

  const denied = await page.evaluate(async () => {
    const response = await fetch('/admin/users', { credentials: 'include' });
    return { status: response.status, body: await response.text() };
  });
  expect(denied.status).toBe(403);
  expect(denied.body).not.toContain('Seed Administrator');
});

async function signIn(page: import('@playwright/test').Page, email: string, userPassword: string) {
  await page.getByLabel('Email').fill(email);
  await page.getByLabel('Password').fill(userPassword);
  await page.getByRole('button', { name: 'Sign in' }).click();
}

async function createSalesUser(page: import('@playwright/test').Page, email: string) {
  await page.evaluate(async ({ email, password }) => {
    const response = await fetch('/admin/users', {
      method: 'POST',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, displayName: 'Sales Denied User Admin', password, role: 'Sales' })
    });
    if (!response.ok) throw new Error(`create sales user failed: ${response.status}`);
  }, { email, password });
}
