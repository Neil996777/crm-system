import { expect, test } from '@playwright/test';

const adminEmail = process.env.E2E_ADMIN_EMAIL ?? 'admin@example.com';
const adminPassword = process.env.E2E_ADMIN_PASSWORD ?? 'AdminChangeMe-001!';
const salesPassword = 'OppSales-001!';

test.beforeEach(async ({ page }) => {
  await page.goto('/');
  await page.getByLabel('邮箱').fill(adminEmail);
  await page.getByLabel('密码').fill(adminPassword);
  await page.getByRole('button', { name: '登录' }).click();
  await expect(page.getByRole('heading', { name: '工作台' })).toBeVisible();
});

test('TEST-OPP-STAGE-002 shows backend-backed blocked transition alert', async ({ page }) => {
  const title = `E2E Blocked Stage ${Date.now()}`;
  await createOpportunity(page, title);

  await page.getByLabel('商机阶段').getByRole('button', { name: '报价', exact: true }).click();
  await expect(page.getByRole('alert')).toContainText('不允许该阶段流转。');
  await expect(page.getByText('当前阶段：新商机')).toBeVisible();
});

test('TEST-OPP-CLOSE-002 blocks Won until related contract is Signed', async ({ page }) => {
  const title = `E2E Early Won ${Date.now()}`;
  await createOpportunity(page, title);

  await page.getByRole('button', { name: '关闭为赢单' }).click();
  await expect(page.getByRole('button', { name: '确认赢单' })).toBeDisabled();
  await page.getByLabel('合同 ID').fill('contract_missing_e2e');
  await page.getByLabel('关闭日期').fill('2027-07-01');
  await expect(page.getByRole('button', { name: '确认赢单' })).toBeEnabled();
  await page.getByRole('button', { name: '确认赢单' }).click();

  await expect(page.getByRole('alert')).toContainText('赢单需要有已签署的关联合同。');
  await expect(page.getByText('当前阶段：新商机')).toBeVisible();
});

test('TEST-OPP-CLOSE-003 closes Lost with reason and terminal detail is read-only', async ({ page }) => {
  const title = `E2E Lost ${Date.now()}`;
  await createOpportunity(page, title);

  await page.getByRole('button', { name: '关闭为丢单' }).click();
  await expect(page.getByRole('button', { name: '确认丢单' })).toBeDisabled();
  await page.getByLabel('关闭日期').fill('2027-07-02');
  await page.getByLabel('丢单原因').selectOption('PRICE');
  await expect(page.getByRole('button', { name: '确认丢单' })).toBeDisabled();
  await page.getByLabel('原因详情').fill('Competitor pricing');
  await expect(page.getByRole('button', { name: '确认丢单' })).toBeEnabled();
  await page.getByRole('button', { name: '确认丢单' }).click();

  await expect(page.getByText('当前阶段：丢单')).toBeVisible();
  await expect(page.getByText('已关闭记录')).toBeVisible();
  await expect(page.getByText('终态商机只读，阶段、关闭操作和工作记录新增均已停用。')).toBeVisible();
  await expect(page.getByRole('button', { name: '编辑' })).toHaveCount(0);
  await expect(page.getByRole('button', { name: '转移负责人' })).toHaveCount(0);
  await expect(page.getByRole('button', { name: '归档', exact: true })).toHaveCount(0);
  await expect(page.getByRole('button', { name: '关闭为丢单' })).toBeDisabled();
  await expect(page.getByRole('button', { name: '保存备注' })).toHaveCount(0);
  await expect(page.getByRole('button', { name: '保存动态' })).toHaveCount(0);
  await expect(page.getByRole('button', { name: '保存任务' })).toHaveCount(0);
});

test('TEST-UIUX-P1-001 manager detail edit transfer and archive actions are live', async ({ page }) => {
  const title = `E2E Detail Actions ${Date.now()}`;
  await createOpportunity(page, title);

  await page.getByRole('button', { name: '编辑' }).click();
  await expect(page.getByRole('heading', { name: '编辑商机' })).toBeVisible();
  await page.getByLabel('标题').fill(`${title} Edited`);
  await page.getByRole('button', { name: '保存编辑' }).click();
  await expect(page.getByLabel('商机详情').getByRole('heading', { name: `${title} Edited` })).toBeVisible();

  await page.getByRole('button', { name: '转移负责人' }).click();
  await page.getByLabel('新负责人 ID').fill('sales-2');
  await page.getByRole('button', { name: '确认转移负责人' }).click();
  await expect(page.getByText('负责人 sales-2')).toBeVisible();

  await page.getByRole('button', { name: '归档', exact: true }).click();
  await expect(page.getByText('已归档')).toBeVisible();
  await expect(page.getByRole('button', { name: '归档', exact: true })).toBeDisabled();
  await restoreOpportunityProjection(page, `${title} Edited`);
});

