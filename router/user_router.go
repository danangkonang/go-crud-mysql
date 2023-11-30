package router

import (
	"github.com/danangkonang/go-crud-mysql/database"
	"github.com/danangkonang/go-crud-mysql/service"
	"github.com/gorilla/mux"
)

func UserRoutes(db *database.Database, r *mux.Router, s *service.UserService) {
	v1 := r.PathPrefix("/api/v1").Subrouter()

	v1.HandleFunc("/users", s.FindUsers).Methods("GET")
	v1.HandleFunc("/user/{user_id:[0-9]+}", s.FindUserById).Methods("GET")
	v1.HandleFunc("/user", s.StoreUser).Methods("POST")
	v1.HandleFunc("/user", s.UpdateUser).Methods("PUT")
	v1.HandleFunc("/user", s.DestroyUser).Methods("DELETE")
}
