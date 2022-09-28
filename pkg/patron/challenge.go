package patron

import "gorm.io/gorm"

type Challenge struct {
	gorm.Model

	Challenger string `gorm:"primaryKey;unique"`
	Contender  string
	Escrow     int
}
