package user

import (
	"context"
	"fmt"
	"log"
	"time"
	"timesheet/commons/res"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Repository interface {
	SelectUserByLoginName(ctx context.Context, loginName string) (*User, error)

	SelectUsers(ctx context.Context) ([]*SelectUser, error)

	//InsertUser is using the db variable contacting the database to create a new user.
	//If there is any error in the flow it will return the error
	InsertUser(ctx context.Context, user *User) (*uuid.UUID, error)

	UpdatePassword(ctx context.Context, psswd []byte, loginName string) (string, error)

	UpdateUser(ctx context.Context, loginName string, user *User) (string, error)
}

type repository struct {
	db *pgxpool.Pool
}

// NewRepository is a function which holds Database variable. Which is used to perform different actions on DB.
func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (repo *repository) SelectUserByLoginName(ctx context.Context, loginName string) (*User, error) {

	user := &User{}

	if err := pgxscan.Get(
		ctx, repo.db, user, "SELECT * FROM users where login_name=$1", loginName,
	); err != nil {
		// Handle query or rows processing error.
		if pgxscan.NotFound(err) {
			//return nil, &res.AppError{ResponseCode: UserDoesNotExist, Cause: err}
			//No error, but no user either
			return nil, nil
		}
		return nil, &res.AppError{ResponseCode: res.DatabaseError, Cause: err}
	}
	// users variable now contains data from all rows.
	return user, nil
}

// InsertUser is using the db variable contacting the database to create a new user.
// If there is any error in the flow it will return the error
func (repo *repository) InsertUser(ctx context.Context, user *User) (*uuid.UUID, error) {
	log.Println("Insert user into DB")

	insertqry := `INSERT INTO users (id, login_name, "password", "name", address, department, 
	security_no, dob, city, state, job_title, is_perm, gender, passport, reporting_mngr, work_mail,
	personal_mail,phone_number) 
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18);
	`
	user.Id = uuid.New()

	var tag pgconn.CommandTag
	var err error
	if tag, err = repo.db.Exec(ctx, insertqry, user.Id, user.LoginName, user.Password,
		user.Name, user.Address, user.Department,
		user.SocailSecurityNumber, user.DOB, user.City, user.State, user.JobTitle,
		user.IsPermanent, user.Gender, user.Passport, user.ReportingManager, user.WorkMail,
		user.PersonalMail, user.PhoneNumber); err != nil {
		return nil, &res.AppError{ResponseCode: res.DatabaseError, Cause: err}
	}

	log.Printf("userID[%v]\n", user.Id)

	log.Printf("Rows affectd [%d]\n", tag.RowsAffected())

	return &user.Id, nil
}

func (repo *repository) UpdatePassword(ctx context.Context, psswd []byte, loginName string) (string, error) {
	var err error
	updPswdQry := `update users set password=$2
				  where login_name=$1`

	if _, err = repo.db.Exec(ctx, updPswdQry, loginName, psswd); err != nil {
		log.Printf("Unable to perform update password. Error is [%v]\n", err)
		return "", err
	}

	return loginName, nil
}

func (repo *repository) UpdateUser(ctx context.Context, loginName string, user *User) (string, error) {
	var err error

	var updateStr string

	updateUserQuery := `UPDATE public.users
	SET "name" = $1, dob = $2, city = $3, state = $4,
	address = $5, job_title = $6, gender = $7, passport = $8, 
	work_mail = $9, personal_mail = $10, phone_number = $11, updated_at = $12
	WHERE login_name = $13;
	`

	if _, err = repo.db.Exec(ctx, updateUserQuery, user.Name, user.DOB, user.City, user.State,
		user.Address, user.JobTitle, user.Gender, user.Passport, user.WorkMail,
		user.PersonalMail, user.PhoneNumber, time.Now(), loginName); err != nil {
		log.Printf("Error whil performing update user %v\n", err)
		return "", err
	}

	updateStr = fmt.Sprintf("User details are updated successfully for the given loginName %s", loginName)

	return updateStr, nil
}

func (repo *repository) SelectUsers(ctx context.Context) ([]*SelectUser, error) {
	users := []*SelectUser{}

	var err error

	var rows pgx.Rows

	selectUsers := `select name,login_name , department ,security_no ,dob , city ,state ,
	address ,job_title ,is_perm ,gender ,passport ,
	reporting_mngr ,work_mail ,personal_mail ,phone_number 
	from users order by created_at desc;`

	rows, err = repo.db.Query(ctx, selectUsers)
	if err != nil {
		log.Println("Error while performing selectusers")
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {

		user := &SelectUser{}

		if err = rows.Scan(&user.Name, &user.LoginName, &user.Department, &user.SocailSecurityNumber,
			&user.DOB, &user.City,
			&user.State, &user.Address, &user.JobTitle, &user.IsPermanent, &user.Gender, &user.Passport,
			&user.ReportingManager,
			&user.WorkMail, &user.PersonalMail, &user.PhoneNumber); err != nil {

			log.Printf("Error while scanning the data into user struct %v\n", err)
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}
