package repo

import (
	"fmt"

	domain "main/internal/domain/models"
	"main/internal/infra/db/interfaces"
	"time"
)

type BurndownRow struct {
	Date         string `gorm:"column:date"`
	TotalScope   uint   `gorm:"column:total_scope"`
	RemainingDev uint   `gorm:"column:remaining_dev"`
	TotalScopeQA uint   `gorm:"column:total_scope_qa"`
	RemainingQA  uint   `gorm:"column:remaining_qa"`
}

type ReportRepository struct {
	conn interfaces.ConnectionInterface `di.inject:"db"`
}

func (r *ReportRepository) GetBurndownData(sprintId uint) ([]domain.BurndownDataPoint, error) {
	var rows []BurndownRow

	query := `
		WITH RECURSIVE sprint_dates AS (
			SELECT date(s.start_date) AS date, date(s.end_date) AS end_date
			FROM sprints s
			WHERE s.id = ?
			UNION ALL
			SELECT date(sd.date, '+1 day'), sd.end_date
			FROM sprint_dates sd
			WHERE sd.date < sd.end_date
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
				AND date(h.created_at) <= sd.date
		)
		SELECT
			lh.date,
			COALESCE(SUM(CASE WHEN lh.removed_at IS NULL OR date(lh.removed_at) > lh.date THEN COALESCE(lh.estimate_dev, 0) ELSE 0 END), 0) AS total_scope,
			COALESCE(SUM(CASE WHEN (lh.removed_at IS NULL OR date(lh.removed_at) > lh.date) AND lh.bind_status != 'done' THEN COALESCE(lh.estimate_dev, 0) ELSE 0 END), 0) AS remaining_dev,
			COALESCE(SUM(CASE WHEN lh.removed_at IS NULL OR date(lh.removed_at) > lh.date THEN COALESCE(lh.estimate_qa, 0) ELSE 0 END), 0) AS total_scope_qa,
			COALESCE(SUM(CASE WHEN (lh.removed_at IS NULL OR date(lh.removed_at) > lh.date) AND lh.bind_status != 'done' THEN COALESCE(lh.estimate_qa, 0) ELSE 0 END), 0) AS remaining_qa
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

type BurnupRow struct {
	Date         string `gorm:"column:date"`
	ScopeDev     uint   `gorm:"column:scope_dev"`
	CompletedDev uint   `gorm:"column:completed_dev"`
	ScopeQA      uint   `gorm:"column:scope_qa"`
	CompletedQA  uint   `gorm:"column:completed_qa"`
}

func (r *ReportRepository) GetBurnupData(sprintId uint) ([]domain.BurnupDataPoint, error) {
	var rows []BurnupRow

	query := `
		WITH RECURSIVE sprint_dates AS (
			SELECT date(s.start_date) AS date, date(s.end_date) AS end_date
			FROM sprints s
			WHERE s.id = ?
			UNION ALL
			SELECT date(sd.date, '+1 day'), sd.end_date
			FROM sprint_dates sd
			WHERE sd.date < sd.end_date
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
				AND date(h.created_at) <= sd.date
		),
		daily_scope AS (
			SELECT
				lh.date,
				COALESCE(SUM(CASE WHEN lh.removed_at IS NULL OR date(lh.removed_at) > lh.date THEN COALESCE(lh.estimate_dev, 0) ELSE 0 END), 0) AS scope_dev,
				COALESCE(SUM(CASE WHEN lh.removed_at IS NULL OR date(lh.removed_at) > lh.date THEN COALESCE(lh.estimate_qa, 0) ELSE 0 END), 0) AS scope_qa
			FROM latest_history lh
			WHERE lh.rn = 1
			GROUP BY lh.date
		),
		completed_tasks AS (
			SELECT
				sd.date,
				COALESCE(SUM(CASE WHEN h.bind_status = 'done' AND date(h.created_at) <= sd.date THEN COALESCE(h.estimate_dev, 0) ELSE 0 END), 0) AS completed_dev,
				COALESCE(SUM(CASE WHEN h.bind_status = 'done' AND date(h.created_at) <= sd.date THEN COALESCE(h.estimate_qa, 0) ELSE 0 END), 0) AS completed_qa
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
