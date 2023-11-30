package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/danangkonang/go-crud-mysql/database"
	"github.com/danangkonang/go-crud-mysql/repository"
	"github.com/danangkonang/go-crud-mysql/router"
	"github.com/danangkonang/go-crud-mysql/service"
	"github.com/danangkonang/go-crud-mysql/util"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

func main() {
	flag.Parse()
	db, err := database.DB()
	if err != nil {
		panic(err.Error())
	}

	_, err = db.Mysql.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			user_id int NOT NULL AUTO_INCREMENT,
			name varchar(255) NOT NULL,
			email varchar(255) NOT NULL,
			phone varchar(255) NOT NULL,
			PRIMARY KEY (user_id)
		) engine=InnoDB charset=UTF8;
	`)
	if err != nil {
		panic(err.Error())
	}

	r := mux.NewRouter().StrictSlash(false)

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		util.Json(w, 404, "Not Found", nil)
	})

	r.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		util.Json(w, 405, "Method NotAllowed", nil)
	})

	usr := service.NewUserService(
		repository.NewUserRepository(db),
	)
	router.UserRoutes(db, r, usr)

	serverloging := fmt.Sprintf("local server started at http://localhost:%s", os.Getenv("APP_PORT"))
	fmt.Println(serverloging)

	srv := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%s", os.Getenv("APP_PORT")),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler: handlers.CORS(
			handlers.AllowedHeaders([]string{
				"X-Requested-With",
				"Content-Type",
				"Authorization",
				"Access-Control-Allow-Credentials",
				"Access-Control-Allow-Origin",
			}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
			handlers.AllowedOrigins([]string{
				"http://127.0.0.1:8080",
			}),
			handlers.AllowCredentials(),
		)(r),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	srv.Shutdown(ctx)
	log.Println("shutting down")
	os.Exit(0)
}
