import type {
  ButtonHTMLAttributes,
  CSSProperties,
  HTMLAttributes,
  InputHTMLAttributes,
  ReactNode,
  SelectHTMLAttributes,
  TextareaHTMLAttributes
} from 'react';
import {
  AlertCircle,
  CheckCircle2,
  ChevronDown,
  ChevronLeft,
  ChevronRight,
  ChevronUp,
  ChevronsUpDown,
  Loader2,
  Search,
  X
} from 'lucide-react';

function cx(...values: Array<string | false | null | undefined>) {
  return values.filter(Boolean).join(' ');
}

export type BadgeTone = 'neutral' | 'primary' | 'success' | 'warning' | 'danger';
export type AccentTone = 'primary' | 'sky' | 'mint' | 'peach' | 'purple' | 'success' | 'warning' | 'danger';

export function Card({ className, ...props }: HTMLAttributes<HTMLDivElement>) {
  return <div className={cx('card', className)} {...props} />;
}

export function Panel({ className, ...props }: HTMLAttributes<HTMLElement>) {
  return <section className={cx('panel', className)} {...props} />;
}

export function PanelHeader({
  title,
  description,
  meta,
  actions,
  className
}: {
  title: ReactNode;
  description?: ReactNode;
  meta?: ReactNode;
  actions?: ReactNode;
  className?: string;
}) {
  return (
    <div className={cx('sectionHeader', className)}>
      <div>
        <h2>{title}</h2>
        {description ? <p>{description}</p> : null}
      </div>
      {meta || actions ? (
        <div className="headerActions">
          {meta ? <span className="badge">{meta}</span> : null}
          {actions}
        </div>
      ) : null}
    </div>
  );
}

export function PageHeader({
  title,
  description,
  actions,
  className
}: {
  title: string;
  description?: ReactNode;
  actions?: ReactNode;
  className?: string;
}) {
  return (
    <section className={cx('pageHeader', className)}>
      <div>
        <h1>{title}</h1>
        {description ? <p>{description}</p> : null}
      </div>
      {actions ? <div className="pageActions">{actions}</div> : null}
    </section>
  );
}

export function Button({
  variant = 'secondary',
  className,
  ...props
}: ButtonHTMLAttributes<HTMLButtonElement> & { variant?: 'primary' | 'secondary' | 'ghost' }) {
  const variantClass = variant === 'primary' ? 'primaryButton' : variant === 'ghost' ? 'ghostButton' : 'secondaryButton';
  return <button className={cx(variantClass, className)} type="button" {...props} />;
}

export function Badge({
  tone = 'neutral',
  className,
  ...props
}: HTMLAttributes<HTMLSpanElement> & { tone?: BadgeTone }) {
  return <span className={cx('badge', tone !== 'neutral' && tone, className)} {...props} />;
}

export function StatusBadge({ children, tone = 'primary' }: { children: ReactNode; tone?: Exclude<BadgeTone, 'neutral'> }) {
  return <span className={cx('statusPill', tone)}>{children}</span>;
}

export function MetricCard({
  label,
  value,
  icon,
  tone = 'sky',
  delta,
  className
}: {
  label: ReactNode;
  value: ReactNode;
  icon?: ReactNode;
  tone?: 'sky' | 'mint' | 'peach' | 'purple';
  delta?: ReactNode;
  className?: string;
}) {
  return (
    <article className={cx('metricTile', className)}>
      <div>
        <span>{label}</span>
        <strong>{value}</strong>
        {delta ? <div className="panelMeta">{delta}</div> : null}
      </div>
      {icon ? <span className={cx('metricIcon', tone)}>{icon}</span> : null}
    </article>
  );
}

export function EmptyState({ children = '暂无数据。', className }: { children?: ReactNode; className?: string }) {
  return <p className={cx('emptyState', className)}>{children}</p>;
}

export function ErrorState({ children = '加载失败，请稍后重试。', className }: { children?: ReactNode; className?: string }) {
  return (
    <div className={cx('errorBanner', className)} role="alert">
      <AlertCircle size={16} aria-hidden="true" />
      <span>{children}</span>
    </div>
  );
}

export function PermissionDenied({ children = '当前账号无权查看该页面。' }: { children?: ReactNode }) {
  return (
    <section className="permissionDenied" role="status" aria-live="polite">
      {children}
    </section>
  );
}

