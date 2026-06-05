import { expect, test } from '@playwright/test';

const adminEmail = process.env.E2E_ADMIN_EMAIL ?? 'admin@example.com';
const adminPassword = process.env.E2E_ADMIN_PASSWORD ?? 'AdminChangeMe-001!';

test.beforeEach(async ({ page }) => {
  await page.goto('/');
  await page.getByLabel('邮箱').fill(adminEmail);
  await page.getByLabel('密码').fill(adminPassword);
  await page.getByRole('button', { name: '登录' }).click();
  await expect(page.getByRole('heading', { name: '工作台' })).toBeVisible();
});

test('TEST-QUOTE-LIFECYCLE-002 validates quote required fields', async ({ page }) => {
  await page.getByRole('button', { name: '报价', exact: true }).click();
  await expect(page.getByRole('heading', { name: '报价', exact: true })).toBeVisible();

  await page.getByRole('button', { name: '新建报价' }).click();
  await page.getByLabel('商机 ID').fill(`opp_quote_missing_${Date.now()}`);
  await page.getByLabel('客户 ID').fill(`acct_quote_missing_${Date.now()}`);
  await page.getByLabel('负责人 ID').fill('sales-1');
  await page.getByRole('button', { name: '保存报价' }).click();

  await expect(page.getByRole('alert')).toContainText('The quote input is invalid.');
});

test('TEST-QUOTE-LIFECYCLE-002 shows expired quote warning and blocks contract link', async ({ page }) => {
  const opportunityId = `opp_quote_expire_${Date.now()}`;

  await page.getByRole('button', { name: '报价', exact: true }).click();
  await page.getByRole('button', { name: '新建报价' }).click();
  await page.getByLabel('商机 ID').fill(opportunityId);
  await page.getByLabel('客户 ID').fill(`acct_${opportunityId}`);
  await page.getByLabel('金额').fill('8800.00');
  await page.getByLabel('有效期截止日').fill('2027-10-31');
  await page.getByLabel('负责人 ID').fill('sales-1');
  await page.getByRole('button', { name: '保存报价' }).click();

  await page.getByRole('button', { name: opportunityId }).click();
  await page.getByRole('button', { name: '标记过期', exact: true }).click();

  await expect(page.getByText('状态：已过期')).toBeVisible();
  await expect(page.getByRole('alert')).toContainText('已过期报价不能关联合同。');
  await expect(page.getByText('禁止关联合同')).toBeVisible();
});

test('TEST-QUOTE-ACCEPT-001 creates sends and accepts a quote with contract link indicator', async ({ page }) => {
  const opportunityId = `opp_quote_accept_${Date.now()}`;

  await page.getByRole('button', { name: '报价', exact: true }).click();
  await page.getByRole('button', { name: '新建报价' }).click();
  await page.getByLabel('商机 ID').fill(opportunityId);
  await page.getByLabel('客户 ID').fill(`acct_${opportunityId}`);
  await page.getByLabel('金额').fill('12000.00');
  await page.getByLabel('有效期截止日').fill('2027-09-30');
  await page.getByLabel('负责人 ID').fill('sales-1');
  await page.getByRole('button', { name: '保存报价' }).click();

  await page.getByRole('button', { name: opportunityId }).click();
  await expect(page.getByText('状态：草稿')).toBeVisible();
  await page.getByRole('button', { name: '发送', exact: true }).click();
  await expect(page.getByText('状态：已发送')).toBeVisible();
  await page.getByRole('button', { name: '接受', exact: true }).click();
  await expect(page.getByText('状态：已接受')).toBeVisible();
  await expect(page.getByText('可关联合同')).toBeVisible();
  await expect(page.getByRole('button', { name: '再次为此商机创建报价' })).toHaveCount(0);
});
