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
  await expect(page.getByLabel('导出对象类型').locator('.badge')).toHaveCount(1);
  await expect(page.getByLabel('导出对象类型')).toContainText('线索');
  await page.getByRole('checkbox', { name: '确认导出范围并记录审计日志' }).check();
  await page.getByRole('button', { name: '开始导出' }).click();

  await expect(page.getByLabel('导出结果')).toContainText('导出行数');
  await expect(page.getByLabel('导出结果')).toContainText('包含归档');
  await expect(page.getByLabel('导出结果')).toContainText('文件安全');
  await expect(page.getByLabel('导出结果')).toContainText('危险单元格已安全前缀化');
  await expect(page.getByLabel('导出结果')).not.toContainText('dangerous_cells_prefixed');
  await expect(page.getByLabel('审计与清理').last()).toContainText('审计记录状态');
  await expect(page.getByLabel('审计与清理').last()).toContainText('清理状态');
  await expect(page.getByLabel('审计与清理').last()).toContainText('保留至');
});
