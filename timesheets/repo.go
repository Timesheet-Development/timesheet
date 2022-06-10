package timesheets

import (
	"context"
	"fmt"
	"timesheet/commons/res"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"
)

type Repository interface {
	InsertTimesheet(ctx context.Context, ts *Timesheet) (string, error)

	SelectTimesheetByLoginName(ctx context.Context, loginName string, month, year int) (bool, error)

	UpdateTimesheetByGivenCriteria(ctx context.Context, ts *Timesheet, loginName string, month, year int) (string, error)

	SelectAllTimesheetByLoginName(ctx context.Context, loginName string) ([]*GetAllTimesheets, error)

	SelectTimesheetByWeek(ctx context.Context, loginName string, week, month, year int) (*GetAllTimesheets, error)

	DeleteTimesheet(ctx context.Context, loginName string, month, year int) (string, error)

	SelectTimesheetUUID(ctx context.Context, loginName string, month, year int) (string, error)

	UpsertTimesheetNotes(ctx context.Context, notes *AddorUpdateNotes, uuids string) (string, error)
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

	res = fmt.Sprintf("Updated Sucessfully with given criteria %s,%d,%d", loginName, month, year)
	return res, nil

}

func (repo *repository) SelectAllTimesheetByLoginName(ctx context.Context, loginName string) ([]*GetAllTimesheets, error) {
	var err error
	var rows pgx.Rows
	tsArr := []*GetAllTimesheets{}

	selectQry := `select login_name,placement,info,"month","year",total_hours,status,week_hours_info,week_day_info from timesheets t 
				  where t.login_name = $1;`

	rows, err = repo.db.Query(ctx, selectQry, loginName)
	if err != nil {
		log.Error().Err(err).Str("loginName", loginName).Msg("Error while fetching the timesheet data")
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		ts := &GetAllTimesheets{}
		err = rows.Scan(&ts.LoginName, &ts.Placement, &ts.Info, &ts.Month, &ts.Year, &ts.TotalHours,
			&ts.Status, &ts.WeekHrs, &ts.WeekDay)
		if err != nil {
			log.Error().Err(err).Str("loginName", loginName).Msg("Error while scaning each field from the timesheet")
			return nil, err
		}

		tsArr = append(tsArr, ts)
	}
	log.Info().Str("loginName", loginName).Msg("Successfully return timesheet info")

	return tsArr, nil
}

func (repo *repository) SelectTimesheetByWeek(ctx context.Context, loginName string, week, month, year int) (*GetAllTimesheets, error) {
	var err error
	ts := &GetAllTimesheets{}

	selectQry := `select login_name,placement,info,"month","year",total_hours,status,week_hours_info,week_day_info from timesheets t 
				  where t.login_name = $1
				  and t."month" = $2
				  and t."year" = $3;`

	if err = pgxscan.Get(
		ctx, repo.db, ts, selectQry, loginName, month, year,
	); err != nil {
		// Handle query or rows processing error.
		if pgxscan.NotFound(err) {
			//return nil, &res.AppError{ResponseCode: UserDoesNotExist, Cause: err}
			//No error, but no user either
			return nil, nil
		}
		return nil, &res.AppError{ResponseCode: res.DatabaseError, Cause: err}
	}

	return ts, nil
}

func (repo *repository) DeleteTimesheet(ctx context.Context, loginName string, month, year int) (string, error) {
	var err error
	var response string

	deletQry := `delete from timesheets t
				where t.login_name = $1
				and t.month = $2
				and t.year = $3`
	if _, err = repo.db.Exec(ctx, deletQry, loginName, month, year); err != nil {
		log.Error().Err(err).Str("loginName", loginName).Msg("Error while deleting the data")
		return "", err
	}
	response = fmt.Sprintf("Successfully delete the record for the given criteria %s %d %d", loginName, month, year)

	return response, nil
}

func (repo *repository) SelectTimesheetUUID(ctx context.Context, loginName string, month, year int) (string, error) {
	var err error
	var uuid string
	selectQry := `select id from timesheets where login_name = $1 and month = $2 and year = $3;`

	err = repo.db.QueryRow(ctx, selectQry, loginName, month, year).Scan(&uuid)
	if err != nil {
		return "", nil
	}

	log.Info().Msgf("UUID %s", uuid)
	return uuid, nil
}

func (repo *repository) UpsertTimesheetNotes(ctx context.Context, notes *AddorUpdateNotes, uuids string) (string, error) {

	var err error
	var res string

	if uuids == "" {
		uuids = uuid.New().String()
	}

	upsertQry := `insert into timesheets(id,login_name,month,year,info) values($1,$2,$3,$4,$5)
				 on conflict(id,month,year)
				 do update set info = excluded.info`

	if _, err = repo.db.Exec(ctx, upsertQry, uuids, notes.LoginName, notes.Month, notes.Year, notes.Info); err != nil {
		return "", err
	}

	res = "Added notes successfully"
	return res, nil
}
