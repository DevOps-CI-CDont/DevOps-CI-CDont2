package models

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	Author_id   int    `json:"author_id"`
	Text        string `json:"text"`
	Pub_date    int    `json:"pub_date"`
	Flagged     int    `json:"flagged"`
	Author_name string `json:"author_name"`
}
