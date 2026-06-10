import type { CSSProperties, KeyboardEvent as ReactKeyboardEvent, ReactNode } from 'react';
import { useEffect, useLayoutEffect, useMemo, useRef, useState } from 'react';
import {
  Activity,
  Bell,
  BriefcaseBusiness,
  CreditCard,
  Expand,
  ListChecks,
  RefreshCcw,
  Target,
  TrendingUp,
  Trophy,
} from 'lucide-react';
import { CurrentUser } from '../api/auth';
import { Contract, listContracts } from '../api/contracts';
import { Lead, listLeads } from '../api/leads';
import { Opportunity, listOpportunities } from '../api/opportunities';
import { listPaymentContracts } from '../api/payments';
import { Quote, listQuotes } from '../api/quotes';
import { BasicReport, GroupRow, ManagerOverview as ManagerOverviewData, PaymentGroupRow, getBasicReport, getManagerOverview } from '../api/reports';
import { ReminderRow, listReminders } from '../api/reminders';
import { Activity as WorkActivity, WorkTask, listActivities, listTasks } from '../api/work';
import {
  Badge,
  Button,
  DataTable,
  EmptyState,
  ErrorState,
  FocusSideCard,
  FocusStage,
  FunnelBars,
  Leaderboard,
  MetricCard,
  SkeletonBlock,
  StageDonut
} from '../components/ui';
import {
  contractStatusLabel,
  labelFor,
  localizeError,
  objectTypeLabel,
  opportunityStageLabel,
  paymentStatusLabel,
  reminderTypeLabel,
  taskStatusLabel
} from '../i18n/labels';

type DashboardCardKey = 'funnel' | 'stage' | 'trend' | 'leaderboard' | 'todo' | 'payments' | 'key-opportunities' | 'activity';
type AccentTone = 'sky' | 'mint' | 'peach' | 'purple';
type DashboardMotionPhase = 'idle' | 'entering' | 'exiting' | 'switching' | 'reduced-entering' | 'reduced-exiting' | 'reduced-switching';
type DashboardRect = { left: number; top: number; width: number; height: number };

const focusEnterMs = 320;
const focusExitMs = 220;
const focusSwitchMs = 220;
const focusReducedMs = 80;

type DashboardSnapshot = {
  overview: ManagerOverviewData | null;
  report: BasicReport | null;
  opportunities: Opportunity[];
  leads: Lead[];
  quotes: Quote[];
  contracts: Contract[];
  paymentContracts: Contract[];
  tasks: WorkTask[];
  reminders: ReminderRow[];
  activities: WorkActivity[];
  errors: string[];
  loadedAt: Date;
  businessDate: string;
};

type StageAggregate = {
  key: string;
  label: string;
  count: number;
  amount: number;
  tone: 'primary' | 'sky' | 'mint' | 'peach' | 'purple' | 'success' | 'danger';
};

type DashboardModel = {
  scopeWord: '团队' | '我的';
  scopeDescription: string;
  scopeBadge: string;
  currency: string;
  stageRows: StageAggregate[];
  funnelRows: StageAggregate[];
  trendPoints: Array<{ label: string; value: number }>;
  leaderboardRows: Array<{ ownerId: string; wonCount: number; amount: number }>;
  todoRows: Array<{ id: string; title: string; meta: string; badge: string; tone: 'primary' | 'warning' | 'danger' }>;
  paymentRows: Array<{ id: string; title: string; meta: string; amount: number; status: string; tone: 'success' | 'warning' | 'danger' | 'primary' }>;
  keyOpportunities: Opportunity[];
  activities: WorkActivity[];
  metrics: {
    monthLeadCount: number;
    activeOpportunityAmount: number;
    monthWonCount: number;
    taskCount: number;
    opportunityCount: number;
    paymentAmount: number;
  };
  live: {
    leadCount: number;
    opportunityCount: number;
    paymentAmount: number;
    taskCount: number;
  };
};

const stageOrder = ['New Opportunity', 'Needs Confirmed', 'Quote', 'Contract Negotiation', 'Won', 'Lost'];
const funnelStageOrder = ['New Opportunity', 'Needs Confirmed', 'Quote', 'Contract Negotiation', 'Won'];
const nonTerminalStages = new Set(['New Opportunity', 'Needs Confirmed', 'Quote', 'Contract Negotiation']);

