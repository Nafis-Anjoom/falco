package main

import (
	"chat/messaging"
	"chat/messaging/idGenerator"
	"chat/database"
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
    config config
    messageService *messaging.MessageService
    userService *UserService
}

type config struct {
    port int
    databaseURL string
    maxDBOpenConns string
}

func (app *application) serve() {
    srv := &http.Server{
        Addr: fmt.Sprintf(":%d", app.config.port),
        Handler: app.routes(),
        IdleTimeout: time.Minute,
        ReadTimeout: time.Minute,
        WriteTimeout: time.Minute,
    }

    go app.messageService.Run()

    log.Printf("starting server on localhost:%d", app.config.port)
    log.Fatal(srv.ListenAndServe())
}

func main() {
    var config config

    config.databaseURL = os.Getenv("DB_URL")
    port, err := strconv.ParseInt(os.Getenv("PORT"), 10, 64)
    if err != nil {
        log.Fatalln("invalid port")
    }
    config.port = int(port)

    dbPool, err := pgxpool.New(context.Background(), config.databaseURL)
    if err != nil {
        log.Fatalln("unable to open database connection")
    }

    models := database.NewModels(dbPool)
    idGenerator, err := idGenerator.NewIdGenerator(0)
    if err != nil {
        log.Fatalln(err)
    }
    messageService := messaging.NewMessageService(models, idGenerator)
    userService := NewUserService(models)

    app := &application{
        config: config,
        messageService: messageService,
        userService: userService,
    }

    app.serve()
}
