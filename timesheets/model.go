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
	WeekHrs    sql.JSONText `db:"week_hours_info"`
	WeekDay    sql.JSONText `db:"week_day_info"`
}

type GetTimesheet struct {
	LoginName  string
	Status     string
	Placement  string
	Info       string
	TotalHours float64
	Month      int
	Year       int
	WeekData   WeekHrs
}

type timesheetStatus string

const (
	timesheetStatusSubmitted timesheetStatus = "Submitted"
	timesheetStatusViewed    timesheetStatus = "Viewed"
	timesheetStatusApproved  timesheetStatus = "Approved"
	timesheetStatusRejected  timesheetStatus = "Rejected"
)

type WeekHrs struct {
	WeekInfo int `json:"week_info"`
	Day1     float64
	Day2     float64
	Day3     float64
	Day4     float64
	Day5     float64
}

type WeekDay struct {
	WeekInfo int `json:"week_info"`
	Day1     string
	Day2     string
	Day3     string
	Day4     string
	Day5     string
}

type AddorUpdateNotes struct {
	LoginName string
	Month     int
	Year      int
	Info      string
}