export function SkeletonBlock({ lines = 3, label = '加载中' }: { lines?: number; label?: string }) {
  return (
    <div className="loadingState skeleton" role="status" aria-label={label}>
      {Array.from({ length: lines }, (_, index) => (
        <span className="loadingLine" key={index} style={{ width: `${92 - index * 12}%` }} />
      ))}
    </div>
  );
}

export type DataTableSortDirection = 'asc' | 'desc' | null | undefined;
export type DataTableRowKey = string | number;
export type DataTableColumn<T> = {
  key: string;
  header: ReactNode;
  render: (row: T, index: number) => ReactNode;
  align?: 'left' | 'center' | 'right';
  className?: string;
  headerClassName?: string;
  width?: string;
  sortable?: boolean;
  sortDirection?: DataTableSortDirection;
  onSort?: () => void;
};

export type DataTableProps<T> = {
  caption?: string;
  children?: ReactNode;
  className?: string;
  columns?: Array<DataTableColumn<T>>;
  rows?: T[];
  rowKey?: (row: T, index: number) => DataTableRowKey;
  selectedRowKeys?: DataTableRowKey[];
  onToggleRow?: (row: T, checked: boolean) => void;
  onToggleAll?: (checked: boolean) => void;
  getRowClassName?: (row: T, index: number) => string | undefined;
  empty?: ReactNode;
  actions?: (row: T, index: number) => ReactNode;
};

function sortIcon(direction: DataTableSortDirection) {
  if (direction === 'asc') return <ChevronUp size={14} aria-hidden="true" />;
  if (direction === 'desc') return <ChevronDown size={14} aria-hidden="true" />;
  return <ChevronsUpDown size={14} aria-hidden="true" />;
}

function ariaSort(direction: DataTableSortDirection) {
  if (direction === 'asc') return 'ascending';
  if (direction === 'desc') return 'descending';
  return 'none';
}

export function DataTable<T>({
  caption,
  children,
  className,
  columns,
  rows = [],
  rowKey,
  selectedRowKeys = [],
  onToggleRow,
  onToggleAll,
  getRowClassName,
  empty,
  actions
}: DataTableProps<T>) {
  if (children || !columns) {
    return (
      <div className={cx('tableWrap', className)}>
        <table className="dataTable">
          {caption ? <caption className="srOnly">{caption}</caption> : null}
          {children}
        </table>
      </div>
    );
  }

  const selected = new Set<DataTableRowKey>(selectedRowKeys);
  const keyForRow = rowKey ?? ((_row: T, index: number) => index);
  const allSelected = rows.length > 0 && rows.every((row, index) => selected.has(keyForRow(row, index)));
  const someSelected = rows.some((row, index) => selected.has(keyForRow(row, index)));
  const hasSelection = Boolean(onToggleRow);
  const colSpan = columns.length + (hasSelection ? 1 : 0) + (actions ? 1 : 0);

  return (
    <div className={cx('tableWrap', className)}>
      <table className="dataTable">
        {caption ? <caption className="srOnly">{caption}</caption> : null}
        <thead>
          <tr>
            {hasSelection ? (
              <th className="selectCell" scope="col">
                <input
                  aria-label="选择全部"
                  aria-checked={someSelected && !allSelected ? 'mixed' : allSelected}
                  checked={allSelected}
                  className="rowCheckbox"
                  disabled={!onToggleAll && !onToggleRow}
                  type="checkbox"
                  onChange={(event) => onToggleAll?.(event.currentTarget.checked)}
                />
              </th>
            ) : null}
            {columns.map((column) => (
              <th
                aria-sort={column.sortable ? ariaSort(column.sortDirection) : undefined}
                className={cx(column.align && `align-${column.align}`, column.headerClassName)}
                key={column.key}
                scope="col"
                style={column.width ? ({ width: column.width } as CSSProperties) : undefined}
              >
                {column.sortable ? (
                  <button className="sortButton" type="button" onClick={column.onSort} disabled={!column.onSort}>
                    <span>{column.header}</span>
                    {sortIcon(column.sortDirection)}
                  </button>
                ) : (
                  column.header
                )}
              </th>
            ))}
            {actions ? (
              <th className="rowActionsCell" scope="col">
                操作
              </th>
            ) : null}
          </tr>
        </thead>
        <tbody>
          {rows.length === 0 ? (
            <tr>
              <td className="emptyCell" colSpan={colSpan}>
                {empty ?? '暂无数据。'}
              </td>
            </tr>
          ) : (
            rows.map((row, index) => {
              const currentKey = keyForRow(row, index);
              return (
                <tr className={getRowClassName?.(row, index)} key={currentKey}>
                  {hasSelection ? (
                    <td className="selectCell">
                      <input
                        aria-label={`选择第 ${index + 1} 行`}
                        checked={selected.has(currentKey)}
                        className="rowCheckbox"
                        type="checkbox"
                        onChange={(event) => onToggleRow?.(row, event.currentTarget.checked)}
                      />
                    </td>
                  ) : null}
                  {columns.map((column) => (
                    <td className={cx(column.align && `align-${column.align}`, column.className)} key={column.key}>
                      {column.render(row, index)}
                    </td>
                  ))}
                  {actions ? <td className="rowActionsCell">{actions(row, index)}</td> : null}
                </tr>
              );
            })
          )}
        </tbody>
      </table>
    </div>
  );
}

