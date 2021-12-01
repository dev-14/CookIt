package models

import (
	"time"
)

// type Hstore struct {
// 	Map map[string]sql.NullString
// }

type Recipe struct {
	Id          int
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Ingredients []string `gorm:"type:text[]" json:"ingredients"`

	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
