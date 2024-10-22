package dto

import (
	"music-echo/utils"
)

type TrackPostRequest struct {
	Artist struct {
		Name string `json:"name"`
	} `json:"artist"`
	Title    string         `validate:"required,min=1" json:"title"`
	Duration utils.Duration `validate:"required" json:"duration"`
	Year     int64          `validate:"required,number,min=1900,max=2024" json:"year"`
	Genre    []string       `validate:"required" json:"genre"`
}

type TrackUpdateRequest struct {
	Artist *struct {
		Name *string `validate:"omitempty" json:"name"`
	} `json:"artist"`
	Title    *string         `validate:"omitempty,min=1" json:"title"`
	Duration *utils.Duration `validate:"omitempty" json:"duration"`
	Year     *int64          `validate:"omitempty,number,min=1900,max=2024" json:"year"`
	Genre    *[]string       `validate:"omitempty" json:"genre"`
}

type UserRegisterActivated struct {
	Token string `validate:"required" json:"token"`
}

type UserPostRequest struct {
	Email    string `validate:"required,email" json:"email"`
	Password string `validate:"required,min=8,max=72" json:"password"`
	Name     string `validate:"required,min=2,max=500" json:"name"`
}
