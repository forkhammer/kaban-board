package repo

import (
	"fmt"
	"strings"

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
				ON h.sprint_id = ?
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
				ON h.sprint_id = ?
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

type SprintStatsRow struct {
	PlannedDev  uint `gorm:"column:planned_dev"`
	PlannedQA   uint `gorm:"column:planned_qa"`
	VelocityDev uint `gorm:"column:velocity_dev"`
	VelocityQA  uint `gorm:"column:velocity_qa"`
}

func (r *ReportRepository) GetSprintStats(sprintId uint, assigneeId *uint, groupId *uint) (*domain.SprintStats, error) {
	var row SprintStatsRow
	query := `
		SELECT
			COALESCE(SUM(ib.estimate_dev), 0) AS planned_dev,
			COALESCE(SUM(ib.estimate_qa), 0)  AS planned_qa,
			COALESCE(SUM(CASE WHEN ib.bind_status = 'done' THEN ib.estimate_dev ELSE 0 END), 0) AS velocity_dev,
			COALESCE(SUM(CASE WHEN ib.bind_status = 'done' THEN ib.estimate_qa  ELSE 0 END), 0) AS velocity_qa
		FROM issue_bindings ib
		WHERE ib.sprint_id = ?
			AND ib.deleted_at IS NULL`
	args := []any{sprintId}
	if assigneeId != nil {
		query += " AND ib.assignee_id = ?"
		args = append(args, *assigneeId)
	}
	if groupId != nil {
		query += " AND ib.id IN (SELECT ib2.id FROM issue_bindings ib2 JOIN user_groups ug ON ug.user_id = ib2.assignee_id WHERE ug.group_id = ?)"
		args = append(args, *groupId)
	}
	if err := r.conn.GetEngine().Raw(query, args...).Scan(&row).Error; err != nil {
		return nil, err
	}
	return &domain.SprintStats{
		PlannedDev:  row.PlannedDev,
		PlannedQA:   row.PlannedQA,
		VelocityDev: row.VelocityDev,
		VelocityQA:  row.VelocityQA,
	}, nil
}

func (r *ReportRepository) GetActiveUserIdsByTeam(teamId uint, groupId *uint) ([]uint, error) {
	var ids []uint
	query := `
		SELECT DISTINCT u.id
		FROM users u
		JOIN user_groups ug ON u.id = ug.user_id
		JOIN team_groups tg ON tg.group_id = ug.group_id
		WHERE tg.team_id = ? AND u.is_active = true`
	args := []any{teamId}
	if groupId != nil {
		query += " AND ug.group_id = ?"
		args = append(args, *groupId)
	}
	if err := r.conn.GetEngine().Raw(query, args...).Scan(&ids).Error; err != nil {
		return nil, err
	}
	return ids, nil
}

type WipRow struct {
	Date     string `gorm:"column:date"`
	WipCount int    `gorm:"column:wip_count"`
}

func wipBucketExprSQLite(interval string) string {
	switch interval {
	case "week":
		return "strftime('%Y-%W', date)"
	case "2weeks":
		return "CAST((julianday(date) - julianday('2000-01-03')) / 14 AS INTEGER)"
	case "month":
		return "strftime('%Y-%m', date)"
	default: // day
		return "date"
	}
}

func (r *ReportRepository) GetWipData(startDate, endDate time.Time, interval string, teamId *uint, userId *uint) ([]domain.WipDataPoint, error) {
	var rows []WipRow

	var joinClause string
	var args []any
	args = append(args, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))

	joinClause = "JOIN issue_binding_histories h ON date(h.created_at) <= dr.date"
	if teamId != nil {
		joinClause += " JOIN sprints s ON h.sprint_id = s.id AND s.team_id = ?"
		args = append(args, *teamId)
	}
	if userId != nil {
		joinClause += " AND h.assignee_id = ?"
		args = append(args, *userId)
	}

	bucket := wipBucketExprSQLite(interval)

	query := strings.Join([]string{
		`WITH RECURSIVE date_range AS (`,
		`    SELECT date(?) AS date, date(?) AS end_date`,
		`    UNION ALL`,
		`    SELECT date(dr.date, '+1 day'), dr.end_date`,
		`    FROM date_range dr WHERE dr.date < dr.end_date`,
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
		`    SELECT lh.date,`,
		`        COALESCE(COUNT(CASE WHEN lh.bind_status = 'in_progress'`,
		`                    AND (lh.removed_at IS NULL OR date(lh.removed_at) > lh.date)`,
		`              THEN 1 END), 0) AS wip_count`,
		`    FROM latest_history lh WHERE lh.rn = 1`,
		`    GROUP BY lh.date`,
		`)`,
		`SELECT MIN(date) AS date, CAST(ROUND(AVG(wip_count)) AS INTEGER) AS wip_count`,
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
