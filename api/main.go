package main

import (
	"chat/chat"
	"chat/chat/idGenerator"
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
    messageService *chat.MessageService
    models *database.Models
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

    // there should be a better way to initialize models and messageService
    // models should not be a top level field in application. Then all services should have a models field
    models := database.NewModels(dbPool)
    idGenerator, err := idGenerator.NewIdGenerator(0)
    if err != nil {
        log.Fatalln(err)
    }
    messageService := chat.NewMessageService(models, idGenerator)

    app := &application{
        config: config,
        messageService: messageService,
        models: models,
    }

    app.serve()
}
