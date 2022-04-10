package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func addRoutes(r *chi.Mux) {
	log.Println("Registering routes")
	addIAMRoutes(r)

	log.Println("Registering routes .. done")
}
func addIAMRoutes(r *chi.Mux) {

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
