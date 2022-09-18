package bank

import (
	"errors"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	ErrUserNotRegistered       = errors.New("user not found")
	ErrorUserAlreadyRegistered = errors.New("this user is already registered")
	ErrInvalidAmount           = errors.New("invalid monies amount")
	ErrUnhandledError          = errors.New("didn't bother to catch it")
)

type BankDB struct {
	db *gorm.DB
}

type bank struct {
	gorm.Model
	UserID string `gorm:"primaryKey"`
	Funds  int
}

func LoadBankDB() (*BankDB, error) {
	db, err := gorm.Open(sqlite.Open("bank.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&bank{})

	return &BankDB{
		db: db,
	}, nil
}

func (bdb *BankDB) AddUser(userID string) error {
	result := bdb.db.FirstOrCreate(&bank{
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

func (bdb *BankDB) CheckFunds(userID string) (int, error) {
	var bank bank
	result := bdb.db.First(&bank, "user_id = ?", userID)

	if result.Error != nil {
		switch result.Error {
		case gorm.ErrRecordNotFound:
			return 0, ErrUserNotRegistered
		default:
			return 0, ErrUnhandledError
		}
	}

	return bank.Funds, nil
}

func (bdb *BankDB) AddFunds(userID string, amount int) (int, error) {
	if amount < 0 {
		return 0, ErrInvalidAmount
	}

	var bank bank
	result := bdb.db.First(&bank, "user_id = ?", userID)

	if result.Error != nil {
		return 0, result.Error
	}

	result = bdb.db.Model(&bank).Update("funds", bank.Funds+amount)
	if result.Error != nil {
		return 0, result.Error
	}

	return bank.Funds, nil
}

func (bdb *BankDB) TakeFunds(userID string, amount int) (int, error) {
	if amount < 0 {
		return 0, ErrInvalidAmount
	}

	var bank bank
	result := bdb.db.First(&bank, "user_id = ?", userID)

	if result.Error != nil {
		return 0, result.Error
	}

	if bank.Funds-amount < 0 {
		return 0, errors.New("Funds cannot be negative")
	}

	result = bdb.db.Model(&bank).Update("funds", bank.Funds-amount)
	if result.Error != nil {
		return 0, result.Error
	}

	return bank.Funds, nil
}
