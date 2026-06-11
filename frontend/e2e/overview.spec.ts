import { expect, test, type Page } from '@playwright/test';

const adminEmail = process.env.E2E_ADMIN_EMAIL ?? 'admin@example.com';
const adminPassword = process.env.E2E_ADMIN_PASSWORD ?? 'AdminChangeMe-001!';
const dashboardPassword = 'Dashboard-001!';
const managerFocusRailKeys = ['funnel', 'stage', 'trend', 'leaderboard', 'todo', 'payments', 'key-opportunities', 'activity'];
const salesFocusRailKeys = ['funnel', 'todo', 'stage', 'trend', 'payments', 'activity'];

test.beforeEach(async ({ page }) => {
  await page.goto('/');
  await page.getByLabel('邮箱').fill(adminEmail);
  await page.getByLabel('密码').fill(adminPassword);
  await page.getByRole('button', { name: '登录' }).click();
  await expect(page.getByRole('heading', { name: '工作台' })).toBeVisible();
});

test('TEST-TEAM-OVERVIEW-003 renders manager overview empty state through gateway', async ({ page }) => {
  const overview = await page.evaluate(async () => {
    const response = await fetch('/api/reports/team-overview', { credentials: 'include' });
    const body = await response.json();
    return { status: response.status, scope: body.data?.scope, hasMetrics: Boolean(body.data?.metrics) };
  });
  expect(overview.status).toBe(200);
  expect(overview.scope).toBe('all');
  expect(overview.hasMetrics).toBe(true);

  await page.getByLabel('主导航').getByRole('button', { name: '报表' }).click();
  await expect(page.getByRole('heading', { name: '团队报表' })).toBeVisible();
  await expect(page.getByRole('heading', { name: '经理团队总览' })).toHaveCount(0);
  await expect(page.locator('[data-uiux="reports-team"]')).toBeVisible();
  await expect(page.locator('.reportsKpiStrip .metricTile')).toHaveCount(9);
  await expect(page.getByLabel('管道分析')).toBeVisible();
  await expect(page.getByLabel('负责人分组')).toBeVisible();
});

test('TEST-UIUX-DASHBOARD-001 renders manager dashboard 8 cards and per-card focus stage', async ({ page }) => {
  await expect(page.locator('[data-uiux="dashboard"]')).toBeVisible();
  await expect(page.locator('section[aria-label="今日实时战报"]')).toBeVisible();
  await expect(page.locator('.dashboardKpis .metricTile')).toHaveCount(4);
  await expect(page.locator('[aria-label="管端工作台数据卡"] [data-dashboard-card]')).toHaveCount(8);
  await expect(page.getByLabel('管端工作台数据卡')).toContainText('团队销售漏斗');
  await expect(page.getByLabel('管端工作台数据卡')).toContainText('商机阶段构成');
  await expect(page.getByLabel('管端工作台数据卡')).toContainText('团队赢单金额趋势');
  await expect(page.getByLabel('管端工作台数据卡')).toContainText('销售业绩榜');
  await expect(page.getByLabel('管端工作台数据卡')).toContainText('团队待办与预警');
  await expect(page.getByLabel('管端工作台数据卡')).toContainText('团队回款到账');
  await expect(page.getByLabel('管端工作台数据卡')).toContainText('重点商机');
  await expect(page.getByLabel('管端工作台数据卡')).toContainText('团队最近活动');
  await expect(page.getByRole('button', { name: '专注模式' })).toHaveCount(0);

  const managerFunnelCard = page.locator('[data-dashboard-card="funnel"]');
  await expect(managerFunnelCard).toHaveAttribute('role', 'button');
  await expect(managerFunnelCard).toHaveAttribute('tabindex', '0');
  await managerFunnelCard.click();
  await expect(page.locator('[data-uiux="dashboard-focus"]')).toBeVisible();
  await expect(page.locator('.shell.focusMode')).toBeVisible();
  await expect(page.getByRole('heading', { name: '团队销售漏斗' })).toBeVisible();
  await expectSingleFocusExitControl(page);
  await expect(page.getByLabel('看板选择器').locator('.sideCard')).toHaveCount(8);
  await expect.poll(() => focusRailKeys(page)).toEqual(managerFocusRailKeys);
  await expectFocusRailSelection(page, 'funnel');
  await expect(page.getByRole('button', { name: /商机阶段构成/ })).toBeVisible();
  await page.getByRole('button', { name: /商机阶段构成/ }).click();
  await expect(page.getByRole('heading', { name: '团队商机阶段构成' })).toBeVisible();
  await expect.poll(() => focusRailKeys(page)).toEqual(managerFocusRailKeys);
  await expectFocusRailSelection(page, 'stage');
  await page.keyboard.press('Escape');
  await expect(page.locator('[data-uiux="dashboard"]')).toBeVisible();
  await expect(page.locator('.shell.focusMode')).toHaveCount(0);
  await page.locator('[data-dashboard-card="activity"]').click();
  await expect(page.locator('[data-uiux="dashboard-focus"]')).toBeVisible();
  await expectSingleFocusExitControl(page);
  await page.getByRole('button', { name: '返回' }).click();
  await expect(page.locator('[data-uiux="dashboard"]')).toBeVisible();
});

