import { expect, test, type Page } from '@playwright/test';

const adminEmail = process.env.E2E_ADMIN_EMAIL ?? 'admin@example.com';
const adminPassword = process.env.E2E_ADMIN_PASSWORD ?? 'AdminChangeMe-001!';
const timelineRefreshTimeoutMs = 15_000;

type E2EWorkTask = {
  id: string;
  title: string;
  relatedType: string;
  relatedId: string;
  ownerId: string;
  dueDate: string;
  status: string;
  version: number;
};

type WorkContextFixture = {
  relatedType: 'Lead' | 'Contract' | 'Payment';
  relatedId: string;
  label: string;
  taskTitle: string;
  recordSearch: string;
  recordHeading: string;
  actualPaymentId?: string;
};

type Quote = {
  id: string;
  opportunityId: string;
  customerId: string;
  amount: string;
  version: number;
};

type Contract = {
  id: string;
  opportunityId: string;
  customerId: string;
  amount: string;
  version: number;
};

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
  const createdTask = await waitForTaskCreate(page, () => page.getByRole('button', { name: '保存任务', exact: true }).click());
  expect(createdTask.title).toBe(taskTitle);
  await expect(page.getByLabel('活动时间线').getByText(taskTitle)).toBeVisible({ timeout: timelineRefreshTimeoutMs });
  await expect(page.getByLabel('活动时间线').getByText('待处理')).toBeVisible({ timeout: timelineRefreshTimeoutMs });

  await waitForTaskListContaining(page, createdTask.id, taskTitle, () => page.getByRole('button', { name: '任务', exact: true }).click());
  await expect(page.getByRole('heading', { name: '任务', exact: true })).toBeVisible();
  await expect(page.getByLabel('任务列表')).toBeVisible();
  await page.getByLabel('搜索任务、关联记录或负责人').fill(taskTitle);
  await expect(page.getByRole('button', { name: '应用筛选' })).toBeVisible();
  await waitForTaskListContaining(page, createdTask.id, taskTitle, () => page.getByRole('button', { name: '应用筛选' }).click());
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
  const viewTaskButton = page.getByRole('button', { name: `查看任务 ${taskTitle}` });
  await expect(viewTaskButton).toBeVisible();
  await expect(viewTaskButton).toBeEnabled();
  await viewTaskButton.click();
  await expect(page.getByRole('heading', { name: taskTitle })).toBeVisible();
  await expect(page.getByRole('button', { name: '完成任务', exact: true })).toBeEnabled();
  const completedTask = await waitForTaskCompletion(page, createdTask.id, () => page.getByRole('button', { name: '完成任务', exact: true }).click());
  expect(completedTask.status).toBe('Completed');
  await expect(page.getByRole('heading', { name: taskTitle })).toBeVisible();
  await expect(page.getByLabel('任务详情').getByText('已完成').first()).toBeVisible();
});