test('TEST-UIUX-P1-002 create form selects only the four non-terminal stages', async ({ page }) => {
  const title = `E2E Stage Choice ${Date.now()}`;
  await page.getByRole('button', { name: '商机', exact: true }).click();
  await page.getByRole('button', { name: '新建商机', exact: true }).click();

  const stageGroup = page.getByRole('group', { name: '阶段' });
  await expect(stageGroup.getByRole('button', { name: '新商机' })).toHaveAttribute('aria-pressed', 'true');
  await expect(stageGroup.getByRole('button', { name: '需求已确认' })).toBeVisible();
  await expect(stageGroup.getByRole('button', { name: '报价' })).toBeVisible();
  await expect(stageGroup.getByRole('button', { name: '合同谈判' })).toBeVisible();
  await expect(stageGroup.getByRole('button', { name: '赢单' })).toHaveCount(0);
  await expect(stageGroup.getByRole('button', { name: '丢单' })).toHaveCount(0);
  await expect(page.getByLabel('预计关闭日期')).toHaveClass(/dateControl/);

  await stageGroup.getByRole('button', { name: '报价' }).click();
  await page.getByLabel('标题').fill(title);
  await page.getByLabel('客户 ID').fill(`acct_${Date.now()}`);
  await page.getByLabel('负责人 ID').fill('sales-1');
  await page.getByLabel('预计金额').fill('12000.00');
  await page.getByLabel('预计关闭日期').fill('2027-08-01');
  await page.getByRole('button', { name: '保存商机' }).click();
  await expect(page.getByText('当前阶段：报价')).toBeVisible();
});

test('TEST-UIUX-A4-OPP-001 sales hides bulk actions and locks owner field to self', async ({ page }) => {
  const suffix = Date.now();
  const sales = await createUser(page, `opp-sales-${suffix}@example.com`, `Opp Sales ${suffix}`, 'Sales');

  await page.getByRole('button', { name: '退出登录' }).click();
  await signIn(page, sales.email, salesPassword);
  await expect(page.locator('.topbar').getByText(sales.displayName)).toBeVisible();

  await page.getByRole('button', { name: '商机', exact: true }).click();
  await expect(page.getByText('批量转移负责人')).toHaveCount(0);
  await expect(page.getByText('批量归档')).toHaveCount(0);
  await page.getByRole('button', { name: '新建商机', exact: true }).click();
  await expect(page.getByLabel('负责人 ID')).toBeDisabled();
  await expect(page.getByLabel('负责人 ID')).toHaveValue(sales.id);
  await page.getByLabel('标题').fill(`Sales Restricted Detail ${suffix}`);
  await page.getByLabel('客户 ID').fill(`acct_sales_${suffix}`);
  await page.getByLabel('预计金额').fill('9000.00');
  await page.getByLabel('预计关闭日期').fill('2027-09-01');
  await page.getByRole('button', { name: '保存商机' }).click();
  await expect(page.getByLabel('商机详情').getByRole('heading', { name: `Sales Restricted Detail ${suffix}` })).toBeVisible();
  await expect(page.getByRole('button', { name: '编辑' })).toHaveCount(0);
  await expect(page.getByRole('button', { name: '转移负责人' })).toHaveCount(0);
  await expect(page.getByRole('button', { name: '归档', exact: true })).toHaveCount(0);
});

async function createOpportunity(page: import('@playwright/test').Page, title: string) {
  await page.getByRole('button', { name: '商机', exact: true }).click();
  await expect(page.getByRole('heading', { name: '商机', exact: true })).toBeVisible();
  await page.getByRole('button', { name: '新建商机', exact: true }).click();
  await page.getByLabel('标题').fill(title);
  await page.getByLabel('客户 ID').fill(`acct_${Date.now()}`);
  await page.getByLabel('负责人 ID').fill('sales-1');
  await page.getByLabel('预计金额').fill('10000.00');
  await expect(page.getByLabel('预计关闭日期')).toHaveClass(/dateControl/);
  await page.getByLabel('预计关闭日期').fill('2027-06-30');
  await page.getByRole('button', { name: '保存商机' }).click();
  await expect(page.getByLabel('商机详情').getByRole('heading', { name: title })).toBeVisible();
  await expect(page.getByLabel('活动时间线')).toBeVisible();
}

async function restoreOpportunityProjection(page: import('@playwright/test').Page, title: string) {
  await page.evaluate(async ({ title }) => {
    const query = new URLSearchParams({ search: title, includeArchived: 'true' });
    const listResponse = await fetch(`/api/opportunities?${query.toString()}`, { credentials: 'include' });
    const listBody = await listResponse.json();
    if (!listResponse.ok) throw new Error(JSON.stringify(listBody));
    const opportunity = listBody.data.items.find((item: { title: string }) => item.title === title);
    if (!opportunity) throw new Error(`opportunity not found: ${title}`);
    const updateResponse = await fetch(`/api/opportunities/${opportunity.id}`, {
      method: 'PATCH',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        expectedVersion: opportunity.version,
        customerId: opportunity.customerId,
        ownerId: opportunity.ownerId,
        stage: opportunity.stage,
        expectedAmount: opportunity.expectedAmount,
        expectedCloseDate: opportunity.expectedCloseDate,
        title: opportunity.title
      })
    });
    const updateBody = await updateResponse.json();
    if (!updateResponse.ok) throw new Error(JSON.stringify(updateBody));
  }, { title });
}

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
    if (!response.ok) throw new Error(JSON.stringify(body));
    return body.user as { id: string; email: string; displayName: string; role: string; status: string };
  }, { email, displayName, password: salesPassword, role });
}
