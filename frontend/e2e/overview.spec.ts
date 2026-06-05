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

test('TEST-TEAM-OVERVIEW-003 renders manager overview empty state through gateway', async ({ page }) => {
  await page.getByRole('button', { name: '报表' }).click();
  await expect(page.getByRole('heading', { name: '经理团队总览' })).toBeVisible();
  await expect(page.locator('section[aria-label="团队指标"]')).toBeVisible();
  await expect(page.getByText('销售管道状态')).toBeVisible();
});
