import { expect, type Page, test } from '@playwright/test';

const adminEmail = process.env.E2E_ADMIN_EMAIL ?? 'admin@example.com';
const adminPassword = process.env.E2E_ADMIN_PASSWORD ?? 'AdminChangeMe-001!';

type FixtureRecord = {
  id: string;
  name: string;
  search: string;
  titleButton: string;
  nav: string;
  heading: string;
  table: string;
  detailLabel: string;
  detailText: string;
  backButton: string;
  menuButton?: RegExp;
};

test.beforeEach(async ({ page }) => {
  await page.goto('/');
  await page.getByLabel('邮箱').fill(adminEmail);
  await page.getByLabel('密码').fill(adminPassword);
  await page.getByRole('button', { name: '登录' }).click();
  await expect(page.getByRole('heading', { name: '工作台' })).toBeVisible();
});

test('TEST-UIUX-G12-018 row action menu portal stays topmost and clickable on first middle last rows', async ({ page }) => {
  const prefix = `BLK018 Menu ${Date.now()}`;
  await createLeadBatch(page, prefix, 30);
  await openList(page, '线索', '线索');
  await filterCurrentList(page, '搜索线索名、公司名', prefix);

  const dataRows = leadRows(page, prefix);
  await expect(dataRows).toHaveCount(25);

  for (const index of [0, 12, 24]) {
    const row = dataRows.nth(index);
    await row.scrollIntoViewIfNeeded();
    const titleButton = row.getByRole('button', { name: /^打开线索 / });
    const label = await titleButton.getAttribute('aria-label');
    const companyName = label?.replace('打开线索 ', '') ?? prefix;

    await row.getByRole('button', { name: /打开 .* 的行操作菜单/ }).click();
    await expectMenuTopmost(page);
    await page.getByRole('menuitem', { name: '查看' }).click();
    await expect(page.getByLabel('线索详情').getByRole('heading', { name: companyName })).toBeVisible();
    await page.getByRole('button', { name: '返回线索列表' }).click();
    await expect(page.getByRole('heading', { name: '线索', exact: true })).toBeVisible();
  }
});

test('TEST-UIUX-G12-018 seven lists use clickable record names and no legacy arrow view buttons', async ({ page }) => {
  const suffix = String(Date.now());
  const fixtures = await createListFixtures(page, suffix);

  for (const record of fixtures) {
    await openList(page, record.nav, record.heading);
    await filterCurrentList(page, searchPlaceholderFor(record.nav), record.search);
    const row = page.getByRole('table', { name: record.table }).locator('tbody tr').filter({ hasText: record.search }).first();
    await expect(row).toBeVisible();
    await expect(row.getByRole('button', { name: /^(查看|查看报价|查看合同|查看回款)/ })).toHaveCount(0);

    if (record.menuButton) {
      await row.getByRole('button', { name: record.menuButton }).click();
      await expectMenuTopmost(page);
      await page.getByRole('menuitem', { name: '查看' }).click();
      await expect(page.getByLabel(record.detailLabel)).toContainText(record.detailText);
      await page.getByRole('button', { name: record.backButton }).click();
      await expect(page.getByRole('heading', { name: record.heading, exact: true })).toBeVisible();
    }

    const refreshedRow = page.getByRole('table', { name: record.table }).locator('tbody tr').filter({ hasText: record.search }).first();
    await refreshedRow.getByRole('button', { name: record.titleButton, exact: true }).click();
    await expect(page.getByLabel(record.detailLabel)).toContainText(record.detailText);
  }
});

async function openList(page: Page, nav: string, heading: string) {
  await page.getByLabel('主导航').getByRole('button', { name: nav, exact: true }).click();
  await expect(page.getByRole('heading', { name: heading, exact: true })).toBeVisible();
}

async function filterCurrentList(page: Page, placeholder: string, value: string) {
  await page.getByPlaceholder(placeholder).fill(value);
  await page.getByRole('button', { name: '应用筛选' }).click();
}