export function WorkOverview({
  user,
  onFocusChange
}: {
  user: CurrentUser;
  onFocusChange: (enabled: boolean) => void;
}) {
  const [snapshot, setSnapshot] = useState<DashboardSnapshot | null>(null);
  const [loading, setLoading] = useState(true);
  const [activeCard, setActiveCard] = useState<DashboardCardKey | null>(null);
  const [motionPhase, setMotionPhase] = useState<DashboardMotionPhase>('idle');
  const [motionRects, setMotionRects] = useState<Partial<Record<DashboardCardKey, DashboardRect>>>({});
  const [liveMessage, setLiveMessage] = useState('');
  const focusRootRef = useRef<HTMLElement | null>(null);
  const motionTimerRef = useRef<number | null>(null);
  const restoreCardKeyRef = useRef<DashboardCardKey | null>(null);
  const prefersReducedMotion = usePrefersReducedMotion();
  const isManagerView = user.role !== 'Sales';
  const model = useMemo(() => (snapshot ? buildModel(snapshot, user) : null), [snapshot, user]);
  const cards = useMemo(() => (model ? dashboardCards(model, isManagerView) : []), [model, isManagerView]);

  useEffect(() => {
    void refresh();
  }, [user.id, user.role]);

  useEffect(() => {
    onFocusChange(Boolean(activeCard));
    return () => onFocusChange(false);
  }, [activeCard, onFocusChange]);

  useEffect(() => {
    if (!activeCard) return;
    const onKeyDown = (event: KeyboardEvent) => {
      if (event.key === 'Escape') {
        requestExitFocus();
      }
    };
    window.addEventListener('keydown', onKeyDown);
    return () => window.removeEventListener('keydown', onKeyDown);
  }, [activeCard, motionPhase, prefersReducedMotion]);

  useLayoutEffect(() => {
    if (!activeCard || !focusRootRef.current) return;
    const root = focusRootRef.current;
    const stage = root.querySelector<HTMLElement>('.stage');
    const activeRect = motionRects[activeCard];
    if (stage && activeRect) {
      applyMotionOrigin(stage, activeRect, stage, 'hero');
    }

    root.querySelectorAll<HTMLElement>('[data-focus-side-card]').forEach((sideCard) => {
      const key = sideCard.dataset.focusSideCard as DashboardCardKey | undefined;
      if (!key) return;
      const originRect = motionRects[key];
      if (!originRect) return;
      applyMotionOrigin(sideCard, originRect, sideCard, 'strip');
    });
  }, [activeCard, motionPhase, motionRects]);

  useEffect(() => {
    if (!activeCard || motionPhase === 'idle') return;
    if (motionTimerRef.current) {
      window.clearTimeout(motionTimerRef.current);
      motionTimerRef.current = null;
    }

    const isReducedPhase = motionPhase.startsWith('reduced');
    const duration = isReducedPhase
      ? focusReducedMs
      : motionPhase === 'entering'
        ? focusEnterMs
        : motionPhase === 'exiting'
          ? focusExitMs
          : focusSwitchMs;

    motionTimerRef.current = window.setTimeout(() => {
      motionTimerRef.current = null;
      if (motionPhase === 'exiting' || motionPhase === 'reduced-exiting') {
        const restoreKey = restoreCardKeyRef.current;
        setActiveCard(null);
        setMotionPhase('idle');
        setMotionRects({});
        window.requestAnimationFrame(() => {
          if (!restoreKey) return;
          document.querySelector<HTMLElement>(`[data-dashboard-card="${restoreKey}"]`)?.focus();
        });
        return;
      }
      setMotionPhase('idle');
      focusStageHeading();
    }, duration);

    return () => {
      if (motionTimerRef.current) {
        window.clearTimeout(motionTimerRef.current);
        motionTimerRef.current = null;
      }
    };
  }, [activeCard, motionPhase]);

  useEffect(() => () => {
    if (motionTimerRef.current) {
      window.clearTimeout(motionTimerRef.current);
    }
  }, []);

  async function refresh() {
    setLoading(true);
    const errors: string[] = [];
    const businessDate = today();
    const safe = async <T,>(label: string, promise: Promise<T>, fallback: T): Promise<T> => {
      try {
        return await promise;
      } catch (err) {
        errors.push(`${label}：${localizeError(err as { safeMessage?: string }, '加载失败。')}`);
        return fallback;
      }
    };

    const [
      overview,
      report,
      opportunities,
      leads,
      quotes,
      contracts,
      paymentContracts,
      tasks,
      reminders,
      activities
    ] = await Promise.all([
      isManagerView ? safe('团队总览', getManagerOverview(), null) : Promise.resolve(null),
      isManagerView ? safe('基础报表', getBasicReport(), null) : Promise.resolve(null),
      safe('商机', listOpportunities('', '', false).then((response) => response.items), [] as Opportunity[]),
      safe('线索', listLeads('', false).then((response) => response.items), [] as Lead[]),
      safe('报价', listQuotes('').then((response) => response.items), [] as Quote[]),
      safe('合同', listContracts('', false).then((response) => response.items), [] as Contract[]),
      safe('回款合同', listPaymentContracts('', false).then((response) => response.items), [] as Contract[]),
      safe('任务', listTasks({ activeOnly: true, businessDate }).then((response) => response.items), [] as WorkTask[]),
      safe('提醒', listReminders(businessDate).then((response) => response.rows), [] as ReminderRow[]),
      safe('最近活动', listActivities('', '').then((response) => response.items), [] as WorkActivity[])
    ]);

    setSnapshot({
      overview,
      report,
      opportunities,
      leads,
      quotes,
      contracts,
      paymentContracts,
      tasks,
      reminders,
      activities,
      errors,
      loadedAt: new Date(),
      businessDate
    });
    setLoading(false);
  }

  function beginCardFocus(cardKey: DashboardCardKey, sourceElement: HTMLElement) {
    const rects = captureDashboardRects();
    rects[cardKey] = rectFromDom(sourceElement);
    restoreCardKeyRef.current = cardKey;
    setMotionRects(rects);
    setActiveCard(cardKey);
    setMotionPhase(prefersReducedMotion ? 'reduced-entering' : 'entering');
    const targetCard = cards.find((card) => card.key === cardKey);
    setLiveMessage(`已进入${targetCard?.title ?? '卡片'}聚焦视图。`);
  }

  function switchFocus(cardKey: DashboardCardKey) {
    if (cardKey === activeCard) return;
    restoreCardKeyRef.current = cardKey;
    setActiveCard(cardKey);
    setMotionPhase(prefersReducedMotion ? 'reduced-switching' : 'switching');
    const targetCard = cards.find((card) => card.key === cardKey);
    setLiveMessage(`已切换到${targetCard?.title ?? '卡片'}聚焦视图。`);
  }

  function requestExitFocus() {
    if (!activeCard || motionPhase === 'exiting' || motionPhase === 'reduced-exiting') return;
    const targetCard = cards.find((card) => card.key === activeCard);
    restoreCardKeyRef.current = activeCard;
    setMotionPhase(prefersReducedMotion ? 'reduced-exiting' : 'exiting');
    setLiveMessage(`已返回工作台，焦点回到${targetCard?.title ?? '原卡片'}。`);
  }

  function focusStageHeading() {
    window.requestAnimationFrame(() => {
      focusRootRef.current?.querySelector<HTMLElement>('[data-focus-heading]')?.focus();
    });
  }

  if (loading && !snapshot) {
    return (
      <main className="content dashboardPage" data-uiux="dashboard">
        <DashboardHeader isManagerView={isManagerView} loadedAt={null} onRefresh={() => void refresh()} />
        <SkeletonBlock lines={8} label="加载工作台" />
      </main>
    );
  }

  if (!snapshot || !model) {
    return (
      <main className="content dashboardPage" data-uiux="dashboard">
        <DashboardHeader isManagerView={isManagerView} loadedAt={null} onRefresh={() => void refresh()} />
        <ErrorState>无法加载工作台。</ErrorState>
      </main>
    );
  }

  const active = activeCard ? cards.find((card) => card.key === activeCard) : null;

  if (active) {
    const sideCards: FocusSideCard[] = cards
      .filter((card) => card.key !== active.key)
      .map((card, index) => ({
        key: card.key,
        title: card.title,
        metric: card.metric,
        meta: card.meta,
        icon: card.sideIcon,
        motionIndex: index,
        onSelect: () => switchFocus(card.key)
      }));
    const focusClassName = [
      'content',
      'dashboardFocusPage',
      focusTransitionClass(motionPhase),
      prefersReducedMotion ? 'dashboardMotionReduced' : 'dashboardMotionFull'
    ].filter(Boolean).join(' ');

    return (
      <main
        className={focusClassName}
        data-motion-mode={prefersReducedMotion ? 'reduced' : 'full'}
        data-transition-phase={motionPhase}
        data-uiux="dashboard-focus"
        ref={focusRootRef}
        style={focusRootStyle(motionPhase)}
      >
        <FocusStage
          title={active.title}
          subtitle={<><span className="liveDot livePulse" aria-hidden="true" />实时 · {active.focusSubtitle}</>}
          icon={active.sideIcon}
          sideCards={sideCards}
          onBack={requestExitFocus}
          backLabel="返回"
          escapeHint="Esc 返回"
          tools={<Badge tone="primary">{model.scopeBadge}</Badge>}
        >
          <div className="dashboardStageContent" key={active.key}>
            {active.focus}
          </div>
        </FocusStage>
        <span className="srOnly" role="status" aria-live="polite">{liveMessage}</span>
      </main>
    );
  }

  return (
    <main className="content dashboardPage motionFadeIn" data-uiux="dashboard">
      <DashboardHeader isManagerView={isManagerView} loadedAt={snapshot.loadedAt} onRefresh={() => void refresh()} />
      {snapshot.errors.length ? <ErrorState className="dashboardError">{snapshot.errors.join('；')}</ErrorState> : null}
      <LiveReport model={model} />
      <section className="dashboardKpis" aria-label={isManagerView ? '团队关键指标' : '个人关键指标'}>
        <MetricCard label={`${model.scopeWord}本月新增线索`} value={model.metrics.monthLeadCount} icon={<ListChecks size={18} aria-hidden="true" />} delta="授权记录" />
        <MetricCard label={`${model.scopeWord}进行中商机金额`} value={formatMoney(model.metrics.activeOpportunityAmount, model.currency)} tone="sky" icon={<BriefcaseBusiness size={18} aria-hidden="true" />} delta="非终态阶段" />
        <MetricCard label={`${model.scopeWord}本月赢单`} value={model.metrics.monthWonCount} tone="mint" icon={<Trophy size={18} aria-hidden="true" />} delta="已关闭赢单" />
        <MetricCard label={`${model.scopeWord}待办任务`} value={model.metrics.taskCount} tone="peach" icon={<Bell size={18} aria-hidden="true" />} delta={`${model.todoRows.filter((row) => row.tone === 'danger').length} 条预警`} />
      </section>
      <section className={isManagerView ? 'roleGrid managerDashboardGrid' : 'roleGrid salesDashboardGrid'} aria-label={isManagerView ? '管端工作台数据卡' : '销售工作台数据卡'}>
        {cards.map((card) => (
          <DashboardPanel card={card} key={card.key} onFocus={(element) => beginCardFocus(card.key, element)} />
        ))}
      </section>
      <span className="srOnly" role="status" aria-live="polite">{liveMessage}</span>
    </main>
  );
}

