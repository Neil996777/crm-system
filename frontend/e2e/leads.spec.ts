import { expect, test } from '@playwright/test';

const adminEmail = process.env.E2E_ADMIN_EMAIL ?? 'admin@example.com';
const adminPassword = process.env.E2E_ADMIN_PASSWORD ?? 'AdminChangeMe-001!';
const generatedUserPassword = 'LeadTransfer-001!';

type LeadRecord = {
  id: string;
  companyName: string;
  ownerId: string;
  status: string;
  version: number;
};

type ManagedUser = {
  id: string;
  email: string;
  displayName: string;
  role: string;
  status: string;
};

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

  await expect(page.getByRole('alert')).toContainText('线索输入无效。');
});

test('TEST-LEAD-QUALIFY-003 converts a valid lead through the UI', async ({ page }) => {
  const companyName = `E2E Convert ${Date.now()}`;

  await page.getByRole('button', { name: '线索', exact: true }).click();
  await page.getByRole('button', { name: '新建线索' }).click();
  await page.getByLabel('公司名称').fill(companyName);
  await page.getByLabel('来源').fill('Website');
  await page.getByLabel('负责人 ID').fill('sales-1');
  await page.getByRole('button', { name: '保存线索' }).click();

  await expect(page.getByLabel('线索详情').getByRole('heading', { name: companyName })).toBeVisible();
  await page.getByRole('button', { name: '标记有效' }).click();
  await expect(page.getByLabel('线索详情').getByText('有效', { exact: true }).first()).toBeVisible();

  await page.getByRole('button', { name: '转换线索' }).click();
  await page.getByLabel('预计金额').fill('99000.00');
  await page.getByLabel('预计关闭日期').fill('2026-12-15');
  await page.getByRole('button', { name: '转换', exact: true }).click();

  await expect(page.getByLabel('线索详情').getByText('已转为商机', { exact: true }).first()).toBeVisible();
  await expect(page.getByLabel('线索详情').getByText(/opp_/).first()).toBeVisible();
});

test('TEST-LEAD-QUALIFY-004 shows Unassigned qualification as unavailable and backend-denied', async ({ page }) => {
  const companyName = `E2E Unassigned ${Date.now()}`;

  await page.getByRole('button', { name: '线索', exact: true }).click();
  await page.getByRole('button', { name: '新建线索' }).click();
  await page.getByLabel('公司名称').fill(companyName);
  await page.getByLabel('来源').fill('Referral');
  await page.getByRole('button', { name: '保存线索' }).click();

  await expect(page.getByLabel('线索详情').getByRole('heading', { name: companyName })).toBeVisible();
  await expect(page.getByRole('button', { name: '标记有效' })).toBeDisabled();
  await expect(page.getByText('未分配线索不能确认或转换。')).toBeVisible();
});

test('TEST-UIUX-FUNC-ROWMENU-001 lead row menu and row click open the record detail', async ({ page }) => {
  const companyName = `E2E Row Menu ${Date.now()}`;

  await page.getByRole('button', { name: '线索', exact: true }).click();
  await page.getByRole('button', { name: '新建线索' }).click();
  await page.getByLabel('公司名称').fill(companyName);
  await page.getByLabel('来源').fill('Website');
  await page.getByLabel('负责人 ID').fill('sales-1');
  await page.getByRole('button', { name: '保存线索' }).click();
  await expect(page.getByLabel('线索详情').getByRole('heading', { name: companyName })).toBeVisible();

  await page.getByRole('button', { name: '返回线索列表' }).click();
  await page.getByPlaceholder('搜索线索名、公司名').fill(companyName);
  await page.getByRole('button', { name: '应用筛选' }).click();
  const row = page.getByRole('row', { name: new RegExp(companyName) });
  await expect(row).toBeVisible();
  await row.getByRole('button', { name: new RegExp(`打开 ${companyName} 的行操作菜单`) }).click();
  await expect(page.getByRole('menu')).toBeVisible();
  await page.getByRole('menuitem', { name: '查看' }).click();
  await expect(page.getByLabel('线索详情').getByRole('heading', { name: companyName })).toBeVisible();

  await page.getByRole('button', { name: '返回线索列表' }).click();
  await row.click();
  await expect(page.getByLabel('线索详情').getByRole('heading', { name: companyName })).toBeVisible();
});