export type ActiveFilter = {
  label: ReactNode;
  value?: ReactNode;
  tone?: BadgeTone;
};

export function Toolbar({
  className,
  searchValue,
  onSearchChange,
  searchPlaceholder = '搜索',
  filters,
  activeFilters,
  onClearFilters,
  summary,
  actions,
  children,
  ...props
}: HTMLAttributes<HTMLDivElement> & {
  searchValue?: string;
  onSearchChange?: (value: string) => void;
  searchPlaceholder?: string;
  filters?: ReactNode;
  activeFilters?: ActiveFilter[];
  onClearFilters?: () => void;
  summary?: ReactNode;
  actions?: ReactNode;
}) {
  const hasStructuredContent = onSearchChange || filters || activeFilters?.length || onClearFilters || summary || actions;

  if (!hasStructuredContent) {
    return (
      <div className={cx('toolbar', className)} {...props}>
        {children}
      </div>
    );
  }

  return (
    <div className={cx('toolbar', 'toolbarStructured', className)} {...props}>
      {onSearchChange ? (
        <label className="toolbarSearch">
          <span className="srOnly">{searchPlaceholder}</span>
          <Search size={16} aria-hidden="true" />
          <input value={searchValue ?? ''} placeholder={searchPlaceholder} onChange={(event) => onSearchChange(event.currentTarget.value)} />
        </label>
      ) : null}
      {filters ? <div className="toolbarFilters">{filters}</div> : null}
      {actions ? <div className="toolbarActions">{actions}</div> : null}
      {children}
      {activeFilters?.length || summary || onClearFilters ? (
        <div className="activeFilterRow">
          {summary ? <span className="filterSummary">{summary}</span> : null}
          {activeFilters?.map((filter, index) => (
            <Badge key={index} tone={filter.tone ?? 'neutral'}>
              {filter.label}
              {filter.value ? `：${filter.value}` : null}
            </Badge>
          ))}
          {onClearFilters ? (
            <button className="clearFilterButton" type="button" onClick={onClearFilters}>
              <X size={14} aria-hidden="true" />
              清除筛选
            </button>
          ) : null}
        </div>
      ) : null}
    </div>
  );
}

type FormFieldBase = {
  label: string;
  hint?: ReactNode;
  error?: ReactNode;
  className?: string;
};

export function TextField({ label, hint, error, className, type, ...props }: FormFieldBase & InputHTMLAttributes<HTMLInputElement>) {
  return (
    <label className={cx(className, type === 'date' && 'dateField')}>
      <span>{label}</span>
      <input className={cx(type === 'date' && 'dateControl')} type={type} {...props} />
      {error ? <span className="dangerText">{error}</span> : hint ? <span className="panelMeta">{hint}</span> : null}
    </label>
  );
}

export function SelectField({ label, hint, error, className, children, ...props }: FormFieldBase & SelectHTMLAttributes<HTMLSelectElement>) {
  return (
    <label className={className}>
      <span>{label}</span>
      <select {...props}>{children}</select>
      {error ? <span className="dangerText">{error}</span> : hint ? <span className="panelMeta">{hint}</span> : null}
    </label>
  );
}

