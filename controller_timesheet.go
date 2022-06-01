package main

import (
	"encoding/json"
	"net/http"
	"timesheet/commons/res"
	"timesheet/timesheets"

	"github.com/rs/zerolog/log"
)

func createTimesheet(w http.ResponseWriter, r *http.Request) {
	var err error
	var loginName string
	t := &timesheets.Timesheet{}

	if err = json.NewDecoder(r.Body).Decode(t); err != nil {
		log.Error().Err(err).Str("loginname", t.LoginName).Msg("Unable to parse timesheet json to struct")
		res.SendError(w, r, err, config.Debug.PrintRootCause)
	}
	loginName = t.LoginName

	loginName, err = timesheetService.CreateTimesheet(r.Context(), t)
	if err != nil {
		res.SendError(w, r, err, config.Debug.PrintRootCause)
	}
	res.SendResponse(w, r, res.OK, loginName)

}
