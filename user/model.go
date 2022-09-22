package user

import (
	"time"

	"github.com/google/uuid"
)

/*
<<<<---User--->>>>
Contains fields which discuss about a employee or a person in a organization.
*/
type User struct {
	Id                   uuid.UUID `db:"id"`
	Name                 string    `db:"name"`
	LoginName            string    `db:"login_name"`
	Password             string    `db:"password"`
	Department           string    `db:"department"`
	SocailSecurityNumber string    `db:"security_no"`
	DOB                  time.Time `db:"dob"`
	City                 string    `db:"city"`
	State                string    `db:"state"`
	Address              string    `db:"address"`
	JobTitle             string    `db:"job_title"`
	IsPermanent          bool      `db:"is_perm"`
	Gender               string    `db:"gender"`
	Passport             string    `db:"passport"`
	ReportingManager     uuid.UUID `db:"reporting_mngr"`
	WorkMail             string    `db:"work_mail"`
	PersonalMail         string    `db:"personal_mail"`
	PhoneNumber          string    `db:"phone_number"`
	CreatedAt            time.Time `db:"created_at"`
	UpdatedAt            time.Time `db:"updated_at"`
}

type SelectUser struct {
	Name                 string    `db:"name"`
	LoginName            string    `db:"login_name"`
	Department           string    `db:"department"`
	SocailSecurityNumber string    `db:"security_no"`
	DOB                  time.Time `db:"dob"`
	City                 string    `db:"city"`
	State                string    `db:"state"`
	Address              string    `db:"address"`
	JobTitle             string    `db:"job_title"`
	IsPermanent          bool      `db:"is_perm"`
	Gender               string    `db:"gender"`
	Passport             string    `db:"passport"`
	ReportingManager     uuid.UUID `db:"reporting_mngr"`
	WorkMail             string    `db:"work_mail"`
	PersonalMail         string    `db:"personal_mail"`
	PhoneNumber          string    `db:"phone_number"`
}

type UpdatePassword struct {
	OldPassword string
	NewPassword string
}
