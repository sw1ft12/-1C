package models

type Dish struct {
    Name         string `json:"name,omitempty" db:"name" default:"-"`
    Calorie      int    `json:"calorie,omitempty" db:"calorie" default:"-1"`
    Protein      int    `json:"protein,omitempty" db:"protein" default:"-1"`
    Fat          int    `json:"fat,omitempty" db:"fat" default:"-1"`
    Carbohydrate int    `json:"carbohydrate,omitempty" db:"carbohydrate" default:"-1"`
}