export function TextAreaField({ label, hint, error, className, ...props }: FormFieldBase & TextareaHTMLAttributes<HTMLTextAreaElement>) {
  return (
    <label className={className}>
      <span>{label}</span>
      <textarea {...props} />
      {error ? <span className="dangerText">{error}</span> : hint ? <span className="panelMeta">{hint}</span> : null}
    </label>
  );
}

function pageWindow(page: number, totalPages: number) {
  const safeTotal = Math.max(totalPages, 1);
  const start = Math.max(1, Math.min(page - 2, safeTotal - 4));
  const end = Math.min(safeTotal, start + 4);
  return Array.from({ length: end - start + 1 }, (_, index) => start + index);
}

export function Pagination({
  page,
  totalPages,
  totalItems,
  pageSize,
  pageSizeOptions = [10, 20, 50],
  onPageChange,
  onPageSizeChange,
  onPrevious,
  onNext,
  className
}: {
  page: number;
  totalPages: number;
  totalItems?: number;
  pageSize?: number;
  pageSizeOptions?: number[];
  onPageChange?: (page: number) => void;
  onPageSizeChange?: (pageSize: number) => void;
  onPrevious?: () => void;
  onNext?: () => void;
  className?: string;
}) {
  const safeTotalPages = Math.max(totalPages, 1);
  const safePage = Math.min(Math.max(page, 1), safeTotalPages);
  const previous = () => (onPrevious ? onPrevious() : onPageChange?.(Math.max(1, safePage - 1)));
  const next = () => (onNext ? onNext() : onPageChange?.(Math.min(safeTotalPages, safePage + 1)));

  return (
    <nav className={cx('pagination', className)} aria-label="分页">
      {typeof totalItems === 'number' ? <span className="pageMeta">共 {totalItems} 条</span> : null}
      <div className="pageNumberGroup">
        <button className="pageButton" type="button" onClick={previous} disabled={safePage <= 1}>
          <ChevronLeft size={15} aria-hidden="true" />
          上一页
        </button>
        {pageWindow(safePage, safeTotalPages).map((pageNumber) => (
          <button
            aria-current={pageNumber === safePage ? 'page' : undefined}
            className={cx('pageButton', pageNumber === safePage && 'selected')}
            key={pageNumber}
            type="button"
            onClick={() => onPageChange?.(pageNumber)}
            disabled={!onPageChange && pageNumber !== safePage}
          >
            {pageNumber}
          </button>
        ))}
        <button className="pageButton" type="button" onClick={next} disabled={safePage >= safeTotalPages}>
          下一页
          <ChevronRight size={15} aria-hidden="true" />
        </button>
      </div>
      {pageSize && onPageSizeChange ? (
        <label className="pageSizeControl">
          <span>每页</span>
          <select value={pageSize} onChange={(event) => onPageSizeChange(Number(event.currentTarget.value))}>
            {pageSizeOptions.map((option) => (
              <option key={option} value={option}>
                {option} 条
              </option>
            ))}
          </select>
        </label>
      ) : (
        <span className="chip" aria-live="polite">
          第 {safePage} / {safeTotalPages} 页
        </span>
      )}
    </nav>
  );
}

export function BulkActionBar({ children }: { children: ReactNode }) {
  return (
    <div className="bulkBar" aria-live="polite">
      {children}
    </div>
  );
}

export function LiveToggle({ active, onToggle }: { active: boolean; onToggle: () => void }) {
  return (
    <button className="secondaryButton" type="button" aria-pressed={active} onClick={onToggle}>
      <span className={cx('liveDot', active && 'livePulse')} aria-hidden="true" />
      {active ? '实时开启' : '实时暂停'}
    </button>
  );
}

export function FunnelBars({
  rows,
  max
}: {
  rows: Array<{ label: ReactNode; value: number; suffix?: ReactNode }>;
  max?: number;
}) {
  const computedMax = max ?? Math.max(1, ...rows.map((row) => row.value));
  return (
    <div className="pipelineViz">
      {rows.map((row, index) => {
        const width = computedMax > 0 ? Math.max(4, Math.round((row.value / computedMax) * 100)) : 4;
        return (
          <div className="funnelRow" key={index}>
            <span className="funnelLabel">{row.label}</span>
            <span className="funnelTrack" aria-hidden="true">
              <span className="funnelFill" style={{ width: `${width}%` }} />
            </span>
            <span className="funnelValue">
              {row.value}
              {row.suffix}
            </span>
          </div>
        );
      })}
    </div>
  );
}

