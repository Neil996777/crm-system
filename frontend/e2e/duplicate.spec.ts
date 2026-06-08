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

test('TEST-DUPLICATE-WARN-001/005 account duplicate warning proceeds without merge', async ({ page }) => {
  const companyName = `E2E Duplicate Account ${Date.now()}`;

  await page.getByRole('button', { name: '公司/客户' }).click();
  await page.getByRole('button', { name: '新建客户' }).click();
  await page.getByLabel('公司名称').fill(companyName);
  await page.getByLabel('客户状态').fill('Prospect');
  await page.getByLabel('负责人 ID').fill('sales-1');
  await page.getByRole('button', { name: '保存客户' }).click();
  await expect(page.getByLabel('客户详情').getByRole('heading', { name: companyName })).toBeVisible();
  await page.getByRole('button', { name: '返回客户列表' }).click();

  await page.getByRole('button', { name: '新建客户' }).click();
  await page.getByLabel('公司名称').fill(`  ${companyName.toUpperCase()}  `);
  await page.getByLabel('客户状态').fill('Prospect');
  await page.getByLabel('负责人 ID').fill('sales-1');
  await page.getByRole('button', { name: '保存客户' }).click();

  await expect(page.getByRole('alert')).toContainText('可能重复');
  await page.getByRole('button', { name: '仍然创建' }).click();
  await page.getByRole('button', { name: '返回客户列表' }).click();
  await page.getByLabel('搜索').fill(companyName);
  await page.getByRole('button', { name: '应用筛选' }).click();
  await expect(page.getByRole('table', { name: '客户结果表' }).getByText(new RegExp(companyName, 'i'))).toHaveCount(2);
});

test('TEST-DUPLICATE-WARN-004 lead duplicate warning and unique no-warning path', async ({ page }) => {
  const suffix = Date.now();
  const companyName = `E2E Duplicate Lead ${suffix}`;
  const uniqueCompany = `E2E Unique Lead ${suffix}`;
  const email = `leaddup-${suffix}@example.com`;
  const phone = `139${String(suffix).slice(-8)}`;

  await page.getByRole('button', { name: '线索' }).click();
  await page.getByRole('button', { name: '新建线索' }).click();
  await page.getByLabel('公司名称').fill(companyName);
  await page.getByLabel('来源').fill('Website');
  await page.getByLabel('负责人 ID').fill('sales-1');
  await page.getByLabel('邮箱').fill(email.toUpperCase());
  await page.getByLabel('电话').fill(`+86 ${phone}`);
  await page.getByRole('button', { name: '保存线索' }).click();
  await expect(page.getByLabel('线索详情').getByRole('heading', { name: companyName })).toBeVisible();
  await page.getByRole('button', { name: '返回线索列表' }).click();

  await page.getByRole('button', { name: '新建线索' }).click();
  await page.getByLabel('公司名称').fill(` ${companyName.toLowerCase()} `);
  await page.getByLabel('来源').fill('Website');
  await page.getByLabel('负责人 ID').fill('sales-1');
  await page.getByLabel('邮箱').fill(email);
  await page.getByLabel('电话').fill(phone);
  await page.getByRole('button', { name: '保存线索' }).click();

  await expect(page.getByRole('alert')).toContainText('可能重复');
  await page.getByRole('button', { name: '仍然创建' }).click();
  await page.getByRole('button', { name: '返回线索列表' }).click();
  await page.getByLabel('搜索').fill(companyName);
  await page.getByRole('button', { name: '应用筛选' }).click();
  await expect(page.getByRole('table', { name: '线索结果表' }).getByText(new RegExp(companyName, 'i'))).toHaveCount(2);

  await page.getByRole('button', { name: '新建线索' }).click();
  await page.getByLabel('公司名称').fill(uniqueCompany);
  await page.getByLabel('来源').fill('Referral');
  await page.getByLabel('负责人 ID').fill('sales-1');
  await page.getByRole('button', { name: '保存线索' }).click();
  await expect(page.getByRole('alert')).toHaveCount(0);
  await expect(page.getByLabel('线索详情').getByRole('heading', { name: uniqueCompany })).toBeVisible();
});
