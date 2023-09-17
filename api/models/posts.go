package models

import "time"

type Post struct {
    ID        string    `json:"id"`
    File      string    `json:"file"`
    Title     string    `json:"title"`
    Author    string    `json:"author"`
    Content   string    `json:"content"`
    Comments  []string  `json:"comments"`
    Likes     int64     `json:"likes"`
    CreatedAt time.Time `json:"createdAt"`
}
