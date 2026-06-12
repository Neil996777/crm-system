import type { ReactNode } from 'react';
import { ArrowLeft, Download } from 'lucide-react';
import { Badge, Button, Card, Pagination } from './ui';

export const DEFAULT_PAGE_SIZE = 25;

export type PageSlice<T> = {
  page: number;
  pageSize: number;
  totalPages: number;
  totalItems: number;
  items: T[];
};

export function paginate<T>(items: T[], page: number, pageSize: number): PageSlice<T> {
  const totalItems = items.length;
  const totalPages = Math.max(1, Math.ceil(totalItems / pageSize));
  const safePage = Math.min(Math.max(page, 1), totalPages);
  const start = (safePage - 1) * pageSize;
  return {
    page: safePage,
    pageSize,
    totalPages,
    totalItems,
    items: items.slice(start, start + pageSize)
  };
}

export function exportRows(filename: string, rows: Array<Record<string, ReactNode>>) {
  const headers = Object.keys(rows[0] ?? {});
  const escapeCell = (value: ReactNode) => {
    const text = String(value ?? '').replaceAll('"', '""');
    return `"${text}"`;
  };
  const csv = [headers.map(escapeCell).join(','), ...rows.map((row) => headers.map((header) => escapeCell(row[header])).join(','))].join('\n');
  const blob = new Blob([`\uFEFF${csv}`], { type: 'text/csv;charset=utf-8' });
  const url = URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.href = url;
  link.download = filename;
  link.click();
  URL.revokeObjectURL(url);
}

export function CrudListShell({
  title,
  description,
  scope,
  actions,
  toolbar,
  activeFilters,
  bulkBar,
  table,
  pagination,
  className
}: {
  title: string;
  description: ReactNode;
  scope?: ReactNode;
  actions?: ReactNode;
  toolbar: ReactNode;
  activeFilters?: ReactNode;
  bulkBar?: ReactNode;
  table: ReactNode;
  pagination: ReactNode;
  className?: string;
}) {
  return (
    <main className="content crudPage">
      <section className="pageHeader">
        <div>
          <h1>{title}</h1>
          <p>{description}</p>
        </div>
        {actions ? <div className="pageActions">{actions}</div> : null}
      </section>
      <Card className={className ? `listShell ${className}` : 'listShell'} aria-label={`${title}列表`}>
        <div className="listToolbarSlot">{toolbar}</div>
        {activeFilters ? <div className="applied">{activeFilters}</div> : null}
        {bulkBar}
        <div className="listTableSlot">{table}</div>
        <div className="listPaginationSlot">
          {scope ? <span className="pageMeta">{scope}</span> : null}
          {pagination}
        </div>
      </Card>
    </main>
  );
}

export function ActiveFilterSummary({
  children,
  onClear
}: {
  children: ReactNode;
  onClear: () => void;
}) {
  return (
    <>
      <span className="appliedLabel">已筛选</span>
      {children}
      <button className="clearFilterButton clearLink" type="button" onClick={onClear}>
        清除筛选
      </button>
    </>
  );
}

export function RecordIdentity({
  icon,
  title,
  subtitle,
  tone = 'sky',
  onTitleClick,
  titleAriaLabel
}: {
  icon?: ReactNode;
  title: ReactNode;
  subtitle?: ReactNode;
  tone?: 'sky' | 'mint' | 'peach' | 'purple' | 'primary';
  onTitleClick?: () => void;
  titleAriaLabel?: string;
}) {
  return (
    <div className="recordIdentity">
      {icon ? <span className={tone === 'primary' ? 'icon' : `icon ${tone}`}>{icon}</span> : null}
      <div>
        {onTitleClick ? (
          <button
            aria-label={titleAriaLabel}
            className="recordLinkButton"
            data-row-interactive="true"
            type="button"
            onClick={(event) => {
              event.stopPropagation();
              onTitleClick();
            }}
          >
            <span className="primaryLink">{title}</span>
          </button>
        ) : (
          <span className="primaryLink">{title}</span>
        )}
        {subtitle ? <span className="subMeta">{subtitle}</span> : null}
      </div>
    </div>
  );
}

