package schema

import "time"

type Wine struct {
	ID        int       `json:"id"`
	ImgUrl    string    `json:"img_url"`
	Name      string    `json:"name"`
	Price     string    `json:"price"`
	Category  string    `json:"category"`
	Region    string    `json:"region"`
	Flag      string    `json:"flag"`
	Rating    string    `json:"rating"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