test('TEST-UIUX-A6-001 dashboard remains stable on narrow desktop viewport', async ({ page }) => {
  const dashboardWidths = [
    { width: 1700, height: 900 },
    { width: 1680, height: 900 },
    { width: 1600, height: 900 },
    { width: 1512, height: 900 },
    { width: 1440, height: 900 },
    { width: 1280, height: 900 },
    { width: 1180, height: 900 },
    { width: 1024, height: 768 },
    { width: 900, height: 720 }
  ];

  for (const viewport of dashboardWidths) {
    await page.setViewportSize(viewport);
    await expect(page.locator('[data-uiux="dashboard"]')).toBeVisible();
    await expect(page.locator('section[aria-label="今日实时战报"]')).toBeVisible();
    await expect(page.locator('[data-dashboard-card]')).toHaveCount(8);
    await expectDashboardCardsNotClipped(page);
    await expectDashboardInlineContentStable(page);
    if (viewport.width >= 1024) {
      await expectDashboardReadableTextNotTruncated(page);
      await expectDashboardFlowRowsDoNotOverlap(page);
    }
    const overflow = await page.evaluate(() => Math.max(
      document.documentElement.scrollWidth - document.documentElement.clientWidth,
      document.body.scrollWidth - document.body.clientWidth
    ));
    expect(overflow).toBeLessThanOrEqual(1);
  }
});

test('TEST-UIUX-DASHBOARD-002 renders sales personal dashboard variant without manager cards', async ({ page }) => {
  const suffix = Date.now();
  const sales = await createUser(page, `dashboard-sales-${suffix}@example.com`, `Dashboard Sales ${suffix}`, 'Sales');

  await page.getByRole('button', { name: '退出登录' }).click();
  await signIn(page, sales.email, dashboardPassword);
  await expect(page.getByRole('heading', { name: '我的工作台' })).toBeVisible();
  await expect(page.locator('[data-uiux="dashboard"]')).toBeVisible();
  await expect(page.locator('.dashboardKpis .metricTile')).toHaveCount(4);
  await expect(page.locator('[aria-label="销售工作台数据卡"] [data-dashboard-card]')).toHaveCount(6);
  await expect(page.getByLabel('销售工作台数据卡')).toContainText('我的销售漏斗');
  await expect(page.getByLabel('销售工作台数据卡')).toContainText('我的待办与预警');
  await expect(page.getByLabel('销售工作台数据卡')).toContainText('我的商机阶段构成');
  await expect(page.getByLabel('销售工作台数据卡')).toContainText('我的赢单金额趋势');
  await expect(page.getByLabel('销售工作台数据卡')).toContainText('我的回款到账');
  await expect(page.getByLabel('销售工作台数据卡')).toContainText('我的最近活动');
  await expect(page.getByLabel('销售工作台数据卡')).not.toContainText('销售业绩榜');

  await page.locator('[data-dashboard-card="funnel"]').click();
  await expect(page.locator('[data-uiux="dashboard-focus"]')).toBeVisible();
  await expect(page.getByRole('heading', { name: '我的销售漏斗' })).toBeVisible();
  await expectSingleFocusExitControl(page);
  await expect(page.getByLabel('看板选择器').locator('.sideCard')).toHaveCount(6);
  await expect.poll(() => focusRailKeys(page)).toEqual(salesFocusRailKeys);
  await expectFocusRailSelection(page, 'funnel');
  await page.keyboard.press('Escape');
  await expect(page.locator('[data-uiux="dashboard"]')).toBeVisible();
});

