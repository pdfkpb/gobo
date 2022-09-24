package patron

import "gorm.io/gorm"

type Challenge struct {
	gorm.Model

	ID        int    `gorm:"primaryKey,autoIncrement"`
	Contender string `gorm:"unique"`

	Escrow int
}
