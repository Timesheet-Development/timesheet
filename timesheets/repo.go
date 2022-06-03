package timesheets

import (
	"context"
	"fmt"
	"timesheet/commons/res"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"
)

type Repository interface {
	InsertTimesheet(ctx context.Context, ts *Timesheet) (string, error)

	SelectTimesheetByLoginName(ctx context.Context, loginName string, month, year int) (bool, error)

	UpdateTimesheetByGivenCriteria(ctx context.Context, ts *Timesheet, loginName string, month, year int) (string, error)
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

func (repo *repository) SelectTimesheetByLoginName(ctx context.Context, loginName string, month, year int) (bool, error) {
	var err error
	var isExisting bool
	var count int

	selectQry := `select count(*) from timesheets t
	 where t.login_name = $1 and t."month" = $2 and t."year" = $3;`
	if err = pgxscan.Get(
		ctx, repo.db, &count, selectQry, loginName, month, year,
	); err != nil {
		// Handle query or rows processing error.
		if pgxscan.NotFound(err) {
			//return nil, &res.AppError{ResponseCode: UserDoesNotExist, Cause: err}
			//No error, but no user either
			return false, nil
		}
		return false, &res.AppError{ResponseCode: res.DatabaseError, Cause: err}
	}

	if count > 0 {
		isExisting = true
	}
	return isExisting, nil
}

func (repo *repository) UpdateTimesheetByGivenCriteria(ctx context.Context, ts *Timesheet, loginName string, month, year int) (string, error) {

	var err error
	var res string
	UpdateQry := `UPDATE public.timesheets
	SET placement=$1, info=$2, total_hours=$3, week_hours_info=$4
	WHERE login_name =$5 AND "year"=$6 AND "month"=$7;
	`
	if _, err = repo.db.Exec(ctx, UpdateQry, ts.Placement, ts.Info, ts.TotalHours, ts.WeekHrs,
		loginName, year, month); err != nil {
		return "", err
	}

	res = fmt.Sprintf("Updated Sucessfully with given criteria %s,%d,%d \n", loginName, month, year)
	return res, nil

}
