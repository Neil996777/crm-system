import { expect, test } from '@playwright/test';

const adminEmail = process.env.E2E_ADMIN_EMAIL ?? 'admin@example.com';
const adminPassword = process.env.E2E_ADMIN_PASSWORD ?? 'AdminChangeMe-001!';
const password = 'UserAdmin-001!';

test.beforeEach(async ({ page }) => {
  await page.goto('/');
  await signIn(page, adminEmail, adminPassword);
  await expect(page.getByRole('heading', { name: '工作台' })).toBeVisible();
});

test('TEST-USER-ADMIN-001 creates user and changes role/status with confirmation', async ({ page }) => {
  const suffix = Date.now();
  const email = `user-admin-${suffix}@example.com`;
  const displayName = `User Admin Evidence ${suffix}`;

  await page.getByRole('button', { name: '管理：用户与角色' }).click();
  await expect(page.getByRole('heading', { name: '用户与角色' })).toBeVisible();
  await expect(page.getByLabel('末位管理员保护')).toBeVisible();
  await expect(page.getByLabel('分页')).toBeVisible();
  await expect(page.getByText(/角色仅：管理员 \/ 销售经理 \/ 销售/)).toBeVisible();
  const downloadPromise = page.waitForEvent('download');
  await page.getByRole('button', { name: '导出', exact: true }).click();
  const download = await downloadPromise;
  expect(download.suggestedFilename()).toBe('users-filtered.csv');

  await page.getByRole('button', { name: '新建用户' }).click();
  const createForm = page.locator('form.createPanel');
  await createForm.getByLabel('邮箱').fill(email);
  await createForm.getByLabel('显示名称').fill(displayName);
  await createForm.getByLabel('密码').fill(password);
  await createForm.getByLabel('角色').selectOption('Sales');
  await page.getByRole('button', { name: '创建用户' }).click();
  await page.getByPlaceholder('搜索显示名或邮箱').fill(displayName);
  const createdRow = page.getByRole('row', { name: displayName });
  await expect(createdRow).toContainText('销售');
  await expect(createdRow.getByRole('button', { name: `编辑 ${displayName}` })).toBeVisible();
  await expect(createdRow.getByRole('button', { name: `停用 ${displayName}` })).toBeVisible();
  await expect(createdRow.getByRole('button', { name: `改角色 ${displayName}` })).toBeVisible();

  await page.getByRole('button', { name: `编辑 ${displayName}` }).click();
  await page.getByLabel('新角色').selectOption('Sales Manager');
  await page.getByRole('button', { name: '复核角色/状态变更' }).click();
  await expect(page.getByRole('dialog')).toContainText('原角色：销售');
  await expect(page.getByRole('dialog')).toContainText('新角色：销售经理');
  await expect(page.getByRole('dialog')).toContainText('访问影响');
  await expect(page.getByRole('dialog')).toContainText('操作日志');
  await page.getByRole('button', { name: '确认变更' }).click();
  await expect(page.getByRole('row', { name: displayName })).toContainText('销售经理');

  await page.getByRole('button', { name: `编辑 ${displayName}` }).click();
  await page.getByLabel('新状态').selectOption('Disabled');
  await page.getByRole('button', { name: '复核角色/状态变更' }).click();
  await expect(page.getByRole('dialog')).toContainText('新状态：停用');
  await page.getByRole('button', { name: '确认变更' }).click();
  await expect(page.getByRole('row', { name: displayName })).toContainText('停用');
});

test('TEST-INV-LASTADMIN-001 blocks disabling or downgrading the last active Administrator', async ({ page }) => {
  await page.getByRole('button', { name: '管理：用户与角色' }).click();
  await expect(page.getByLabel('末位管理员保护')).toContainText('唯一启用管理员');
  await page.getByPlaceholder('搜索显示名或邮箱').fill('Seed Administrator');
  await expect(page.getByRole('row', { name: 'Seed Administrator' }).getByRole('button', { name: '停用 Seed Administrator' })).toBeDisabled();
  await page.getByRole('button', { name: '编辑 Seed Administrator' }).click();
  await expect(page.getByLabel('新角色').locator('option[value="Sales"]')).toHaveAttribute('disabled', '');
  await expect(page.getByLabel('新状态').locator('option[value="Disabled"]')).toHaveAttribute('disabled', '');
  await expect(page.getByText('末位启用管理员受保护，不能降级或停用。')).toBeVisible();
  await expect(page.getByRole('button', { name: '复核角色/状态变更' })).toBeDisabled();
});

test('TEST-PERM-USERADMIN-002/003 sales is denied user administration', async ({ page }) => {
  const suffix = Date.now();
  const email = `user-admin-sales-${suffix}@example.com`;
  await createSalesUser(page, email);

  await page.getByRole('button', { name: '退出登录' }).click();
  await signIn(page, email, password);
  await expect(page.locator('.topbar').getByText('Sales Denied User Admin')).toBeVisible();
  await expect(page.getByRole('button', { name: '管理：用户与角色' })).toHaveCount(0);

  const denied = await page.evaluate(async () => {
    const response = await fetch('/admin/users', { credentials: 'include' });
    return { status: response.status, body: await response.text() };
  });
  expect(denied.status).toBe(403);
  expect(denied.body).not.toContain('Seed Administrator');
});

async function signIn(page: import('@playwright/test').Page, email: string, userPassword: string) {
  await page.getByLabel('邮箱').fill(email);
  await page.getByLabel('密码').fill(userPassword);
  await page.getByRole('button', { name: '登录' }).click();
}

async function createSalesUser(page: import('@playwright/test').Page, email: string) {
  await page.evaluate(async ({ email, password }) => {
    const response = await fetch('/admin/users', {
      method: 'POST',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, displayName: 'Sales Denied User Admin', password, role: 'Sales' })
    });
    if (!response.ok) throw new Error(`create sales user failed: ${response.status}`);
  }, { email, password });
}
