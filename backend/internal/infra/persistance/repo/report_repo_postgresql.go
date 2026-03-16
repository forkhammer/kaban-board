package repo

import (
	"fmt"
	"strings"

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
				ON h.sprint_id = ?
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
		date, err := time.Parse(time.RFC3339, row.Date)
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
				ON h.sprint_id = ?
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
				ON h.sprint_id = ?
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
		date, err := time.Parse(time.RFC3339, row.Date)
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

func wipBucketExprPostgresql(interval string) string {
	switch interval {
	case "week":
		return "DATE_TRUNC('week', date::date)::text"
	case "2weeks":
		return "FLOOR((date::date - '2000-01-03'::date) / (14 * 86400))::bigint"
	case "month":
		return "DATE_TRUNC('month', date::date)::text"
	default: // day
		return "date::text"
	}
}

func (r *ReportRepositoryPostgresql) GetWipData(startDate, endDate time.Time, interval string, teamId *uint, userId *uint) ([]domain.WipDataPoint, error) {
	var rows []WipRow

	var joinClause string
	var args []any
	args = append(args, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))

	joinClause = "INNER JOIN issue_binding_histories h ON h.created_at::date <= dr.date"
	if teamId != nil {
		joinClause += " INNER JOIN sprints s ON h.sprint_id = s.id AND s.team_id = ?"
		args = append(args, *teamId)
	}
	if userId != nil {
		joinClause += " AND h.assignee_id = ?"
		args = append(args, *userId)
	}

	bucket := wipBucketExprPostgresql(interval)

	query := strings.Join([]string{
		`WITH date_range AS (`,
		`    SELECT d::date AS date`,
		`    FROM generate_series(?::date, ?::date, '1 day'::interval) d`,
		`),`,
		`latest_history AS (`,
		`    SELECT h.issue_binding_id, dr.date, h.bind_status, h.removed_at,`,
		`        ROW_NUMBER() OVER (`,
		`            PARTITION BY h.issue_binding_id, dr.date`,
		`            ORDER BY h.created_at DESC, h.id DESC`,
		`        ) AS rn`,
		`    FROM date_range dr`,
		joinClause,
		`),`,
		`daily_wip AS (`,
		`    SELECT lh.date::text AS date,`,
		`        COALESCE(COUNT(CASE WHEN lh.bind_status = 'in_progress'`,
		`                    AND (lh.removed_at IS NULL OR lh.removed_at::date > lh.date)`,
		`              THEN 1 END), 0) AS wip_count`,
		`    FROM latest_history lh WHERE lh.rn = 1`,
		`    GROUP BY lh.date`,
		`)`,
		`SELECT MIN(date)::text AS date, ROUND(AVG(wip_count))::int AS wip_count`,
		`FROM daily_wip`,
		`GROUP BY ` + bucket,
		`ORDER BY MIN(date)`,
	}, "\n")

	if err := r.conn.GetEngine().Raw(query, args...).Scan(&rows).Error; err != nil {
		return nil, err
	}

	result := make([]domain.WipDataPoint, len(rows))
	for i, row := range rows {
		date, err := time.Parse("2006-01-02", row.Date)
		if err != nil {
			return nil, fmt.Errorf("cannot parse date: %w", err)
		}
		result[i] = domain.WipDataPoint{
			Date:     date,
			WipCount: row.WipCount,
		}
	}

	return result, nil
}
