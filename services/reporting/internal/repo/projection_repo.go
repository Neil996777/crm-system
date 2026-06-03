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

type OverviewData struct {
	Metrics  OverviewMetrics
	Pipeline []GroupRow
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
