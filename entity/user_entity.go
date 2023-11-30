package entity

type StoreUser struct {
	Id    int64  `json:"user_id"`
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required"`
	Phone string `json:"phone" validate:"required"`
}

type DestroyeUser struct {
	Id int64 `json:"user_id" validate:"required"`
}

type UpdateUser struct {
	Id    int64  `json:"user_id" validate:"required"`
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required"`
	Phone string `json:"phone" validate:"required"`
}

type ResponseUser struct {
	Id    int64  `json:"user_id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}
