package requestparams

import (
	"time"
)

type UserProfileListParams struct {
	Name        string    `json:"name"`
	Email       string    `json:"email" valid:"length(3|1000)~Email at least 3 character"`
	PhoneNumber string    `json:"phone_number"`
	Status      int       `json:"status"`
	DateFrom    time.Time `json:"date_from"`
	DateTo      time.Time `json:"date_to"`
	CurrentPage int       `json:"current_page" valid:"-"`
	RowPerPage  int       `json:"row_per_page"`
}

type UserInfoParams struct {
	UserID int `json:"user_id" valid:"required"`
}

type UpdateProfileParams struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	Birthday    string `json:"birthday"`
	Avatar      string `json:"avatar"`
}

type AllUserName struct {
	UserID            int       `json:"user_id"`
	FullName          string    `json:"full_name"`
	Avatar            string    `json:"avatar"`
	Email             string    `json:"email"`
	CompanyJoinedDate time.Time `json:"company_joined_date"`
	Birthday          string    `json:"birthday"`
}

type DeleteUserParams struct {
	UserId int `json:"user_id" valid:"required"`
}

type EmployeeDetailParams struct {
	UserID int `json:"user_id" valid:"required"`
}

type RegisterParams struct {
	FirstName         string `json:"first_name" validate:"required"`
	LastName          string `json:"last_name" validate:"required"`
	Email             string `json:"email" validate:"required,email"`
	Password          string `json:"password" validate:"required"`
	PhoneNumber       string `json:"phone_number"`
	PersonnalEmail    string `json:"personnal_email" validate:"omitempty,email"`
	Birthday          string `json:"birthday"`
	Department        string `json:"department" validate:"required"`
	Avatar            string `json:"avatar"`
	Gender            int    `json:"gender" validate:"required"`
	RoleID            int    `json:"role_id" validate:"required"`
	CompanyJoinedDate string `json:"company_joined_date"`
}
