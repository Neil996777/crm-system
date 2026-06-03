package repo

import (
	"context"
	"database/sql"
	"fmt"
)

type ProjectionRepo struct {
	db *sql.DB
}

func NewProjectionRepo(db *sql.DB) *ProjectionRepo {
	return &ProjectionRepo{db: db}
}

type OverviewMetrics struct {
	LeadCount        int
	OpportunityCount int
	TaskCount        int
	WonCount         int
	LostCount        int
	QuoteAmount      string
	ContractAmount   string
	PaidAmount       string
	ReceivableAmount string
}

type GroupRow struct {
	Key    string `json:"key"`
	Label  string `json:"label"`
	Count  int    `json:"count"`
	Amount string `json:"amount"`
}

type PaymentGroupRow struct {
	Key        string `json:"key"`
	Label      string `json:"label"`
	Count      int    `json:"count"`
	Amount     string `json:"amount"`
	DueAmount  string `json:"dueAmount"`
	PaidAmount string `json:"paidAmount"`
}

type OverviewData struct {
	Metrics  OverviewMetrics
	Pipeline []GroupRow
}

type ReportBreakdowns struct {
	LeadsByStatus        []GroupRow        `json:"leadsByStatus"`
	OpportunitiesByStage []GroupRow        `json:"opportunitiesByStage"`
	QuotesByStatus       []GroupRow        `json:"quotesByStatus"`
	ContractsByStatus    []GroupRow        `json:"contractsByStatus"`
	PaymentsByStatus     []PaymentGroupRow `json:"paymentsByStatus"`
}

type SalesReportData struct {
	Metrics    OverviewMetrics
	Breakdowns ReportBreakdowns
	Groups     []GroupRow
}

