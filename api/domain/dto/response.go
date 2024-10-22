package dto

import (
	"music-echo/api/domain/dao"
	"time"
)

type MetadataResponse struct {
	CurrentPage int64 `json:"current_page,omitempty"`
	PageSize    int64 `json:"page_size,omitempty"`
	FirstPage   int64 `json:"first_page,omitempty"`
	LastPage    int64 `json:"last_page,omitempty"`
	TotalRecord int64 `json:"total_record,omitempty"`
}

type TrackInsertResponse struct {
	Id        int64     `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
}

type TrackGetResponse struct {
	Track  *dao.Tracks  `json:"track"`
	Artist *dao.Artists `json:"artist"`
	Likes  *int64       `json:"likes"`
}

type TrackUpdateResponse struct {
	Track  *dao.Tracks  `json:"track"`
	Artist *dao.Artists `json:"artist"`
	Likes  *int64       `json:"likes"`
}

type TrackGetAllResponse struct {
	Track  *dao.Tracks  `json:"track"`
	Artist *dao.Artists `json:"artist"`
	Likes  int64        `json:"likes"`
}

type UsersCreateResponse struct {
	Id        int64     `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Version   int       `json:"version"`
}

type WebResponse struct {
	Message  string           `json:"message,omitempty"`
	Metadata MetadataResponse `json:"metadata,omitempty"`
	Data     interface{}      `json:"data,omitempty"`
}