function leadRows(page: Page, prefix: string) {
  return page.getByRole('table', { name: '线索结果表' }).locator('tbody tr').filter({ hasText: prefix });
}

async function expectMenuTopmost(page: Page) {
  const menu = page.getByRole('menu');
  await expect(menu).toBeVisible();
  const result = await menu.evaluate((element) => {
    const rect = element.getBoundingClientRect();
    const x = rect.left + rect.width / 2;
    const y = rect.top + rect.height / 2;
    const top = document.elementFromPoint(x, y);
    return {
      insideMenu: top === element || element.contains(top),
      tagName: top?.tagName ?? '',
      className: top instanceof HTMLElement ? top.className : '',
      role: top instanceof HTMLElement ? top.getAttribute('role') : null
    };
  });
  expect(result.insideMenu, JSON.stringify(result)).toBe(true);
}

async function createLeadBatch(page: Page, prefix: string, count: number) {
  await page.evaluate(async ({ prefix, count }) => {
    for (let index = 0; index < count; index += 1) {
      await request('/api/leads', {
        method: 'POST',
        body: JSON.stringify({
          leadName: `${prefix} ${String(index).padStart(2, '0')}`,
          companyName: `${prefix} ${String(index).padStart(2, '0')}`,
          source: 'Website',
          ownerId: 'sales-1'
        })
      });
    }

    async function request<T>(path: string, init?: RequestInit): Promise<T> {
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
      return body.data as T;
    }
  }, { prefix, count });
}

