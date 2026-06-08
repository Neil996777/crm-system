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
  await page.getByLabel('邮箱').fill(adminEmail);
  await page.getByLabel('密码').fill(adminPassword);
  await page.getByRole('button', { name: '登录' }).click();
  await expect(page.getByRole('heading', { name: '工作台' })).toBeVisible();
});

test('TEST-CONTRACT-CREATE-002 validates required fields amount difference reason and pending reminder', async ({ page }) => {
  const quote = await createAcceptedQuote(page, `contract_create_${Date.now()}`, '14000.00');

  await openContracts(page);
  await page.getByRole('button', { name: '新建合同', exact: true }).click();
  await page.getByRole('button', { name: '保存合同', exact: true }).click();
  await expect(page.getByRole('alert')).toContainText('合同输入无效。');

  await page.getByLabel('报价 ID').fill(quote.id);
  await page.getByLabel('商机 ID').fill(quote.opportunityId);
  await page.getByLabel('客户 ID').fill(quote.customerId);
  await page.getByLabel('金额', { exact: true }).fill('14400.00');
  await page.getByLabel('预计签署日期').fill('2026-01-01');
  await page.getByLabel('合同备注').fill('Commercial note required for contract creation');
  await page.getByLabel('负责人 ID').fill('sales-1');
  await page.getByRole('button', { name: '保存合同', exact: true }).click();
  await expect(page.getByRole('alert')).toContainText('合同金额与报价金额不一致时必须填写原因。');

  await page.getByLabel('金额差异原因').fill('Negotiated implementation services');
  await page.getByRole('button', { name: '保存合同', exact: true }).click();

  await expect(page.getByText('状态：待签署')).toBeVisible();
  await expect(page.getByRole('alert')).toContainText('待签署合同的预计签署日期已过。');
  await expect(page.getByText('Negotiated implementation services')).toBeVisible();
});

test('TEST-CONTRACT-LIFECYCLE-002 rejects signing without signed effective date then completes lifecycle', async ({ page }) => {
  const quote = await createAcceptedQuote(page, `contract_lifecycle_${Date.now()}`, '19000.00');

  await openContracts(page);
  await page.getByRole('button', { name: '新建合同', exact: true }).click();
  await page.getByLabel('报价 ID').fill(quote.id);
  await page.getByLabel('商机 ID').fill(quote.opportunityId);
  await page.getByLabel('客户 ID').fill(quote.customerId);
  await page.getByLabel('金额', { exact: true }).fill(quote.amount);
  await page.getByLabel('预计签署日期').fill('2027-12-01');
  await page.getByLabel('合同备注').fill('Lifecycle contract note');
  await page.getByLabel('负责人 ID').fill('sales-1');
  await page.getByRole('button', { name: '保存合同', exact: true }).click();

  await page.getByRole('button', { name: '签署', exact: true }).click();
  await expect(page.getByRole('alert')).toContainText('该合同状态需要填写签署或生效日期。');

  await page.getByLabel('签署/生效日期').fill('2027-12-15');
  await page.getByRole('button', { name: '签署', exact: true }).click();
  await expect(page.getByText('状态：已签署')).toBeVisible();
  await expect(page.locator('dl.detailGrid dd', { hasText: '2027-12-15' }).first()).toBeVisible();

  await page.getByRole('button', { name: '启用', exact: true }).click();
  await expect(page.getByText('状态：启用')).toBeVisible();

  await page.getByRole('button', { name: '完成', exact: true }).click();
  await expect(page.getByText('状态：已完成')).toBeVisible();
});

async function openContracts(page: import('@playwright/test').Page) {
  await page.getByRole('button', { name: '合同' }).click();
  await expect(page.getByRole('heading', { name: '合同' })).toBeVisible();
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