test('TEST-UIUX-B3-001 dashboard live polling applies buffered updates without full reload', async ({ page }) => {
  await page.addInitScript(() => {
    (window as Window & { __dashboardLivePollMs?: number; __dashboardLiveCoalesceMs?: number }).__dashboardLivePollMs = 1_000;
    (window as Window & { __dashboardLivePollMs?: number; __dashboardLiveCoalesceMs?: number }).__dashboardLiveCoalesceMs = 100;
  });

  let activityVersion = 0;
  let delayNextActivities = false;
  await page.route('**/api/activities*', async (route) => {
    if (delayNextActivities) {
      delayNextActivities = false;
      await new Promise((resolve) => setTimeout(resolve, 350));
    }
    const response = await route.fetch();
    const body = await response.json();
    const items = body.data?.items;
    if (Array.isArray(items) && activityVersion > 0) {
      body.data.items = [
        {
          id: `e2e-live-activity-${activityVersion}`,
          relatedType: 'Opportunity',
          relatedId: `OPP-LIVE-${activityVersion}`,
          activityType: 'note',
          content: `轮询动态 ${activityVersion}`,
          ownerId: 'live-e2e',
          occurredAt: new Date(Date.now() + activityVersion * 1_000).toISOString()
        },
        ...items.filter((item: { id?: string }) => !item.id?.startsWith('e2e-live-activity-'))
      ];
    }
    await route.fulfill({ response, json: body });
  });

  await page.reload();
  await expect(page.locator('[data-uiux="dashboard"]')).toBeVisible();
  await expect(page.getByRole('button', { name: '本月' })).toHaveCount(0);
  await expect(page.getByText('自动合并')).toHaveCount(0);
  await expect(page.getByRole('button', { name: '实时更新' })).toHaveAttribute('aria-pressed', 'true');
  const activityCard = page.locator('[data-dashboard-card="activity"]');
  await expect(activityCard.locator('.event.arrived')).toHaveCount(0);
  await page.evaluate(() => {
    sessionStorage.setItem('__dashboardReloaded', 'false');
    window.addEventListener('beforeunload', () => sessionStorage.setItem('__dashboardReloaded', 'true'));
  });

  activityVersion = 1;
  await expect(activityCard).toContainText('轮询动态 1', { timeout: 4_000 });
  await expect(activityCard.locator('.event.arrived')).toHaveCount(1);
  await expect(activityCard.locator('.event.arrived')).toContainText('轮询动态 1');
  expect(await page.evaluate(() => sessionStorage.getItem('__dashboardReloaded'))).toBe('false');

  await page.getByRole('button', { name: '实时更新' }).click();
  await expect(page.getByRole('button', { name: '暂停' })).toHaveAttribute('aria-pressed', 'false');
  await expect(page.locator('.reportLead .liveDot.paused')).toBeVisible();
  activityVersion = 2;
  await expect(page.getByRole('button', { name: /有 \d+ 条新更新 · 点击刷新/ })).toBeVisible({ timeout: 4_000 });
  await expect(activityCard).not.toContainText('轮询动态 2');
  await page.getByRole('button', { name: '暂停' }).click();
  await expect(page.getByRole('button', { name: '实时更新' })).toHaveAttribute('aria-pressed', 'true');
  await expect(activityCard).toContainText('轮询动态 2');
  await expect(activityCard.locator('.event.arrived')).toHaveCount(1);
  await expect(activityCard.locator('.event.arrived')).toContainText('轮询动态 2');

  activityVersion = 3;
  delayNextActivities = true;
  await page.getByRole('button', { name: '刷新数据' }).click();
  await expect(page.getByRole('button', { name: '刷新中' })).toBeVisible();
  await expect(page.locator('.updateMeta')).toContainText(/已刷新/, { timeout: 4_000 });
  await expect(activityCard).toContainText('轮询动态 3');
});

