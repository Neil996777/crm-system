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

test('TEST-BASIC-REPORT-002 renders basic reports empty state through gateway', async ({ page }) => {
  await page.getByRole('button', { name: '报表' }).click();
  await expect(page.getByRole('heading', { name: '基础销售报表' })).toBeVisible();
  await expect(page.locator('section[aria-label="基础报表指标"]')).toBeVisible();
  await expect(page.getByText('按状态统计线索')).toBeVisible();
  await expect(page.getByText('按状态统计回款')).toBeVisible();
});
