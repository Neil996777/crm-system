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
  await page.getByLabel('邮箱').fill(adminEmail);
  await page.getByLabel('密码').fill(adminPassword);
  await page.getByRole('button', { name: '登录' }).click();
  await expect(page.getByRole('heading', { name: '工作台' })).toBeVisible();
});

test('TEST-PAYMENT-RECORD-002 creates overdue plan and records partial payment with remaining amount', async ({ page }) => {
  const contract = await createContract(page, `payment_record_${Date.now()}`, '10000.00');

  await openPayments(page);
  await page.getByLabel('搜索').fill(contract.opportunityId);
  await page.getByRole('button', { name: '搜索', exact: true }).click();
  await page.getByRole('button', { name: contract.opportunityId }).click();

  await page.getByLabel('计划金额').fill('10000.00');
  await page.getByLabel('计划到期日').fill('2026-01-01');
  await page.getByRole('button', { name: '保存回款计划', exact: true }).click();
  await expect(page.getByText('计划状态：未回款')).toBeVisible();
  await expect(page.getByRole('alert')).toContainText('回款计划已逾期。');

  await page.getByLabel('回款金额').fill('4000.00');
  await page.getByLabel('回款日期').fill('2026-06-02');
  await page.getByLabel('幂等键').fill(`pay_${Date.now()}`);
  await page.getByLabel('回款备注').fill('Partial payment collected');
  await page.getByRole('button', { name: '登记回款', exact: true }).click();

  await expect(page.getByText('回款状态：部分回款')).toBeVisible();
  await expect(page.getByText('剩余金额：6000.00')).toBeVisible();
});

test('TEST-PAYMENT-GUARD-003 blocks zero amount and contract overpayment', async ({ page }) => {
  const contract = await createContract(page, `payment_guard_${Date.now()}`, '10000.00');

  await openPayments(page);
  await page.getByLabel('搜索').fill(contract.opportunityId);
  await page.getByRole('button', { name: '搜索', exact: true }).click();
  await page.getByRole('button', { name: contract.opportunityId }).click();

  await page.getByLabel('计划金额').fill('10000.00');
  await page.getByLabel('计划到期日').fill('2027-08-01');
  await page.getByRole('button', { name: '保存回款计划', exact: true }).click();

  await page.getByLabel('回款金额').fill('0.00');
  await page.getByLabel('回款日期').fill('2027-08-05');
  await page.getByLabel('幂等键').fill(`pay_zero_${Date.now()}`);
  await page.getByRole('button', { name: '登记回款', exact: true }).click();
  await expect(page.getByRole('alert')).toContainText('Payment amount must be greater than zero.');

  await page.getByLabel('回款金额').fill('9000.00');
  await page.getByLabel('幂等键').fill(`pay_ok_${Date.now()}`);
  await page.getByRole('button', { name: '登记回款', exact: true }).click();
  await expect(page.getByText('剩余金额：1000.00')).toBeVisible();

  await page.getByLabel('回款金额').fill('1000.01');
  await page.getByLabel('幂等键').fill(`pay_over_${Date.now()}`);
  await page.getByRole('button', { name: '登记回款', exact: true }).click();
  await expect(page.getByRole('alert')).toContainText('Payment exceeds the remaining contract amount.');
});

async function openPayments(page: import('@playwright/test').Page) {
  await page.getByRole('button', { name: '回款' }).click();
  await expect(page.getByRole('heading', { name: '回款' })).toBeVisible();
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