test('TEST-UIUX-A7-001 card focus respects reduced-motion mode and still snaps between states', async ({ page }) => {
  await page.emulateMedia({ reducedMotion: 'no-preference' });
  await installDashboardAnimationRecorder(page);
  await expect(page.locator('[data-uiux="dashboard"]')).toBeVisible();
  const paymentsCard = page.locator('[data-dashboard-card="payments"]');
  await paymentsCard.focus();
  await expect(paymentsCard).toBeFocused();
  await resetDashboardAnimationRecorder(page);
  await page.keyboard.press(' ');
  await expect(page.locator('[data-uiux="dashboard-focus"]')).toBeVisible();
  await expect(page.locator('[data-uiux="dashboard-focus"]')).toHaveAttribute('data-motion-mode', 'full');
  await expect(page.locator('.shell.focusMode')).toBeVisible();
  await expectSingleFocusExitControl(page);
  await expect(page.getByLabel('看板选择器').locator('.sideCard')).toHaveCount(8);
  await expect.poll(() => focusRailKeys(page)).toEqual(managerFocusRailKeys);
  await expectFocusRailSelection(page, 'payments');
  await expectDashboardAnimationStarted(page, 'dashboardStageEnter');
  await expectDashboardAnimationTotal(page, 'dashboardStageEnter', 450);
  await expectDashboardAnimationStarted(page, 'dashboardStripEnter');
  await expect(page.getByRole('heading', { name: '团队回款到账' })).toBeFocused({ timeout: 2_000 });

  await resetDashboardAnimationRecorder(page);
  const railOrderBeforeSwitch = await focusRailKeys(page);
  await page.getByRole('button', { name: /商机阶段构成/ }).click();
  await expect(page.getByRole('heading', { name: '团队商机阶段构成' })).toBeVisible();
  await expect.poll(() => focusRailKeys(page)).toEqual(railOrderBeforeSwitch);
  await expectFocusRailSelection(page, 'stage');
  await expectDashboardAnimationStarted(page, 'dashboardStageSwitch');
  await expectDashboardAnimationTotal(page, 'dashboardStageSwitch', 220);
  const switchAnimations = await dashboardAnimationNames(page);
  expect(switchAnimations).not.toContain('dashboardStageExit');

  await resetDashboardAnimationRecorder(page);
  await page.keyboard.press('Escape');
  await expectDashboardAnimationStarted(page, 'dashboardStageExit');
  await expectDashboardAnimationTotal(page, 'dashboardStageExit', 310);
  await expectDashboardAnimationStarted(page, 'dashboardStripExit');
  await expect(page.locator('[data-uiux="dashboard"]')).toBeVisible();
  await expect(page.locator('[data-dashboard-card="stage"]')).toBeFocused();

  await page.emulateMedia({ reducedMotion: 'reduce' });
  await resetDashboardAnimationRecorder(page);
  const reducedPaymentsCard = page.locator('[data-dashboard-card="payments"]');
  await reducedPaymentsCard.focus();
  await page.keyboard.press('Enter');
  await expect(page.locator('[data-uiux="dashboard-focus"]')).toBeVisible();
  await expect(page.locator('[data-uiux="dashboard-focus"]')).toHaveAttribute('data-motion-mode', 'reduced');
  await expect(page.getByRole('heading', { name: '团队回款到账' })).toBeVisible();
  await expectSingleFocusExitControl(page);
  await expect(page.getByLabel('看板选择器').locator('.sideCard')).toHaveCount(8);
  await expect.poll(() => focusRailKeys(page)).toEqual(managerFocusRailKeys);
  await expectFocusRailSelection(page, 'payments');
  const reducedTransform = await page.locator('.stage').evaluate((stage) => getComputedStyle(stage).transform);
  expect(reducedTransform === 'none' || reducedTransform === 'matrix(1, 0, 0, 1, 0, 0)').toBeTruthy();
  await expectDashboardAnimationStarted(page, 'dashboardReducedFocusAppear');
  const reducedAnimations = await dashboardAnimationNames(page);
  expect(reducedAnimations).not.toContain('dashboardStageEnter');
  expect(reducedAnimations).not.toContain('dashboardStripEnter');
});

test('TEST-UIUX-FOCUS-LAYOUT-001 focus stage keeps real desktop width after rail collapse', async ({ page }) => {
  const desktopFocusWidths = [1280, 1366, 1440];

  for (const width of desktopFocusWidths) {
    await page.setViewportSize({ width, height: 900 });
    if (await page.locator('[data-uiux="dashboard-focus"]').count() === 0) {
      await page.locator('[data-dashboard-card="funnel"]').click();
      await expect(page.locator('[data-uiux="dashboard-focus"]')).toBeVisible();
    }

    await expectFocusStageRealWidth(page, 500);
    await expect(page.getByLabel('看板选择器').locator('.sideCard')).toHaveCount(8);
    await expectFocusRailSelection(page, 'funnel');
  }
});

