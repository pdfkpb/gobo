package patron

import (
	"errors"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	ErrUserNotRegistered       = errors.New("user not found")
	ErrorUserAlreadyRegistered = errors.New("this user is already registered")
	ErrInvalidAmount           = errors.New("invalid monies amount")
	ErrAlreadyLotteryRolled    = errors.New("this user already rolled for the lottery")
	ErrUnhandledError          = errors.New("didn't bother to catch it")
)

type PatronDB struct {
	db *gorm.DB
}

type Patron struct {
	gorm.Model
	UserID      string `gorm:"primaryKey"`
	Funds       int
	LotteryRoll int
}

func LoadPatronDB() (*PatronDB, error) {
	db, err := gorm.Open(sqlite.Open("patron.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&Patron{})

	return &PatronDB{
		db: db,
	}, nil
}

func (pdb *PatronDB) AddUser(userID string) error {
	result := pdb.db.FirstOrCreate(&Patron{
		UserID: userID,
		Funds:  1000,
	})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected > 0 {
		return nil
	} else {
		return ErrorUserAlreadyRegistered
	}
}

func (pdb *PatronDB) CheckFunds(userID string) (int, error) {
	var patron Patron
	result := pdb.db.First(&patron, "user_id = ?", userID)

	if result.Error != nil {
		switch result.Error {
		case gorm.ErrRecordNotFound:
			return 0, ErrUserNotRegistered
		default:
			return 0, ErrUnhandledError
		}
	}

	return patron.Funds, nil
}

func (pdb *PatronDB) AddFunds(userID string, amount int) (int, error) {
	if amount < 0 {
		return 0, ErrInvalidAmount
	}

	var patron Patron
	result := pdb.db.First(&patron, "user_id = ?", userID)

	if result.Error != nil {
		return 0, result.Error
	}

	result = pdb.db.Model(&patron).Update("funds", patron.Funds+amount)
	if result.Error != nil {
		return 0, result.Error
	}

	return patron.Funds, nil
}

func (pdb *PatronDB) TakeFunds(userID string, amount int) (int, error) {
	if amount < 0 {
		return 0, ErrInvalidAmount
	}

	var patron Patron
	result := pdb.db.First(&patron, "user_id = ?", userID)

	if result.Error != nil {
		return 0, result.Error
	}

	if patron.Funds-amount < 0 {
		return 0, errors.New("Funds cannot be negative")
	}

	result = pdb.db.Model(&patron).Update("funds", patron.Funds-amount)
	if result.Error != nil {
		return 0, result.Error
	}

	return patron.Funds, nil
}

// Lottery Functions

func (pdb *PatronDB) SetLotteryRoll(userID string, roll int) error {
	var patron Patron
	result := pdb.db.Where("lottery_roll = ?", 0).First(&patron, "user_id = ?", userID)
	if result.Error != nil {
		return ErrAlreadyLotteryRolled
	}

	result = pdb.db.Model(&patron).Update("lottery_roll", roll)
	if result.Error != nil {
		return ErrUnhandledError
	}

	return nil
}

func (pdb *PatronDB) GetLotteryWinner() ([]string, int, error) {
	var patrons []Patron
	result := pdb.db.Order("lottery_roll desc").Where("lottery_roll > ?", 0).First(&patrons)
	if result.Error != nil {
		return []string{}, 0, ErrUnhandledError
	}

	var winners []string
	for _, patron := range patrons {
		winners = append(winners, patron.UserID)
	}

	return winners, patrons[0].LotteryRoll, nil
}

func (pdb *PatronDB) ClearLottery() error {
	result := pdb.db.Model(Patron{}).Update("lottery_roll", 0)
	if result.Error != nil {
		return ErrUnhandledError
	}

	return nil
}