export function StatusPill({
  children,
  tone = 'primary'
}: {
  children: ReactNode;
  tone?: 'primary' | 'success' | 'warning' | 'danger' | 'neutral';
}) {
  if (tone === 'neutral') return <Badge>{children}</Badge>;
  return <Badge tone={tone}>{children}</Badge>;
}

export function CrudPagination({
  slice,
  onPageChange,
  onPageSizeChange
}: {
  slice: PageSlice<unknown>;
  onPageChange: (page: number) => void;
  onPageSizeChange: (pageSize: number) => void;
}) {
  return (
    <Pagination
      page={slice.page}
      totalPages={slice.totalPages}
      totalItems={slice.totalItems}
      pageSize={slice.pageSize}
      onPageChange={onPageChange}
      onPageSizeChange={onPageSizeChange}
    />
  );
}

export function ExportSelectedButton({
  disabled,
  onExport
}: {
  disabled: boolean;
  onExport: () => void;
}) {
  return (
    <Button className="bulkButton" disabled={disabled} onClick={onExport}>
      <Download size={14} aria-hidden="true" />
      导出所选
    </Button>
  );
}

export function DetailHero({
  eyebrow,
  title,
  subtitle,
  icon,
  status,
  actions,
  stats,
  onBack
}: {
  eyebrow: ReactNode;
  title: ReactNode;
  subtitle?: ReactNode;
  icon?: ReactNode;
  status?: ReactNode;
  actions?: ReactNode;
  stats?: ReactNode;
  onBack: () => void;
}) {
  return (
    <article className="card detailHero">
      <button className="crumb" type="button" onClick={onBack}>
        <ArrowLeft size={15} aria-hidden="true" />
        {eyebrow}
      </button>
      <div className="heroTop">
        <div className="titleBlock">
          {icon ? <span className="icon sky">{icon}</span> : null}
          <div>
            <h1>{title}</h1>
            {subtitle || status ? (
              <div className="heroMeta">
                {status}
                {subtitle ? <span>{subtitle}</span> : null}
              </div>
            ) : null}
          </div>
        </div>
        {actions ? <div className="actions">{actions}</div> : null}
      </div>
      {stats ? <div className="detailStats">{stats}</div> : null}
    </article>
  );
}

export function DetailStat({
  label,
  value,
  icon,
  tone = 'primary'
}: {
  label: ReactNode;
  value: ReactNode;
  icon?: ReactNode;
  tone?: 'primary' | 'sky' | 'mint' | 'peach' | 'purple';
}) {
  return (
    <div className="detailStat">
      {icon ? <span className={tone === 'primary' ? 'statIcon' : `statIcon ${tone}`}>{icon}</span> : null}
      <div>
        <label>{label}</label>
        <strong>{value}</strong>
      </div>
    </div>
  );
}

export function FormShell({
  title,
  description,
  badge,
  children,
  side,
  actions,
  onCancel
}: {
  title: ReactNode;
  description?: ReactNode;
  badge?: ReactNode;
  children: ReactNode;
  side?: ReactNode;
  actions?: ReactNode;
  onCancel: () => void;
}) {
  return (
    <main className="content crudPage">
      <section className="pageHeader">
        <div>
          <h1>{title}</h1>
          {description ? <p>{description}</p> : null}
        </div>
        <div className="pageActions">
          <Button onClick={onCancel}>取消</Button>
          {actions}
        </div>
      </section>
      <div className="formShell">
        <article className="card formCard">
          <div className="formHead">
            <div className="formHeadTitle">
              <span className="icon sky" aria-hidden="true" />
              <div>
                <strong>{title}</strong>
                {description ? <div className="helper">{description}</div> : null}
              </div>
            </div>
            {badge ? <Badge tone="primary">{badge}</Badge> : null}
          </div>
          {children}
        </article>
        {side ? <aside className="formSide">{side}</aside> : null}
      </div>
    </main>
  );
}

export function FormSection({ title, children }: { title: ReactNode; children: ReactNode }) {
  return (
    <section className="formSection">
      <h2>{title}</h2>
      {children}
    </section>
  );
}
