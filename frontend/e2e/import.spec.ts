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

test('TEST-CSV-IMPORT-001/002 imports valid CSV rows and shows row errors', async ({ page }) => {
  const suffix = Date.now();
  const companyName = `Import E2E ${suffix}`;
  const csv = `companyName,leadName,source,ownerId\n${companyName},Imported Lead,Website,sales-1\nBroken Lead,, ,sales-1\n`;

  await page.getByRole('button', { name: '导入/导出' }).click();
  await expect(page.getByRole('heading', { name: '导入/导出' })).toBeVisible();
  await expect(page.getByLabel('导入对象类型').locator('.badge')).toHaveCount(1);
  await expect(page.getByLabel('导入对象类型')).toContainText('线索');
  await expect(page.getByRole('button', { name: '最近批次' })).toHaveCount(0);
  await page.getByRole('button', { name: '新建导入' }).click();
  await expect(page.locator('form.importForm').getByLabel('CSV 文件')).toBeFocused();
  const importForm = page.locator('form.importForm');
  await importForm.getByLabel('对象类型').selectOption('lead');
  await expect(importForm.getByRole('button', { name: '开始导入' })).toBeDisabled();
  await importForm.getByLabel('CSV 文件').setInputFiles({
    name: 'leads.csv',
    mimeType: 'text/csv',
    buffer: Buffer.from(csv)
  });
  await expect(importForm.getByRole('button', { name: '开始导入' })).toBeEnabled();
  await importForm.getByRole('button', { name: '开始导入' }).click();

  const importResult = page.getByRole('region', { name: '导入结果' });
  await expect(importResult).toContainText('总行数');
  await expect(importResult).toContainText('成功数');
  await expect(importResult).toContainText('失败数');
  await expect(page.getByLabel('导入结果字段')).toContainText('2');
  await expect(page.getByRole('table', { name: '导入逐行错误表' })).toBeVisible();
  await expect(page.getByText('第 3 行')).toBeVisible();
  await expect(page.getByLabel('审计与清理').first()).toContainText('审计记录状态');
  await expect(page.getByLabel('审计与清理').first()).toContainText('清理状态');
  await expect(page.getByLabel('审计与清理').first()).toContainText('保留至');

  await page.getByLabel('主导航').getByRole('button', { name: '线索', exact: true }).click();
  await expect(page.getByText(companyName)).toBeVisible();
});
