export const appName = 'CRM 系统';

export const navLabels = {
  overview: '工作台',
  leads: '线索',
  accounts: '公司/客户',
  contacts: '联系人',
  opportunities: '商机',
  quotes: '报价',
  contracts: '合同',
  payments: '回款',
  tasks: '任务',
  reminders: '提醒中心',
  managerOverview: '报表',
  importExport: '导入/导出',
  userManagement: '管理：用户与角色',
  operationLogs: '操作日志'
} as const;

export const roleLabel: Record<string, string> = {
  Administrator: '管理员',
  'Sales Manager': '销售经理',
  Sales: '销售'
};

export const userStatusLabel: Record<string, string> = {
  Active: '启用',
  Disabled: '停用',
  Inactive: '停用'
};

export const accountStatusLabel: Record<string, string> = {
  Prospect: '潜在客户',
  Active: '活跃客户',
  Inactive: '停用客户'
};

export const leadStatusLabel: Record<string, string> = {
  Unassigned: '未分配',
  'Pending Qualification': '待确认',
  Valid: '有效',
  Invalid: '无效',
  'Converted To Opportunity': '已转为商机'
};

export const opportunityStageLabel: Record<string, string> = {
  'New Opportunity': '新商机',
  'Needs Confirmed': '需求已确认',
  Quote: '报价',
  'Contract Negotiation': '合同谈判',
  Won: '赢单',
  Lost: '丢单'
};

export const quoteStatusLabel: Record<string, string> = {
  Draft: '草稿',
  Sent: '已发送',
  Accepted: '已接受',
  Rejected: '已拒绝',
  Expired: '已过期'
};

export const contractStatusLabel: Record<string, string> = {
  'Pending Signature': '待签署',
  Signed: '已签署',
  Active: '启用',
  Completed: '已完成',
  Terminated: '已终止'
};

export const paymentStatusLabel: Record<string, string> = {
  'No plan': '无计划',
  Unpaid: '未回款',
  Pending: '待回款',
  PartiallyPaid: '部分回款',
  Paid: '已回款',
  Overdue: '已逾期',
  Cancelled: '已取消'
};

export const taskStatusLabel: Record<string, string> = {
  Open: '待处理',
  Completed: '已完成',
  Cancelled: '已取消',
  Overdue: '已逾期'
};

export const reminderTypeLabel: Record<string, string> = {
  task_due: '任务到期',
  task_overdue: '任务逾期',
  contract_pending_signature: '合同待签署',
  payment_due: '回款到期',
  payment_overdue: '回款逾期'
};

export const objectTypeLabel: Record<string, string> = {
  lead: '线索',
  Lead: '线索',
  Opportunity: '商机',
  Account: '客户',
  Contact: '联系人',
  Contract: '合同',
  Quote: '报价',
  Payment: '回款',
  Task: '任务',
  User: '用户'
};

export const lostReasonLabel: Record<string, string> = {
  PRICE: '价格',
  COMPETITOR: '竞争对手',
  NO_BUDGET: '无预算',
  NO_DECISION: '未决策',
  OTHER: '其他'
};

export const resultLabel: Record<string, string> = {
  SUCCESS: '成功',
  FAILED: '失败',
  Success: '成功',
  Failed: '失败',
  success: '成功',
  failed: '失败'
};

export const runStatusLabel: Record<string, string> = {
  Pending: '待处理',
  Running: '运行中',
  Completed: '已完成',
  Succeeded: '成功',
  Failed: '失败',
  Cancelled: '已取消',
  Retained: '已留存',
  Cleaned: '已清理'
};

export const archiveStatusLabel: Record<string, string> = {
  Archived: '已归档'
};

export const reportScopeLabel: Record<string, string> = {
  all: '全部',
  team: '团队',
  owned: '本人'
};

export const reportArchiveFilterLabel: Record<string, string> = {
  active_default: '默认仅活动记录',
  include_archived: '包含已归档',
  archived_only: '仅已归档'
};

export function labelFor(labels: Record<string, string>, value: string | null | undefined) {
  if (!value) return '';
  return labels[value] ?? value;
}

export function summaryTextZh(summary: Record<string, unknown> | undefined) {
  if (!summary || Object.keys(summary).length === 0) return '无';
  return Object.entries(summary)
    .filter(([, value]) => value !== '' && value !== null && value !== undefined)
    .map(([key, value]) => `${key}: ${String(value)}`)
    .join(', ');
}
