package models

import "time"

type Person struct {
	ID        int64
	Name      string
	Birthdate time.Time
}
