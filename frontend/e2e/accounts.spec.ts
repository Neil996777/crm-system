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

test('TEST-CUSTOMER-CRUD-002 validates required customer fields', async ({ page }) => {
  await page.getByRole('button', { name: '公司/客户' }).click();
  await expect(page.getByRole('heading', { name: '公司/客户' })).toBeVisible();

  await page.getByRole('button', { name: '新建客户' }).click();
  await page.getByLabel('公司名称').fill('E2E Missing Status');
  await page.getByRole('button', { name: '保存客户' }).click();

  await expect(page.getByRole('alert')).toContainText('客户输入无效。');
});

test('TEST-CONTACT-LINK-003 creates two contacts visible in customer context', async ({ page }) => {
  const suffix = Date.now();
  const companyName = `E2E Account ${suffix}`;

  await page.getByRole('button', { name: '公司/客户' }).click();
  await page.getByRole('button', { name: '新建客户' }).click();
  await page.getByLabel('公司名称').fill(companyName);
  await page.getByLabel('客户状态').fill('Prospect');
  await page.getByLabel('负责人 ID').fill('sales-1');
  await page.getByRole('button', { name: '保存客户' }).click();

  await page.getByRole('button', { name: companyName }).click();
  await page.getByRole('button', { name: '添加联系人', exact: true }).click();
  await expect(page.getByLabel('联系人姓名')).toBeVisible();
  await page.getByLabel('联系人姓名').fill('Primary Buyer');
  await page.getByLabel('邮箱').fill(`buyer-${suffix}@example.com`);
  await page.getByRole('button', { name: '保存联系人' }).click();
  await expect(page.getByRole('table', { name: '联系人' }).getByText('Primary Buyer')).toBeVisible();
  await expect(page.getByLabel('联系人姓名')).toHaveCount(0);

  await page.getByRole('button', { name: '添加联系人', exact: true }).click();
  await expect(page.getByLabel('联系人姓名')).toBeVisible();
  await page.getByLabel('联系人姓名').fill('Technical Reviewer');
  await page.getByLabel('邮箱').fill(`technical-${suffix}@example.com`);
  await page.getByLabel('角色备注').fill('Technical review');
  await page.getByRole('button', { name: '保存联系人' }).click();

  await expect(page.getByRole('table', { name: '联系人' }).getByText('Primary Buyer')).toBeVisible();
  await expect(page.getByRole('table', { name: '联系人' }).getByText('Technical Reviewer')).toBeVisible();
});
