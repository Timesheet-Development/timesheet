package user

import (
	"context"
	"net/http"
	"strings"
	"timesheet/commons/res"
	"timesheet/commons/validate"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	IsUserAlreadyExisting(ctx context.Context, loginName string) (bool, error)

	GetUser(ctx context.Context, loginName string) (*User, error)

	//CreateUser will create a new record in database if exists it will raise a error
	CreateUser(ctx context.Context, iam *User) (*uuid.UUID, error)

	ForgotPassword(ctx context.Context, loginName string, updPswd *UpdatePassword) (string, error)
}

var UserAlreadyExists = &res.ResponseCode{Code: "UserAlreadyExists", Message: "User already exists", HttpStatus: http.StatusBadRequest}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateUser(ctx context.Context, iam *User) (*uuid.UUID, error) {
	log.Info().Msg("User service initialized")

	var err error

	//Transforming fields
	iam.LoginName = strings.ToUpper(iam.LoginName)

	//Validate fields
	log.Info().Msgf("Validating user creation request for [%s]\n", iam.LoginName)

	ve := validate.New().
		IsRequired("LoginName", iam.LoginName).
		IsRequired("Password", iam.Password).
		IsSizeInRange("Password", iam.Password, 3, 20)
	if ve.HasErrors() {
		return nil, ve
	}

	log.Info().Msgf("Checking if user details [%v] exists", iam)

	var userExists = false
	if userExists, err = s.IsUserAlreadyExisting(ctx, iam.LoginName); err != nil {
		return nil, err
	}

	if userExists {
		return nil, &res.AppError{ResponseCode: UserAlreadyExists, Cause: nil}
	}

	log.Info().Msgf("User [%s] does not exist. Create it now\n", iam.LoginName)

	//Generate a password hash.
	log.Info().Msgf("Generating password hash for [%s]\n", iam.LoginName)

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(iam.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err).Msgf("Error occurred while hashing password for user [%s]. Error is [%s]\n", iam.LoginName, err.Error())
		return nil, err
	}
	iam.Password = string(passwordHash)

	if err = bcrypt.CompareHashAndPassword(passwordHash, []byte(iam.Password)); err != nil {
		log.Error().Err(err).Msg("Password Hash generated does not match the password! ")
	}

	ID := &iam.Id

	log.Info().Msgf("Inserting User [%s]\n", iam.LoginName)

	if ID, err = s.repo.InsertUser(ctx, iam); err != nil {
		log.Error().Err(err).Msgf("Creating user [%s] failed. Error is [%s]\n", iam.LoginName, err.Error())
		return nil, err
	}
	return ID, nil
}

func (svc *service) IsUserAlreadyExisting(ctx context.Context, loginName string) (bool, error) {
	//Transform fields
	loginName = strings.ToUpper(loginName)
	log.Info().Msgf("Checking if user [%s] exists", loginName)

	var user *User
	var err error
	if user, err = svc.GetUser(ctx, loginName); err != nil {
		return false, err
	}

	return user != nil, nil
}

func (svc *service) GetUser(ctx context.Context, loginName string) (*User, error) {
	//Transform fields
	loginName = strings.ToUpper(loginName)
	log.Info().Msgf("Checking if user [%s] exists", loginName)

	var user *User
	var err error

	if user, err = svc.repo.SelectUserByLoginName(ctx, loginName); err != nil {
		log.Printf("Get user failed. Error is [%v]\n", err)
		return nil, err
	}

	if user == nil {
		log.Info().Msgf("User [%s] not found", loginName)
		return nil, nil
	}

	log.Info().Msgf("User [%s] found", loginName)
	return user, nil
}

func (s *service) ForgotPassword(ctx context.Context, loginName string, updPswd *UpdatePassword) (string, error) {
	var err error
	user := &User{}

	user, err = s.GetUser(ctx, loginName)
	if err != nil {
		log.Error().Err(err).Str("loginName", loginName).Msg("User details not found for the given loginName")
		return "", err
	}

	if user != nil {

		if updPswd.OldPassword != "" {
			if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(updPswd.OldPassword)); err != nil {
				log.Info().Msg("Password Hash generated does not match the password! ")
			}
		}

		//Generating new hash for the user
		log.Info().Msgf("Generating password hash for [%s]\n", loginName)

		passwordHash, err := bcrypt.GenerateFromPassword([]byte(updPswd.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Error occurred while hashing password for user [%s]. Error is [%s]\n", user.LoginName, err.Error())
			return "", err
		}
		//log.Info().Msgf("Password generated is [%s]\n", passwordHash)

		log.Info().Msg("Contacting repo to update password")

		if loginName, err = s.repo.UpdatePassword(ctx, passwordHash, loginName); err != nil {
			log.Error().Err(err).Msg("Error while updating new password")
			return "", err
		}
	} else {
		log.Error().Err(err).Msgf("User info is not available with the given loginName %s\n", loginName)
		return "", err
	}

	return loginName, nil
}
