package bank

import (
	"errors"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	ErrUserNotRegistered = errors.New("user not found")
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

func (bdb *BankDB) AddFunds(userID string, amount int) error {
	if amount < 0 {
		return errors.New("")
	}

	var bank Bank
	result := bdb.db.First(&bank, userID)

	if result.Error != nil {
		return result.Error
	}

	bank.funds += amount

	return bdb.db.Save(bank).Error
}

func (bdb *BankDB) TakeFunds(userID string, amount int) error {
	if amount < 0 {
		return errors.New("")
	}

	var bank Bank
	result := bdb.db.First(&bank, userID)

	if result.Error != nil {
		return result.Error
	}

	if bank.funds-amount < 0 {
		return errors.New("funds cannot be negative")
	}

	bank.funds -= amount

	return bdb.db.Save(bank).Error
}
