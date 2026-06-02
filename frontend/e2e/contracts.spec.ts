import { expect, test } from '@playwright/test';

const adminEmail = process.env.E2E_ADMIN_EMAIL ?? 'admin@example.com';
const adminPassword = process.env.E2E_ADMIN_PASSWORD ?? 'AdminChangeMe-001!';

type Quote = {
  id: string;
  opportunityId: string;
  customerId: string;
  amount: string;
  version: number;
};

test.beforeEach(async ({ page }) => {
  await page.goto('/');
  await page.getByLabel('Email').fill(adminEmail);
  await page.getByLabel('Password').fill(adminPassword);
  await page.getByRole('button', { name: 'Sign in' }).click();
  await expect(page.getByRole('heading', { name: 'Work Overview' })).toBeVisible();
});

test('TEST-CONTRACT-CREATE-002 validates required fields amount difference reason and pending reminder', async ({ page }) => {
  const quote = await createAcceptedQuote(page, `contract_create_${Date.now()}`, '14000.00');

  await openContracts(page);
  await page.getByRole('button', { name: 'New contract', exact: true }).click();
  await page.getByRole('button', { name: 'Save contract', exact: true }).click();
  await expect(page.getByRole('alert')).toContainText('The contract input is invalid.');

  await page.getByLabel('Quote ID').fill(quote.id);
  await page.getByLabel('Opportunity ID').fill(quote.opportunityId);
  await page.getByLabel('Customer ID').fill(quote.customerId);
  await page.getByLabel('Amount', { exact: true }).fill('14400.00');
  await page.getByLabel('Expected signed date').fill('2026-01-01');
  await page.getByLabel('Contract note').fill('Commercial note required for contract creation');
  await page.getByLabel('Owner ID').fill('sales-1');
  await page.getByRole('button', { name: 'Save contract', exact: true }).click();
  await expect(page.getByRole('alert')).toContainText('A reason is required when contract amount differs from quote amount.');

  await page.getByLabel('Amount difference reason').fill('Negotiated implementation services');
  await page.getByRole('button', { name: 'Save contract', exact: true }).click();

  await expect(page.getByText('Status: Pending Signature')).toBeVisible();
  await expect(page.getByRole('alert')).toContainText('Pending signature expected date has passed.');
  await expect(page.getByText('Negotiated implementation services')).toBeVisible();
});

test('TEST-CONTRACT-LIFECYCLE-002 rejects signing without signed effective date then completes lifecycle', async ({ page }) => {
  const quote = await createAcceptedQuote(page, `contract_lifecycle_${Date.now()}`, '19000.00');

  await openContracts(page);
  await page.getByRole('button', { name: 'New contract', exact: true }).click();
  await page.getByLabel('Quote ID').fill(quote.id);
  await page.getByLabel('Opportunity ID').fill(quote.opportunityId);
  await page.getByLabel('Customer ID').fill(quote.customerId);
  await page.getByLabel('Amount', { exact: true }).fill(quote.amount);
  await page.getByLabel('Expected signed date').fill('2027-12-01');
  await page.getByLabel('Contract note').fill('Lifecycle contract note');
  await page.getByLabel('Owner ID').fill('sales-1');
  await page.getByRole('button', { name: 'Save contract', exact: true }).click();

  await page.getByRole('button', { name: 'Sign', exact: true }).click();
  await expect(page.getByRole('alert')).toContainText('Signed or effective date is required for this contract status.');

  await page.getByLabel('Signed/effective date').fill('2027-12-15');
  await page.getByRole('button', { name: 'Sign', exact: true }).click();
  await expect(page.getByText('Status: Signed')).toBeVisible();
  await expect(page.getByText('2027-12-15')).toBeVisible();

  await page.getByRole('button', { name: 'Activate', exact: true }).click();
  await expect(page.getByText('Status: Active')).toBeVisible();

  await page.getByRole('button', { name: 'Complete', exact: true }).click();
  await expect(page.getByText('Status: Completed')).toBeVisible();
});

async function openContracts(page: import('@playwright/test').Page) {
  await page.getByRole('button', { name: 'Contracts' }).click();
  await expect(page.getByRole('heading', { name: 'Contracts' })).toBeVisible();
}

async function createAcceptedQuote(page: import('@playwright/test').Page, key: string, amount: string): Promise<Quote> {
  return page.evaluate(async ({ key, amount }) => {
    async function request<T>(path: string, init?: RequestInit): Promise<T> {
      const response = await fetch(path, {
        ...init,
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
          ...(init?.headers ?? {})
        }
      });
      const body = await response.json();
      if (!response.ok) {
        throw new Error(body.error?.safeMessage ?? 'Request failed.');
      }
      return body.data as T;
    }
    const quote = await request<Quote>('/api/quotes', {
      method: 'POST',
      body: JSON.stringify({
        opportunityId: `opp_${key}`,
        customerId: `acct_${key}`,
        amount,
        status: 'Draft',
        validityEnd: '2027-12-31',
        ownerId: 'sales-1'
      })
    });
    const sent = await request<Quote>(`/api/quotes/${quote.id}/status`, {
      method: 'POST',
      body: JSON.stringify({ expectedVersion: quote.version, toStatus: 'Sent' })
    });
    return request<Quote>(`/api/quotes/${quote.id}/status`, {
      method: 'POST',
      body: JSON.stringify({ expectedVersion: sent.version, toStatus: 'Accepted' })
    });
  }, { key, amount });
}
