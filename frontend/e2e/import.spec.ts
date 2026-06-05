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
  const importForm = page.locator('form.importForm');
  await importForm.getByLabel('对象类型').selectOption('lead');
  await importForm.getByLabel('CSV 文件').setInputFiles({
    name: 'leads.csv',
    mimeType: 'text/csv',
    buffer: Buffer.from(csv)
  });
  await page.getByRole('button', { name: '开始导入' }).click();

  await expect(page.getByText('已导入 1 / 2 行')).toBeVisible();
  await expect(page.getByText('第 3 行')).toBeVisible();

  await page.getByRole('button', { name: '线索', exact: true }).click();
  await expect(page.getByText(companyName)).toBeVisible();
});
