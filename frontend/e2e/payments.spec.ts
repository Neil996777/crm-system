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

type Contract = {
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

test('TEST-PAYMENT-RECORD-002 creates overdue plan and records partial payment with remaining amount', async ({ page }) => {
  const contract = await createContract(page, `payment_record_${Date.now()}`, '10000.00');

  await openPayments(page);
  await page.getByLabel('Search').fill(contract.opportunityId);
  await page.getByRole('button', { name: 'Search', exact: true }).click();
  await page.getByRole('button', { name: contract.opportunityId }).click();

  await page.getByLabel('Plan amount').fill('10000.00');
  await page.getByLabel('Plan due date').fill('2026-01-01');
  await page.getByRole('button', { name: 'Save payment plan', exact: true }).click();
  await expect(page.getByText('Plan status: Unpaid')).toBeVisible();
  await expect(page.getByRole('alert')).toContainText('Payment plan is overdue.');

  await page.getByLabel('Payment amount').fill('4000.00');
  await page.getByLabel('Payment date').fill('2026-06-02');
  await page.getByLabel('Idempotency key').fill(`pay_${Date.now()}`);
  await page.getByLabel('Payment note').fill('Partial payment collected');
  await page.getByRole('button', { name: 'Record payment', exact: true }).click();

  await expect(page.getByText('Payment status: PartiallyPaid')).toBeVisible();
  await expect(page.getByText('Remaining amount: 6000.00')).toBeVisible();
});

test('TEST-PAYMENT-GUARD-003 blocks zero amount and contract overpayment', async ({ page }) => {
  const contract = await createContract(page, `payment_guard_${Date.now()}`, '10000.00');

  await openPayments(page);
  await page.getByLabel('Search').fill(contract.opportunityId);
  await page.getByRole('button', { name: 'Search', exact: true }).click();
  await page.getByRole('button', { name: contract.opportunityId }).click();

  await page.getByLabel('Plan amount').fill('10000.00');
  await page.getByLabel('Plan due date').fill('2027-08-01');
  await page.getByRole('button', { name: 'Save payment plan', exact: true }).click();

  await page.getByLabel('Payment amount').fill('0.00');
  await page.getByLabel('Payment date').fill('2027-08-05');
  await page.getByLabel('Idempotency key').fill(`pay_zero_${Date.now()}`);
  await page.getByRole('button', { name: 'Record payment', exact: true }).click();
  await expect(page.getByRole('alert')).toContainText('Payment amount must be greater than zero.');

  await page.getByLabel('Payment amount').fill('9000.00');
  await page.getByLabel('Idempotency key').fill(`pay_ok_${Date.now()}`);
  await page.getByRole('button', { name: 'Record payment', exact: true }).click();
  await expect(page.getByText('Remaining amount: 1000.00')).toBeVisible();

  await page.getByLabel('Payment amount').fill('1000.01');
  await page.getByLabel('Idempotency key').fill(`pay_over_${Date.now()}`);
  await page.getByRole('button', { name: 'Record payment', exact: true }).click();
  await expect(page.getByRole('alert')).toContainText('Payment exceeds the remaining contract amount.');
});

async function openPayments(page: import('@playwright/test').Page) {
  await page.getByRole('button', { name: 'Payments' }).click();
  await expect(page.getByRole('heading', { name: 'Payments' })).toBeVisible();
}

async function createContract(page: import('@playwright/test').Page, key: string, amount: string): Promise<Contract> {
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
    const accepted = await request<Quote>(`/api/quotes/${quote.id}/status`, {
      method: 'POST',
      body: JSON.stringify({ expectedVersion: sent.version, toStatus: 'Accepted' })
    });
    return request<Contract>('/api/contracts', {
      method: 'POST',
      body: JSON.stringify({
        quoteId: accepted.id,
        opportunityId: accepted.opportunityId,
        customerId: accepted.customerId,
        amount,
        status: 'Pending Signature',
        contractNote: 'Payment UI contract note',
        expectedSignedDate: '2027-08-01',
        ownerId: 'sales-1'
      })
    });
  }, { key, amount });
}
