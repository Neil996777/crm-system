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

export const priorityLabel: Record<string, string> = {
  P0: 'P0',
  P1: 'P1',
  P2: 'P2',
  P3: 'P3'
};

export const objectTypeLabel: Record<string, string> = {
  lead: '线索',
  opportunity: '商机',
  account: '客户',
  contact: '联系人',
  contract: '合同',
  quote: '报价',
  payment: '回款',
  task: '任务',
  user: '用户',
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

export const fileSafetyLabel: Record<string, string> = {
  dangerous_cells_prefixed: '危险单元格已安全前缀化',
  safe: '安全',
  Safe: '安全'
};

export const errorMessageZh: Record<string, string> = {
  'The request is invalid.': '请求无效。',
  'The account input is invalid.': '客户输入无效。',
  'The archive input is invalid.': '归档输入无效。',
  'The contact input is invalid.': '联系人输入无效。',
  'The lead input is invalid.': '线索输入无效。',
  'The qualification input is invalid.': '资质评估输入无效。',
  'The conversion input is invalid.': '转化输入无效。',
  'The opportunity input is invalid.': '商机输入无效。',
  'The stage transition input is invalid.': '阶段流转输入无效。',
  'The close-won input is invalid.': '赢单输入无效。',
  'The close-lost input is invalid.': '丢单输入无效。',
  'The quote input is invalid.': '报价输入无效。',
  'The quote status input is invalid.': '报价状态输入无效。',
  'The contract input is invalid.': '合同输入无效。',
  'The contract status input is invalid.': '合同状态输入无效。',
  'The contract quote link is invalid.': '合同与报价的关联无效。',
  'The payment input is invalid.': '回款输入无效。',
  'The payment plan input is invalid.': '回款计划输入无效。',
  'The task input is invalid.': '任务输入无效。',
  'The task status input is invalid.': '任务状态输入无效。',
  'The work item input is invalid.': '工作项输入无效。',
  'The owner transfer input is invalid.': '负责人转移输入无效。',
  'The duplicate check input is invalid.': '查重输入无效。',
  'The duplicate warning confirmation is invalid.': '重复提醒确认无效。',
  'The duplicate warning confirmation was already used.': '该重复提醒确认已被使用。',
  'The obligation query input is invalid.': '未结事项查询输入无效。',
  'The reminder query input is invalid.': '提醒查询输入无效。',
  'The filter is invalid.': '筛选条件无效。',
  'The projection input is invalid.': '报表数据输入无效。',
  'The import input is invalid.': '导入输入无效。',
  'The export input is invalid.': '导出输入无效。',
  'The CSV content is invalid.': 'CSV 内容无效。',
  'Only CSV import is supported.': '仅支持 CSV 导入。',
  'The object type is not supported for import.': '该对象类型不支持导入。',
  'The object type is not supported for export.': '该对象类型不支持导出。',
  'Export confirmation is required.': '需要确认后才能导出。',
  'The import run could not be saved.': '导入记录保存失败。',
  'The export could not be completed.': '导出未能完成。',
  'The requested stage transition is not allowed.': '不允许该阶段流转。',
  'The requested quote status transition is not allowed.': '不允许该报价状态流转。',
  'The requested contract status transition is not allowed.': '不允许该合同状态流转。',
  'The requested task status transition is not allowed.': '不允许该任务状态流转。',
  'Terminal opportunity records cannot be closed again.': '已终结的商机不能再次关闭。',
  'Terminal opportunity records cannot be edited.': '已终结的商机不能编辑。',
  'Terminal opportunity records cannot change stage.': '已终结的商机不能变更阶段。',
  'A quote already exists for this opportunity.': '该商机已存在报价。',
  'A contract already exists for this quote.': '该报价已存在合同。',
  'A reason is required when contract amount differs from quote amount.': '合同金额与报价金额不一致时必须填写原因。',
  'Signed or effective date is required for this contract status.': '该合同状态需要填写签署或生效日期。',
  'Won requires a Signed related contract.': '赢单需要有已签署的关联合同。',
  'Lost reason is required.': '必须填写丢单原因。',
  'Payment amount must be greater than zero.': '回款金额必须大于零。',
  'Payment exceeds the remaining contract amount.': '回款金额超过合同剩余应收。',
  'Payments use the committed single currency.': '回款只能使用约定的单一币种。',
  'The lead cannot be converted in its current state.': '当前状态的线索无法转化。',
  'The lead has already been converted.': '该线索已转化。',
  'The record changed after it was opened.': '记录在你打开后已被他人修改，请刷新重试。',
  'The requested resource was not found.': '未找到所请求的资源。',
  'Permission denied.': '没有权限执行该操作。',
  'A required service is unavailable.': '依赖的服务暂不可用，请稍后重试。',
  'Service authentication failed.': '服务认证失败。',
  'Audit log failed.': '审计日志写入失败。',
  'The audit event could not be persisted.': '审计事件未能持久化。',
  'Authentication failed.': '认证失败。',
  'Request failed.': '请求失败。',
  'A required dependency is unavailable.': '依赖的服务暂不可用，请稍后重试。',
  'A required service returned an invalid response.': '依赖服务返回了无效响应。',
  'A required service returned an error.': '依赖服务返回错误。',
  'Choose a supported status.': '请选择支持的状态。',
  'Active obligations must be resolved before archive.': '归档前必须先处理未完成事项。',
  'Pending signature contracts must be signed or terminated before archive.': '待签署合同必须先签署或终止后才能归档。',
  'Unpaid payment plans must be paid before archive.': '未回款计划必须完成回款后才能归档。',
  'The last active Administrator cannot be disabled or downgraded.': '不能停用或降级最后一个启用的管理员。',
  'Company name or lead name is required.': '必须填写公司名称或线索名称。',
  'Source is required.': '必须填写来源。',
  'CSV cell is not safe to import.': 'CSV 单元格存在安全风险，不能导入。',
  'Row could not be imported.': '该行未能导入。',
  'Lead archived.': '线索已归档。',
  'Payment recorded.': '回款已登记。',
  'Lead created': '线索已创建',
  'Lead qualified': '线索已确认',
  'Lead disqualified': '线索已标记无效',
  'Lead converted': '线索已转化',
  'Record archived': '记录已归档',
  'Duplicate warning raised': '已触发重复提醒',
  'Account created': '客户已创建',
  'Account updated': '客户已更新',
  'Contact created': '联系人已创建',
  'Opportunity created': '商机已创建',
  'Opportunity updated': '商机已更新',
  'Opportunity stage changed': '商机阶段已变更',
  'Opportunity closed won': '商机已赢单关闭',
  'Opportunity closed lost': '商机已丢单关闭',
  'Opportunity won': '商机已赢单',
  'Opportunity lost': '商机已丢单',
  'Quote created': '报价已创建',
  'Quote changed': '报价已变更',
  'Quote accepted': '报价已接受',
  'Contract created': '合同已创建',
  'Contract signed': '合同已签署',
  'Contract terminated': '合同已终止',
  'Contract status changed': '合同状态已变更',
  'Payment plan created': '回款计划已创建',
  'Open work owner transferred': '待处理工作负责人已转移',
  'Work item created': '工作项已创建',
  'Task status changed': '任务状态已变更',
  'Report access denied': '报表访问被拒绝',
  'User signed in': '用户已登录',
  'User signed out': '用户已退出',
  'User access denied': '用户访问被拒绝',
  'User admin changed': '用户管理已变更'
};

export const actionLabel: Record<string, string> = {
  create_user: '新建用户',
  change_role: '变更角色',
  change_status: '变更状态',
  user_admin_changed: '用户管理变更',
  last_admin_blocked: '阻止最后一个管理员变更',
  access_denied: '访问被拒绝',
  login_failed: '登录失败',
  sign_in: '登录',
  sign_out: '退出登录',
  archive: '归档',
  payment_recorded: '登记回款',
  csv_import: 'CSV 导入',
  csv_export: 'CSV 导出',
  'Lead created': '线索已创建',
  'Lead qualified': '线索已确认',
  'Lead disqualified': '线索已标记无效',
  'Lead converted': '线索已转化',
  'Record archived': '记录已归档',
  'Duplicate warning raised': '已触发重复提醒',
  'Account created': '客户已创建',
  'Account updated': '客户已更新',
  'Contact created': '联系人已创建',
  'Opportunity created': '商机已创建',
  'Opportunity updated': '商机已更新',
  'Opportunity won': '商机已赢单',
  'Opportunity lost': '商机已丢单',
  'Quote accepted': '报价已接受',
  'Quote changed': '报价已变更',
  'Contract created': '合同已创建',
  'Contract signed': '合同已签署',
  'Contract terminated': '合同已终止',
  'Contract status changed': '合同状态已变更',
  'Payment recorded': '回款已登记',
  'Open work owner transferred': '待处理工作负责人已转移',
  'Work item created': '工作项已创建',
  'Task status changed': '任务状态已变更',
  'Report access denied': '报表访问被拒绝',
  'Owner changed': '负责人已变更',
  'Stage changed': '阶段已变更',
  'Access denied': '访问被拒绝'
};

export type LocalizableError = {
  safeMessage?: string;
  fieldErrors?: Array<{ field: string; code: string; safeMessage?: string }>;
};

export function labelFor(labels: Record<string, string>, value: string | null | undefined) {
  if (!value) return '';
  return labels[value] ?? value;
}

export function localizeMessage(message: string | null | undefined, fallback = '请求失败。') {
  if (!message) return fallback;
  return errorMessageZh[message] ?? message;
}

export function localizeError(error: LocalizableError | undefined, fallback = '请求失败。') {
  if (!error) return fallback;
  return localizeMessage(error.safeMessage, fallback);
}

export function localizeFieldErrors(error: LocalizableError | undefined) {
  return error?.fieldErrors?.map((fieldError) => ({
    ...fieldError,
    safeMessage: localizeMessage(fieldError.safeMessage, fieldError.code)
  })) ?? [];
}

const summaryKeyLabel: Record<string, string> = {
  traceability: '追溯',
  action: '动作',
  value: '值',
  before: '原值',
  after: '新值',
  status: '状态',
  beforeStatus: '原状态',
  fromStatus: '原状态',
  toStatus: '新状态',
  paymentStatus: '回款状态',
  stage: '阶段',
  fromStage: '原阶段',
  toStage: '新阶段',
  role: '角色',
  fromRole: '原角色',
  toRole: '新角色',
  actorId: '操作者',
  actorRole: '操作者角色',
  actorDisplay: '操作者',
  correlationId: '关联 ID',
  ownerId: '负责人',
  fromOwnerId: '原负责人',
  toOwnerId: '新负责人',
  oldOwnerId: '原负责人',
  newOwnerId: '新负责人',
  teamId: '团队',
  amount: '金额',
  dueAmount: '计划金额',
  paidAmount: '已回款',
  remainingAmount: '剩余金额',
  result: '结果',
  objectType: '对象类型',
  totalRows: '总行数',
  successCount: '成功行数',
  failureCount: '失败行数',
  includeArchived: '包含已归档',
  exportedCount: '导出行数',
  companyName: '公司名称',
  leadName: '线索名称',
  leadId: '线索',
  accountId: '客户',
  customerId: '客户',
  contactId: '联系人',
  contactIds: '联系人',
  opportunityId: '商机',
  quoteId: '报价',
  contractId: '合同',
  paymentId: '回款',
  source: '来源',
  title: '标题',
  reason: '原因',
  invalidReason: '无效原因',
  lostReason: '丢单原因',
  lostReasonCode: '丢单原因',
  signedEffectiveDate: '签署/生效日期',
  expectedSignedDate: '预计签署日期',
  expectedCloseDate: '预计关闭日期',
  closeDate: '关闭日期',
  dueDate: '到期日'
};

function summaryValueZh(key: string, value: unknown): string {
  if (typeof value === 'boolean') return value ? '是' : '否';
  if (Array.isArray(value)) return value.length === 0 ? '无' : value.map((item) => summaryValueZh(key, item)).join(', ');
  if (typeof value === 'object' && value !== null) return JSON.stringify(value);
  const text = String(value);
  if (key === 'action') return labelFor(actionLabel, text);
  if (key === 'objectType' || key.toLowerCase().endsWith('type')) return labelFor(objectTypeLabel, text);
  if (key.toLowerCase().includes('role')) return labelFor(roleLabel, text);
  if (key.toLowerCase().includes('status')) {
    return labelFor({ ...leadStatusLabel, ...quoteStatusLabel, ...contractStatusLabel, ...paymentStatusLabel, ...taskStatusLabel, ...userStatusLabel }, text);
  }
  if (key.toLowerCase().includes('stage')) return labelFor(opportunityStageLabel, text);
  if (key.toLowerCase().includes('reason')) return labelFor(lostReasonLabel, text);
  if (key === 'result') return labelFor(resultLabel, text);
  return localizeMessage(text, text);
}

export function summaryTextZh(summary: Record<string, unknown> | undefined) {
  if (!summary || Object.keys(summary).length === 0) return '无';
  return Object.entries(summary)
    .filter(([, value]) => value !== '' && value !== null && value !== undefined)
    .map(([key, value]) => `${labelFor(summaryKeyLabel, key)}: ${summaryValueZh(key, value)}`)
    .join(', ');
}
