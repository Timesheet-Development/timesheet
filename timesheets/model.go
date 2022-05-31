package timesheets

import (
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
}
