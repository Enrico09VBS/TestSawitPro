// This file contains types that are used in the repository layer.
package repository

type RegistrationParam struct {
	PhoneNumber string `json:"phone_number"`
	FullName    string `json:"full_name"`
	Password    string `json:"password"`
}

type LoginParam struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type UserData struct {
	PhoneNumber string `json:"phone_number" db:"phone_number"`
	FullName    string `json:"full_name" db:"full_name"`
}

type UpdateProfileParam struct {
	ID          int64
	PhoneNumber string `json:"phone_number"`
	FullName    string `json:"full_name"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}
