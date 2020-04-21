package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/todo/handler"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const defaultPort = "7771"

func port() string {
	p := os.Getenv("PORT")
	if p != "" {
		return ":" + p
	}
	return ":" + defaultPort
}

func main() {
	connStr := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True",
		"root",
		"{}",
		"localhost:3306",
		"todo",
	)

	db, err := gorm.Open("mysql", connStr)
	if err != nil {
		log.Fatal(err)
	}

	h := handler.New(db)

	server := &http.Server{
		Addr:    port(),
		Handler: h,
	}

	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt)

		<-stop
		log.Println("Shutting down...")

		if err := server.Shutdown(context.Background()); err != nil {
			log.Println("Unable to shutdown:", err)
		}

		log.Println("Server stopped")
	}()

	log.Println("Listening on http://localhost" + port())
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
