package models

import "gorm.io/gorm"

type Follower struct {
	gorm.Model
	Who_id  int `json:"who_id"`
	Whom_id int `json:"whom_id"`
}