export type TrendPoint = {
  label: string;
  value: number;
};

function trendPath(points: TrendPoint[]) {
  if (points.length === 0) return { line: '', area: '' };
  const values = points.map((point) => point.value);
  const min = Math.min(...values);
  const max = Math.max(...values);
  const span = Math.max(1, max - min);
  const coords = points.map((point, index) => {
    const x = points.length === 1 ? 120 : 12 + (index * 216) / (points.length - 1);
    const y = 84 - ((point.value - min) / span) * 64;
    return [x, y] as const;
  });
  const line = coords.map(([x, y], index) => `${index === 0 ? 'M' : 'L'} ${x.toFixed(2)} ${y.toFixed(2)}`).join(' ');
  const area = `${line} L ${coords[coords.length - 1][0].toFixed(2)} 90 L ${coords[0][0].toFixed(2)} 90 Z`;
  return { line, area };
}

export function TrendPanel({
  title,
  children,
  meta,
  points,
  valueLabel
}: {
  title: ReactNode;
  children?: ReactNode;
  meta?: ReactNode;
  points?: TrendPoint[];
  valueLabel?: (point: TrendPoint) => ReactNode;
}) {
  const paths = trendPath(points ?? []);
  return (
    <Panel className="trendPanel">
      <PanelHeader title={title} meta={meta} />
      {points?.length ? (
        <div className="trendChart" role="img" aria-label="趋势折线">
          <svg viewBox="0 0 240 96" aria-hidden="true" focusable="false">
            <path className="trendArea" d={paths.area} />
            <path className="trendLine" d={paths.line} />
            {points.map((point, index) => {
              const values = points.map((current) => current.value);
              const min = Math.min(...values);
              const max = Math.max(...values);
              const span = Math.max(1, max - min);
              const cxValue = points.length === 1 ? 120 : 12 + (index * 216) / (points.length - 1);
              const cyValue = 84 - ((point.value - min) / span) * 64;
              return <circle className="trendPoint" cx={cxValue} cy={cyValue} key={point.label} r="3.5" />;
            })}
          </svg>
          <div className="trendLegend">
            {points.map((point) => (
              <span key={point.label}>
                <small>{point.label}</small>
                <strong>{valueLabel ? valueLabel(point) : point.value}</strong>
              </span>
            ))}
          </div>
        </div>
      ) : (
        children
      )}
    </Panel>
  );
}

export type DonutSegment = {
  label: ReactNode;
  value: number;
  tone?: AccentTone;
};

const donutToneOrder: AccentTone[] = ['primary', 'sky', 'mint', 'peach', 'purple', 'success', 'warning', 'danger'];

export function DonutChart({
  segments,
  label,
  center,
  className
}: {
  segments: DonutSegment[];
  label: string;
  center?: ReactNode;
  className?: string;
}) {
  const total = Math.max(0, segments.reduce((sum, segment) => sum + Math.max(0, segment.value), 0));
  const circumference = 2 * Math.PI * 50;
  let offset = 0;

  return (
    <div className={cx('donutShell', className)} role="img" aria-label={label}>
      <svg viewBox="0 0 128 128" aria-hidden="true" focusable="false">
        <circle className="donutTrack" cx="64" cy="64" r="50" />
        {segments.map((segment, index) => {
          const tone = segment.tone ?? donutToneOrder[index % donutToneOrder.length];
          const value = Math.max(0, segment.value);
          const dash = total > 0 ? (value / total) * circumference : 0;
          const currentOffset = offset;
          offset += dash;
          return (
            <circle
              className={cx('donutSegment', `tone-${tone}`)}
              cx="64"
              cy="64"
              key={`${index}-${String(segment.label)}`}
              r="50"
              strokeDasharray={`${Math.max(0, dash - 2).toFixed(2)} ${circumference.toFixed(2)}`}
              strokeDashoffset={(-currentOffset).toFixed(2)}
            />
          );
        })}
        <text className="donutCenter" dominantBaseline="middle" textAnchor="middle" x="64" y="64">
          {center ?? total}
        </text>
      </svg>
      <div className="legend">
        {segments.map((segment, index) => {
          const tone = segment.tone ?? donutToneOrder[index % donutToneOrder.length];
          const percent = total > 0 ? Math.round((Math.max(0, segment.value) / total) * 100) : 0;
          return (
            <span className="legendItem" key={`${index}-${String(segment.label)}`}>
              <span className={cx('legendSwatch', `tone-${tone}`)} aria-hidden="true" />
              <span>{segment.label}</span>
              <strong>{percent}%</strong>
            </span>
          );
        })}
      </div>
    </div>
  );
}

