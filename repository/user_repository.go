package repository

import (
	"github.com/danangkonang/go-crud-mysql/database"
	"github.com/danangkonang/go-crud-mysql/entity"
)

type UserRepository interface {
	FindUserById(id int64) (*entity.ResponseUser, error)
	FindUsers() ([]*entity.ResponseUser, error)
	StoreUser(*entity.StoreUser) (int64, error)
	DestroyUser(id int64) (int64, error)
	UpdateUser(*entity.UpdateUser) (int64, error)
}

type repository struct {
	db *database.Database
}

func NewUserRepository(db *database.Database) UserRepository {
	return &repository{db: db}
}

func (r repository) FindUserById(id int64) (*entity.ResponseUser, error) {
	user := new(entity.ResponseUser)
	row := r.db.Mysql.QueryRow("SELECT user_id, name, email, phone FROM users WHERE user_id=?", id)
	if err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Phone); err != nil {
		return user, err
	}
	return user, nil
}

func (r repository) FindUsers() ([]*entity.ResponseUser, error) {
	var users []*entity.ResponseUser
	rows, err := r.db.Mysql.Query("SELECT user_id, name, email, phone FROM users")
	if err != nil {
		return users, err
	}
	defer rows.Close()
	for rows.Next() {
		u := new(entity.ResponseUser)
		err := rows.Scan(
			&u.Id,
			&u.Name,
			&u.Email,
			&u.Phone,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r repository) StoreUser(u *entity.StoreUser) (int64, error) {
	query := "INSERT INTO users (name, email, phone) VALUES (?,?,?)"
	row, err := r.db.Mysql.Exec(query, u.Name, u.Email, u.Phone)
	if err != nil {
		return 0, err
	}
	return row.LastInsertId()
}

func (r repository) DestroyUser(id int64) (int64, error) {
	query := "DELETE FROM users WHERE user_id=?"
	row, err := r.db.Mysql.Exec(query, id)
	if err != nil {
		return 0, err
	}
	return row.RowsAffected()
}

func (r repository) UpdateUser(u *entity.UpdateUser) (int64, error) {
	query := "UPDATE users SET name=?, email=?, phone=? WHERE user_id=?"
	row, err := r.db.Mysql.Exec(query, u.Name, u.Email, u.Phone, u.Id)
	if err != nil {
		return 0, err
	}
	return row.RowsAffected()
}