test('TEST-UIUX-NAV-01 focus rail keeps nav reachable and flyout above the stage', async ({ page }) => {
  await page.setViewportSize({ width: 1440, height: 900 });
  await expect(page.locator('[data-uiux="dashboard"]')).toBeVisible();
  await page.locator('[data-dashboard-card="funnel"]').click();
  await expect(page.locator('[data-uiux="dashboard-focus"]')).toBeVisible();
  await expect(page.locator('.shell.focusMode')).toBeVisible();
  await expectFocusStageRealWidth(page, 500);
  await expectFocusNavIconsReachable(page, 14);

  const contactNav = page.getByLabel('主导航').getByRole('button', { name: '联系人' });
  await contactNav.hover();
  await expect.poll(async () => (await navFlyoutState(page, '联系人')).opacity, { timeout: 1_200 }).toBeGreaterThan(0.95);
  const hoverState = await navFlyoutState(page, '联系人');
  expect(hoverState.afterContent).toContain('联系人');
  expect(hoverState.visibility).toBe('visible');
  expect(hoverState.topIsStage).toBe(false);
  expect(hoverState.topIsNavButton).toBe(true);

  const quoteNav = page.getByLabel('主导航').getByRole('button', { name: '报价' });
  await page.keyboard.press('Tab');
  await quoteNav.focus();
  await expect(quoteNav).toBeFocused();
  await expect.poll(async () => (await navFlyoutState(page, '报价')).opacity, { timeout: 500 }).toBeGreaterThan(0.95);
  const focusState = await navFlyoutState(page, '报价');
  expect(focusState.afterContent).toContain('报价');
  expect(focusState.topIsStage).toBe(false);
  expect(focusState.topIsNavButton).toBe(true);

  await page.keyboard.press('Escape');
  await expect(page.locator('[data-uiux="dashboard"]')).toBeVisible();
  await page.setViewportSize({ width: 1440, height: 720 });
  await page.evaluate(() => window.scrollTo(0, 0));
  await page.locator('[data-dashboard-card="funnel"]').click();
  await expect(page.locator('[data-uiux="dashboard-focus"]')).toBeVisible();
  await expectFocusStageRealWidth(page, 500);
  await expectFocusNavIconsReachable(page, 14);
});

test('TEST-UIUX-A5-001 main navigation is keyboard reachable', async ({ page }) => {
  const leadsNav = page.getByLabel('主导航').getByRole('button', { name: '线索' });
  await leadsNav.focus();
  await expect(leadsNav).toBeFocused();
  await page.keyboard.press('Enter');
  await expect(page.getByRole('heading', { name: '线索', exact: true })).toBeVisible();
});

async function signIn(page: import('@playwright/test').Page, email: string, password: string) {
  await page.getByLabel('邮箱').fill(email);
  await page.getByLabel('密码').fill(password);
  await page.getByRole('button', { name: '登录' }).click();
}

async function expectDashboardCardsNotClipped(page: Page) {
  const clippedCards = await page.locator('[data-dashboard-card], .dashboardKpis .metricTile').evaluateAll((cards) => (
    cards
      .map((card) => {
        const element = card as HTMLElement;
        return {
          key: element.dataset.dashboardCard ?? element.textContent?.trim().slice(0, 24) ?? '',
          viewportWidth: window.innerWidth,
          verticalOverflow: element.scrollHeight - element.clientHeight,
          horizontalOverflow: element.scrollWidth - element.clientWidth
        };
      })
      .filter((card) => card.verticalOverflow > 1 || card.horizontalOverflow > 1)
  ));
  expect(clippedCards).toEqual([]);
}

async function expectDashboardInlineContentStable(page: Page) {
  const layoutIssues = await page.evaluate(() => {
    const issues: Array<{ type: string; text: string }> = [];
    document.querySelectorAll<HTMLElement>('.dashboardKpis .metricTile').forEach((tile) => {
      const icon = tile.querySelector<HTMLElement>('.metricIcon');
      const value = tile.querySelector<HTMLElement>('strong');
      if (!icon || !value) return;
      const iconRect = icon.getBoundingClientRect();
      const valueRect = value.getBoundingClientRect();
      const overlapsX = valueRect.right > iconRect.left - 1 && iconRect.right > valueRect.left - 1;
      const overlapsY = valueRect.bottom > iconRect.top - 1 && iconRect.bottom > valueRect.top - 1;
      if (overlapsX && overlapsY) {
        issues.push({ type: 'metric-icon-overlap', text: tile.textContent?.trim() ?? '' });
      }
    });

    document.querySelectorAll<HTMLElement>('[data-dashboard-card] .legendItem span:not(.legendSwatch)').forEach((label) => {
      const style = window.getComputedStyle(label);
      const rect = label.getBoundingClientRect();
      const fontSize = Number.parseFloat(style.fontSize) || 12;
      const lineHeight = Number.parseFloat(style.lineHeight) || fontSize * 1.2;
      if (style.whiteSpace !== 'nowrap' || rect.height > lineHeight * 1.5) {
        issues.push({ type: 'legend-wrap', text: label.textContent?.trim() ?? '' });
      }
    });

    return issues;
  });
  expect(layoutIssues).toEqual([]);
}

