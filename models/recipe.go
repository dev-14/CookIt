package models

import (
	"time"

	"github.com/lib/pq"
)

// type Hstore struct {
// 	Map map[string]sql.NullString
// }

type Recipe struct {
	Id          int
	Name        string
	Description string
	Ingredients pq.StringArray `gorm:"type:text[]"`

	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
