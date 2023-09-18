package models

import (
    "time"

    "github.com/google/uuid"
)

type Post struct {
    ID             string    `json:"id"`
    File           string    `json:"file"`
    Title          string    `json:"title"`
    Author         string    `json:"author"`
    Content        string    `json:"content"`
    Comments       []string  `json:"comments"`
    Likes          int64     `json:"likes"`
    CreatedAt      time.Time `json:"createdAt"`
    IsTerminalFile bool      `json:"isTerminalFile"`
}

func NewPost(file string, title string, author string, content string) *Post {
    return &Post{
        ID:             uuid.New().String(),
        File:           file,
        Title:          title,
        Content:        content,
        Likes:          0,
        CreatedAt:      time.Now(),
        IsTerminalFile: true,
    }
}

func (p *Post) Create() error {
    return nil
}
