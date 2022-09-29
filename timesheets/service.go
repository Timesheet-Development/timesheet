package timesheets

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"timesheet/user"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type Service interface {
	CreateTimesheet(ctx context.Context, ts *Timesheet) (string, error)

	UpdateTimesheet(ctx context.Context, ts *Timesheet, loginName string, month, year int) (string, error)

	GetListofTimesheets(ctx context.Context, loginName string) ([]*GetAllTimesheets, error)

	GetTimesheetsByWeek(ctx context.Context, loginName string, week, month, year int) (*GetTimesheet, error)

	DeleteTimesheet(ctx context.Context, loginName string, month, year int) (string, error)

	AddorUpdatenotes(ctx context.Context, notes *Notes) (string, error)

	ListSubmittedTimesheets(ctx context.Context, status string) ([]*GetAllTimesheets, error)

	//GenerateCSV(ctx context.Context, req []*GetTimesheet) (string, error)
}

type service struct {
	repo     Repository
	userRepo user.Repository
}

func NewService(repo Repository, userRepo user.Repository) Service {
	return &service{repo: repo,
		userRepo: userRepo}
}

func (s *service) CreateTimesheet(ctx context.Context, ts *Timesheet) (string, error) {
	var err error
	var loginName string
	userRes := &user.User{}

	if ts.LoginName == "" {
		err = errors.New("loginName is empty")
		return "", err
	}

	//Transforming fields
	ts.LoginName = strings.ToUpper(ts.LoginName)

	userRes, err = s.userRepo.SelectUserByLoginName(ctx, ts.LoginName)
	if err != nil {
		log.Error().Err(err).Str("loginName", ts.LoginName).Msg("User details not found for the given loginName")
		return "", err
	}

	log.Info().Str("loginName", userRes.LoginName).Msg("logging the timesheet info")

	if userRes.UserType == user.UserTypeApprover || userRes.UserType == user.UserTypeEmployee || userRes.UserType == user.UserTypeManager {

		var sb strings.Builder
		sb.WriteString(userRes.Department)
		sb.WriteString(" ")
		sb.WriteString(userRes.JobTitle)

		ts.ID = uuid.New()

		ts.LoginName = userRes.LoginName
		ts.Placement = sb.String()
		ts.Status = string(timesheetStatusSubmitted)

		//Unmarshalling the week hrs json
		wArr := []*WeekHrs{}

		if err = json.Unmarshal(ts.WeekHrs, &wArr); err != nil {
			log.Error().Err(err).Msg("Error while unmarshalling week hrs json")
		}

		//Unmarshalling week day json
		wDayArr := []*WeekDay{}
		if err = json.Unmarshal(ts.WeekDay, &wDayArr); err != nil {
			log.Error().Err(err).Msg("Error while unmarshalling week day json")
		}

		for _, eachDayHrs := range wArr {
			ts.TotalHours += eachDayHrs.Day1 + eachDayHrs.Day2 + eachDayHrs.Day3 + eachDayHrs.Day4 + eachDayHrs.Day5
		}

		if loginName, err = s.repo.InsertTimesheet(ctx, ts, wArr, wDayArr); err != nil {
			log.Error().Err(err).Str("loginName", loginName).Msg("Error while calling repo in timesheet service")
			return "", err
		}
	} else {
		log.Error().Err(err).Str("loginName", loginName).Msg("Operator is not allowed to log timesheets")
		err = errors.New("Operator is not allowed to log timesheets")
		return "", err
	}

	return loginName, nil
}

func (s *service) UpdateTimesheet(ctx context.Context, ts *Timesheet, loginName string, month, year int) (string, error) {
	/*
		Step1: Check  loginName,month,Year are not empty.
		Step2 : check whether timesheets is having a record with this login
		Step3: Call repo from service.
	*/
	var err error
	var tsResponse *Timesheet
	var res string

	if loginName == "" {
		err = errors.New("loginName is empty")
		return "", err
	}

	if month == 0 {
		err = errors.New("month cannot be zero")
		return "", err
	}

	if year == 0 {
		err = errors.New("year cannot be zero")
		return "", err
	}

	loginName = strings.ToUpper(loginName)

	tsResponse, err = s.repo.SelectTimesheetByLoginName(ctx, loginName, month, year)
	if err != nil {
		log.Error().Err(err).Msgf("Error while fetching Timesheet with given LoginName")
		return "", err
	}

	if tsResponse != nil {
		//Unmarshalling the week hrs json
		wArr := []*WeekHrs{}

		if err = json.Unmarshal(ts.WeekHrs, &wArr); err != nil {
			log.Error().Err(err).Msg("Error while unmarshalling week hrs json")
		}

		for _, eachDayHrs := range wArr {
			ts.TotalHours += eachDayHrs.Day1 + eachDayHrs.Day2 + eachDayHrs.Day3 + eachDayHrs.Day4 + eachDayHrs.Day5
		}

		if ts.Placement == "" {
			ts.Placement = tsResponse.Placement
		}
		if ts.Info == "" {
			ts.Info = tsResponse.Info
		}
		if ts.WeekHrs == nil {
			ts.WeekHrs = tsResponse.WeekHrs
		}

		res, err = s.repo.UpdateTimesheetByGivenCriteria(ctx, ts, loginName, month, year, wArr)
		if err != nil {
			log.Error().Err(err).Msgf("update Timesheet is failed with given criteria %s,%d,%d ", loginName, month, year)
			return "", err
		}
	}
	return res, nil
}

