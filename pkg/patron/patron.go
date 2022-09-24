package patron

import (
	"errors"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	ErrUserNotRegistered       = errors.New("user not found")
	ErrChallengeNotFound       = errors.New("patron has not outstanding challenges")
	ErrorUserAlreadyRegistered = errors.New("this user is already registered")
	ErrInvalidAmount           = errors.New("invalid monies amount")
	ErrFundsCannotBeNeg        = errors.New("Funds cannot be negative")
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
	Challenge   Challenge `gorm:"foreignKey:ID"`
}

func LoadPatronDB() (*PatronDB, error) {
	db, err := gorm.Open(sqlite.Open("patron.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&Patron{})
	if err != nil {

	}

	err = db.AutoMigrate(&Challenge{})
	if err != nil {

	}

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
		return 0, ErrFundsCannotBeNeg
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
		switch result.Error {
		case gorm.ErrRecordNotFound:
			return []string{}, 0, nil
		default:
			return []string{}, 0, ErrUnhandledError
		}
	}

	var winners []string
	for _, patron := range patrons {
		winners = append(winners, patron.UserID)
	}

	return winners, patrons[0].LotteryRoll, nil
}

func (pdb *PatronDB) ClearLottery() error {
	result := pdb.db.Model(&Patron{}).Where("1 = 1").Update("lottery_roll", 0)
	if result.Error != nil {
		return ErrUnhandledError
	}

	return nil
}

// Challenge Functions

func (pdb *PatronDB) CreateChallenge(userID string, contender string, amount int) error {
	if amount < 0 {
		return ErrInvalidAmount
	}

	var patron Patron
	result := pdb.db.First(&patron, "user_id = ?", userID)

	if result.Error != nil {
		return result.Error
	}

	if patron.Funds-amount < 0 {
		return errors.New("Funds cannot be negative")
	}

	patron.Funds -= amount
	patron.Challenge = Challenge{
		Contender: contender,
		Escrow:    amount,
	}

	result = pdb.db.Model(&Patron{}).Where("1 = 1").Updates(patron)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (pdb *PatronDB) GetChallenge(userID string) (int, error) {
	var patron Patron
	result := pdb.db.First(&patron, "user_id = ?", userID)

	if result.Error != nil {
		switch result.Error {
		case gorm.ErrRecordNotFound:
			return 0, ErrChallengeNotFound
		default:
			return 0, ErrUnhandledError
		}
	}

	if patron.Challenge.Contender == "" {
		return 0, ErrChallengeNotFound
	}

	return patron.Challenge.Escrow, nil
}

func (pdb *PatronDB) ClearChallenge(userID string) error {
	var patron Patron
	result := pdb.db.First(&patron, "user_id = ?", userID)

	if result.Error != nil {
		switch result.Error {
		case gorm.ErrRecordNotFound:
			return ErrChallengeNotFound
		default:
			return ErrUnhandledError
		}
	}

	if patron.Challenge.Contender == "" {
		return ErrChallengeNotFound
	}

	patron.Challenge = Challenge{}

	result = pdb.db.Model(&Patron{}).Where("1 = 1").Updates(patron)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
