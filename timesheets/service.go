package timesheets

import (
	"context"
	"encoding/json"
	"strings"
	"timesheet/user"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type Service interface {
	CreateTimesheet(ctx context.Context, ts *Timesheet) (string, error)
	UpdateTimesheet(ctx context.Context, ts *Timesheet, loginName string, month, year int) (string, error)
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
	user := &user.User{}

	if ts.LoginName == "" {
		err = errors.New("loginName is empty")
		return "", err
	}

	//Transforming fields
	ts.LoginName = strings.ToUpper(ts.LoginName)

	user, err = s.userRepo.SelectUserByLoginName(ctx, ts.LoginName)
	if err != nil {
		log.Error().Err(err).Str("loginName", ts.LoginName).Msg("User details not found for the given loginName")
		return "", err
	}

	log.Info().Str("loginName", user.LoginName).Msg("logging the timesheet info")

	ts.ID = uuid.New()

	ts.LoginName = user.LoginName
	ts.Placement = user.Department + " " + user.JobTitle
	ts.Status = string(timesheetStatusSubmitted)

	//Unmarshalling the week hrs json
	wArr := []WeekHrs{}

	if err = json.Unmarshal(ts.WeekHrs, &wArr); err != nil {
		log.Error().Err(err).Msg("Error while unmarshalling week hrs json")
	}

	for _, eachDayHrs := range wArr {
		ts.TotalHours = eachDayHrs.Day1 + eachDayHrs.Day2 + eachDayHrs.Day3 + eachDayHrs.Day4 + eachDayHrs.Day5
	}

	if loginName, err = s.repo.InsertTimesheet(ctx, ts); err != nil {
		log.Error().Err(err).Str("loginName", loginName).Msg("Error while calling repo in timesheet service")
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
	var isExisting bool
	var res string

	if loginName == "" {
		err = errors.New("loginName is empty")
		return "", err
	}

	if month == 0 {
		err = errors.New("loginName is empty")
		return "", err
	}

	if year == 0 {
		err = errors.New("loginName is empty")
		return "", err
	}

	isExisting, err = s.repo.SelectTimesheetByLoginName(ctx, loginName, month, year)
	if err != nil {
		log.Error().Err(err).Msgf("Error while fetching Timesheet with given LoginName")
		return "", err
	}

	if isExisting {
		res, err = s.repo.UpdateTimesheetByGivenCriteria(ctx, ts, loginName, month, year)
		if err != nil {
			log.Error().Err(err).Msgf("update Timesheet is failed with given criteria %s,%d,%d ", loginName, month, year)
			return "", err
		}
	}
	return res, nil
}