function captureDashboardRects(): Partial<Record<DashboardCardKey, DashboardRect>> {
  const rects: Partial<Record<DashboardCardKey, DashboardRect>> = {};
  document.querySelectorAll<HTMLElement>('[data-dashboard-card]').forEach((element) => {
    const key = element.dataset.dashboardCard as DashboardCardKey | undefined;
    if (!key) return;
    rects[key] = rectFromDom(element);
  });
  return rects;
}

function rectFromDom(element: HTMLElement): DashboardRect {
  const rect = element.getBoundingClientRect();
  return {
    left: rect.left,
    top: rect.top,
    width: rect.width,
    height: rect.height
  };
}

function applyMotionOrigin(element: HTMLElement, origin: DashboardRect, target: HTMLElement, prefix: 'hero' | 'strip') {
  const targetRect = target.getBoundingClientRect();
  const safeWidth = Math.max(1, targetRect.width);
  const safeHeight = Math.max(1, targetRect.height);
  element.style.setProperty(`--${prefix}-start-x`, `${origin.left - targetRect.left}px`);
  element.style.setProperty(`--${prefix}-start-y`, `${origin.top - targetRect.top}px`);
  element.style.setProperty(`--${prefix}-start-scale-x`, `${Math.max(0.08, origin.width / safeWidth).toFixed(4)}`);
  element.style.setProperty(`--${prefix}-start-scale-y`, `${Math.max(0.08, origin.height / safeHeight).toFixed(4)}`);
}

function focusTransitionClass(phase: DashboardMotionPhase) {
  if (phase === 'entering' || phase === 'reduced-entering') return 'dashboardTransitionEntering';
  if (phase === 'exiting' || phase === 'reduced-exiting') return 'dashboardTransitionExiting';
  if (phase === 'switching' || phase === 'reduced-switching') return 'dashboardTransitionSwitching';
  return 'dashboardTransitionIdle';
}

function focusRootStyle(phase: DashboardMotionPhase): CSSProperties {
  return {
    '--dashboard-transition-reverse': phase === 'exiting' || phase === 'reduced-exiting' ? 'reverse' : 'normal'
  } as CSSProperties;
}

