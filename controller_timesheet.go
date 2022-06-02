package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"timesheet/commons/res"
	"timesheet/timesheets"

	"github.com/go-chi/chi/v5"
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
func updateTimesheet(w http.ResponseWriter, r *http.Request) {
	var err error
	loginName := chi.URLParam(r, "loginName")
	month := chi.URLParam(r, "month")
	year := chi.URLParam(r, "year")
	var m, y int

	m, err = strconv.Atoi(month)
	if err != nil {
		log.Error().Err(err).Msg("month conversion of string to int is failed")
	}
	y, err = strconv.Atoi(year)
	if err != nil {
		log.Error().Err(err).Msg("year conversion of string to int is failed")
	}

	t := &timesheets.Timesheet{}

	if err = json.NewDecoder(r.Body).Decode(t); err != nil {
		log.Error().Err(err).Str("loginname", t.LoginName).Msg("Unable to parse timesheet json to struct")
		res.SendError(w, r, err, config.Debug.PrintRootCause)
	}
	var response string
	response, err = timesheetService.UpdateTimesheet(r.Context(), t, loginName, m, y)
	if err != nil {
		res.SendError(w, r, err, config.Debug.PrintRootCause)
	}
	res.SendResponse(w, r, res.OK, response)

}
