import { expect, test } from '@playwright/test';

const adminEmail = process.env.E2E_ADMIN_EMAIL ?? 'admin@example.com';
const adminPassword = process.env.E2E_ADMIN_PASSWORD ?? 'AdminChangeMe-001!';

test('TEST-AUTH-LOGIN-001/005 signs in through gateway and persists session', async ({ page }) => {
  await page.goto('/');
  await expect(page.getByRole('heading', { name: 'CRM 系统' })).toBeVisible();

  await page.getByLabel('邮箱').fill(adminEmail);
  await page.getByLabel('密码').fill(adminPassword);
  await page.getByRole('button', { name: '登录' }).click();

  await expect(page.getByRole('heading', { name: '工作台' })).toBeVisible();
  await expect(page.getByText('管理员', { exact: true })).toBeVisible();
  await expect(page.getByRole('navigation').getByText('管理：用户与角色')).toBeVisible();
  await expect(page.getByRole('navigation').getByText('操作日志')).toBeVisible();

  await page.reload();
  await expect(page.getByRole('heading', { name: '工作台' })).toBeVisible();
  await expect(page.getByText('Seed Administrator')).toBeVisible();
});

test('TEST-AUTH-LOGIN-002 shows one generic sign-in failure', async ({ page }) => {
  await page.goto('/');
  await page.getByLabel('邮箱').fill('missing@example.com');
  await page.getByLabel('密码').fill('wrong-password');
  await page.getByRole('button', { name: '登录' }).click();

  await expect(page.getByRole('alert')).toHaveText('认证失败。');
});