export const StageDonut = DonutChart;

export type LeaderboardItem = {
  label: ReactNode;
  value: number;
  meta?: ReactNode;
  suffix?: ReactNode;
  tone?: AccentTone;
};

export function Leaderboard({
  items,
  title,
  className
}: {
  items: LeaderboardItem[];
  title?: ReactNode;
  className?: string;
}) {
  const max = Math.max(1, ...items.map((item) => item.value));
  return (
    <div className={cx('leaderboard', className)}>
      {title ? <h3>{title}</h3> : null}
      <ol>
        {items.map((item, index) => {
          const width = Math.max(4, Math.round((item.value / max) * 100));
          return (
            <li className="leaderboardRow" key={index}>
              <span className="leaderRank">{index + 1}</span>
              <div className="leaderMain">
                <strong>{item.label}</strong>
                {item.meta ? <span>{item.meta}</span> : null}
                <span className="leaderTrack" aria-hidden="true">
                  <span className={cx('leaderFill', item.tone && `tone-${item.tone}`)} style={{ width: `${width}%` }} />
                </span>
              </div>
              <strong className="leaderValue">
                {item.value}
                {item.suffix}
              </strong>
            </li>
          );
        })}
      </ol>
    </div>
  );
}

export function ReminderRowCard({
  title,
  description,
  icon,
  meta,
  time,
  badges,
  actions,
  tone = 'primary',
  overdue = false,
  className
}: {
  title: ReactNode;
  description?: ReactNode;
  icon?: ReactNode;
  meta?: ReactNode;
  time?: ReactNode;
  badges?: ReactNode;
  actions?: ReactNode;
  tone?: AccentTone;
  overdue?: boolean;
  className?: string;
}) {
  return (
    <article className={cx('reminderCard', `tone-${tone}`, overdue && 'overdue', className)}>
      {icon ? <span className={cx('flowIcon', tone)}>{icon}</span> : null}
      <div className="reminderBody">
        <div className="reminderTitleRow">
          <strong>{title}</strong>
          {badges ? <div className="badgeRow">{badges}</div> : null}
        </div>
        {description ? <p>{description}</p> : null}
        {meta || time ? (
          <div className="reminderMeta">
            {meta ? <span>{meta}</span> : null}
            {time ? <time>{time}</time> : null}
          </div>
        ) : null}
      </div>
      {actions ? <div className="reminderActions">{actions}</div> : null}
    </article>
  );
}

export type FocusSideCard = {
  key: string;
  title: ReactNode;
  metric?: ReactNode;
  meta?: ReactNode;
  icon?: ReactNode;
  selected?: boolean;
  motionIndex?: number;
  onSelect?: () => void;
};

