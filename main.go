package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/mitchan/go_gym_bro/internal/app"
	"github.com/mitchan/go_gym_bro/internal/routes"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "server port")
	flag.Parse()

	app, err := app.NewApplication()
	if err != nil {
		panic(err)
	}

	app.Logger.Println("running")

	r := routes.SetupRoutes(app)
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	fmt.Printf("we are running on port :%d", port)

	err = server.ListenAndServe()
	if err != nil {
		app.Logger.Fatal(err)
	}
}