function usePrefersReducedMotion() {
  const [reduced, setReduced] = useState(false);

  useEffect(() => {
    const query = window.matchMedia('(prefers-reduced-motion: reduce)');
    const sync = () => setReduced(query.matches);
    sync();
    query.addEventListener('change', sync);
    return () => query.removeEventListener('change', sync);
  }, []);

  return reduced;
}

function DashboardHeader({
  isManagerView,
  loadedAt,
  onRefresh
}: {
  isManagerView: boolean;
  loadedAt: Date | null;
  onRefresh: () => void;
}) {
  return (
    <section className="dashboardTitle">
      <div>
        <h1>{isManagerView ? '团队工作台' : '我的工作台'}</h1>
        <p>{currentMonthLabel()} · {isManagerView ? '团队销售运营总览' : '我的销售工作台'} · 数据更新时间 {loadedAt ? formatClock(loadedAt) : '加载中'}</p>
      </div>
      <div className="dashboardActions">
        <Button>本月</Button>
        <Button variant="primary" onClick={onRefresh}>
          <RefreshCcw size={16} aria-hidden="true" />
          刷新数据
        </Button>
        <div className="liveSegment" aria-label="实时更新状态">
          <span className="segmentOn"><span className="liveDot livePulse" aria-hidden="true" />实时更新</span>
          <span className="segmentOff">自动合并</span>
        </div>
        <span className="updateMeta">更新于 {loadedAt ? formatClock(loadedAt) : '加载中'}</span>
      </div>
    </section>
  );
}

function LiveReport({ model }: { model: DashboardModel }) {
  return (
    <section className="liveReport" aria-label="今日实时战报">
      <div className="reportLead">
        <span className="liveDot livePulse" aria-hidden="true" />
        今日实时战报
      </div>
      <ReportItem label="今日活跃线索" value={model.live.leadCount} tone="lavender" />
      <ReportItem label="今日推进商机" value={model.live.opportunityCount} tone="mint" />
      <ReportItem label="回款金额" value={formatMoney(model.live.paymentAmount, model.currency)} tone="peach" />
      <ReportItem label="今日待办" value={model.live.taskCount} tone="purple" />
    </section>
  );
}

function ReportItem({ label, value, tone }: { label: string; value: ReactNode; tone: 'lavender' | 'mint' | 'peach' | 'purple' }) {
  return (
    <div className="reportItem">
      <i className={`reportDot tone-${tone}`} aria-hidden="true" />
      <span>{label}</span>
      <strong>{value}</strong>
    </div>
  );
}

type DashboardCard = {
  key: DashboardCardKey;
  title: string;
  meta: ReactNode;
  metric: ReactNode;
  sideIcon: ReactNode;
  focusSubtitle: ReactNode;
  children: ReactNode;
  focus: ReactNode;
  tone?: AccentTone;
};

function dashboardCards(model: DashboardModel, isManagerView: boolean): DashboardCard[] {
  const scope = model.scopeWord;
  const base: DashboardCard[] = [
    {
      key: 'funnel',
      title: `${scope}销售漏斗`,
      meta: '数量 · 金额',
      metric: model.metrics.opportunityCount,
      sideIcon: <TrendingUp size={18} aria-hidden="true" />,
      focusSubtitle: '数量 · 金额 · 重点商机',
      tone: 'sky',
      children: <FunnelOverview model={model} />,
      focus: <FunnelFocus model={model} />
    },
    {
      key: 'stage',
      title: `${scope}商机阶段构成`,
      meta: `${scope} ${model.metrics.opportunityCount} 个商机`,
      metric: `${model.metrics.opportunityCount}个`,
      sideIcon: <Target size={18} aria-hidden="true" />,
      focusSubtitle: '阶段占比 · 金额',
      tone: 'purple',
      children: <StageOverview model={model} />,
      focus: <StageFocus model={model} />
    },
    {
      key: 'trend',
      title: `${scope}赢单金额趋势`,
      meta: '近 6 月',
      metric: formatMoney(model.trendPoints.reduce((sum, point) => sum + point.value, 0), model.currency),
      sideIcon: <Activity size={18} aria-hidden="true" />,
      focusSubtitle: '近 6 月 · CNY',
      children: <TrendOverview model={model} />,
      focus: <TrendFocus model={model} />
    },
    {
      key: 'todo',
      title: `${scope}待办与预警`,
      meta: `今天 ${model.metrics.taskCount} 项`,
      metric: `${model.metrics.taskCount}`,
      sideIcon: <Bell size={18} aria-hidden="true" />,
      focusSubtitle: '任务 · 提醒 · 预警',
      tone: 'peach',
      children: <TodoOverview model={model} />,
      focus: <TodoFocus model={model} />
    },
    {
      key: 'payments',
      title: `${scope}回款到账`,
      meta: <><span className="liveDot" aria-hidden="true" />实时</>,
      metric: formatMoney(model.metrics.paymentAmount, model.currency),
      sideIcon: <CreditCard size={18} aria-hidden="true" />,
      focusSubtitle: '回款状态 · 授权记录',
      tone: 'peach',
      children: <PaymentsOverview model={model} />,
      focus: <PaymentsFocus model={model} />
    },
    {
      key: 'activity',
      title: `${scope}最近活动`,
      meta: <><span className="liveDot" aria-hidden="true" />更新于 {formatClock(new Date())}</>,
      metric: model.activities.length ? formatTime(model.activities[0].occurredAt) : '暂无',
      sideIcon: <ListChecks size={18} aria-hidden="true" />,
      focusSubtitle: '动态 · 最近更新',
      children: <ActivityOverview model={model} />,
      focus: <ActivityFocus model={model} />
    }
  ];

  if (!isManagerView) {
    return [base[0], base[3], base[1], base[2], base[4], base[5]];
  }

  return [
    base[0],
    base[1],
    base[2],
    {
      key: 'leaderboard',
      title: '销售业绩榜',
      meta: '赢单数 · 金额',
      metric: model.leaderboardRows[0]?.ownerId ?? '暂无',
      sideIcon: <Trophy size={18} aria-hidden="true" />,
      focusSubtitle: '负责人排行 · 赢单金额',
      tone: 'mint',
      children: <LeaderboardOverview model={model} />,
      focus: <LeaderboardFocus model={model} />
    },
    base[3],
    base[4],
    {
      key: 'key-opportunities',
      title: '重点商机',
      meta: '团队',
      metric: model.keyOpportunities.length,
      sideIcon: <BriefcaseBusiness size={18} aria-hidden="true" />,
      focusSubtitle: '高金额 · 临近预计关闭',
      tone: 'sky',
      children: <KeyOpportunityOverview model={model} />,
      focus: <KeyOpportunityFocus model={model} />
    },
    base[5]
  ];
}

