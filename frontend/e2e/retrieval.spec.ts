import { expect, test } from '@playwright/test';

const adminEmail = process.env.E2E_ADMIN_EMAIL ?? 'admin@example.com';
const adminPassword = process.env.E2E_ADMIN_PASSWORD ?? 'AdminChangeMe-001!';
const salesPassword = 'Retrieval-001!';

test.beforeEach(async ({ page }) => {
  await page.goto('/');
  await signIn(page, adminEmail, adminPassword);
  await expect(page.getByRole('heading', { name: 'Work Overview' })).toBeVisible();
});

test('TEST-NAV-RETRIEVE-001 lists and details contacts from the primary navigation', async ({ page }) => {
  const suffix = Date.now();
  const companyName = `Retrieval Contacts ${suffix}`;
  const contactName = `Retrieval Person ${suffix}`;

  await page.getByRole('button', { name: 'Companies/Customers' }).click();
  await page.getByRole('button', { name: 'New customer' }).click();
  await page.getByLabel('Company name').fill(companyName);
  await page.getByLabel('Customer status').fill('Prospect');
  await page.getByLabel('Owner ID').fill('sales-1');
  await page.getByRole('button', { name: 'Save customer' }).click();
  await page.getByRole('button', { name: companyName }).click();
  await page.getByRole('button', { name: 'Add contact' }).click();
  await page.getByLabel('Contact name').fill(contactName);
  await page.getByLabel('Email').fill(`retrieval-${suffix}@example.com`);
  await page.getByLabel('Phone').fill(`138${String(suffix).slice(-8)}`);
  await page.getByRole('button', { name: 'Save contact' }).click();
  await expect(page.getByText(contactName)).toBeVisible();

  await page.getByRole('navigation', { name: 'Primary' }).getByRole('button', { name: 'Contacts', exact: true }).click();
  await expect(page.getByRole('heading', { name: 'Contacts' })).toBeVisible();
  await page.getByLabel('Search').fill(contactName);
  await page.getByRole('button', { name: 'Search' }).click();
  await page.getByRole('button', { name: contactName }).click();
  await expect(page.getByLabel('Contact detail')).toContainText(contactName);
  await expect(page.getByLabel('Contact detail')).toContainText(companyName);
});

test('TEST-NAV-RETRIEVE-003/004 shows empty state and invalid filter feedback', async ({ page }) => {
  await page.getByRole('button', { name: 'Leads' }).click();
  await page.getByLabel('Search').fill(`missing-${Date.now()}`);
  await page.getByRole('button', { name: 'Search' }).click();
  await expect(page.getByText('No leads found.')).toBeVisible();

  const invalid = await page.evaluate(async () => {
    const response = await fetch('/api/opportunities?stage=NotAStage', { credentials: 'include' });
    return { status: response.status, body: await response.text() };
  });
  expect(invalid.status).toBe(400);
  expect(invalid.body).toContain('INVALID_FILTER');
});

test('TEST-NAV-RETRIEVE-005 hides unauthorized records from sales lists', async ({ page }) => {
  const suffix = Date.now();
  const companyName = `Restricted Retrieval ${suffix}`;
  const salesEmail = `retrieval-sales-${suffix}@example.com`;
  await createSalesUser(page, salesEmail);

  await page.getByRole('button', { name: 'Leads' }).click();
  await page.getByRole('button', { name: 'New lead' }).click();
  await page.getByLabel('Company name').fill(companyName);
  await page.getByLabel('Source').fill('Website');
  await page.getByLabel('Owner ID').fill('sales-1');
  await page.getByRole('button', { name: 'Save lead' }).click();
  await expect(page.getByRole('button', { name: companyName })).toBeVisible();

  await page.getByRole('button', { name: 'Sign out' }).click();
  await signIn(page, salesEmail, salesPassword);
  await page.getByRole('button', { name: 'Leads' }).click();
  await page.getByLabel('Search').fill(companyName);
  await page.getByRole('button', { name: 'Search' }).click();
  await expect(page.getByText('No leads found.')).toBeVisible();
  await expect(page.getByText(companyName)).toHaveCount(0);
});

async function signIn(page: import('@playwright/test').Page, email: string, password: string) {
  await page.getByLabel('Email').fill(email);
  await page.getByLabel('Password').fill(password);
  await page.getByRole('button', { name: 'Sign in' }).click();
}

async function createSalesUser(page: import('@playwright/test').Page, email: string) {
  await page.evaluate(async ({ email, password }) => {
    const response = await fetch('/admin/users', {
      method: 'POST',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, displayName: 'Retrieval Sales', password, role: 'Sales' })
    });
    if (!response.ok) throw new Error(`create sales user failed: ${response.status}`);
  }, { email, password: salesPassword });
}
