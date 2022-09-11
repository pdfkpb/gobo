package bank

import (
	"errors"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type BankDB struct {
	db *gorm.DB
}

type Bank struct {
	gorm.Model
	userID int `gorm:"primarykey"`
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

func (bdb *BankDB) AddUser(userID int) {
	bdb.db.Create(&Bank{
		userID: userID,
		funds:  0,
	})
}

func (bdb *BankDB) CheckFunds(userID int) (int, error) {
	var bank Bank
	result := bdb.db.First(&bank, userID)

	if result.Error != nil {
		return 0, result.Error
	}

	return bank.funds, nil
}

func (bdb *BankDB) AddFunds(userID int, amount int) error {
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

func (bdb *BankDB) TakeFunds(userID int, amount int) error {
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