function DashboardPanel({ card, onFocus }: { card: DashboardCard; onFocus: (element: HTMLElement) => void }) {
  const handleKeyDown = (event: ReactKeyboardEvent<HTMLElement>) => {
    if (event.key !== 'Enter' && event.key !== ' ') return;
    event.preventDefault();
    onFocus(event.currentTarget);
  };

  return (
    <article
      className="card dashboardPanel"
      data-dashboard-card={card.key}
      aria-label={`展开${card.title}`}
      role="button"
      tabIndex={0}
      onClick={(event) => onFocus(event.currentTarget)}
      onKeyDown={handleKeyDown}
    >
      <div className="dashboardPanelHeader">
        <div className="titleGroup">
          <span className={`panelIcon ${card.tone ?? ''}`}>{card.sideIcon}</span>
          <div>
            <div className="panelTitle">{card.title}</div>
            <div className="panelMeta">{card.meta}</div>
          </div>
        </div>
        <span className="expand" aria-hidden="true">
          <Expand size={16} aria-hidden="true" />
        </span>
      </div>
      {card.children}
    </article>
  );
}

function FunnelOverview({ model }: { model: DashboardModel }) {
  return (
    <FunnelBars
      rows={model.funnelRows.map((row) => ({
        label: row.label,
        value: row.count,
        suffix: <> · {formatMoney(row.amount, model.currency)}</>
      }))}
    />
  );
}

function FunnelFocus({ model }: { model: DashboardModel }) {
  const max = Math.max(1, ...model.funnelRows.map((row) => row.count));
  return (
    <div className="dashboardFocusStack">
      <div className="focusFunnelRows">
        {model.funnelRows.map((row) => (
          <div className="focusFunnelRow" key={row.key}>
            <div className="funnelLabel">{row.label}</div>
            <div className="funnelTrack" aria-hidden="true">
              <span className="funnelFill" style={{ width: `${Math.max(4, Math.round((row.count / max) * 100))}%` }} />
            </div>
            <div className="funnelValue">
              {row.count} · {formatMoney(row.amount, model.currency)}
              <span className="rate">转化率 {max > 0 ? Math.round((row.count / max) * 100) : 0}%</span>
            </div>
          </div>
        ))}
      </div>
      <OpportunityTable rows={model.keyOpportunities} caption="聚焦商机明细" />
    </div>
  );
}

function StageOverview({ model }: { model: DashboardModel }) {
  return (
    <StageDonut
      label="商机阶段构成"
      center={model.metrics.opportunityCount}
      segments={model.stageRows.map((row) => ({ label: row.label, value: row.count, tone: row.tone }))}
    />
  );
}

function StageFocus({ model }: { model: DashboardModel }) {
  return (
    <div className="dashboardFocusStack twoColumnFocus">
      <StageOverview model={model} />
      <DataTable
        caption="阶段构成明细"
        rows={model.stageRows}
        rowKey={(row) => row.key}
        columns={[
          { key: 'stage', header: '阶段', render: (row) => <Badge tone={stageBadgeTone(row.key)}>{row.label}</Badge> },
          { key: 'count', header: '数量', align: 'right', render: (row) => row.count },
          { key: 'amount', header: '金额', align: 'right', render: (row) => <span className="money">{formatMoney(row.amount, model.currency)}</span> }
        ]}
      />
    </div>
  );
}

function TrendOverview({ model }: { model: DashboardModel }) {
  return <TrendVisual points={model.trendPoints} currency={model.currency} compact />;
}

function TrendFocus({ model }: { model: DashboardModel }) {
  return (
    <div className="dashboardFocusStack">
      <TrendVisual points={model.trendPoints} currency={model.currency} />
      <DataTable
        caption="赢单趋势明细"
        rows={model.trendPoints}
        rowKey={(row) => row.label}
        columns={[
          { key: 'month', header: '月份', render: (row) => row.label },
          { key: 'amount', header: '赢单金额', align: 'right', render: (row) => <span className="money">{formatMoney(row.value, model.currency)}</span> }
        ]}
      />
    </div>
  );
}

function TrendVisual({ points, currency, compact = false }: { points: Array<{ label: string; value: number }>; currency: string; compact?: boolean }) {
  const path = trendPath(points);
  return (
    <div className={compact ? 'dashboardTrendChart compact' : 'dashboardTrendChart'} role="img" aria-label="赢单金额趋势">
      <svg viewBox="0 0 520 184" aria-hidden="true" focusable="false">
        <defs>
          <linearGradient id={compact ? 'dashboardTrendCompact' : 'dashboardTrendFocus'} x1="0" x2="0" y1="0" y2="1">
            <stop offset="0" stopColor="#2563EB" stopOpacity=".18" />
            <stop offset="1" stopColor="#2563EB" stopOpacity="0" />
          </linearGradient>
        </defs>
        <g className="trendGrid">
          <line x1="42" y1="30" x2="492" y2="30" />
          <line x1="42" y1="70" x2="492" y2="70" />
          <line x1="42" y1="110" x2="492" y2="110" />
          <line x1="42" y1="150" x2="492" y2="150" />
        </g>
        <path d={path.area} fill={`url(#${compact ? 'dashboardTrendCompact' : 'dashboardTrendFocus'})`} />
        <path className="trendLine" d={path.line} />
      </svg>
      <div className="trendLegend">
        {points.map((point) => (
          <span key={point.label}>
            <small>{point.label}</small>
            <strong>{formatCompactMoney(point.value, currency)}</strong>
          </span>
        ))}
      </div>
    </div>
  );
}

