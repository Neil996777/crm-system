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

test('TEST-OPP-STAGE-002 shows backend-backed blocked transition alert', async ({ page }) => {
  const title = `E2E Blocked Stage ${Date.now()}`;
  await createOpportunity(page, title);

  await page.getByLabel('商机阶段').getByRole('button', { name: '报价', exact: true }).click();
  await expect(page.getByRole('alert')).toContainText('The requested stage transition is not allowed.');
  await expect(page.getByText('当前阶段：新商机')).toBeVisible();
});

test('TEST-OPP-CLOSE-002 blocks Won until related contract is Signed', async ({ page }) => {
  const title = `E2E Early Won ${Date.now()}`;
  await createOpportunity(page, title);

  await page.getByRole('button', { name: '关闭为赢单' }).click();
  await page.getByLabel('合同 ID').fill('contract_missing_e2e');
  await page.getByLabel('关闭日期').fill('2027-07-01');
  await page.getByRole('button', { name: '确认赢单' }).click();

  await expect(page.getByRole('alert')).toContainText('Won requires a Signed related contract.');
  await expect(page.getByText('当前阶段：新商机')).toBeVisible();
});

test('TEST-OPP-CLOSE-003 closes Lost with reason and terminal detail is read-only', async ({ page }) => {
  const title = `E2E Lost ${Date.now()}`;
  await createOpportunity(page, title);

  await page.getByRole('button', { name: '关闭为丢单' }).click();
  await page.getByLabel('关闭日期').fill('2027-07-02');
  await page.getByLabel('丢单原因').selectOption('PRICE');
  await page.getByLabel('原因详情').fill('Competitor pricing');
  await page.getByRole('button', { name: '确认丢单' }).click();

  await expect(page.getByText('当前阶段：丢单')).toBeVisible();
  await expect(page.getByText('已关闭记录')).toBeVisible();
  await expect(page.getByRole('button', { name: '关闭为丢单' })).toBeDisabled();
});

async function createOpportunity(page: import('@playwright/test').Page, title: string) {
  await page.getByRole('button', { name: '商机', exact: true }).click();
  await expect(page.getByRole('heading', { name: '商机', exact: true })).toBeVisible();
  await page.getByRole('button', { name: '新建商机', exact: true }).click();
  await page.getByLabel('标题').fill(title);
  await page.getByLabel('客户 ID').fill(`acct_${Date.now()}`);
  await page.getByLabel('负责人 ID').fill('sales-1');
  await page.getByLabel('预计金额').fill('10000.00');
  await page.getByLabel('预计关闭日期').fill('2027-06-30');
  await page.getByRole('button', { name: '保存商机' }).click();
  await expect(page.getByLabel('商机详情').getByRole('heading', { name: title })).toBeVisible();
}
