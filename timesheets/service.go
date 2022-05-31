package timesheets

import (
	"context"
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

	/* check whether loginName is existing or not
	if not throw error,
	checking if the loginName is in db
	ToDo: total hrs calculation, define status
	*/

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

	if loginName, err = s.repo.InsertTimesheet(ctx, ts); err != nil {
		log.Error().Err(err).Str("loginName", loginName).Msg("Error while calling repo in timesheet service")
		return "", err
	}

	return loginName, nil
}