function LeaderboardOverview({ model }: { model: DashboardModel }) {
  if (model.leaderboardRows.length === 0) return <EmptyState>暂无赢单排行。</EmptyState>;
  return (
    <Leaderboard
      items={model.leaderboardRows.slice(0, 4).map((row) => ({
        label: row.ownerId,
        value: row.wonCount,
        meta: formatMoney(row.amount, model.currency),
        suffix: ' 单',
        tone: 'mint'
      }))}
    />
  );
}

function LeaderboardFocus({ model }: { model: DashboardModel }) {
  return (
    <DataTable
      caption="销售业绩榜明细"
      rows={model.leaderboardRows}
      rowKey={(row) => row.ownerId}
      empty="暂无赢单排行。"
      columns={[
        { key: 'owner', header: '负责人', render: (row) => row.ownerId },
        { key: 'won', header: '赢单数', align: 'right', render: (row) => row.wonCount },
        { key: 'amount', header: '赢单金额', align: 'right', render: (row) => <span className="money">{formatMoney(row.amount, model.currency)}</span> }
      ]}
    />
  );
}

function TodoOverview({ model }: { model: DashboardModel }) {
  if (model.todoRows.length === 0) return <EmptyState>暂无待办与预警。</EmptyState>;
  return (
    <div className="dashboardList">
      {model.todoRows.slice(0, 4).map((row) => (
        <DashboardListRow key={row.id} title={row.title} meta={row.meta} badges={<Badge tone={row.tone}>{row.badge}</Badge>} />
      ))}
    </div>
  );
}

function TodoFocus({ model }: { model: DashboardModel }) {
  return (
    <DataTable
      caption="待办与预警明细"
      rows={model.todoRows}
      rowKey={(row) => row.id}
      empty="暂无待办与预警。"
      columns={[
        { key: 'title', header: '事项', render: (row) => row.title },
        { key: 'meta', header: '范围', render: (row) => row.meta },
        { key: 'badge', header: '状态', render: (row) => <Badge tone={row.tone}>{row.badge}</Badge> }
      ]}
    />
  );
}

function PaymentsOverview({ model }: { model: DashboardModel }) {
  if (model.paymentRows.length === 0) return <EmptyState>暂无回款记录。</EmptyState>;
  return (
    <>
      <div className="paymentRows">
        {model.paymentRows.slice(0, 3).map((row, index) => (
          <div className={index === 0 ? 'paymentRow arrived' : 'paymentRow'} key={row.id}>
            <span className="flowIcon pay" aria-hidden="true"><CreditCard size={13} /></span>
            <div>
              <strong>{row.title}</strong>
              <small>{row.meta}</small>
            </div>
            <div className="paymentRight">
              <span className="money">{formatMoney(row.amount, model.currency)}</span>
              <Badge tone={row.tone}>{row.status}</Badge>
            </div>
          </div>
        ))}
      </div>
      <div className="footer"><span className="paymentSummary">{model.scopeWord}回款金额 <span className="money">{formatMoney(model.metrics.paymentAmount, model.currency)}</span></span></div>
    </>
  );
}

function PaymentsFocus({ model }: { model: DashboardModel }) {
  return (
    <DataTable
      caption="回款到账明细"
      rows={model.paymentRows}
      rowKey={(row) => row.id}
      empty="暂无回款记录。"
      columns={[
        { key: 'record', header: '记录', render: (row) => row.title },
        { key: 'meta', header: '说明', render: (row) => row.meta },
        { key: 'status', header: '状态', render: (row) => <Badge tone={row.tone}>{row.status}</Badge> },
        { key: 'amount', header: '金额', align: 'right', render: (row) => <span className="money">{formatMoney(row.amount, model.currency)}</span> }
      ]}
    />
  );
}

function KeyOpportunityOverview({ model }: { model: DashboardModel }) {
  if (model.keyOpportunities.length === 0) return <EmptyState>暂无重点商机。</EmptyState>;
  return (
    <div className="dashboardList">
      {model.keyOpportunities.slice(0, 4).map((opportunity) => (
        <DashboardListRow
          key={opportunity.id}
          title={opportunity.title}
          meta={`${opportunity.customerId} · ${opportunity.ownerId}`}
          badges={<><Badge tone={stageBadgeTone(opportunity.stage)}>{labelFor(opportunityStageLabel, opportunity.stage)}</Badge><Badge>{formatMoney(opportunity.expectedAmount, 'CNY')}</Badge></>}
        />
      ))}
    </div>
  );
}

function KeyOpportunityFocus({ model }: { model: DashboardModel }) {
  return <OpportunityTable rows={model.keyOpportunities} caption="重点商机明细" />;
}

function ActivityOverview({ model }: { model: DashboardModel }) {
  if (model.activities.length === 0) return <EmptyState>暂无最近活动。</EmptyState>;
  return (
    <div className="timeline">
      {model.activities.slice(0, 4).map((activity, index) => (
        <div className={index === 0 ? 'event arrived' : 'event'} key={activity.id}>
          <span className="flowIcon" aria-hidden="true"><Activity size={13} /></span>
          <div>
            <p>{activity.content || '动态已记录'}</p>
            <small><Badge>{labelFor(objectTypeLabel, activity.relatedType)}</Badge> {activity.ownerId}</small>
          </div>
          <span className="eventTime">{formatTime(activity.occurredAt)}</span>
        </div>
      ))}
    </div>
  );
}

function ActivityFocus({ model }: { model: DashboardModel }) {
  return (
    <DataTable
      caption="最近活动明细"
      rows={model.activities}
      rowKey={(row) => row.id}
      empty="暂无最近活动。"
      columns={[
        { key: 'content', header: '活动', render: (row) => row.content || '动态已记录' },
        { key: 'related', header: '对象', render: (row) => `${labelFor(objectTypeLabel, row.relatedType)} ${row.relatedId}` },
        { key: 'owner', header: '负责人', render: (row) => row.ownerId },
        { key: 'time', header: '时间', align: 'right', render: (row) => formatDate(row.occurredAt) }
      ]}
    />
  );
}

