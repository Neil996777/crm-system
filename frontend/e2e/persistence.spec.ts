import { execSync } from 'node:child_process';
import { expect, test } from '@playwright/test';

const adminEmail = process.env.E2E_ADMIN_EMAIL ?? 'admin@example.com';
const adminPassword = process.env.E2E_ADMIN_PASSWORD ?? 'AdminChangeMe-001!';

async function signIn(page: import('@playwright/test').Page, options: { navigate?: boolean } = {}) {
  if (options.navigate !== false) {
    await page.goto('/', { waitUntil: 'domcontentloaded' });
  }
  const authState = await waitForInteractiveAuth(page);
  if (authState === 'session') {
    return;
  }
  const emailInput = page.locator('input[type="email"]');
  const passwordInput = page.locator('input[type="password"]');
  await expect(emailInput).toBeEditable({ timeout: 30_000 });
  await expect(passwordInput).toBeEditable({ timeout: 30_000 });
  await emailInput.fill(adminEmail);
  await passwordInput.fill(adminPassword);
  await page.getByRole('button', { name: '登录' }).click();
  await expect(page.getByRole('heading', { name: '工作台' })).toBeVisible();
  await waitForLeadGatewayAuthenticated(page);
}

async function waitForInteractiveAuth(page: import('@playwright/test').Page) {
  await expect.poll(() => currentAuthState(page), { timeout: 30_000, intervals: [250, 500, 1_000, 2_000] }).not.toBe('loading');
  return currentAuthState(page);
}

async function currentAuthState(page: import('@playwright/test').Page): Promise<'session' | 'login' | 'loading'> {
  const logoutVisible = await page.getByRole('button', { name: '退出登录' }).isVisible().catch(() => false);
  if (logoutVisible) return 'session';
  const emailReady = await page.locator('input[type="email"]').isEnabled().catch(() => false);
  const passwordReady = await page.locator('input[type="password"]').isEnabled().catch(() => false);
  if (emailReady && passwordReady) return 'login';
  return 'loading';
}

async function waitForLoginForm(page: import('@playwright/test').Page) {
  await expect.poll(() => currentAuthState(page), { timeout: 30_000, intervals: [250, 500, 1_000, 2_000] }).toBe('login');
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

async function waitForLeadServiceReady(page: import('@playwright/test').Page, leadId: string, companyName: string) {
  await expect.poll(async () => {
    return page.evaluate(async ({ id, name }) => {
      try {
        const response = await fetch('/api/leads', { credentials: 'include' });
        if (!response.ok) return `status:${response.status}`;
        const body = await response.json();
        const found = Boolean(body.data?.items?.some((item: { id?: string; companyName?: string }) => item.id === id && item.companyName === name));
        return found ? 'ready' : 'missing';
      } catch (error) {
        return error instanceof Error ? `error:${error.message}` : 'error:unknown';
      }
    }, { id: leadId, name: companyName });
  }, { timeout: 75_000, intervals: [500, 1_000, 2_000, 5_000] }).toBe('ready');
}

async function waitForLeadGatewayAuthenticated(page: import('@playwright/test').Page) {
  await expect.poll(async () => {
    return page.evaluate(async () => {
      try {
        const response = await fetch('/api/leads', { credentials: 'include' });
        return response.status;
      } catch {
        return 0;
      }
    });
  }, { timeout: 30_000, intervals: [250, 500, 1_000, 2_000] }).toBe(200);
}

test('TEST-PERSISTENCE-001..005 lead data survives refresh, re-login, and service restart; failed save is surfaced', async ({ page }) => {
  test.setTimeout(180_000);
  const suffix = Date.now();
  const companyName = `Persist E2E ${suffix}`;
  const failedCompanyName = `Persist Failed ${suffix}`;

  await signIn(page);
  const leadId = await createLead(page, companyName);
  await expectLeadPersisted(page, leadId, companyName);

  await page.reload();
  await expectLeadPersisted(page, leadId, companyName);

  await page.getByRole('button', { name: '退出登录' }).click();
  await waitForLoginForm(page);
  await signIn(page, { navigate: false });
  await expectLeadPersisted(page, leadId, companyName);

  execSync('docker compose restart lead', { cwd: '..', stdio: 'pipe', timeout: 90_000 });
  await waitForLeadServiceReady(page, leadId, companyName);
  await page.reload({ waitUntil: 'domcontentloaded' });
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