async function expectDashboardReadableTextNotTruncated(page: Page) {
  const truncated = await page.evaluate(() => {
    const selectors = [
      '[data-dashboard-card="payments"] .paymentRow small',
      '[data-dashboard-card="payments"] .paymentRight .money',
      '[data-dashboard-card="payments"] .paymentRight .badge',
      '[data-dashboard-card="stage"] .legendItem span:not(.legendSwatch)',
      '[data-dashboard-card="activity"] .event p',
      '[data-dashboard-card="activity"] .event small',
      '[data-dashboard-card="activity"] .eventTime'
    ];

    return selectors.flatMap((selector) => (
      Array.from(document.querySelectorAll<HTMLElement>(selector))
        .map((element) => ({
          selector,
          text: element.textContent?.trim() ?? '',
          viewportWidth: window.innerWidth,
          scrollWidth: element.scrollWidth,
          clientWidth: element.clientWidth
        }))
        .filter((item) => item.text.length > 0 && item.scrollWidth > item.clientWidth + 1)
    ));
  });
  expect(truncated).toEqual([]);
}

async function expectDashboardFlowRowsDoNotOverlap(page: Page) {
  const overlaps = await page.evaluate(() => {
    const intersects = (a: DOMRect, b: DOMRect) => (
      a.right > b.left + 1
      && b.right > a.left + 1
      && a.bottom > b.top + 1
      && b.bottom > a.top + 1
    );

    const describeRect = (rect: DOMRect) => ({
      x: Math.round(rect.x),
      y: Math.round(rect.y),
      width: Math.round(rect.width),
      height: Math.round(rect.height)
    });

    return Array.from(document.querySelectorAll<HTMLElement>(
      '[data-dashboard-card="payments"] .paymentRow, [data-dashboard-card="activity"] .event'
    )).flatMap((row) => {
      const icon = row.querySelector<HTMLElement>('.flowIcon');
      if (!icon) return [];
      const iconRect = icon.getBoundingClientRect();
      return Array.from(row.children)
        .filter((child): child is HTMLElement => child instanceof HTMLElement && !child.classList.contains('flowIcon'))
        .map((element) => ({
          card: row.closest<HTMLElement>('[data-dashboard-card]')?.dataset.dashboardCard ?? '',
          rowClass: row.className,
          text: element.textContent?.trim() ?? '',
          viewportWidth: window.innerWidth,
          gridColumns: window.getComputedStyle(row).gridTemplateColumns,
          icon: describeRect(iconRect),
          textRect: describeRect(element.getBoundingClientRect()),
          overlaps: intersects(iconRect, element.getBoundingClientRect())
        }))
        .filter((item) => item.text.length > 0 && item.overlaps);
    });
  });
  expect(overlaps).toEqual([]);
}

async function installDashboardAnimationRecorder(page: Page) {
  await page.evaluate(() => {
    type DashboardAnimationRecord = { name: string; duration: number | null; delay: number | null; total: number | null };
    const global = window as Window & {
      __dashboardAnimationNames?: string[];
      __dashboardAnimationRecords?: DashboardAnimationRecord[];
      __dashboardAnimationRecorderInstalled?: boolean;
    };
    global.__dashboardAnimationNames = [];
    global.__dashboardAnimationRecords = [];
    if (global.__dashboardAnimationRecorderInstalled) return;
    document.addEventListener('animationstart', (event) => {
      const target = event.target;
      if (!(target instanceof Element)) return;
      const isDashboardFocus = target.matches('[data-uiux="dashboard-focus"]') || Boolean(target.closest('[data-uiux="dashboard-focus"]'));
      if (!isDashboardFocus) return;
      const animationName = (event as AnimationEvent).animationName;
      const timing = target
        .getAnimations()
        .find((animation) => (
          'animationName' in animation && (animation as CSSAnimation).animationName === animationName
        ))
        ?.effect
        ?.getTiming();
      const duration = typeof timing?.duration === 'number' ? timing.duration : null;
      const delay = typeof timing?.delay === 'number' ? timing.delay : null;
      global.__dashboardAnimationNames?.push(animationName);
      global.__dashboardAnimationRecords?.push({
        name: animationName,
        duration,
        delay,
        total: duration !== null && delay !== null ? duration + delay : null
      });
    }, true);
    global.__dashboardAnimationRecorderInstalled = true;
  });
}