async function createListFixtures(page: Page, suffix: string): Promise<FixtureRecord[]> {
  return page.evaluate(async (suffix) => {
    const lead = await request<{ id: string; companyName: string }>('/api/leads', {
      method: 'POST',
      body: JSON.stringify({
        leadName: `BLK018 Lead ${suffix}`,
        companyName: `BLK018 Lead ${suffix}`,
        source: 'Website',
        ownerId: 'sales-1'
      })
    });
    const account = await request<{ id: string; companyName: string }>('/api/accounts', {
      method: 'POST',
      body: JSON.stringify({ companyName: `BLK018 Account ${suffix}`, customerStatus: 'Active', ownerId: 'sales-1' })
    });
    const contact = await request<{ id: string; contactName: string }>(`/api/accounts/${account.id}/contacts`, {
      method: 'POST',
      body: JSON.stringify({
        contactName: `BLK018 Contact ${suffix}`,
        email: `blk018-${suffix}@example.com`,
        phone: `139${suffix.slice(-8)}`,
        roleNote: '采购联系人'
      })
    });
    const opportunity = await request<{ id: string; title: string; customerId: string }>('/api/opportunities', {
      method: 'POST',
      body: JSON.stringify({
        customerId: account.id,
        ownerId: 'sales-1',
        stage: 'New Opportunity',
        expectedAmount: '18000.00',
        expectedCloseDate: '2027-09-01',
        title: `BLK018 Opportunity ${suffix}`
      })
    });
    const quote = await request<{ id: string; opportunityId: string; customerId: string; amount: string; version: number }>('/api/quotes', {
      method: 'POST',
      body: JSON.stringify({
        opportunityId: opportunity.id,
        customerId: account.id,
        amount: '18000.00',
        status: 'Draft',
        validityEnd: '2027-12-31',
        ownerId: 'sales-1'
      })
    });
    const sentQuote = await request<{ id: string; opportunityId: string; customerId: string; amount: string; version: number }>(`/api/quotes/${quote.id}/status`, {
      method: 'POST',
      body: JSON.stringify({ expectedVersion: quote.version, toStatus: 'Sent' })
    });
    const acceptedQuote = await request<{ id: string; opportunityId: string; customerId: string; amount: string; version: number }>(`/api/quotes/${sentQuote.id}/status`, {
      method: 'POST',
      body: JSON.stringify({ expectedVersion: sentQuote.version, toStatus: 'Accepted' })
    });
    const contract = await request<{ id: string; opportunityId: string }>('/api/contracts', {
      method: 'POST',
      body: JSON.stringify({
        quoteId: acceptedQuote.id,
        opportunityId: acceptedQuote.opportunityId,
        customerId: acceptedQuote.customerId,
        amount: acceptedQuote.amount,
        status: 'Pending Signature',
        contractNote: 'BLK018 合同',
        expectedSignedDate: '2027-12-01',
        amountDifferenceReason: '',
        ownerId: 'sales-1'
      })
    });

    return [
      {
        id: lead.id,
        name: lead.companyName,
        search: lead.companyName,
        titleButton: `打开线索 ${lead.companyName}`,
        nav: '线索',
        heading: '线索',
        table: '线索结果表',
        detailLabel: '线索详情',
        detailText: lead.companyName,
        backButton: '返回线索列表',
        menuButton: /打开 .* 的行操作菜单/
      },
      {
        id: account.id,
        name: account.companyName,
        search: account.companyName,
        titleButton: `打开客户 ${account.companyName}`,
        nav: '公司/客户',
        heading: '公司/客户',
        table: '客户结果表',
        detailLabel: '客户详情',
        detailText: account.companyName,
        backButton: '返回客户列表',
        menuButton: /打开 .* 的行操作菜单/
      },
      {
        id: opportunity.id,
        name: opportunity.title,
        search: opportunity.title,
        titleButton: `打开商机 ${opportunity.title}`,
        nav: '商机',
        heading: '商机',
        table: '商机结果表',
        detailLabel: '商机详情',
        detailText: opportunity.title,
        backButton: '返回商机列表',
        menuButton: /打开 .* 的行操作菜单/
      },
      {
        id: acceptedQuote.id,
        name: acceptedQuote.opportunityId,
        search: acceptedQuote.opportunityId,
        titleButton: `打开报价 ${acceptedQuote.id}`,
        nav: '报价',
        heading: '报价',
        table: '报价结果表',
        detailLabel: '报价详情',
        detailText: acceptedQuote.opportunityId,
        backButton: '返回报价列表',
        menuButton: /打开报价 .* 的行操作菜单/
      },
      {
        id: contract.id,
        name: contract.opportunityId,
        search: contract.opportunityId,
        titleButton: `打开合同 ${contract.id}`,
        nav: '合同',
        heading: '合同',
        table: '合同结果表',
        detailLabel: '合同详情',
        detailText: contract.id,
        backButton: '返回合同列表',
        menuButton: /打开合同 .* 的行操作菜单/
      },
      {
        id: contract.id,
        name: contract.opportunityId,
        search: contract.opportunityId,
        titleButton: `打开回款合同 ${contract.id}`,
        nav: '回款',
        heading: '回款',
        table: '回款合同结果表',
        detailLabel: '回款详情',
        detailText: contract.opportunityId,
        backButton: '返回回款列表',
        menuButton: /打开回款 .* 的行操作菜单/
      },
      {
        id: contact.id,
        name: contact.contactName,
        search: contact.contactName,
        titleButton: `打开联系人 ${contact.contactName}`,
        nav: '联系人',
        heading: '联系人',
        table: '联系人结果表',
        detailLabel: '联系人详情',
        detailText: contact.contactName,
        backButton: '返回联系人列表'
      }
    ] satisfies FixtureRecord[];

    async function request<T>(path: string, init?: RequestInit): Promise<T> {
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
      return body.data as T;
    }
  }, suffix);
}

function searchPlaceholderFor(nav: string) {
  if (nav === '线索') return '搜索线索名、公司名';
  if (nav === '公司/客户') return '搜索公司客户';
  if (nav === '联系人') return '搜索联系人或客户';
  if (nav === '商机') return '搜索商机名、客户名';
  if (nav === '报价') return '搜索报价、商机或客户';
  if (nav === '合同') return '搜索合同、报价或客户';
  return '搜索合同、商机或客户';
}