func (s *service) GetListofTimesheets(ctx context.Context, loginName string) ([]*GetAllTimesheets, error) {
	var err error
	ts := []*GetAllTimesheets{}

	loginName = strings.ToUpper(loginName)

	if loginName == "" {
		err = errors.New("loginName is empty")
		return nil, err
	}
	if ts, err = s.repo.SelectAllTimesheetByLoginName(ctx, loginName); err != nil {
		log.Error().Err(err).Msgf("Error while fetching timesheet info by the given login Name %s", loginName)
		return nil, err
	}
	return ts, nil
}

func (s *service) GetTimesheetsByWeek(ctx context.Context, loginName string, week, month, year int) (*GetTimesheet, error) {
	var err error
	ts := &GetAllTimesheets{}

	w := &WeekHrs{}

	loginName = strings.ToUpper(loginName)

	if loginName == "" {
		err = errors.New("loginName is empty")
		return nil, err
	}

	if week == 0 {
		err = errors.New("invalid week")
		return nil, err
	}

	if month == 0 {
		err = errors.New("invalid month")
		return nil, err
	}

	if year == 0 {
		err = errors.New("invalid year")
		return nil, err
	}

	if ts, err = s.repo.SelectTimesheetByGivenCriteria(ctx, loginName, month, year); err != nil {
		log.Error().Err(err).Str("loginName", loginName).Msgf("Error while fetching timesheet info by the given week %d", week)
		return nil, err
	}

	if w, err = s.repo.SelectWeekHoursByWeek(ctx, loginName, week, month, year); err != nil {
		log.Error().Err(err).Str("loginName", loginName).Msg("Error while fetching timesheet info by the given week from week hours table")
		return nil, err
	}

	// log.Info().Msgf("json data from db %v", ts.WeekHrs)

	// //Step1 : Unmarshal the week hrs info data

	// wArr := []WeekHrs{}

	// if err = json.Unmarshal(ts.WeekHrs, &wArr); err != nil {
	// 	log.Error().Err(err).Msg("Error while unmarshalling week hrs json")
	// }

	// log.Info().Msgf("un marshalled data into golang struct %v", wArr)

	// w := WeekHrs{}

	// for _, wI := range wArr {

	// 	if wI.WeekInfo == week {
	// 		w.WeekInfo = wI.WeekInfo
	// 		w.Day1 = wI.Day1
	// 		w.Day2 = wI.Day2
	// 		w.Day3 = wI.Day3
	// 		w.Day4 = wI.Day4
	// 		w.Day5 = wI.Day5

	// 		break
	// 	}
	// }

	timesheet := &GetTimesheet{
		LoginName:  ts.LoginName,
		Status:     ts.Status,
		Placement:  ts.Placement,
		Info:       ts.Info,
		TotalHours: ts.TotalHours,
		Month:      ts.Month,
		Year:       ts.Year,
		WeekData:   w,
	}

	return timesheet, nil
}

func (s *service) DeleteTimesheet(ctx context.Context, loginName string, month, year int) (string, error) {
	var err error
	var response string
	var ts *Timesheet

	loginName = strings.ToUpper(loginName)

	if loginName == "" {
		err = errors.New("Login name is invalid")
		return "", err
	}

	ts, err = s.repo.SelectTimesheetByLoginName(ctx, loginName, month, year)
	if err != nil {
		log.Error().Err(err).Msgf("Error while fetching Timesheet with given LoginName")
		return "", err
	}
	if ts == nil {

		response = fmt.Sprintf("We don't have data with the given login Name %s", loginName)
		return response, nil
	}

	if ts != nil {
		response, err = s.repo.DeleteTimesheet(ctx, loginName, month, year)
		if err != nil {
			log.Error().Err(err).Str("loginname", loginName).Msg("Error while calling repo DeleteTimesheet")
			return "", err
		}
	}
	return response, nil
}

func (s *service) AddorUpdatenotes(ctx context.Context, notes *Notes) (string, error) {
	var err error
	var res string

	notes.LoginName = strings.ToUpper(notes.LoginName)

	if notes.LoginName == "" && notes.Month == 0 && notes.Year == 0 {
		err = errors.New("Criteria is not valid")
		return "", err
	}

	res, err = s.repo.UpsertTimesheetNotes(ctx, notes)
	if err != nil {
		return "", err
	}

	return res, nil
}

// func (s *service) GenerateCSV(ctx context.Context, req []*GetTimesheet) (string, error) {
// 	var err error
// 	var res string
// 	if res ,err = s.repo.SelectTimesheetByLoginName()
// }
func (s *service) ListSubmittedTimesheets(ctx context.Context, status string) ([]*GetAllTimesheets, error) {
	var err error

	if status == "" {
		err = errors.New("status is mandatory")
		return nil, err
	}

	timesheets := []*GetAllTimesheets{}

	if timesheets, err = s.repo.SelectTimesheetsByStatus(ctx, status); err != nil {
		log.Error().Err(err).Str("status", status).Msg("Error while calling SelectTimesheetsByStatus ")
		return nil, err
	}

	return timesheets, nil

}
