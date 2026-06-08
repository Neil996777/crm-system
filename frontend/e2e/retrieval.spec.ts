import { expect, test } from '@playwright/test';

const adminEmail = process.env.E2E_ADMIN_EMAIL ?? 'admin@example.com';
const adminPassword = process.env.E2E_ADMIN_PASSWORD ?? 'AdminChangeMe-001!';
const salesPassword = 'Retrieval-001!';

test.beforeEach(async ({ page }) => {
  await page.goto('/');
  await signIn(page, adminEmail, adminPassword);
  await expect(page.getByRole('heading', { name: '工作台' })).toBeVisible();
});

test('TEST-NAV-RETRIEVE-001 lists and details contacts from the primary navigation', async ({ page }) => {
  const suffix = Date.now();
  const companyName = `Retrieval Contacts ${suffix}`;
  const contactName = `Retrieval Person ${suffix}`;

  await page.getByRole('button', { name: '公司/客户' }).click();
  await page.getByRole('button', { name: '新建客户' }).click();
  await page.getByLabel('公司名称').fill(companyName);
  await page.getByLabel('客户状态').fill('Prospect');
  await page.getByLabel('负责人 ID').fill('sales-1');
  await page.getByRole('button', { name: '保存客户' }).click();
  await expect(page.getByLabel('客户详情').getByRole('heading', { name: companyName })).toBeVisible();
  await page.getByRole('button', { name: '添加联系人' }).click();
  await page.getByLabel('联系人姓名').fill(contactName);
  await page.getByLabel('邮箱').fill(`retrieval-${suffix}@example.com`);
  await page.getByLabel('电话').fill(`138${String(suffix).slice(-8)}`);
  await page.getByRole('button', { name: '保存联系人' }).click();
  await expect(page.getByText(contactName)).toBeVisible();

  await page.getByRole('navigation', { name: '主导航' }).getByRole('button', { name: '联系人', exact: true }).click();
  await expect(page.getByRole('heading', { name: '联系人' })).toBeVisible();
  await page.getByLabel('搜索').fill(contactName);
  await page.getByRole('button', { name: '应用筛选' }).click();
  await page.getByRole('button', { name: new RegExp(`查看 ${escapeRegExp(contactName)}`) }).click();
  await expect(page.getByLabel('联系人详情')).toContainText(contactName);
  await expect(page.getByLabel('联系人详情')).toContainText(companyName);
});

test('TEST-NAV-RETRIEVE-003/004 shows empty state and invalid filter feedback', async ({ page }) => {
  await page.getByRole('button', { name: '线索', exact: true }).click();
  await page.getByLabel('搜索').fill(`missing-${Date.now()}`);
  await page.getByRole('button', { name: '应用筛选' }).click();
  await expect(page.getByText('没有符合当前筛选条件的线索。')).toBeVisible();

  const invalid = await page.evaluate(async () => {
    const response = await fetch('/api/opportunities?stage=NotAStage', { credentials: 'include' });
    return { status: response.status, body: await response.text() };
  });
  expect(invalid.status).toBe(400);
  expect(invalid.body).toContain('INVALID_FILTER');
});

test('TEST-NAV-RETRIEVE-005 hides unauthorized records from sales lists', async ({ page }) => {
  const suffix = Date.now();
  const companyName = `Restricted Retrieval ${suffix}`;
  const salesEmail = `retrieval-sales-${suffix}@example.com`;
  await createSalesUser(page, salesEmail);

  await page.getByRole('button', { name: '线索', exact: true }).click();
  await page.getByRole('button', { name: '新建线索' }).click();
  await page.getByLabel('公司名称').fill(companyName);
  await page.getByLabel('来源').fill('Website');
  await page.getByLabel('负责人 ID').fill('sales-1');
  await page.getByRole('button', { name: '保存线索' }).click();
  await expect(page.getByLabel('线索详情').getByRole('heading', { name: companyName })).toBeVisible();

  await page.getByRole('button', { name: '退出登录' }).click();
  await signIn(page, salesEmail, salesPassword);
  await page.getByRole('button', { name: '线索', exact: true }).click();
  await page.getByLabel('搜索').fill(companyName);
  await page.getByRole('button', { name: '应用筛选' }).click();
  await expect(page.getByText('没有符合当前筛选条件的线索。')).toBeVisible();
  await expect(page.getByText(companyName)).toHaveCount(0);
});

async function signIn(page: import('@playwright/test').Page, email: string, password: string) {
  await page.getByLabel('邮箱').fill(email);
  await page.getByLabel('密码').fill(password);
  await page.getByRole('button', { name: '登录' }).click();
}

async function createSalesUser(page: import('@playwright/test').Page, email: string) {
  await page.evaluate(async ({ email, password }) => {
    const response = await fetch('/admin/users', {
      method: 'POST',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, displayName: 'Retrieval Sales', password, role: 'Sales' })
    });
    if (!response.ok) throw new Error(`create sales user failed: ${response.status}`);
  }, { email, password: salesPassword });
}

function escapeRegExp(value: string) {
  return value.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
}
