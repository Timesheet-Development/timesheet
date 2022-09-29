package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
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

func getListofTimesheets(w http.ResponseWriter, r *http.Request) {
	loginName := chi.URLParam(r, "loginName")

	stms, err := timesheetService.GetListofTimesheets(r.Context(), loginName)

	if err != nil {
		res.SendError(w, r, err, config.Debug.PrintRootCause)
	} else if stms == nil {
		res.SendResponse(w, r, res.RecordNotFound, nil)
	} else {
		res.SendResponse(w, r, res.OK, stms)
	}
}

func getTimesheetsByWeek(w http.ResponseWriter, r *http.Request) {
	var err error
	var weekInt, monthInt, yearInt int

	loginName := chi.URLParam(r, "loginName")
	week := chi.URLParam(r, "week")

	weekInt, err = strconv.Atoi(week)
	if err != nil {
		log.Error().Err(err).Msg("error while converting week datatype string to int")
	}

	month := chi.URLParam(r, "month")
	monthInt, err = strconv.Atoi(month)
	if err != nil {
		log.Error().Err(err).Msg("error while converting month datatype string to int")
	}

	year := chi.URLParam(r, "year")
	yearInt, err = strconv.Atoi(year)
	if err != nil {
		log.Error().Err(err).Msg("error while converting year datatype string to int")
	}

	var timesheet *timesheets.GetTimesheet
	timesheet, err = timesheetService.GetTimesheetsByWeek(r.Context(), loginName, weekInt, monthInt, yearInt)
	if err != nil {
		log.Error().Err(err).Msg("error while calling GetTimesshetsByWeek")
		res.SendError(w, r, err, config.Debug.PrintRootCause)
	} else if timesheet == nil {
		res.SendResponse(w, r, res.RecordNotFound, nil)
	} else {
		res.SendResponse(w, r, res.OK, timesheet)
	}
}

func deleteTimesheet(w http.ResponseWriter, r *http.Request) {
	var err error
	var monthInt, yearInt int
	var response string

	loginName := chi.URLParam(r, "loginName")

	month := chi.URLParam(r, "month")
	monthInt, err = strconv.Atoi(month)
	if err != nil {
		log.Error().Err(err).Msg("error while converting month datatype string to int")
	}

	year := chi.URLParam(r, "year")
	yearInt, err = strconv.Atoi(year)
	if err != nil {
		log.Error().Err(err).Msg("error while converting year datatype string to int")
	}

	if response, err = timesheetService.DeleteTimesheet(r.Context(), loginName, monthInt, yearInt); err != nil {
		log.Error().Err(err).Str("loginName", loginName).Msg("Error while deleting timesheet record")
		res.SendError(w, r, err, config.Debug.PrintRootCause)
	}
	res.SendResponse(w, r, res.OK, &response)
}

func addorUpdateNotes(w http.ResponseWriter, r *http.Request) {
	var err error
	notes := &timesheets.Notes{}
	var response string

	if err = json.NewDecoder(r.Body).Decode(notes); err != nil {
		log.Error().Err(err).Str("loginname", notes.LoginName).Msg("Unable to parse timesheet json to struct")
		res.SendError(w, r, err, config.Debug.PrintRootCause)
	}

	response, err = timesheetService.AddorUpdatenotes(r.Context(), notes)
	if err != nil {
		log.Error().Err(err).Msg("error msg")
		res.SendError(w, r, err, config.Debug.PrintRootCause)
	}
	res.SendResponse(w, r, res.OK, response)

}

func generateCSV(w http.ResponseWriter, r *http.Request) {
	var err error
	// var weekInt, monthInt, yearInt int

	// loginName := chi.URLParam(r, "loginName")
	// week := chi.URLParam(r, "week")

	// weekInt, err = strconv.Atoi(week)
	// if err != nil {
	// 	log.Error().Err(err).Msg("error while converting week datatype string to int")
	// }

	// month := chi.URLParam(r, "month")
	// monthInt, err = strconv.Atoi(month)
	// if err != nil {
	// 	log.Error().Err(err).Msg("error while converting month datatype string to int")
	// }

	// year := chi.URLParam(r, "year")
	// yearInt, err = strconv.Atoi(year)
	// if err != nil {
	// 	log.Error().Err(err).Msg("error while converting year datatype string to int")
	// }

	var data string
	req := []*timesheets.GetTimesheet{}

	length := int(r.ContentLength)
	datalen := make([]byte, length)
	r.Body.Read(datalen)
	json.Unmarshal(datalen, &req)

	strTime := time.Now().Format("2006_01_02_15:04:05")

	//data, err = timesheetService.GenerateCSV(r.Context(), req)
	if err != nil {
		log.Error().Err(err).Msg("error while calling GetTimesshetsByWeek")
		res.SendError(w, r, err, config.Debug.PrintRootCause)
	} else if data == "" {
		res.SendResponse(w, r, res.RecordNotFound, nil)
	} else {
		filename := fmt.Sprintf("\"Timesheet-Report-%s.csv\"", strTime)
		w.Header().Add("Content-Disposition", "attachment; filename= "+filename)
		w.Header().Add("x-export-filename", filename)
		res.WriteCSV(w, r, data)
	}
}
func listSubmittedTimesheets(w http.ResponseWriter, r *http.Request) {
	var err error
	status := chi.URLParam(r, "status")

	timesheets := []*timesheets.GetAllTimesheets{}

	if timesheets, err = timesheetService.ListSubmittedTimesheets(r.Context(), status); err != nil {
		log.Error().Err(err).Str("status", status).Msg("Failed retriving submitted timesheets")
		res.SendError(w, r, err, config.Debug.PrintRootCause)
	}
	res.SendResponse(w, r, res.OK, timesheets)
}