async function resetDashboardAnimationRecorder(page: Page) {
  await page.evaluate(() => {
    const global = window as Window & {
      __dashboardAnimationNames?: string[];
      __dashboardAnimationRecords?: Array<{ name: string; duration: number | null; delay: number | null; total: number | null }>;
    };
    global.__dashboardAnimationNames = [];
    global.__dashboardAnimationRecords = [];
  });
}

async function dashboardAnimationNames(page: Page) {
  return page.evaluate(() => (window as Window & { __dashboardAnimationNames?: string[] }).__dashboardAnimationNames ?? []);
}

async function expectDashboardAnimationStarted(page: Page, animationName: string) {
  await expect.poll(() => dashboardAnimationNames(page), { timeout: 2_000 }).toContain(animationName);
}

async function dashboardAnimationRecords(page: Page) {
  return page.evaluate(() => (
    (window as Window & {
      __dashboardAnimationRecords?: Array<{ name: string; duration: number | null; delay: number | null; total: number | null }>;
    }).__dashboardAnimationRecords ?? []
  ));
}

async function expectDashboardAnimationTotal(page: Page, animationName: string, expectedMs: number, toleranceMs = 24) {
  await expect.poll(async () => {
    const record = (await dashboardAnimationRecords(page)).find((item) => item.name === animationName && item.total !== null);
    return record?.total ?? 0;
  }, { timeout: 2_000 }).toBeGreaterThan(0);
  const record = (await dashboardAnimationRecords(page)).find((item) => item.name === animationName && item.total !== null);
  expect(record?.total).toBeGreaterThanOrEqual(expectedMs - toleranceMs);
  expect(record?.total).toBeLessThanOrEqual(expectedMs + toleranceMs);
}

async function focusRailKeys(page: Page) {
  return page.getByLabel('看板选择器').locator('[data-focus-side-card]').evaluateAll((cards) => (
    cards.map((card) => card.getAttribute('data-focus-side-card') ?? '')
  ));
}

async function expectFocusRailSelection(page: Page, selectedKey: string) {
  await expect(page.getByLabel('看板选择器').locator('[aria-current="true"]')).toHaveCount(1);
  await expect(page.getByLabel('看板选择器').locator(`[data-focus-side-card="${selectedKey}"]`)).toHaveAttribute('aria-current', 'true');
}

async function focusStageLayout(page: Page) {
  return page.locator('[data-uiux="dashboard-focus"]').evaluate((rootElement) => {
    const root = rootElement as HTMLElement;
    const focus = root.querySelector<HTMLElement>('.focus');
    const stage = root.querySelector<HTMLElement>('.stage');
    const side = root.querySelector<HTMLElement>('.side');
    const shell = root.closest<HTMLElement>('.shell');
    const workspace = root.closest<HTMLElement>('.workspace');
    const sidebar = shell?.querySelector<HTMLElement>('.sidebar') ?? null;
    if (!focus || !stage || !side || !workspace || !sidebar || !shell) {
      throw new Error('focus stage layout nodes are missing');
    }

    const focusRect = focus.getBoundingClientRect();
    const stageRect = stage.getBoundingClientRect();
    const sideRect = side.getBoundingClientRect();
    const workspaceRect = workspace.getBoundingClientRect();
    const sidebarRect = sidebar.getBoundingClientRect();
    return {
      viewportWidth: window.innerWidth,
      focusColumns: getComputedStyle(focus).gridTemplateColumns,
      shellColumns: getComputedStyle(shell).gridTemplateColumns,
      sidebarPosition: getComputedStyle(sidebar).position,
      sidebarWidth: Math.round(sidebarRect.width),
      workspaceLeft: Math.round(workspaceRect.left),
      workspaceWidth: Math.round(workspaceRect.width),
      focusWidth: Math.round(focusRect.width),
      stageLeft: Math.round(stageRect.left),
      stageWidth: Math.round(stageRect.width),
      sideLeft: Math.round(sideRect.left),
      sideWidth: Math.round(sideRect.width)
    };
  });
}

