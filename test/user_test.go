package test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/danangkonang/go-crud-mysql/database"
	"github.com/danangkonang/go-crud-mysql/entity"
	"github.com/danangkonang/go-crud-mysql/repository"
	"github.com/danangkonang/go-crud-mysql/service"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

type Database struct {
	Mysql *sql.DB
}

func NewMock() (*sql.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		return nil, nil, err
	}
	return db, mock, nil
}

func UserServer(db *database.Database) *mux.Router {
	r := mux.NewRouter()
	s := service.NewUserService(
		repository.NewUserRepository((*database.Database)(db)),
	)

	v1 := r.PathPrefix("/api/v1").Subrouter()

	v1.HandleFunc("/users", s.FindUsers).Methods("GET")
	v1.HandleFunc("/user/{user_id:[0-9]+}", s.FindUserById).Methods("GET")
	v1.HandleFunc("/user", s.StoreUser).Methods("POST")
	v1.HandleFunc("/user", s.UpdateUser).Methods("PUT")
	v1.HandleFunc("/user", s.DestroyUser).Methods("DELETE")
	return r
}

func TestFindUsers(t *testing.T) {
	sqlDB, mock, err := NewMock()
	if err != nil {
		t.Fatal(err)
	}
	con := &Database{
		Mysql: sqlDB,
	}
	defer con.Mysql.Close()
	rows := sqlmock.NewRows(
		[]string{"user_id", "name", "email", "phone"},
	).AddRow(1, "bar", "foo@gmail.com", "08123123123")
	query := "SELECT user_id, name, email, phone FROM users"
	mock.ExpectQuery(query).WillReturnRows(rows)
	request, _ := http.NewRequest("GET", "/api/v1/users", nil)
	response := httptest.NewRecorder()
	UserServer(&database.Database{Mysql: con.Mysql}).ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code)
	assert.JSONEq(t, response.Body.String(), `{"data":[{"email":"foo@gmail.com", "name":"bar", "phone":"08123123123", "user_id":1}], "message":"Successfuly", "status":200}`)
}

func TestFindUsersById(t *testing.T) {
	sqlDB, mock, err := NewMock()
	if err != nil {
		t.Fatal(err)
	}
	con := &Database{
		Mysql: sqlDB,
	}
	defer con.Mysql.Close()
	rows := sqlmock.NewRows(
		[]string{"user_id", "name", "email", "phone"},
	).AddRow(1, "bar", "foo@gmail.com", "08123123123")
	query := "SELECT user_id, name, email, phone FROM users WHERE user_id=?"
	mock.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)
	request, _ := http.NewRequest("GET", "/api/v1/user/1", nil)
	response := httptest.NewRecorder()
	UserServer(&database.Database{Mysql: con.Mysql}).ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code)
}

func TestFindUsersByIdNotFound(t *testing.T) {
	sqlDB, mock, err := NewMock()
	if err != nil {
		t.Fatal(err)
	}
	con := &Database{
		Mysql: sqlDB,
	}
	defer con.Mysql.Close()
	rows := sqlmock.NewRows(
		[]string{"user_id", "name", "email", "phone"},
	)
	query := "SELECT user_id, name, email, phone FROM users WHERE user_id=?"
	mock.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)
	request, _ := http.NewRequest("GET", "/api/v1/user/1", nil)
	response := httptest.NewRecorder()
	UserServer(&database.Database{Mysql: con.Mysql}).ServeHTTP(response, request)
	assert.Equal(t, 404, response.Code)
}

