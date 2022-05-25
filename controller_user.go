package main

import (
	"encoding/json"
	"net/http"
	"timesheet/commons/res"
	"timesheet/user"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

//
func createUser(w http.ResponseWriter, r *http.Request) {
	var err error
	userReq := &user.User{}

	if err = json.NewDecoder(r.Body).Decode(userReq); err != nil {
		log.Error().Err(err).Str("user", userReq.LoginName).Msg("Unable to parse user json to struct")
		res.SendError(w, r, err, config.Debug.PrintRootCause)
	}

	Id, err := userService.CreateUser(r.Context(), userReq)
	if err != nil {
		res.SendError(w, r, err, config.Debug.PrintRootCause)
	}
	res.SendResponse(w, r, res.OK, Id)

}

func forgotPassword(w http.ResponseWriter, r *http.Request) {
	var err error
	loginName := chi.URLParam(r, "loginName")

	updPswd := &user.UpdatePassword{}

	//TODO: When there is no oldpassword. Implement this logic

	if err = json.NewDecoder(r.Body).Decode(updPswd); err != nil {
		log.Error().Err(err).Str("loginName", loginName).Msg("Unable to parse update password json to struct")
		res.SendError(w, r, err, config.Debug.PrintRootCause)
	}

	loginName, err = userService.ForgotPassword(r.Context(), loginName, updPswd)
	if err != nil {
		res.SendError(w, r, err, config.Debug.PrintRootCause)
	}
	res.SendResponse(w, r, res.OK, loginName)
}

func loginUser(w http.ResponseWriter, r *http.Request) {
	var err error

	var jwtStr string

	user := &user.User{}
	if err = json.NewDecoder(r.Body).Decode(user); err != nil {
		log.Error().Err(err).Str("loginName", user.LoginName).Msg("Unable to parse json to user struct")
		res.SendError(w, r, err, config.Debug.PrintRootCause)
	}
	jwtStr, err = userService.LoginUser(r.Context(), user)
	if err != nil {
		res.SendError(w, r, err, config.Debug.PrintRootCause)
	} else {
		cookie := http.Cookie{
			Name:     "Timesheet",
			Value:    jwtStr,
			Secure:   true,
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)
		res.SendResponse(w, r, res.OK, jwtStr)
	}

}
