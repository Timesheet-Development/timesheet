package timesheets

import (
	"context"
	"fmt"
	"time"
	"timesheet/commons/res"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"
)

type Repository interface {
	InsertTimesheet(ctx context.Context, ts *Timesheet, wArr []*WeekHrs, wDayArr []*WeekDay) (string, error)

	SelectTimesheetByLoginName(ctx context.Context, loginName string, month, year int) (*Timesheet, error)

	UpdateTimesheetByGivenCriteria(ctx context.Context, ts *Timesheet, loginName string, month, year int, wArr []*WeekHrs) (string, error)

	SelectAllTimesheetByLoginName(ctx context.Context, loginName string) ([]*GetAllTimesheets, error)

	SelectTimesheetByGivenCriteria(ctx context.Context, loginName string, month, year int) (*GetAllTimesheets, error)

	SelectWeekHoursByWeek(ctx context.Context, loginName string, week, month, year int) (*WeekHrs, error)

	DeleteTimesheet(ctx context.Context, loginName string, month, year int) (string, error)

	SelectTimesheetUUID(ctx context.Context, loginName string, month, year int) (string, error)

	UpsertTimesheetNotes(ctx context.Context, notes *Notes) (string, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) InsertTimesheet(ctx context.Context, ts *Timesheet, wArr []*WeekHrs, wDayArr []*WeekDay) (string, error) {
	var err error
	var loginName string
	tx, _ := r.db.Begin(ctx)

	insertTimesheetQry := `INSERT INTO public.timesheets
						(id, status, placement, info, total_hours, "month", "year", 
						week_hours_info, week_day_info, login_name)
						VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);`

	if _, err = tx.Exec(ctx, insertTimesheetQry, ts.ID, ts.Status, ts.Placement,
		ts.Info, ts.TotalHours, ts.Month, ts.Year, ts.WeekHrs, ts.WeekDay, ts.LoginName); err != nil {
		tx.Rollback(ctx)
		log.Error().Err(err).Str("loginName", ts.LoginName).Msg("Error while inserting the timesheet data")
		return "", err
	}

	insertIntoWeekHrs := `INSERT INTO public.week_hours_info
	(login_name, "month", "year", week_info, day1, day2, day3, day4, day5)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9);
	`
	for _, weekHrs := range wArr {
		if _, err = tx.Exec(ctx, insertIntoWeekHrs, ts.LoginName, ts.Month, ts.Year,
			weekHrs.WeekInfo, weekHrs.Day1, weekHrs.Day2, weekHrs.Day3, weekHrs.Day4, weekHrs.Day5); err != nil {
			tx.Rollback(ctx)
			log.Error().Err(err).Str("login_name", ts.LoginName).Msg("Error while insering into week hours info table")
			return "", err
		}
	}

	insertIntoWeekDay := `INSERT INTO public.week_day_info
	(login_name, "month", "year", week_info, day1, day2, day3, day4, day5)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9);
	
	`
	for _, weekDay := range wDayArr {
		if _, err = tx.Exec(ctx, insertIntoWeekDay, ts.LoginName, ts.Month, ts.Year,
			weekDay.WeekInfo, weekDay.Day1, weekDay.Day2, weekDay.Day3, weekDay.Day4, weekDay.Day5); err != nil {
			tx.Rollback(ctx)
			log.Error().Err(err).Str("login_name", ts.LoginName).Msg("Error while insering into week hours info table")
			return "", err
		}
	}

	tx.Commit(ctx)

	loginName = ts.LoginName
	return loginName, nil
}

func (repo *repository) SelectTimesheetByLoginName(ctx context.Context, loginName string, month, year int) (*Timesheet, error) {
	var err error

	ts := &Timesheet{}

	selectQry := `select * from timesheets
	 where login_name=$1 and month=$2 and year=$3;`

	if err = repo.db.QueryRow(ctx, selectQry, loginName, month, year).Scan(&ts.ID, &ts.LoginName,
		&ts.Status, &ts.Placement, &ts.Info, &ts.TotalHours, &ts.Month, &ts.Year, &ts.WeekHrs, &ts.WeekDay, &ts.CreatedAt, &ts.UpdatedAt); err != nil {
		log.Error().Err(err).Msg("Error while fetching timesheet information")
	}

	return ts, nil
}

func (repo *repository) UpdateTimesheetByGivenCriteria(ctx context.Context, ts *Timesheet, loginName string, month, year int, wArr []*WeekHrs) (string, error) {

	tx, _ := repo.db.Begin(ctx)
	var err error
	var res string
	UpdateQry := `UPDATE public.timesheets
	SET placement=$1, info=$2, total_hours=$3, week_hours_info=$4
	WHERE login_name=$5 AND "year"=$6 AND "month"=$7;
	`
	if _, err = tx.Exec(ctx, UpdateQry, ts.Placement, ts.Info, ts.TotalHours, ts.WeekHrs,
		loginName, year, month); err != nil {
		tx.Rollback(ctx)
		return "", err
	}

	updateWeekHRS := `UPDATE public.week_hours_info
	SET  day1=$1, day2=$2, day3=$3, day4=$4, day5=$5, updated_at=$6
	WHERE login_name=$7 AND "month"=$8 AND "year"=$9 AND week_info=$10;
	`

	for _, weekhrs := range wArr {
		if _, err = tx.Exec(ctx, updateWeekHRS, weekhrs.Day1, weekhrs.Day2, weekhrs.Day3, weekhrs.Day4,
			weekhrs.Day5, time.Now(), loginName, month, year, weekhrs.WeekInfo); err != nil {
			tx.Rollback(ctx)
			log.Error().Err(err).Msg("Error while updating week hrs info")
			return "", err
		}
	}
	tx.Commit(ctx)
	res = fmt.Sprintf("Updated Sucessfully with given criteria %s,%d,%d", loginName, month, year)
	return res, nil

}

func (repo *repository) SelectAllTimesheetByLoginName(ctx context.Context, loginName string) ([]*GetAllTimesheets, error) {
	var err error
	var rows pgx.Rows
	tsArr := []*GetAllTimesheets{}

	selectQry := `select login_name,placement,info,"month","year",total_hours,status,
				  week_hours_info,week_day_info from timesheets t 
				  where t.login_name = $1
				  order by created_at desc;`

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

func (repo *repository) SelectTimesheetByGivenCriteria(ctx context.Context, loginName string, month, year int) (*GetAllTimesheets, error) {
	var err error
	ts := &GetAllTimesheets{}

	selectQry := `select login_name,placement,info,"month","year",total_hours,
				  status,week_hours_info,week_day_info from timesheets t 
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

func (repo *repository) SelectWeekHoursByWeek(ctx context.Context, loginName string, week, month, year int) (*WeekHrs, error) {

	var err error
	w := &WeekHrs{}

	selectWeekInfo := `select week_info ,day1 ,day2 ,day3 ,day4 ,day5  from week_hours_info whi 
	where login_name = $1 and "month" = $2 and "year" = $3
	and week_info = $4;`

	if err = pgxscan.Get(
		ctx, repo.db, w, selectWeekInfo, loginName, month, year, week,
	); err != nil {
		// Handle query or rows processing error.
		if pgxscan.NotFound(err) {
			//return nil, &res.AppError{ResponseCode: UserDoesNotExist, Cause: err}
			//No error, but no user either
			return nil, nil
		}
		return nil, &res.AppError{ResponseCode: res.DatabaseError, Cause: err}
	}
	return w, nil
}

func (repo *repository) DeleteTimesheet(ctx context.Context, loginName string, month, year int) (string, error) {
	var err error
	var response string

	tx, _ := repo.db.Begin(ctx)

	deletQry := `delete from timesheets t
				where t.login_name = $1
				and t.month = $2
				and t.year = $3`
	if _, err = tx.Exec(ctx, deletQry, loginName, month, year); err != nil {
		log.Error().Err(err).Str("loginName", loginName).Msg("Error while deleting the data")
		tx.Rollback(ctx)
		return "", err
	}

	deleteWeekHrs := `delete from week_hours_info whi 
					  where whi.login_name = $1
					  and whi.month = $2
					  and whi.year = $3;`
	if _, err = tx.Exec(ctx, deleteWeekHrs, loginName, month, year); err != nil {
		log.Error().Err(err).Str("loginName", loginName).Msg("Error while deleting the data")
		tx.Rollback(ctx)
		return "", err
	}

	deleteWeekDay := `delete from week_day_info wdi 
	where wdi.login_name = $1
	and wdi.month = $2
	and wdi.year = $3;`
	if _, err = tx.Exec(ctx, deleteWeekDay, loginName, month, year); err != nil {
		log.Error().Err(err).Str("loginName", loginName).Msg("Error while deleting the data")
		tx.Rollback(ctx)
		return "", err
	}

	tx.Commit(ctx)

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

func (repo *repository) UpsertTimesheetNotes(ctx context.Context, notes *Notes) (string, error) {

	var err error
	var res string

	uuid := uuid.New()

	upsertQry := `insert into timesheets(id,login_name,month,year,info) values($1,$2,$3,$4,$5)
				 on conflict(login_name,month,year)
				 do update set info = excluded.info`

	if _, err = repo.db.Exec(ctx, upsertQry, uuid, notes.LoginName, notes.Month, notes.Year, notes.Info); err != nil {
		return "", err
	}

	res = "Added notes successfully"
	return res, nil
}