test('TEST-LEAD-TRANSFER-001/002 manager transfer persists owner and history; sales transfer is denied', async ({ page }) => {
  const suffix = Date.now();
  const companyName = `E2E Lead Transfer ${suffix}`;
  const managerEmail = `lead-transfer-manager-${suffix}@example.com`;
  const salesEmail = `lead-transfer-sales-${suffix}@example.com`;
  const manager = await createUser(page, managerEmail, `Lead Transfer Manager ${suffix}`, 'Sales Manager');
  const sales = await createUser(page, salesEmail, `Lead Transfer Sales ${suffix}`, 'Sales');
  const lead = await apiRequest<LeadRecord>(page, '/api/leads', {
    method: 'POST',
    body: JSON.stringify({
      companyName,
      source: 'Website',
      ownerId: 'sales-1'
    })
  });

  await page.getByRole('button', { name: '退出登录' }).click();
  await signIn(page, manager.email, generatedUserPassword);
  await expect(page.locator('.topbar').getByText(manager.displayName)).toBeVisible();
  const transferred = await apiRequest<LeadRecord>(page, `/api/leads/${lead.id}/owner-transfer`, {
    method: 'POST',
    body: JSON.stringify({
      expectedVersion: lead.version,
      newOwnerId: sales.id,
      reason: 'E2E manager transfer coverage'
    })
  });
  expect(transferred.ownerId).toBe(sales.id);
  expect(transferred.status).toBe('Pending Qualification');
  expect(transferred.version).toBeGreaterThan(lead.version);

  const persisted = await apiRequest<LeadRecord>(page, `/api/leads/${lead.id}`);
  expect(persisted.ownerId).toBe(sales.id);
  expect(persisted.version).toBe(transferred.version);

  await expect.poll(async () => page.evaluate(async (id) => {
    const response = await fetch(`/api/leads/${id}/history`, { credentials: 'include' });
    const body = await response.json();
    return (body.data?.events ?? [])
      .map((event: { eventId: string; safeSummary: string; afterSummary?: Record<string, unknown> }) => JSON.stringify(event))
      .join('\n');
  }, lead.id), { timeout: 15_000 }).toContain('EVT-OWNER-CHANGED');

  await openLead(page, companyName);
  await expect(page.getByLabel('线索详情')).toContainText(`负责人 ${sales.id}`);
  const history = page.getByLabel('记录历史');
  await expect(history).toContainText('EVT-OWNER-CHANGED');
  await expect(history).toContainText('新负责人');
  await expect(history).toContainText('原负责人');

  await page.getByRole('button', { name: '退出登录' }).click();
  await signIn(page, sales.email, generatedUserPassword);
  await expect(page.locator('.topbar').getByText(sales.displayName)).toBeVisible();
  const denied = await page.evaluate(async ({ id, version }) => {
    const response = await fetch(`/api/leads/${id}/owner-transfer`, {
      method: 'POST',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        expectedVersion: version,
        newOwnerId: 'sales-1',
        reason: 'Sales attempted transfer'
      })
    });
    const body = await response.json();
    return { status: response.status, code: body.error?.code, safeMessage: body.error?.safeMessage };
  }, { id: lead.id, version: transferred.version });
  expect(denied).toEqual({
    status: 403,
    code: 'PERMISSION_DENIED',
    safeMessage: 'Permission denied.'
  });
});

async function signIn(page: import('@playwright/test').Page, email: string, password: string) {
  await page.getByLabel('邮箱').fill(email);
  await page.getByLabel('密码').fill(password);
  await page.getByRole('button', { name: '登录' }).click();
}

async function createUser(page: import('@playwright/test').Page, email: string, displayName: string, role: 'Sales Manager' | 'Sales') {
  return page.evaluate(async ({ email, displayName, password, role }) => {
    const response = await fetch('/admin/users', {
      method: 'POST',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, displayName, password, role })
    });
    const body = await response.json();
    if (!response.ok) {
      throw new Error(body.error?.safeMessage ?? `create user failed: ${response.status}`);
    }
    return body.user as ManagedUser;
  }, { email, displayName, password: generatedUserPassword, role });
}

async function apiRequest<T>(page: import('@playwright/test').Page, path: string, init?: RequestInit): Promise<T> {
  return page.evaluate(async ({ path, init }) => {
    const response = await fetch(path, {
      ...init,
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
        ...(init?.headers ?? {})
      }
    });
    const body = await response.json();
    if (!response.ok) {
      throw new Error(body.error?.safeMessage ?? JSON.stringify(body));
    }
    return body.data;
  }, { path, init }) as Promise<T>;
}

async function openLead(page: import('@playwright/test').Page, companyName: string) {
  await page.getByRole('button', { name: '线索', exact: true }).click();
  await expect(page.getByRole('heading', { name: '线索', exact: true })).toBeVisible();
  await page.getByPlaceholder('搜索线索名、公司名').fill(companyName);
  await page.getByRole('button', { name: '应用筛选' }).click();
  await page.getByRole('button', { name: `打开线索 ${companyName}` }).click();
  await expect(page.getByLabel('线索详情').getByRole('heading', { name: companyName })).toBeVisible();
}
