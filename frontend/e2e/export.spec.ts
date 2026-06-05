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

test('TEST-CSV-EXPORT-001 exports active authorized records after confirmation', async ({ page }) => {
  await page.getByRole('button', { name: '导入/导出' }).click();
  await expect(page.getByRole('heading', { name: '导入/导出' })).toBeVisible();
  await page.getByRole('checkbox', { name: '确认导出范围并记录审计日志' }).check();
  await page.getByRole('button', { name: '开始导出' }).click();

  await expect(page.getByText(/已导出 \d+ 行线索/)).toBeVisible();
  await expect(page.locator('.exportResult').getByText('不包含已归档')).toBeVisible();
});
