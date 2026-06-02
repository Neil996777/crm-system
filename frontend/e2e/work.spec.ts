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

test('TEST-ACTIVITY-NOTE-002 validates missing fields and creates note and activity in record detail', async ({ page }) => {
  const title = `E2E Work Panel ${Date.now()}`;
  await createOpportunity(page, title);

  await expect(page.getByRole('heading', { name: 'Activities, Notes, Tasks' })).toBeVisible();
  await page.getByRole('button', { name: 'Save note', exact: true }).click();
  await expect(page.getByRole('alert')).toContainText('The work item input is invalid.');

  await page.getByLabel('Note content').fill('Decision maker confirmed next step');
  await page.getByRole('button', { name: 'Save note', exact: true }).click();
  await expect(page.getByText('Decision maker confirmed next step')).toBeVisible();

  await page.getByLabel('Activity type').fill('Call');
  await page.getByLabel('Activity content').fill('Introductory call completed');
  await page.getByRole('button', { name: 'Save activity', exact: true }).click();
  await expect(page.getByText('Introductory call completed')).toBeVisible();
});

test('TEST-TASK-LIFECYCLE-002 creates task and completes it from standalone list', async ({ page }) => {
  const title = `E2E Work Task ${Date.now()}`;
  const taskTitle = `Prepare follow up material ${Date.now()}`;
  await createOpportunity(page, title);

  await page.getByLabel('Task title').fill(taskTitle);
  await page.getByLabel('Task due date').fill('2027-03-01');
  await page.getByRole('button', { name: 'Save task', exact: true }).click();
  await expect(page.getByText(taskTitle)).toBeVisible();
  await expect(page.getByText('Open')).toBeVisible();

  await page.getByRole('button', { name: 'Tasks' }).click();
  await expect(page.getByRole('heading', { name: 'Tasks' })).toBeVisible();
  await page.getByRole('button', { name: taskTitle }).click();
  await page.getByRole('button', { name: 'Complete task', exact: true }).click();
  await expect(page.getByRole('heading', { name: taskTitle })).toBeVisible();
  await expect(page.locator('section[aria-label="Task detail"]').getByText('Completed')).toBeVisible();
});

async function createOpportunity(page: import('@playwright/test').Page, title: string) {
  await page.getByRole('button', { name: 'Opportunities' }).click();
  await expect(page.getByRole('heading', { name: 'Opportunities' })).toBeVisible();
  await page.getByRole('button', { name: 'New opportunity', exact: true }).click();
  await page.getByLabel('Title').fill(title);
  await page.getByLabel('Customer ID').fill(`acct_${Date.now()}`);
  await page.getByLabel('Owner ID').fill('sales-1');
  await page.getByLabel('Expected amount').fill('10000.00');
  await page.getByLabel('Expected close date').fill('2027-06-30');
  await page.getByRole('button', { name: 'Save opportunity' }).click();
  await expect(page.getByRole('button', { name: title })).toBeVisible();
  await page.getByRole('button', { name: title }).click();
}
