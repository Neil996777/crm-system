import { expect, test } from '@playwright/test';

const adminEmail = process.env.E2E_ADMIN_EMAIL ?? 'admin@example.com';
const adminPassword = process.env.E2E_ADMIN_PASSWORD ?? 'AdminChangeMe-001!';
const timelineRefreshTimeoutMs = 15_000;

test.beforeEach(async ({ page }) => {
  await page.goto('/');
  await page.getByLabel('邮箱').fill(adminEmail);
  await page.getByLabel('密码').fill(adminPassword);
  await page.getByRole('button', { name: '登录' }).click();
  await expect(page.getByRole('heading', { name: '工作台' })).toBeVisible();
});

test('TEST-ACTIVITY-NOTE-002 validates missing fields and creates note and activity in record detail', async ({ page }) => {
  const title = `E2E Work Panel ${Date.now()}`;
  await createOpportunity(page, title);

  await expect(page.getByRole('heading', { name: '动态、备注、任务' })).toBeVisible();
  await page.getByRole('button', { name: '保存备注', exact: true }).click();
  await expect(page.getByRole('alert')).toContainText('工作项输入无效。');

  await page.getByLabel('备注内容').fill('Decision maker confirmed next step');
  await page.getByRole('button', { name: '保存备注', exact: true }).click();
  await expect(page.getByLabel('活动时间线').getByText('Decision maker confirmed next step')).toBeVisible({ timeout: timelineRefreshTimeoutMs });

  await page.getByLabel('动态类型').fill('Call');
  await page.getByLabel('动态内容').fill('Introductory call completed');
  await page.getByRole('button', { name: '保存动态', exact: true }).click();
  await expect(page.getByLabel('活动时间线').locator('.timelineItem')).toHaveCount(2, { timeout: timelineRefreshTimeoutMs });
  await expect(page.getByLabel('活动时间线').getByText('Introductory call completed')).toBeVisible({ timeout: timelineRefreshTimeoutMs });
});

test('TEST-TASK-LIFECYCLE-002 creates task and completes it from standalone list', async ({ page }) => {
  const title = `E2E Work Task ${Date.now()}`;
  const taskTitle = `Prepare follow up material ${Date.now()}`;
  await createOpportunity(page, title);

  await page.getByLabel('任务标题').fill(taskTitle);
  await expect(page.getByLabel('任务到期日')).toHaveClass(/dateControl/);
  await page.getByLabel('任务到期日').fill('2027-03-01');
  await page.getByRole('button', { name: '保存任务', exact: true }).click();
  await expect(page.getByLabel('活动时间线').getByText(taskTitle)).toBeVisible({ timeout: timelineRefreshTimeoutMs });
  await expect(page.getByLabel('活动时间线').getByText('待处理')).toBeVisible({ timeout: timelineRefreshTimeoutMs });

  await page.getByRole('button', { name: '任务', exact: true }).click();
  await expect(page.getByRole('heading', { name: '任务', exact: true })).toBeVisible();
  await expect(page.getByLabel('任务列表')).toBeVisible();
  await page.getByLabel('搜索任务、关联记录或负责人').fill(taskTitle);
  await expect(page.getByRole('button', { name: '应用筛选' })).toBeVisible();
  await page.getByRole('button', { name: '应用筛选' }).click();
  await expect(page.getByText('已筛选')).toBeVisible();
  const table = page.getByRole('table', { name: '任务结果表' });
  await expect(table).toBeVisible();
  await expect(table.getByRole('columnheader', { name: '任务' })).toBeVisible();
  await expect(table.getByRole('columnheader', { name: '关联记录' })).toBeVisible();
  await expect(table.getByRole('columnheader', { name: '状态' })).toBeVisible();
  await expect(table.getByRole('columnheader', { name: '负责人' })).toBeVisible();
  await expect(table.getByRole('columnheader', { name: '到期日' })).toBeVisible();
  await expect(page.getByLabel('选择全部')).toBeVisible();
  await expect(page.locator('.bulkBar')).toContainText('已选择 0 条');
  await expect(page.getByRole('navigation', { name: '分页' })).toBeVisible();
  await page.getByRole('button', { name: `查看任务 ${taskTitle}` }).click();
  await page.getByRole('button', { name: '完成任务', exact: true }).click();
  await expect(page.getByRole('heading', { name: taskTitle })).toBeVisible();
  await expect(page.getByLabel('任务详情').getByText('已完成').first()).toBeVisible();
});

async function createOpportunity(page: import('@playwright/test').Page, title: string) {
  await page.getByRole('button', { name: '商机', exact: true }).click();
  await expect(page.getByRole('heading', { name: '商机', exact: true })).toBeVisible();
  await page.getByRole('button', { name: '新建商机', exact: true }).click();
  await page.getByLabel('标题').fill(title);
  await page.getByLabel('客户 ID').fill(`acct_${Date.now()}`);
  await page.getByLabel('负责人 ID').fill('sales-1');
  await page.getByLabel('预计金额').fill('10000.00');
  await page.getByLabel('预计关闭日期').fill('2027-06-30');
  await page.getByRole('button', { name: '保存商机' }).click();
  await expect(page.getByLabel('商机详情').getByRole('heading', { name: title })).toBeVisible();
}
