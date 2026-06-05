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

test('TEST-LEAD-CREATE-002 validates create lead required fields', async ({ page }) => {
  await page.getByRole('button', { name: '线索' }).click();
  await expect(page.getByRole('heading', { name: '线索', exact: true })).toBeVisible();

  await page.getByRole('button', { name: '新建线索' }).click();
  await page.getByLabel('公司名称').fill('E2E Missing Source');
  await page.getByRole('button', { name: '保存线索' }).click();

  await expect(page.getByRole('alert')).toContainText('The lead input is invalid.');
});

test('TEST-LEAD-QUALIFY-003 converts a valid lead through the UI', async ({ page }) => {
  const companyName = `E2E Convert ${Date.now()}`;

  await page.getByRole('button', { name: '线索', exact: true }).click();
  await page.getByRole('button', { name: '新建线索' }).click();
  await page.getByLabel('公司名称').fill(companyName);
  await page.getByLabel('来源').fill('Website');
  await page.getByLabel('负责人 ID').fill('sales-1');
  await page.getByRole('button', { name: '保存线索' }).click();

  await page.getByRole('button', { name: companyName }).click();
  await page.getByRole('button', { name: '标记有效' }).click();
  await expect(page.getByLabel('线索详情').getByText('有效', { exact: true })).toBeVisible();

  await page.getByRole('button', { name: '转换线索' }).click();
  await page.getByLabel('预计金额').fill('99000.00');
  await page.getByLabel('预计关闭日期').fill('2026-12-15');
  await page.getByRole('button', { name: '转换', exact: true }).click();

  await expect(page.getByLabel('线索详情').getByText('已转为商机', { exact: true })).toBeVisible();
  await expect(page.getByLabel('线索详情').getByText(/商机：opp_/)).toBeVisible();
});

test('TEST-LEAD-QUALIFY-004 shows Unassigned qualification as unavailable and backend-denied', async ({ page }) => {
  const companyName = `E2E Unassigned ${Date.now()}`;

  await page.getByRole('button', { name: '线索', exact: true }).click();
  await page.getByRole('button', { name: '新建线索' }).click();
  await page.getByLabel('公司名称').fill(companyName);
  await page.getByLabel('来源').fill('Referral');
  await page.getByRole('button', { name: '保存线索' }).click();

  await page.getByRole('button', { name: companyName }).click();
  await expect(page.getByRole('button', { name: '标记有效' })).toBeDisabled();
  await expect(page.getByText('未分配线索不能确认或转换。')).toBeVisible();
});
