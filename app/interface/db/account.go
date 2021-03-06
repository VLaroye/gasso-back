package db

import (
	"github.com/VLaroye/gasso-back/app/domain/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Account struct {
	ID   string
	Name string
}

func NewAccount(id, name string) *Account {
	return &Account{
		ID:   id,
		Name: name,
	}
}

type accountRepository struct {
	db     *gorm.DB
	logger *zap.SugaredLogger
}

func NewAccountRepository(db *gorm.DB, logger *zap.SugaredLogger) *accountRepository {
	logger.Debugw("creating postgres AccountRepository")
	return &accountRepository{
		db:     db,
		logger: logger,
	}
}

func (ar *accountRepository) FindAll() ([]*model.Account, error) {
	var accounts []*Account

	result := ar.db.Find(&accounts)

	if result.Error != nil {
		ar.logger.Errorw("error fetching list accounts from db",
			"error", result.Error,
		)
		return nil, result.Error
	}

	ar.logger.Infow("list accounts fetched from db",
		"nb of accounts fetched", result.RowsAffected,
	)

	response := make([]*model.Account, len(accounts))

	for i, account := range accounts {
		response[i] = model.NewAccount(account.ID, account.Name)
	}

	return response, nil
}

func (ar *accountRepository) FindByName(name string) (*model.Account, error) {
	var account Account
	result := ar.db.Where("name = ?", name).Find(&account)

	if result.Error != nil {
		ar.logger.Errorw("error fechting account by name",
			"name", name,
		)
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		ar.logger.Infow("account not found",
			"name", name,
		)
		return nil, nil
	}

	ar.logger.Infow("account fetched by name",
		"name", name,
	)
	return model.NewAccount(account.ID, account.Name), nil
}

func (ar *accountRepository) FindByID(id string) (*model.Account, error) {
	var account Account
	result := ar.db.Where("id = ?", id).Find(&account)

	if result.Error != nil {
		ar.logger.Errorw("find account by id failed, account not found",
			"id", id,
		)
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		ar.logger.Infow("find account by id failed, account not found",
			"id", id,
		)
		return nil, nil
	}

	ar.logger.Infow("account fetched by id",
		"id", id,
	)
	return model.NewAccount(account.ID, account.Name), nil
}

func (ar *accountRepository) Create(id, name string) error {
	accountToInsert := NewAccount(id, name)

	result := ar.db.Create(&accountToInsert)

	if result.Error != nil {
		ar.logger.Errorw("create account failed",
			"accountId", accountToInsert.ID,
			"error", result.Error,
		)
		return result.Error
	}

	ar.logger.Infow("account created",
		"id", accountToInsert.ID,
		"name", accountToInsert.Name,
	)
	return nil
}

func (ar *accountRepository) Update(id, name string) error {
	accountToUpdate := NewAccount(id, name)

	result := ar.db.Save(&accountToUpdate)

	if result.Error != nil {
		ar.logger.Errorw("update account failed",
			"accountId", id,
			"accountName", name,
			"error", result.Error,
		)
		return result.Error
	}

	ar.logger.Infow("account updated",
		"id", accountToUpdate.ID,
		"name", accountToUpdate.Name,
	)
	return nil
}

func (ar *accountRepository) Delete(id string) error {
	result := ar.db.Where("id = ?", id).Delete(&Account{})

	if result.Error != nil {
		ar.logger.Errorw("delete account failed",
			"accountId", id,
			"error", result.Error,
		)
		return result.Error
	}

	ar.logger.Infow("account deleted",
		"id", id,
	)
	return nil
}
