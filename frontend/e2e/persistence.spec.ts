import { execSync } from 'node:child_process';
import { expect, test } from '@playwright/test';

const adminEmail = process.env.E2E_ADMIN_EMAIL ?? 'admin@example.com';
const adminPassword = process.env.E2E_ADMIN_PASSWORD ?? 'AdminChangeMe-001!';

async function signIn(page: import('@playwright/test').Page) {
  await page.goto('/');
  await page.getByLabel('邮箱').fill(adminEmail);
  await page.getByLabel('密码').fill(adminPassword);
  await page.getByRole('button', { name: '登录' }).click();
  await expect(page.getByRole('heading', { name: '工作台' })).toBeVisible();
}

async function createLead(page: import('@playwright/test').Page, companyName: string) {
  return page.evaluate(async (name) => {
    const response = await fetch('/api/leads', {
      method: 'POST',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ companyName: name, source: 'Website', ownerId: 'sales-1' })
    });
    const body = await response.json();
    if (!response.ok) {
      throw new Error(body.error?.safeMessage ?? 'create failed');
    }
    if (!body.data?.id) {
      throw new Error(JSON.stringify(body));
    }
    return body.data.id as string;
  }, companyName);
}

async function expectLeadPersisted(page: import('@playwright/test').Page, leadId: string, companyName: string) {
  await expect(async () => {
    const result = await page.evaluate(async ({ id, name }) => {
      const response = await fetch('/api/leads', { credentials: 'include' });
      const body = await response.json();
      const found = Boolean(body.data?.items?.some((item: { id?: string; companyName?: string }) => item.id === id && item.companyName === name));
      return { found: response.ok && found, status: response.status, body };
    }, { id: leadId, name: companyName });
    expect(result.found, JSON.stringify(result)).toBe(true);
  }).toPass({ timeout: 30_000 });
}

test('TEST-PERSISTENCE-001..005 lead data survives refresh, re-login, and service restart; failed save is surfaced', async ({ page }) => {
  test.setTimeout(90_000);
  const suffix = Date.now();
  const companyName = `Persist E2E ${suffix}`;
  const failedCompanyName = `Persist Failed ${suffix}`;

  await signIn(page);
  const leadId = await createLead(page, companyName);
  await expectLeadPersisted(page, leadId, companyName);

  await page.reload();
  await expectLeadPersisted(page, leadId, companyName);

  await page.getByRole('button', { name: '退出登录' }).click();
  await signIn(page);
  await expectLeadPersisted(page, leadId, companyName);

  execSync('docker compose restart lead', { cwd: '..', stdio: 'pipe' });
  await page.reload();
  await expectLeadPersisted(page, leadId, companyName);

  const failed = await page.evaluate(async (name) => {
    const response = await fetch('/api/leads', {
      method: 'POST',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ companyName: name, ownerId: 'sales-1' })
    });
    const body = await response.json();
    return { ok: response.ok, safeMessage: body.error?.safeMessage ?? '' };
  }, failedCompanyName);
  expect(failed.ok).toBe(false);
  expect(failed.safeMessage).toContain('invalid');
});
