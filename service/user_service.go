package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/danangkonang/go-crud-mysql/entity"
	"github.com/danangkonang/go-crud-mysql/repository"
	"github.com/danangkonang/go-crud-mysql/util"
	"github.com/danangkonang/validation"
	"github.com/gorilla/mux"
)

type UserService struct {
	userRepository repository.UserRepository
}

func NewUserService(u repository.UserRepository) *UserService {
	return &UserService{userRepository: u}
}

func (s UserService) FindUserById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ID, _ := strconv.ParseInt(vars["user_id"], 0, 64)
	user, err := s.userRepository.FindUserById(ID)
	switch {
	case err == sql.ErrNoRows:
		util.Json(w, 404, "User Not Found", nil)
		return
	case err != nil:
		util.Json(w, 500, "Internal Server Error", nil)
		return

	}
	util.Json(w, 200, "Successfuly", user)
}

func (s UserService) FindUsers(w http.ResponseWriter, r *http.Request) {
	user, err := s.userRepository.FindUsers()
	if err != nil {
		util.Json(w, 500, "Internal Server Error", nil)
		return
	}
	util.Json(w, 200, "Successfuly", user)
}

func (s UserService) StoreUser(w http.ResponseWriter, r *http.Request) {
	var user entity.StoreUser
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		util.Json(w, 400, "Bad Request", nil)
		return
	}
	defer r.Body.Close()

	invalid, err := validation.Check(user)
	if err != nil {
		util.Json(w, 400, "Bad Request", invalid)
		return
	}

	ID, err := s.userRepository.StoreUser(&user)
	if err != nil {
		util.Json(w, 500, "Internal Server Error", nil)
		return
	}
	user.Id = ID
	util.Json(w, 200, "Successfuly", user)
}

func (s UserService) DestroyUser(w http.ResponseWriter, r *http.Request) {
	var user entity.DestroyeUser
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Println(err.Error())
		util.Json(w, 400, "Bad Request", nil)
		return
	}
	defer r.Body.Close()

	invalid, err := validation.Check(user)
	if err != nil {
		util.Json(w, 400, "Bad Request", invalid)
		return
	}

	row, err := s.userRepository.DestroyUser(user.Id)
	if err != nil {
		util.Json(w, 500, "Internal Server Error", nil)
		return
	}
	if row == 0 {
		util.Json(w, 400, "Id Not Found", nil)
		return
	}
	util.Json(w, 200, "Successfuly", nil)
}

func (s UserService) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user entity.UpdateUser
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		util.Json(w, 400, "Bad Request", nil)
		return
	}
	defer r.Body.Close()

	invalid, err := validation.Check(user)
	if err != nil {
		util.Json(w, 400, "Bad Request", invalid)
		return
	}

	row, err := s.userRepository.UpdateUser(&user)
	if err != nil {
		util.Json(w, 500, "Internal Server Error", nil)
		return
	}
	if row == 0 {
		util.Json(w, 400, "Id Not Found", nil)
		return
	}
	util.Json(w, 200, "Successfuly", user)
}
