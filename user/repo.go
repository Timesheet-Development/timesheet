package user

import (
	"context"
	"log"
	"timesheet/commons/res"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Repository interface {
	SelectUserByLoginName(ctx context.Context, loginName string) (*User, error)

	InsertUser(ctx context.Context, user *User) (*uuid.UUID, error)

	UpdatePassword(ctx context.Context, psswd []byte, loginName string) (string, error)
}

type repository struct {
	db *pgxpool.Pool
}

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

func (repo *repository) InsertUser(ctx context.Context, user *User) (*uuid.UUID, error) {
	tx, _ := repo.db.Begin(ctx)
	log.Println("Insert user into DB")

	insertqry := `INSERT INTO users (id, login_name, "password", "name", address, department, 
	security_no, dob, city, state, job_title, is_perm, gender, passport, reporting_mngr) 
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15);
	`
	user.Id = uuid.New()

	var tag pgconn.CommandTag
	var err error
	if tag, err = tx.Exec(ctx, insertqry, user.Id, user.LoginName, user.Password,
		user.Name, user.Address, user.Department,
		user.SocailSecurityNumber, user.DOB, user.City, user.State, user.JobTitle,
		user.IsPermanent, user.Gender, user.Passport, user.ReportingManager); err != nil {
		if err != nil {
			tx.Rollback(ctx)
			return nil, &res.AppError{ResponseCode: res.DatabaseError, Cause: err}
		}
	}
	tx.Commit(ctx)
	//Communicating with database to get user id .
	if err := pgxscan.Get(
		ctx, repo.db, user, "select id from users where login_name=$1", user.LoginName,
	); err != nil {
		// Handle query or rows processing error.
		if pgxscan.NotFound(err) {
			//No error, but no user either
			return nil, nil
		}
		return nil, &res.AppError{ResponseCode: res.DatabaseError, Cause: err}
	}
	log.Printf("userID[%v]\n", user.Id)

	log.Printf("Rows affectd [%d]\n", tag.RowsAffected())

	return &user.Id, nil
}

func (repo *repository) UpdatePassword(ctx context.Context, psswd []byte, loginName string) (string, error) {
	var err error
	updPswdQry := `update users set password=$1
				  where login_name=$2`

	if _, err = repo.db.Exec(ctx, updPswdQry, psswd, loginName); err != nil {
		log.Printf("Unable to perform update password. Error is [%v]\n", err)
		return "", err
	}
	return loginName, nil
}
