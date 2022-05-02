package main

import (
	"encoding/json"
	"net/http"
	"timesheet/commons/res"
	"timesheet/user"

	"github.com/rs/zerolog/log"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
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
