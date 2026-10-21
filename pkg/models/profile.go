package models

import "time"

type Profile struct {
	Level     DifficultyLevel `json:"level"`
	Languages []string        `json:"languages"`
	Interests []string        `json:"interests"`
	MinStars  int             `json:"min_stars,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

func DefaultProfile() *Profile {
	return &Profile{
		Level:     DifficultyBeginner,
		Languages: []string{},
		Interests: []string{},
		MinStars:  0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