test('TEST-ACTIVITY-CONTEXT-001/002/003 persists tasks for lead contract and payment contexts', async ({ page }) => {
  const contexts = await createWorkContextFixtures(page, String(Date.now()));
  await openTasks(page);

  for (const context of contexts) {
    await page.getByRole('button', { name: '新建任务', exact: true }).click();
    await page.getByLabel('关联类型').selectOption(context.relatedType);
    await page.getByLabel('关联记录 ID').fill(context.relatedId);
    await page.getByLabel('任务标题').fill(context.taskTitle);
    await page.getByLabel('任务到期日').fill('2027-04-15');
    await page.getByLabel('负责人 ID').fill('sales-1');

    const created = await waitForTaskCreate(page, () => page.getByRole('button', { name: '保存任务', exact: true }).click());
    expect(created.title).toBe(context.taskTitle);
    expect(created.relatedType).toBe(context.relatedType);
    expect(created.relatedId).toBe(context.relatedId);
    expect(created.ownerId).toBe('sales-1');

    const detail = page.getByLabel('任务详情');
    await expect(detail.getByRole('heading', { name: context.taskTitle })).toBeVisible();
    await expect(detail).toContainText(`${context.label} ${context.relatedId}`);
    await expect(detail).toContainText('负责人 sales-1');

    const persisted = await findPersistedTask(page, context.relatedType, context.relatedId, context.taskTitle);
    expect(persisted?.id).toBe(created.id);
    expect(persisted?.relatedType).toBe(context.relatedType);
    expect(persisted?.relatedId).toBe(context.relatedId);

    await page.getByRole('button', { name: '返回任务列表' }).click();
    await expect(page.getByRole('heading', { name: '任务', exact: true })).toBeVisible();
  }

  for (const context of contexts) {
    await expectRecordDetailVisible(page, context);
  }
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

async function openTasks(page: Page) {
  await page.getByRole('button', { name: '任务', exact: true }).click();
  await expect(page.getByRole('heading', { name: '任务', exact: true })).toBeVisible();
}

async function waitForTaskCreate(page: Page, action: () => Promise<unknown>) {
  const responsePromise = page.waitForResponse((response) => (
    new URL(response.url()).pathname === '/api/tasks'
    && response.request().method() === 'POST'
  ));
  await action();
  const response = await responsePromise;
  expect(response.status()).toBeGreaterThanOrEqual(200);
  expect(response.status()).toBeLessThan(300);
  const body = await response.json();
  return body.data as E2EWorkTask;
}

async function waitForTaskListContaining(page: Page, taskId: string, taskTitle: string, action: () => Promise<unknown>) {
  let matchedTask: E2EWorkTask | undefined;
  const responsePromise = page.waitForResponse(async (response) => {
    if (new URL(response.url()).pathname !== '/api/tasks' || response.request().method() !== 'GET') return false;
    if (response.status() < 200 || response.status() >= 300) return false;
    const body = await response.json().catch(() => null);
    const items = body?.data?.items;
    if (!Array.isArray(items)) return false;
    matchedTask = items.find((task: E2EWorkTask) => task.id === taskId || task.title === taskTitle);
    return Boolean(matchedTask);
  });
  await action();
  await responsePromise;
  expect(matchedTask).toBeTruthy();
  return matchedTask as E2EWorkTask;
}

async function waitForTaskCompletion(page: Page, taskId: string, action: () => Promise<unknown>) {
  let completedTask: E2EWorkTask | undefined;
  const responsePromise = page.waitForResponse(async (response) => {
    if (new URL(response.url()).pathname !== `/api/tasks/${taskId}/status` || response.request().method() !== 'POST') return false;
    if (response.status() < 200 || response.status() >= 300) return false;
    const body = await response.json().catch(() => null);
    const task = body?.data as E2EWorkTask | undefined;
    if (task?.id !== taskId || task.status !== 'Completed') return false;
    completedTask = task;
    return true;
  });
  await action();
  await responsePromise;
  expect(completedTask).toBeTruthy();
  return completedTask as E2EWorkTask;
}

async function createWorkContextFixtures(page: Page, suffix: string): Promise<WorkContextFixture[]> {
  return page.evaluate(async (suffix) => {
    type GatewayEnvelope<T> = { data: T };
    type Lead = { id: string; companyName: string };
    type ActualPayment = { paymentId: string; contractId: string };

    async function request<T>(path: string, init?: RequestInit): Promise<T> {
      const response = await fetch(path, {
        ...init,
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
          ...(init?.headers ?? {})
        }
      });
      const body = await response.json() as GatewayEnvelope<T> & { error?: { safeMessage?: string } };
      if (!response.ok) {
        throw new Error(body.error?.safeMessage ?? JSON.stringify(body));
      }
      return body.data;
    }

    async function createAcceptedQuote(key: string, amount: string) {
      const quote = await request<Quote>('/api/quotes', {
        method: 'POST',
        body: JSON.stringify({
          opportunityId: `opp_${key}`,
          customerId: `acct_${key}`,
          amount,
          status: 'Draft',
          validityEnd: '2027-12-31',
          ownerId: 'sales-1'
        })
      });
      const sent = await request<Quote>(`/api/quotes/${quote.id}/status`, {
        method: 'POST',
        body: JSON.stringify({ expectedVersion: quote.version, toStatus: 'Sent' })
      });
      return request<Quote>(`/api/quotes/${quote.id}/status`, {
        method: 'POST',
        body: JSON.stringify({ expectedVersion: sent.version, toStatus: 'Accepted' })
      });
    }

    async function createContract(key: string, amount: string) {
      const accepted = await createAcceptedQuote(key, amount);
      return request<Contract>('/api/contracts', {
        method: 'POST',
        body: JSON.stringify({
          quoteId: accepted.id,
          opportunityId: accepted.opportunityId,
          customerId: accepted.customerId,
          amount,
          status: 'Pending Signature',
          contractNote: `ACC012 ${key}`,
          expectedSignedDate: '2027-12-01',
          amountDifferenceReason: '',
          ownerId: 'sales-1'
        })
      });
    }

    const lead = await request<Lead>('/api/leads', {
      method: 'POST',
      body: JSON.stringify({
        leadName: `ACC012 Lead ${suffix}`,
        companyName: `ACC012 Lead ${suffix}`,
        source: 'Website',
        ownerId: 'sales-1'
      })
    });
    const contract = await createContract(`acc012_contract_${suffix}`, '21000.00');
    const paymentContract = await createContract(`acc012_payment_${suffix}`, '16000.00');
    await request(`/api/contracts/${paymentContract.id}/payment-plans`, {
      method: 'POST',
      body: JSON.stringify({
        dueAmount: '16000.00',
        dueDate: '2027-04-01',
        currency: 'CNY'
      })
    });
    const payment = await request<ActualPayment>(`/api/contracts/${paymentContract.id}/payments`, {
      method: 'POST',
      body: JSON.stringify({
        idempotencyKey: `acc012_payment_${suffix}`,
        amount: '6000.00',
        paymentDate: '2027-04-05',
        note: 'ACC012 payment context',
        currency: 'CNY'
      })
    });

    return [
      {
        relatedType: 'Lead',
        relatedId: lead.id,
        label: '线索',
        taskTitle: `ACC012 Lead Task ${suffix}`,
        recordSearch: lead.companyName,
        recordHeading: lead.companyName
      },
      {
        relatedType: 'Contract',
        relatedId: contract.id,
        label: '合同',
        taskTitle: `ACC012 Contract Task ${suffix}`,
        recordSearch: contract.opportunityId,
        recordHeading: contract.id
      },
      {
        relatedType: 'Payment',
        relatedId: paymentContract.id,
        label: '回款',
        taskTitle: `ACC012 Payment Task ${suffix}`,
        recordSearch: paymentContract.id,
        recordHeading: paymentContract.opportunityId,
        actualPaymentId: payment.paymentId
      }
    ] satisfies WorkContextFixture[];
  }, suffix);
}

async function findPersistedTask(page: Page, relatedType: string, relatedId: string, taskTitle: string) {
  return page.evaluate(async ({ relatedType, relatedId, taskTitle }) => {
    const params = new URLSearchParams({ relatedType, relatedId });
    const response = await fetch(`/api/tasks?${params.toString()}`, { credentials: 'include' });
    const body = await response.json();
    if (!response.ok) {
      throw new Error(body.error?.safeMessage ?? JSON.stringify(body));
    }
    const items = body.data?.items ?? [];
    return items.find((task: E2EWorkTask) => (
      task.title === taskTitle
      && task.relatedType === relatedType
      && task.relatedId === relatedId
    )) ?? null;
  }, { relatedType, relatedId, taskTitle }) as Promise<E2EWorkTask | null>;
}

async function expectRecordDetailVisible(page: Page, context: WorkContextFixture) {
  if (context.relatedType === 'Lead') {
    await page.getByRole('button', { name: '线索', exact: true }).click();
    await expect(page.getByRole('heading', { name: '线索', exact: true })).toBeVisible();
    await page.getByPlaceholder('搜索线索名、公司名').fill(context.recordSearch);
    await page.getByRole('button', { name: '应用筛选' }).click();
    await page.getByRole('button', { name: `打开线索 ${context.recordHeading}` }).click();
    await expect(page.getByLabel('线索详情').getByRole('heading', { name: context.recordHeading })).toBeVisible();
    return;
  }

  if (context.relatedType === 'Contract') {
    await page.getByRole('button', { name: '合同', exact: true }).click();
    await expect(page.getByRole('heading', { name: '合同', exact: true })).toBeVisible();
    await page.getByPlaceholder('搜索合同、报价或客户').fill(context.recordSearch);
    await page.getByRole('button', { name: '应用筛选' }).click();
    await page.getByRole('button', { name: `打开合同 ${context.relatedId}`, exact: true }).click();
    await expect(page.getByLabel('合同详情')).toContainText(context.relatedId);
    return;
  }

  await page.getByRole('button', { name: '回款', exact: true }).click();
  await expect(page.getByRole('heading', { name: '回款', exact: true })).toBeVisible();
  await page.getByPlaceholder('搜索合同、商机或客户').fill(context.recordSearch);
  await page.getByRole('button', { name: '应用筛选' }).click();
  await page.getByRole('button', { name: `打开回款合同 ${context.relatedId}`, exact: true }).click();
  await expect(page.getByLabel('回款详情')).toContainText(`合同 ${context.relatedId}`);
  if (context.actualPaymentId) {
    expect(context.actualPaymentId).toMatch(/^payment_/);
  }
}