function DashboardListRow({ title, meta, badges }: { title: ReactNode; meta: ReactNode; badges?: ReactNode }) {
  return (
    <div className="dashboardRow">
      <div>
        <strong>{title}</strong>
        <span>{meta}</span>
      </div>
      {badges ? <div className="badges">{badges}</div> : null}
    </div>
  );
}

function OpportunityTable({ rows, caption }: { rows: Opportunity[]; caption: string }) {
  return (
    <DataTable
      caption={caption}
      rows={rows}
      rowKey={(row) => row.id}
      empty="暂无商机明细。"
      columns={[
        { key: 'title', header: '商机', render: (row) => row.title },
        { key: 'customer', header: '客户', render: (row) => row.customerId },
        { key: 'stage', header: '阶段', render: (row) => <Badge tone={stageBadgeTone(row.stage)}>{labelFor(opportunityStageLabel, row.stage)}</Badge> },
        { key: 'owner', header: '负责人', render: (row) => row.ownerId },
        { key: 'date', header: '预计签约', render: (row) => formatDate(row.expectedCloseDate) },
        { key: 'amount', header: '金额', align: 'right', render: (row) => <span className="money">{formatMoney(row.expectedAmount, 'CNY')}</span> }
      ]}
    />
  );
}

function buildModel(snapshot: DashboardSnapshot, user: CurrentUser): DashboardModel {
  const scopeWord = user.role === 'Sales' ? '我的' : '团队';
  const scopeDescription = user.role === 'Administrator' ? '全局 CRM 范围' : user.role === 'Sales Manager' ? '团队范围' : '本人负责和分配范围';
  const scopeBadge = user.role === 'Administrator' ? '全部' : user.role === 'Sales Manager' ? '团队' : '本人';
  const currency = snapshot.overview?.currency ?? snapshot.report?.currency ?? 'CNY';
  const opportunities = snapshot.opportunities.filter((item) => !item.archived);
  const leads = snapshot.leads.filter((item) => !item.archived);
  const contracts = snapshot.contracts.filter((item) => !item.archived);
  const paymentContracts = snapshot.paymentContracts.filter((item) => !item.archived);
  const activeTasks = snapshot.tasks.filter((task) => task.status === 'Open');
  const stageRows = stageOrder.map((stage) => stageAggregate(stage, opportunities, snapshot.report?.breakdowns.opportunitiesByStage));
  const funnelRows = funnelStageOrder.map((stage) => stageAggregate(stage, opportunities, snapshot.overview?.pipeline));
  const wonOpportunities = opportunities.filter((opportunity) => opportunity.stage === 'Won');
  const trendPoints = lastSixMonths().map((month) => ({
    label: month.label,
    value: wonOpportunities
      .filter((opportunity) => month.key === monthKey(opportunity.closeDate || opportunity.expectedCloseDate || opportunity.updatedAt))
      .reduce((sum, opportunity) => sum + numberValue(opportunity.expectedAmount), 0)
  }));
  const leaderboardRows = Array.from(groupWonByOwner(wonOpportunities).values()).sort((a, b) => b.amount - a.amount || b.wonCount - a.wonCount);
  const todoRows = buildTodoRows(snapshot.reminders, activeTasks);
  const paymentRows = buildPaymentRows(snapshot.report, paymentContracts, contracts, currency, user.role === 'Sales');
  const keyOpportunities = opportunities
    .filter((opportunity) => !opportunity.archived)
    .sort((left, right) => Number(nonTerminalStages.has(right.stage)) - Number(nonTerminalStages.has(left.stage)) || numberValue(right.expectedAmount) - numberValue(left.expectedAmount))
    .slice(0, 8);
  const activeOpportunityAmount = opportunities.filter((opportunity) => nonTerminalStages.has(opportunity.stage)).reduce((sum, opportunity) => sum + numberValue(opportunity.expectedAmount), 0);
  const monthLeadCount = leads.filter((lead) => isCurrentMonth(lead.updatedAt)).length;
  const monthWonCount = wonOpportunities.filter((opportunity) => isCurrentMonth(opportunity.closeDate || opportunity.updatedAt)).length;
  const paymentAmount = paymentRows.reduce((sum, row) => sum + row.amount, 0);

  return {
    scopeWord,
    scopeDescription,
    scopeBadge,
    currency,
    stageRows,
    funnelRows,
    trendPoints,
    leaderboardRows,
    todoRows,
    paymentRows,
    keyOpportunities,
    activities: snapshot.activities.slice(0, 12),
    metrics: {
      monthLeadCount,
      activeOpportunityAmount,
      monthWonCount,
      taskCount: activeTasks.length,
      opportunityCount: opportunities.length,
      paymentAmount
    },
    live: {
      leadCount: leads.filter((lead) => isToday(lead.updatedAt)).length,
      opportunityCount: opportunities.filter((opportunity) => isToday(opportunity.updatedAt)).length,
      paymentAmount,
      taskCount: activeTasks.filter((task) => task.dueDate <= snapshot.businessDate).length
    }
  };
}

function stageAggregate(stage: string, opportunities: Opportunity[], reportRows?: GroupRow[]): StageAggregate {
  const reportRow = reportRows?.find((row) => row.label === stage || row.key === stage);
  const fallbackRows = opportunities.filter((opportunity) => opportunity.stage === stage);
  return {
    key: stage,
    label: labelFor(opportunityStageLabel, stage),
    count: reportRow?.count ?? fallbackRows.length,
    amount: reportRow ? numberValue(reportRow.amount) : fallbackRows.reduce((sum, opportunity) => sum + numberValue(opportunity.expectedAmount), 0),
    tone: stageTone(stage)
  };
}

