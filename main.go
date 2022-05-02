package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {

	fmt.Println("Main started")
	//Load configuration
	if err := initConfig(); err != nil {
		log.Fatalf("Failed loading configuration. err=%s\n", err.Error())
	}

	//Setup router and middleware
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	//cors code

	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*", "*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"X-PINGOTHER", "Accept", "Authorization", "Content-Type", "Content-Type: application/json", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	//	log.Println("cors",cors)
	r.Use(cors.Handler)

	//Add custom routes
	addRoutes(r)

	//Print all routes
	printRoutes(r)

	//Initialize various dependent services or Fail if the resources are not bound
	initResourcesOrFail()

	//Start the listener
	log.Printf("Starting service at port %d", config.HTTP.Port)
	if err := http.ListenAndServe(":"+strconv.Itoa(config.HTTP.Port), r); err != nil {
		log.Printf("Error occurred starting HTTP service: %s\n", err.Error())
	}
	log.Printf("HTTP Service shutting down")

}
