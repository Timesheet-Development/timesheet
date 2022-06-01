package timesheets

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"
)

type Repository interface {
	InsertTimesheet(ctx context.Context, ts *Timesheet) (string, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) InsertTimesheet(ctx context.Context, ts *Timesheet) (string, error) {
	var err error
	var loginName string
	insertTimesheetQry := `INSERT INTO public.timesheets
						(id, status, placement, info, total_hours, "month", "year", 
						week_hours_info, week_day_info, login_name)
						VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);`

	if _, err = r.db.Exec(ctx, insertTimesheetQry, ts.ID, ts.Status, ts.Placement,
		ts.Info, ts.TotalHours, ts.Month, ts.Year, ts.WeekHrs, ts.WeekDay, ts.LoginName); err != nil {
		log.Error().Err(err).Str("loginName", ts.LoginName).Msg("Error while inserting the timesheet data")
		return "", err
	}
	loginName = ts.LoginName
	return loginName, nil
}