export function FocusStage({
  title,
  subtitle,
  icon,
  tools,
  children,
  sideCards,
  onBack,
  backLabel = '返回总览'
}: {
  title: ReactNode;
  subtitle?: ReactNode;
  icon?: ReactNode;
  tools?: ReactNode;
  children: ReactNode;
  sideCards?: FocusSideCard[];
  onBack?: () => void;
  backLabel?: string;
}) {
  return (
    <div className="focus">
      <section className="stage" aria-label="聚焦舞台">
        <div className="stageHead">
          <div className="titleBlock">
            {icon ? <span className="panelIcon">{icon}</span> : null}
            <div>
              <h1 data-focus-heading tabIndex={-1}>{title}</h1>
              {subtitle ? <div className="stageSub">{subtitle}</div> : null}
            </div>
          </div>
          <div className="stageTools">
            {tools}
            {onBack ? (
              <button className="secondaryButton" type="button" onClick={onBack} aria-label={backLabel}>
                <ChevronLeft size={15} aria-hidden="true" />
                {backLabel}
              </button>
            ) : null}
          </div>
        </div>
        {children}
      </section>
      {sideCards?.length ? (
        <aside className="side" aria-label="看板选择器">
          {sideCards.map((card) => {
            const motionIndex = card.motionIndex ?? 0;
            const motionStyle = {
              '--strip-index': motionIndex,
              '--strip-enter-delay': `${80 + motionIndex * 24}ms`,
              '--strip-exit-delay': `${motionIndex * 16}ms`
            } as CSSProperties;
            const content = (
              <>
                {card.icon ? <span className="panelIcon">{card.icon}</span> : null}
                <span>
                  <strong>{card.title}</strong>
                  {card.meta ? <small>{card.meta}</small> : null}
                </span>
                {card.metric ? <span className="sideMetric">{card.metric}</span> : null}
              </>
            );
            return card.onSelect ? (
              <button
                className={cx('sideCard', card.selected && 'selected')}
                aria-current={card.selected ? 'true' : undefined}
                data-focus-side-card={card.key}
                key={card.key}
                style={motionStyle}
                type="button"
                onClick={card.onSelect}
              >
                {content}
              </button>
            ) : (
              <article
                className={cx('sideCard', card.selected && 'selected')}
                aria-current={card.selected ? 'true' : undefined}
                data-focus-side-card={card.key}
                key={card.key}
                style={motionStyle}
              >
                {content}
              </article>
            );
          })}
        </aside>
      ) : null}
    </div>
  );
}

export const CardFocusStage = FocusStage;

export function AuditEventCard({
  actor,
  action,
  resource,
  result,
  occurredAt,
  safeSummary,
  eventId,
  correlationId,
  hash,
  badges,
  className
}: {
  actor: ReactNode;
  action: ReactNode;
  resource?: ReactNode;
  result?: ReactNode;
  occurredAt: ReactNode;
  safeSummary: ReactNode;
  eventId?: ReactNode;
  correlationId?: ReactNode;
  hash?: ReactNode;
  badges?: ReactNode;
  className?: string;
}) {
  return (
    <article className={cx('auditEventCard', className)}>
      <div className="auditEventHead">
        <div>
          <strong>{action}</strong>
          <span>
            {actor}
            {resource ? <> · {resource}</> : null}
          </span>
        </div>
        {badges || result ? (
          <div className="badgeRow">
            {badges}
            {result ? <Badge tone="primary">{result}</Badge> : null}
          </div>
        ) : null}
      </div>
      <p className="safeSummary">{safeSummary}</p>
      <div className="auditMeta">
        <time>{occurredAt}</time>
        {eventId ? <span>事件 {eventId}</span> : null}
        {correlationId ? <span>关联 {correlationId}</span> : null}
        {hash ? <span>摘要 {hash}</span> : null}
      </div>
    </article>
  );
}

export function Drawer({
  title,
  open,
  onClose,
  children
}: {
  title: string;
  open: boolean;
  onClose: () => void;
  children: ReactNode;
}) {
  if (!open) return null;
  return (
    <div className="drawerLayer" role="presentation">
      <section className="drawerPanel" role="dialog" aria-modal="true" aria-label={title}>
        <div className="detailHeader">
          <h2>{title}</h2>
          <button className="ghostButton" type="button" onClick={onClose} aria-label="关闭">
            <X size={16} aria-hidden="true" />
          </button>
        </div>
        {children}
      </section>
    </div>
  );
}

export function Toast({
  children,
  tone = 'success'
}: {
  children: ReactNode;
  tone?: 'success' | 'danger';
}) {
  const Icon = tone === 'success' ? CheckCircle2 : AlertCircle;
  return (
    <div className={tone === 'success' ? 'successNotice' : 'errorBanner'} role="status" aria-live="polite">
      <Icon size={16} aria-hidden="true" />
      <span>{children}</span>
    </div>
  );
}

export function InlineLoading({ label = '加载中' }: { label?: string }) {
  return (
    <span className="chip" role="status">
      <Loader2 size={14} aria-hidden="true" />
      {label}
    </span>
  );
}
