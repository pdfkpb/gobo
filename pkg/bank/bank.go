package bank

import (
	"errors"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	ErrUserNotRegistered = errors.New("user not found")
	ErrInvalidAmount     = errors.New("invalid monies amount")
	ErrUnhandledError    = errors.New("didn't bother to catch it")
)

type BankDB struct {
	db *gorm.DB
}

type Bank struct {
	gorm.Model
	userID string `gorm:"primarykey"`
	funds  int
}

func LoadBankDB() (*BankDB, error) {
	db, err := gorm.Open(sqlite.Open("bank.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&Bank{})

	return &BankDB{
		db: db,
	}, nil
}

func (bdb *BankDB) AddUser(userID string) {
	bdb.db.Create(&Bank{
		userID: userID,
		funds:  0,
	})
}

func (bdb *BankDB) CheckFunds(userID string) (int, error) {
	var bank Bank
	result := bdb.db.First(&bank, userID)

	if result.Error != nil {
		switch result.Error {
		case gorm.ErrRecordNotFound:
			return 0, ErrUserNotRegistered
		default:
			return 0, ErrUnhandledError
		}
	}

	return bank.funds, nil
}

func (bdb *BankDB) AddFunds(userID string, amount int) (int, error) {
	if amount < 0 {
		return 0, errors.New("")
	}

	var bank Bank
	result := bdb.db.First(&bank, userID)

	if result.Error != nil {
		return 0, result.Error
	}

	bank.funds += amount

	result = bdb.db.Save(bank)
	if result.Error != nil {
		return 0, result.Error
	}

	return bank.funds, nil
}

func (bdb *BankDB) TakeFunds(userID string, amount int) (int, error) {
	if amount < 0 {
		return 0, errors.New("")
	}

	var bank Bank
	result := bdb.db.First(&bank, userID)

	if result.Error != nil {
		return 0, result.Error
	}

	if bank.funds-amount < 0 {
		return 0, errors.New("funds cannot be negative")
	}

	bank.funds -= amount

	result = bdb.db.Save(bank)
	if result.Error != nil {
		return 0, result.Error
	}

	return bank.funds, nil
}
