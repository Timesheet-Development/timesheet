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
