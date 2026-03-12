package repo

import (
	"fmt"

	domain "main/internal/domain/models"
	"main/internal/infra/db/interfaces"
	"time"
)

type ReportRepositoryPostgresql struct {
	conn interfaces.ConnectionInterface `di.inject:"db"`
}

func (r *ReportRepositoryPostgresql) GetBurndownData(sprintId uint) ([]domain.BurndownDataPoint, error) {
	var rows []BurndownRow

	query := `
		WITH sprint_dates AS (
			SELECT d::date AS date
			FROM sprints s,
				generate_series(s.start_date::date, s.end_date::date, '1 day'::interval) d
			WHERE s.id = ?
		),
		latest_history AS (
			SELECT
				h.issue_binding_id,
				sd.date,
				h.estimate_dev,
				h.estimate_qa,
				h.bind_status,
				h.removed_at,
				ROW_NUMBER() OVER (
					PARTITION BY h.issue_binding_id, sd.date
					ORDER BY h.created_at DESC, h.id DESC
				) AS rn
			FROM sprint_dates sd
			JOIN issue_binding_histories h
				ON h.issue_binding_id IN (
					SELECT ib.id FROM issue_bindings ib WHERE ib.sprint_id = ?
				)
				AND h.created_at::date <= sd.date
		)
		SELECT
			lh.date,
			COALESCE(SUM(CASE WHEN lh.removed_at IS NULL OR lh.removed_at::date > lh.date THEN COALESCE(lh.estimate_dev, 0) ELSE 0 END), 0) AS total_scope,
			COALESCE(SUM(CASE WHEN (lh.removed_at IS NULL OR lh.removed_at::date > lh.date) AND lh.bind_status != 'done' THEN COALESCE(lh.estimate_dev, 0) ELSE 0 END), 0) AS remaining_dev,
			COALESCE(SUM(CASE WHEN lh.removed_at IS NULL OR lh.removed_at::date > lh.date THEN COALESCE(lh.estimate_qa, 0) ELSE 0 END), 0) AS total_scope_qa,
			COALESCE(SUM(CASE WHEN (lh.removed_at IS NULL OR lh.removed_at::date > lh.date) AND lh.bind_status != 'done' THEN COALESCE(lh.estimate_qa, 0) ELSE 0 END), 0) AS remaining_qa
		FROM latest_history lh
		WHERE lh.rn = 1
		GROUP BY lh.date
		ORDER BY lh.date
	`

	if err := r.conn.GetEngine().Raw(query, sprintId, sprintId).Scan(&rows).Error; err != nil {
		return nil, err
	}

	result := make([]domain.BurndownDataPoint, len(rows))
	for i, row := range rows {
		date, err := time.Parse("2006-01-02", row.Date)
		if err != nil {
			return nil, fmt.Errorf("cannot parse date: %w", err)
		}

		result[i] = domain.BurndownDataPoint{
			Date:         date,
			TotalScope:   row.TotalScope,
			RemainingDev: row.RemainingDev,
			TotalScopeQA: row.TotalScopeQA,
			RemainingQA:  row.RemainingQA,
		}
	}

	return result, nil
}

func (r *ReportRepositoryPostgresql) GetBurnupData(sprintId uint) ([]domain.BurnupDataPoint, error) {
	var rows []BurnupRow

	query := `
		WITH sprint_dates AS (
			SELECT d::date AS date
			FROM sprints s,
				generate_series(s.start_date::date, s.end_date::date, '1 day'::interval) d
			WHERE s.id = ?
		),
		latest_history AS (
			SELECT
				h.issue_binding_id,
				sd.date,
				h.estimate_dev,
				h.estimate_qa,
				h.bind_status,
				h.removed_at,
				ROW_NUMBER() OVER (
					PARTITION BY h.issue_binding_id, sd.date
					ORDER BY h.created_at DESC, h.id DESC
				) AS rn
			FROM sprint_dates sd
			JOIN issue_binding_histories h
				ON h.issue_binding_id IN (
					SELECT ib.id FROM issue_bindings ib WHERE ib.sprint_id = ?
				)
				AND h.created_at::date <= sd.date
		),
		daily_scope AS (
			SELECT
				lh.date,
				COALESCE(SUM(CASE WHEN lh.removed_at IS NULL OR lh.removed_at::date > lh.date THEN COALESCE(lh.estimate_dev, 0) ELSE 0 END), 0) AS scope_dev,
				COALESCE(SUM(CASE WHEN lh.removed_at IS NULL OR lh.removed_at::date > lh.date THEN COALESCE(lh.estimate_qa, 0) ELSE 0 END), 0) AS scope_qa
			FROM latest_history lh
			WHERE lh.rn = 1
			GROUP BY lh.date
		),
		completed_tasks AS (
			SELECT
				sd.date,
				COALESCE(SUM(CASE WHEN h.bind_status = 'done' AND h.created_at::date <= sd.date THEN COALESCE(h.estimate_dev, 0) ELSE 0 END), 0) AS completed_dev,
				COALESCE(SUM(CASE WHEN h.bind_status = 'done' AND h.created_at::date <= sd.date THEN COALESCE(h.estimate_qa, 0) ELSE 0 END), 0) AS completed_qa
			FROM sprint_dates sd
			LEFT JOIN issue_binding_histories h
				ON h.issue_binding_id IN (
					SELECT ib.id FROM issue_bindings ib WHERE ib.sprint_id = ?
				)
				AND h.bind_status = 'done'
			GROUP BY sd.date
		)
		SELECT
			ds.date,
			ds.scope_dev,
			ct.completed_dev,
			ds.scope_qa,
			ct.completed_qa
		FROM daily_scope ds
		JOIN completed_tasks ct ON ds.date = ct.date
		ORDER BY ds.date
	`

	if err := r.conn.GetEngine().Raw(query, sprintId, sprintId, sprintId).Scan(&rows).Error; err != nil {
		return nil, err
	}

	result := make([]domain.BurnupDataPoint, len(rows))
	for i, row := range rows {
		date, err := time.Parse("2006-01-02", row.Date)
		if err != nil {
			return nil, fmt.Errorf("cannot parse date: %w", err)
		}

		result[i] = domain.BurnupDataPoint{
			Date:         date,
			ScopeDev:     row.ScopeDev,
			CompletedDev: row.CompletedDev,
			ScopeQA:      row.ScopeQA,
			CompletedQA:  row.CompletedQA,
		}
	}

	return result, nil
}
