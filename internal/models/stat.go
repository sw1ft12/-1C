package models

import "time"

type Stat struct {
    Date        time.Time `json:"date" db:"date"`
    SumCalories int       `json:"sum_calories" db:"sum"`
}
