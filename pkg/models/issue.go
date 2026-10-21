package models

import (
	"strings"
	"time"
)

type Label struct {
	ID          int64  `json:"id,omitempty"`
	Name        string `json:"name"`
	Color       string `json:"color,omitempty"`
	Description string `json:"description,omitempty"`
}

type User struct {
	Login     string `json:"login"`
	ID        int64  `json:"id,omitempty"`
	AvatarURL string `json:"avatar_url,omitempty"`
	HTMLURL   string `json:"html_url,omitempty"`
}

type Issue struct {
	ID            int64      `json:"id"`
	Number        int        `json:"number"`
	Title         string     `json:"title"`
	Body          string     `json:"body"`
	HTMLURL       string     `json:"html_url"`
	State         string     `json:"state"`
	CommentsCount int        `json:"comments"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	ClosedAt      *time.Time `json:"closed_at,omitempty"`
	User          User       `json:"user"`
	Assignee      *User      `json:"assignee,omitempty"`
	Assignees     []User     `json:"assignees,omitempty"`
	Labels        []Label    `json:"labels"`
}