async function expectFocusStageRealWidth(page: Page, minStageWidth: number) {
  await expect.poll(async () => {
    const layout = await focusStageLayout(page);
    return layout.stageWidth;
  }, { timeout: 1_500 }).toBeGreaterThan(minStageWidth);

  const layout = await focusStageLayout(page);
  expect(layout.stageWidth, JSON.stringify(layout)).toBeGreaterThan(minStageWidth);
  expect(layout.stageWidth / layout.focusWidth, JSON.stringify(layout)).toBeGreaterThan(0.45);
  expect(layout.focusColumns, JSON.stringify(layout)).not.toMatch(/^0px\s+300px$/);
  expect(layout.sidebarPosition, JSON.stringify(layout)).not.toBe('fixed');
  expect(layout.sideWidth, JSON.stringify(layout)).toBeGreaterThanOrEqual(280);
  expect(layout.workspaceLeft, JSON.stringify(layout)).toBeGreaterThanOrEqual(layout.sidebarWidth - 1);
}

async function focusNavReachability(page: Page) {
  return page.getByLabel('主导航').evaluate((navElement) => {
    const nav = navElement as HTMLElement;
    const sidebar = nav.closest('.sidebar') as HTMLElement | null;
    const items = Array.from(nav.querySelectorAll<HTMLElement>('.navItem'));
    const last = items[items.length - 1];
    const navRect = nav.getBoundingClientRect();
    const sidebarRect = sidebar?.getBoundingClientRect();
    const lastRect = last.getBoundingClientRect();
    const navStyle = window.getComputedStyle(nav);
    const sidebarStyle = sidebar ? window.getComputedStyle(sidebar) : null;
    return {
      count: items.length,
      viewportHeight: window.innerHeight,
      scrollY: window.scrollY,
      navTop: navRect.top,
      navBottom: navRect.bottom,
      navClientHeight: nav.clientHeight,
      navScrollHeight: nav.scrollHeight,
      navScrollTop: nav.scrollTop,
      navPosition: navStyle.position,
      navMaxHeight: navStyle.maxHeight,
      sidebarTop: sidebarRect?.top ?? null,
      sidebarBottom: sidebarRect?.bottom ?? null,
      sidebarPosition: sidebarStyle?.position ?? null,
      sidebarHeight: sidebarStyle?.height ?? null,
      lastTop: lastRect.top,
      lastBottom: lastRect.bottom,
      canScroll: nav.scrollHeight > nav.clientHeight + 1,
      lastWithinViewport: lastRect.top >= -1 && lastRect.bottom <= window.innerHeight + 1
    };
  });
}

async function expectFocusNavIconsReachable(page: Page, expectedCount: number) {
  const nav = page.getByLabel('主导航');
  await expect(nav.locator('.navItem')).toHaveCount(expectedCount);
  const lastNav = nav.getByRole('button', { name: '操作日志' });
  await expect(lastNav).toBeVisible();

  let metrics = await focusNavReachability(page);
  if (!metrics.lastWithinViewport && metrics.canScroll) {
    await nav.evaluate((navElement) => {
      navElement.scrollTop = navElement.scrollHeight;
    });
    metrics = await focusNavReachability(page);
  }

  expect(metrics.count).toBe(expectedCount);
  expect(metrics.lastTop, JSON.stringify(metrics)).toBeGreaterThanOrEqual(-1);
  expect(metrics.lastBottom, JSON.stringify(metrics)).toBeLessThanOrEqual(metrics.viewportHeight + 1);
}

async function navFlyoutState(page: Page, label: string) {
  return page.getByLabel('主导航').getByRole('button', { name: label }).evaluate((buttonElement) => {
    const button = buttonElement as HTMLElement;
    const style = window.getComputedStyle(button, '::after');
    const rect = button.getBoundingClientRect();
    const x = Math.min(window.innerWidth - 2, rect.right + 24);
    const y = rect.top + rect.height / 2;
    const topElement = document.elementFromPoint(x, y);
    const topButton = topElement?.closest('button');
    return {
      afterContent: style.content,
      opacity: Number.parseFloat(style.opacity),
      visibility: style.visibility,
      zIndex: style.zIndex,
      pointX: Math.round(x),
      pointY: Math.round(y),
      topIsStage: Boolean(topElement?.closest('.stage')),
      topIsNavButton: topButton === button
    };
  });
}

async function expectSingleFocusExitControl(page: Page) {
  const focusStage = page.getByLabel('聚焦舞台');
  await expect(focusStage.locator('.stageTools').getByRole('button')).toHaveCount(1);
  await expect(focusStage.getByRole('button', { name: '返回', exact: true })).toHaveCount(1);
  await expect(focusStage.getByText(/Esc/)).toHaveCount(0);
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
  }, { email, displayName, password: dashboardPassword, role });
}