function groupWonByOwner(opportunities: Opportunity[]) {
  const groups = new Map<string, { ownerId: string; wonCount: number; amount: number }>();
  for (const opportunity of opportunities) {
    const current = groups.get(opportunity.ownerId) ?? { ownerId: opportunity.ownerId || '未分配', wonCount: 0, amount: 0 };
    current.wonCount += 1;
    current.amount += numberValue(opportunity.expectedAmount);
    groups.set(opportunity.ownerId, current);
  }
  return groups;
}

function buildTodoRows(reminders: ReminderRow[], tasks: WorkTask[]) {
  const reminderRows = reminders.map((row) => ({
    id: row.id,
    title: row.relatedRecord.display,
    meta: `${labelFor(reminderTypeLabel, row.type)} · ${row.dueDate}`,
    badge: labelFor(reminderTypeLabel, row.type),
    tone: row.type.includes('overdue') ? 'danger' as const : 'warning' as const
  }));
  const taskRows = tasks.map((task) => ({
    id: task.id,
    title: task.title,
    meta: `${labelFor(objectTypeLabel, task.relatedType)} ${task.relatedId} · ${task.dueDate}`,
    badge: labelFor(taskStatusLabel, task.status),
    tone: task.dueDate < today() ? 'danger' as const : 'primary' as const
  }));
  return [...reminderRows, ...taskRows].slice(0, 12);
}

function buildPaymentRows(report: BasicReport | null, paymentContracts: Contract[], contracts: Contract[], currency: string, salesScope: boolean) {
  if (report?.breakdowns.paymentsByStatus?.length) {
    return report.breakdowns.paymentsByStatus.map((row: PaymentGroupRow) => ({
      id: row.key,
      title: labelFor(paymentStatusLabel, row.label),
      meta: `${row.count} 条 · 应收 ${formatMoney(row.dueAmount, currency)}`,
      amount: numberValue(row.paidAmount || row.amount),
      status: labelFor(paymentStatusLabel, row.label),
      tone: paymentTone(row.label)
    }));
  }

  const source = paymentContracts.length ? paymentContracts : contracts;
  return source.slice(0, 12).map((contract) => ({
    id: contract.id,
    title: contract.opportunityId,
    meta: `${salesScope ? '可见合同/回款记录' : '合同记录'} · ${labelFor(contractStatusLabel, contract.status)}`,
    amount: numberValue(contract.amount),
    status: labelFor(contractStatusLabel, contract.status),
    tone: contract.status === 'Completed' || contract.status === 'Signed' || contract.status === 'Active' ? 'success' as const : 'warning' as const
  }));
}

function stageTone(stage: string): StageAggregate['tone'] {
  if (stage === 'Won') return 'mint';
  if (stage === 'Lost') return 'danger';
  if (stage === 'Contract Negotiation') return 'peach';
  if (stage === 'Quote') return 'purple';
  if (stage === 'Needs Confirmed') return 'sky';
  return 'primary';
}

function stageBadgeTone(stage: string): 'neutral' | 'primary' | 'success' | 'warning' | 'danger' {
  if (stage === 'Won') return 'success';
  if (stage === 'Lost') return 'danger';
  if (stage === 'Contract Negotiation') return 'warning';
  return 'primary';
}

function paymentTone(status: string): 'success' | 'warning' | 'danger' | 'primary' {
  if (status === 'Paid' || status === 'PartiallyPaid') return 'success';
  if (status === 'Overdue') return 'danger';
  if (status === 'No plan') return 'primary';
  return 'warning';
}

function trendPath(points: Array<{ label: string; value: number }>) {
  if (points.length === 0) return { line: '', area: '' };
  const values = points.map((point) => point.value);
  const min = Math.min(...values);
  const max = Math.max(...values);
  const span = Math.max(1, max - min);
  const coords = points.map((point, index) => {
    const x = points.length === 1 ? 260 : 42 + (index * 450) / (points.length - 1);
    const y = 150 - ((point.value - min) / span) * 112;
    return [x, y] as const;
  });
  const line = coords.map(([x, y], index) => `${index === 0 ? 'M' : 'L'} ${x.toFixed(2)} ${y.toFixed(2)}`).join(' ');
  const area = `${line} L ${coords[coords.length - 1][0].toFixed(2)} 166 L ${coords[0][0].toFixed(2)} 166 Z`;
  return { line, area };
}

function lastSixMonths() {
  const now = new Date();
  return Array.from({ length: 6 }, (_, index) => {
    const date = new Date(now.getFullYear(), now.getMonth() - (5 - index), 1);
    return {
      key: `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}`,
      label: `${date.getMonth() + 1}月`
    };
  });
}

function monthKey(value: string) {
  if (!value) return '';
  return value.slice(0, 7);
}

function isCurrentMonth(value: string) {
  return monthKey(value) === monthKey(new Date().toISOString());
}

function isToday(value: string) {
  return Boolean(value) && value.slice(0, 10) === today();
}

function today() {
  return new Date().toISOString().slice(0, 10);
}

function currentMonthLabel() {
  const date = new Date();
  return `${date.getFullYear()} 年 ${date.getMonth() + 1} 月`;
}

function numberValue(value: string | number | undefined) {
  const number = typeof value === 'number' ? value : Number(value ?? 0);
  return Number.isFinite(number) ? number : 0;
}

function formatMoney(value: string | number, currency: string) {
  const number = numberValue(value);
  return new Intl.NumberFormat('zh-CN', { style: 'currency', currency, maximumFractionDigits: 0 }).format(number);
}

function formatCompactMoney(value: string | number, currency: string) {
  const number = numberValue(value);
  if (Math.abs(number) >= 1_000_000) return `${formatMoney(number / 1_000_000, currency)}M`;
  if (Math.abs(number) >= 1_000) return `${formatMoney(number / 1_000, currency)}K`;
  return formatMoney(number, currency);
}

function formatDate(value: string) {
  if (!value) return '未设置';
  return value.length > 10 ? value.slice(0, 10) : value;
}

function formatTime(value: string) {
  if (!value) return '暂无';
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return formatDate(value);
  return date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' });
}

function formatClock(value: Date) {
  return value.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' });
}
