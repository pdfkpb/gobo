package patron

import "gorm.io/gorm"

type Challenge struct {
	gorm.Model

	Contender string `gorm:"unique"`

	Escrow int
}
