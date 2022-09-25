package patron

import "gorm.io/gorm"

type Challenge struct {
	gorm.Model

	SomethingElse int    `gorm:"primaryKey,autoIncrement"`
	Contender     string `gorm:"unique"`

	Escrow int
}
