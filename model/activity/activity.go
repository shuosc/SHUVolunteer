package activity

import "time"

type Activity struct {
	Date time.Time `json:"date"`
	Name string    `json:"name"`
	Team string    `json:"team"`
}
