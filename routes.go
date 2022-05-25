package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

//addRoutes will have different routes  functions and calls those functions...
func addRoutes(r *chi.Mux) {
	log.Println("Registering routes")
	addIAMRoutes(r)
	addTimesheetRoutes(r)

	log.Println("Registering routes .. done")
}

//http://localhost:8085/iam/users

func addIAMRoutes(r *chi.Mux) {
	r.Route("/iam", func(r chi.Router) {

		r.Post("/users", createUser)

		r.Put("/users/{loginName}", forgotPassword)

		//	r.Post("/users/{loginName}", LoginUser)

	})
}

func addTimesheetRoutes(r *chi.Mux) {
	r.Route("/timesheets", func(r chi.Router) {

		r.Get("/{loginName}/{week}", getTimesheetsByWeek)
	})
}

func printRoutes(r *chi.Mux) {

	log.Println("Following routes are supported")

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		fmt.Printf("%s %s\n", method, route)
		return nil
	}

	if err := chi.Walk(r, walkFunc); err != nil {
		fmt.Printf("Logging err: %s\n", err.Error())
	}
}
