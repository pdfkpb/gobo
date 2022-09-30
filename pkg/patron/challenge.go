package patron

type Challenge struct {
	Challenger string `gorm:"primaryKey;unique"`
	Contender  string
	Escrow     int
}
