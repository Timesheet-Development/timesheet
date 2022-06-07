package timesheets

import (
	"time"

	"github.com/google/uuid"
	sql "github.com/jmoiron/sqlx/types"
)

type Timesheet struct {
	ID         uuid.UUID
	LoginName  string
	Status     string
	Placement  string
	Info       string
	TotalHours float64
	Month      int
	Year       int
	WeekHrs    sql.JSONText
	WeekDay    sql.JSONText
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type GetAllTimesheets struct {
	LoginName  string
	Status     string
	Placement  string
	Info       string
	TotalHours float64
	Month      int
	Year       int
	WeekHrs    sql.JSONText
	WeekDay    sql.JSONText
}

type timesheetStatus string

const (
	timesheetStatusSubmitted timesheetStatus = "Submitted"
	timesheetStatusViewed    timesheetStatus = "Viewed"
	timesheetStatusApproved  timesheetStatus = "Approved"
	timesheetStatusRejected  timesheetStatus = "Rejected"
)

type WeekHrs struct {
	WeekInfo int
	Day1     float64
	Day2     float64
	Day3     float64
	Day4     float64
	Day5     float64
}
