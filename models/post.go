package models

import "time"

type Post struct {
	Id          string    `json:"id"`
	UserId      string    `json:"userId"`
	PostContent string    `json:"postContent"`
	CreatedAt   time.Time `json:"createdAt"`
}