func TestSaveUser(t *testing.T) {
	sample := []struct {
		status int
		body   *entity.StoreUser
	}{
		{
			status: 200,
			body: &entity.StoreUser{
				Name:  "bar",
				Email: "foo@email.com",
				Phone: "08123123123",
			},
		},
		{
			status: 400,
			body: &entity.StoreUser{
				Name:  "",
				Email: "foo@email.com",
				Phone: "08123123123",
			},
		},
		{
			status: 400,
			body: &entity.StoreUser{
				Name:  "",
				Email: "foo@email.com",
				Phone: "08123123123",
			},
		},
		{
			status: 400,
			body: &entity.StoreUser{
				Name:  "",
				Email: "",
				Phone: "08123123123",
			},
		},
		{
			status: 400,
			body: &entity.StoreUser{
				Name:  "",
				Email: "",
				Phone: "",
			},
		},
	}
	for _, v := range sample {
		sqlDB, mock, err := NewMock()
		if err != nil {
			t.Fatal(err)
		}
		con := &Database{
			Mysql: sqlDB,
		}
		defer con.Mysql.Close()
		godotenv.Load()

		query := "INSERT INTO users (name, email, phone) VALUES (?,?,?)"
		mock.ExpectExec(query).WithArgs(v.body.Name, v.body.Email, v.body.Phone).WillReturnResult(sqlmock.NewResult(1, 1))
		b, err := json.Marshal(v.body)
		if err != nil {
			t.Error(err.Error())
		}
		request, _ := http.NewRequest("POST", "/api/v1/user", bytes.NewBuffer(b))

		response := httptest.NewRecorder()
		UserServer(&database.Database{Mysql: con.Mysql}).ServeHTTP(response, request)
		assert.Equal(t, v.status, response.Code)
	}
}

func TestUpdateUser(t *testing.T) {
	sample := []struct {
		status int
		body   *entity.UpdateUser
	}{
		{
			status: 200,
			body: &entity.UpdateUser{
				Id:    1,
				Name:  "bar",
				Email: "foo@email.com",
				Phone: "08123123123",
			},
		},
		{
			status: 400,
			body: &entity.UpdateUser{
				Name:  "bar",
				Email: "foo@email.com",
				Phone: "08123123123",
			},
		},
		{
			status: 400,
			body: &entity.UpdateUser{
				Id:    1,
				Name:  "",
				Email: "foo@email.com",
				Phone: "08123123123",
			},
		},
		{
			status: 400,
			body: &entity.UpdateUser{
				Id:    1,
				Name:  "bar",
				Email: "foo@email.com",
				Phone: "",
			},
		},
	}
	for _, v := range sample {
		sqlDB, mock, err := NewMock()
		if err != nil {
			t.Fatal(err)
		}
		con := &Database{
			Mysql: sqlDB,
		}
		defer con.Mysql.Close()

		query := "UPDATE users SET name=?, email=?, phone=? WHERE user_id=?"
		mock.ExpectExec(query).WithArgs(v.body.Name, v.body.Email, v.body.Phone, v.body.Id).WillReturnResult(sqlmock.NewResult(0, 0))
		b, err := json.Marshal(v.body)
		if err != nil {
			t.Error(err.Error())
		}

		request, _ := http.NewRequest("PUT", "/api/v1/user", bytes.NewBuffer([]byte(b)))
		response := httptest.NewRecorder()
		UserServer(&database.Database{Mysql: con.Mysql}).ServeHTTP(response, request)
		assert.Equal(t, v.status, response.Code)

	}
}

func TestDeleteUser(t *testing.T) {
	sample := []struct {
		status int
		body   *entity.DestroyeUser
	}{
		{
			status: 200,
			body: &entity.DestroyeUser{
				Id: 1,
			},
		},
		{
			status: 400,
			body:   &entity.DestroyeUser{},
		},
	}
	for _, v := range sample {
		sqlDB, mock, err := NewMock()
		if err != nil {
			t.Fatal(err)
		}
		con := &Database{
			Mysql: sqlDB,
		}
		defer con.Mysql.Close()

		query := "DELETE FROM users WHERE user_id=?"
		mock.ExpectExec(query).WithArgs(v.body.Id).WillReturnResult(sqlmock.NewResult(0, 0))

		b, err := json.Marshal(v.body)
		if err != nil {
			t.Error(err.Error())
		}

		request, _ := http.NewRequest("DELETE", "/api/v1/user", bytes.NewBuffer(b))
		response := httptest.NewRecorder()
		UserServer(&database.Database{Mysql: con.Mysql}).ServeHTTP(response, request)
		assert.Equal(t, v.status, response.Code)

	}
}
