package requestparams

import (
	"time"
	// m "orientation-training-api/internal/models"
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