func (r *ProjectionRepo) TeamOverview(ctx context.Context, teamID string) (OverviewData, error) {
	var data OverviewData
	data.Pipeline = []GroupRow{}
	err := r.db.QueryRowContext(ctx, `
		SELECT
		  count(*) FILTER (WHERE record_type = 'lead')::int,
		  count(*) FILTER (WHERE record_type = 'opportunity')::int,
		  count(*) FILTER (WHERE record_type = 'task')::int,
		  count(*) FILTER (WHERE record_type = 'opportunity' AND stage = 'Won')::int,
		  count(*) FILTER (WHERE record_type = 'opportunity' AND stage = 'Lost')::int,
		  to_char(COALESCE(sum(amount) FILTER (WHERE record_type = 'quote'), 0), 'FM999999999999990.00'),
		  to_char(COALESCE(sum(amount) FILTER (WHERE record_type = 'contract'), 0), 'FM999999999999990.00'),
		  to_char(COALESCE(sum(amount) FILTER (WHERE record_type = 'payment'), 0), 'FM999999999999990.00'),
		  to_char(COALESCE(sum(amount) FILTER (WHERE record_type = 'contract'), 0) - COALESCE(sum(amount) FILTER (WHERE record_type = 'payment'), 0), 'FM999999999999990.00')
		FROM reporting.record_projections
		WHERE team_id = $1 AND archived_at IS NULL
	`, teamID).Scan(
		&data.Metrics.LeadCount,
		&data.Metrics.OpportunityCount,
		&data.Metrics.TaskCount,
		&data.Metrics.WonCount,
		&data.Metrics.LostCount,
		&data.Metrics.QuoteAmount,
		&data.Metrics.ContractAmount,
		&data.Metrics.PaidAmount,
		&data.Metrics.ReceivableAmount,
	)
	if err != nil {
		return OverviewData{}, err
	}
	rows, err := r.db.QueryContext(ctx, `
		SELECT stage, count(*)::int, to_char(COALESCE(sum(amount), 0), 'FM999999999999990.00')
		FROM reporting.record_projections
		WHERE team_id = $1 AND archived_at IS NULL AND record_type = 'opportunity'
		GROUP BY stage
		ORDER BY stage ASC
	`, teamID)
	if err != nil {
		return OverviewData{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var row GroupRow
		if err := rows.Scan(&row.Key, &row.Count, &row.Amount); err != nil {
			return OverviewData{}, err
		}
		row.Label = row.Key
		data.Pipeline = append(data.Pipeline, row)
	}
	if err := rows.Err(); err != nil {
		return OverviewData{}, err
	}
	return data, nil
}

func (r *ProjectionRepo) SalesOverview(ctx context.Context, teamID string, allTeams bool) (SalesReportData, error) {
	data := SalesReportData{
		Breakdowns: ReportBreakdowns{
			LeadsByStatus:        []GroupRow{},
			OpportunitiesByStage: []GroupRow{},
			QuotesByStatus:       []GroupRow{},
			ContractsByStatus:    []GroupRow{},
			PaymentsByStatus:     []PaymentGroupRow{},
		},
		Groups: []GroupRow{},
	}
	err := r.db.QueryRowContext(ctx, `
		SELECT
		  count(*) FILTER (WHERE record_type = 'lead')::int,
		  count(*) FILTER (WHERE record_type = 'opportunity')::int,
		  count(*) FILTER (WHERE record_type = 'task')::int,
		  count(*) FILTER (WHERE record_type = 'opportunity' AND stage = 'Won')::int,
		  count(*) FILTER (WHERE record_type = 'opportunity' AND stage = 'Lost')::int,
		  to_char(COALESCE(sum(amount) FILTER (WHERE record_type = 'quote'), 0), 'FM999999999999990.00'),
		  to_char(COALESCE(sum(amount) FILTER (WHERE record_type = 'contract'), 0), 'FM999999999999990.00'),
		  to_char(COALESCE(sum(amount) FILTER (WHERE record_type = 'payment' AND lower(status) IN ('paid', 'partiallypaid', 'partially paid')), 0), 'FM999999999999990.00'),
		  to_char(COALESCE(sum(amount) FILTER (WHERE record_type = 'contract'), 0) - COALESCE(sum(amount) FILTER (WHERE record_type = 'payment' AND lower(status) IN ('paid', 'partiallypaid', 'partially paid')), 0), 'FM999999999999990.00')
		FROM reporting.record_projections
		WHERE archived_at IS NULL AND ($2 OR team_id = $1)
	`, teamID, allTeams).Scan(
		&data.Metrics.LeadCount,
		&data.Metrics.OpportunityCount,
		&data.Metrics.TaskCount,
		&data.Metrics.WonCount,
		&data.Metrics.LostCount,
		&data.Metrics.QuoteAmount,
		&data.Metrics.ContractAmount,
		&data.Metrics.PaidAmount,
		&data.Metrics.ReceivableAmount,
	)
	if err != nil {
		return SalesReportData{}, err
	}
	if data.Breakdowns.LeadsByStatus, err = r.groupRows(ctx, teamID, allTeams, "lead", "status", false); err != nil {
		return SalesReportData{}, err
	}
	if data.Breakdowns.OpportunitiesByStage, err = r.groupRows(ctx, teamID, allTeams, "opportunity", "stage", true); err != nil {
		return SalesReportData{}, err
	}
	if data.Breakdowns.QuotesByStatus, err = r.groupRows(ctx, teamID, allTeams, "quote", "status", true); err != nil {
		return SalesReportData{}, err
	}
	if data.Breakdowns.ContractsByStatus, err = r.groupRows(ctx, teamID, allTeams, "contract", "status", true); err != nil {
		return SalesReportData{}, err
	}
	if data.Breakdowns.PaymentsByStatus, err = r.paymentRows(ctx, teamID, allTeams); err != nil {
		return SalesReportData{}, err
	}
	return data, nil
}

func (r *ProjectionRepo) groupRows(ctx context.Context, teamID string, allTeams bool, recordType string, dimension string, includeAmount bool) ([]GroupRow, error) {
	amountExpr := "'0.00'"
	if includeAmount {
		amountExpr = "to_char(COALESCE(sum(amount), 0), 'FM999999999999990.00')"
	}
	rows, err := r.db.QueryContext(ctx, fmt.Sprintf(`
		SELECT COALESCE(%s, ''), count(*)::int, %s
		FROM reporting.record_projections
		WHERE archived_at IS NULL AND ($2 OR team_id = $1) AND record_type = $3
		GROUP BY COALESCE(%s, '')
		ORDER BY COALESCE(%s, '') ASC
	`, dimension, amountExpr, dimension, dimension), teamID, allTeams, recordType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []GroupRow{}
	for rows.Next() {
		var row GroupRow
		if err := rows.Scan(&row.Key, &row.Count, &row.Amount); err != nil {
			return nil, err
		}
		row.Label = row.Key
		result = append(result, row)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *ProjectionRepo) paymentRows(ctx context.Context, teamID string, allTeams bool) ([]PaymentGroupRow, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT
		  COALESCE(status, ''),
		  count(*)::int,
		  to_char(COALESCE(sum(amount), 0), 'FM999999999999990.00'),
		  to_char(COALESCE(sum(amount) FILTER (WHERE lower(status) NOT IN ('paid', 'partiallypaid', 'partially paid')), 0), 'FM999999999999990.00'),
		  to_char(COALESCE(sum(amount) FILTER (WHERE lower(status) IN ('paid', 'partiallypaid', 'partially paid')), 0), 'FM999999999999990.00')
		FROM reporting.record_projections
		WHERE archived_at IS NULL AND ($2 OR team_id = $1) AND record_type = 'payment'
		GROUP BY COALESCE(status, '')
		ORDER BY COALESCE(status, '') ASC
	`, teamID, allTeams)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []PaymentGroupRow{}
	for rows.Next() {
		var row PaymentGroupRow
		if err := rows.Scan(&row.Key, &row.Count, &row.Amount, &row.DueAmount, &row.PaidAmount); err != nil {
			return nil, err
		}
		row.Label = row.Key
		result = append(result, row)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (m OverviewMetrics) Empty() bool {
	return m.LeadCount == 0 &&
		m.OpportunityCount == 0 &&
		m.TaskCount == 0 &&
		m.WonCount == 0 &&
		m.LostCount == 0 &&
		m.QuoteAmount == "0.00" &&
		m.ContractAmount == "0.00" &&
		m.PaidAmount == "0.00"
}

func (m OverviewMetrics) Map() map[string]any {
	return map[string]any{
		"leadCount":        m.LeadCount,
		"opportunityCount": m.OpportunityCount,
		"taskCount":        m.TaskCount,
		"wonCount":         m.WonCount,
		"lostCount":        m.LostCount,
		"quoteAmount":      m.QuoteAmount,
		"contractAmount":   m.ContractAmount,
		"paidAmount":       m.PaidAmount,
		"receivableAmount": m.ReceivableAmount,
	}
}

func ProjectionTableName() string {
	return fmt.Sprintf("%s.%s", "reporting", "record_projections")
}
