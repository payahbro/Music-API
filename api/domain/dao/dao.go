package dao

import (
	"music-echo/utils"
	"time"
)

type Artists struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type Tracks struct {
	Id        int64          `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	IdArtist  int64          `json:"id_artist"`
	Title     string         `json:"title"`
	Duration  utils.Duration `json:"duration"`
	Year      int64          `json:"year"`
	Genre     []string       `json:"genre"`
	Version   int64          `json:"version"`
}

type Likes struct {
	IdUsers  int64 `json:"id_users"`
	IdTracks int64 `json:"id_tracks"`
}

type Users struct {
	Id        int64          `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	Name      string         `json:"name"`
	Email     string         `json:"email"`
	Password  utils.Password `json:"-"`
	Activated bool           `json:"activated"`
	Version   int            `json:"-"`
}

type Token struct {
	Hash   []byte
	UserId int64
	Expiry time.Time
	Scope  string
}
