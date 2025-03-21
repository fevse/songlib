package main

import (
	"context"
	"database/sql"
	"log"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"

	"github.com/fevse/songlib/internal/app"
	"github.com/fevse/songlib/internal/config"
	"github.com/fevse/songlib/internal/server"
	"github.com/fevse/songlib/internal/storage"
)

func main() {
	conf := config.LoadConfig()

	db, err := sql.Open("postgres", conf.DBConnectionString())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	storage := storage.NewStorage(db)

	err = storage.Migrate()
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	app := app.NewSongLibApp(storage, conf.MIURL)
	server := server.NewServer(app, conf.ServHost, conf.ServPort)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			log.Fatalf("Failed to stop server: %v", err)
		}
		log.Println("Server stopped")
	}()

	log.Printf("Server started on %v:%v", conf.ServHost, conf.ServPort)
	if err := server.Start(ctx); err != nil {
		log.Fatalf("!Server closed: %v", err)
	}

}
