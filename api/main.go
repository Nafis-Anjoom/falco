package main

import (
	"chat/chat"
	"chat/database"
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type application struct {
    messageService *chat.MessageService
    models *database.Models
}

func (app *application) serve() {
    srv := &http.Server{
        Addr: ":3000",
        Handler: app.routes(),
        IdleTimeout: time.Minute,
        ReadTimeout: time.Minute,
        WriteTimeout: time.Minute,
    }

    go app.messageService.Run()

    log.Printf("starting server on localhost:%d", 3000)
    log.Fatal(srv.ListenAndServe())
}

func main() {
    postgresUrl := os.Getenv("DATABASE_URL")

    log.Println("database url", postgresUrl)

    dbPool, err := pgxpool.New(context.Background(), postgresUrl)
    if err != nil {
        log.Fatalln("unable to open database connection")
    }

    app := &application{
        messageService: chat.NewMessageService(),
        models: database.NewModels(dbPool),
    }

    app.serve()
}
