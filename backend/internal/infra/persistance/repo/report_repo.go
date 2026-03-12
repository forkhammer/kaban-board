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
			COALESCE(SUM(CASE WHEN (lh.removed_at IS NULL OR date(lh.removed_at) > lh.date) AND lh.bind_status != 'done' THEN COALESCE(lh.estimate_dev, 0) ELSE 0 END), 0) AS remaining_dev
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
		}
	}

	return result, nil
}
