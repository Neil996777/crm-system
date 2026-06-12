import { expect, test } from '@playwright/test';

const adminEmail = process.env.E2E_ADMIN_EMAIL ?? 'admin@example.com';
const adminPassword = process.env.E2E_ADMIN_PASSWORD ?? 'AdminChangeMe-001!';

type AccountRecord = {
  id: string;
  companyName: string;
};

type ContactRecord = {
  id: string;
  accountId: string;
  contactName: string;
  email: string;
  phone: string;
};

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

test('TEST-DUPLICATE-CONTACT-001/002 contact email and phone warnings proceed without merge', async ({ page }) => {
  const suffix = Date.now();
  const sourceAccount = await apiRequest<AccountRecord>(page, '/api/accounts', {
    method: 'POST',
    body: JSON.stringify({
      companyName: `E2E Duplicate Contact Source ${suffix}`,
      customerStatus: 'Active',
      ownerId: 'sales-1'
    })
  });
  const targetAccount = await apiRequest<AccountRecord>(page, '/api/accounts', {
    method: 'POST',
    body: JSON.stringify({
      companyName: `E2E Duplicate Contact Target ${suffix}`,
      customerStatus: 'Active',
      ownerId: 'sales-1'
    })
  });
  const existing = await apiRequest<ContactRecord>(page, `/api/accounts/${sourceAccount.id}/contacts`, {
    method: 'POST',
    body: JSON.stringify({
      contactName: `Existing Buyer ${suffix}`,
      email: `buyer-${suffix}@example.com`,
      phone: `139${String(suffix).slice(-8)}`,
      roleNote: '原始联系人'
    })
  });

  await openAccount(page, targetAccount.companyName);
  const emailDuplicateName = `Email Duplicate ${suffix}`;
  await createContactThroughWarning(page, {
    contactName: emailDuplicateName,
    email: existing.email.toUpperCase(),
    phone: `138${String(suffix).slice(-8)}`,
    roleNote: '邮箱匹配但允许新增'
  });
  await expect(page.getByRole('table', { name: '联系人' }).getByText(emailDuplicateName)).toBeVisible();

  const phoneDuplicateName = `Phone Duplicate ${suffix}`;
  await createContactThroughWarning(page, {
    contactName: phoneDuplicateName,
    email: `phone-unique-${suffix}@example.com`,
    phone: `+86 ${existing.phone}`,
    roleNote: '电话匹配但允许新增'
  });
  await expect(page.getByRole('table', { name: '联系人' }).getByText(phoneDuplicateName)).toBeVisible();

  const sourceContacts = await listContacts(page, sourceAccount.id);
  const targetContacts = await listContacts(page, targetAccount.id);
  expect(sourceContacts.filter((contact) => contact.id === existing.id)).toHaveLength(1);
  expect(sourceContacts.find((contact) => contact.id === existing.id)?.contactName).toBe(existing.contactName);
  const createdEmailDuplicate = targetContacts.find((contact) => contact.contactName === emailDuplicateName);
  const createdPhoneDuplicate = targetContacts.find((contact) => contact.contactName === phoneDuplicateName);
  expect(createdEmailDuplicate?.id).toBeTruthy();
  expect(createdPhoneDuplicate?.id).toBeTruthy();
  expect(new Set([existing.id, createdEmailDuplicate?.id, createdPhoneDuplicate?.id]).size).toBe(3);
});

async function openAccount(page: import('@playwright/test').Page, companyName: string) {
  await page.getByRole('button', { name: '公司/客户', exact: true }).click();
  await expect(page.getByRole('heading', { name: '公司/客户', exact: true })).toBeVisible();
  await page.getByPlaceholder('搜索公司客户').fill(companyName);
  await page.getByRole('button', { name: '应用筛选' }).click();
  await page.getByRole('button', { name: `打开客户 ${companyName}` }).click();
  await expect(page.getByLabel('客户详情').getByRole('heading', { name: companyName })).toBeVisible();
}

async function createContactThroughWarning(
  page: import('@playwright/test').Page,
  input: { contactName: string; email: string; phone: string; roleNote: string }
) {
  await page.getByRole('button', { name: '添加联系人' }).click();
  const form = page.locator('form.createPanel').last();
  await form.getByLabel('联系人姓名').fill(input.contactName);
  await form.getByLabel('邮箱').fill(input.email);
  await form.getByLabel('电话').fill(input.phone);
  await form.getByLabel('角色备注').fill(input.roleNote);
  await form.getByRole('button', { name: '保存联系人' }).click();

  await expect(page.getByRole('alert')).toContainText('可能重复');
  await page.getByRole('button', { name: '仍然创建' }).click();
  await expect(page.getByRole('alert')).toHaveCount(0);
}

async function listContacts(page: import('@playwright/test').Page, accountId: string) {
  const response = await apiRequest<{ items: ContactRecord[] }>(page, `/api/accounts/${accountId}/contacts`);
  return response.items;
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
