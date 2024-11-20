package main

import (
	"chat/auth"
	"chat/database"
	"chat/messaging"
	"chat/messaging/idGenerator"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type application struct {
	config         config
	messageService *messaging.MessageService
	userService    *UserService
	authService    *auth.AuthService
}

type config struct {
	port           int
	databaseURL    string
	maxDBOpenConns string
}

func (app *application) serve() {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.EnableCORS(app.LogRequest(app.authenticateDummy(app.routes()))),
		// Handler:      app.LogRequest(app.authenticateDummy(app.routes())),
		IdleTimeout:  time.Minute,
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
	}

	go app.messageService.Run()

	log.Printf("starting server on localhost:%d", app.config.port)
	log.Fatal(srv.ListenAndServe())
}

func main() {
	var config config

	port, err := strconv.ParseInt(os.Getenv("PORT"), 10, 64)
	if err != nil {
		log.Fatalln("invalid port")
	}
	config.port = int(port)

	config.databaseURL = os.Getenv("DB_URL")
	dbPool, err := pgxpool.New(context.Background(), config.databaseURL)
	if err != nil {
		log.Fatalln("unable to open database connection")
	}
	models := database.NewModels(dbPool)

	machineId, err := strconv.ParseInt(os.Getenv("MACHINE_ID"), 10, 64)
	if err != nil {
		log.Fatalln("machine Id must be an integer between 0 and 2^11")
	}
	idGenerator, err := idGenerator.NewIdGenerator(machineId)
	if err != nil {
		log.Fatalln(err)
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	authService := auth.NewAuthService(jwtSecret)

	messageService := messaging.NewMessageService(models, idGenerator)
	userService := NewUserService(models, authService)

	app := &application{
		config:         config,
		messageService: messageService,
		userService:    userService,
		authService:    authService,
	}

	app.serve()
}
