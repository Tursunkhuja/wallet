package wallet

import (
	"github.com/Tursunkhuja/wallet/pkg/types"
	//"github.com/Tursunkhuja/wallet/pkg/wallet"
)

type Service struct {
	nextAccountID int64
	accounts      []*types.Account
	//payments      []*types.Payment
}

type ErrAccountNotFound string

func (e ErrAccountNotFound) Error() string {
	return string(e)
}

func (service *Service) FindAccountByID(accountID int64) (*types.Account, error) {
	for _, account := range service.accounts {
		if account.ID == accountID {
			return account, nil
		}
	}
	return nil, ErrAccountNotFound("account not found!")
}

type Error string

func (e Error) Error() string {
	return string(e)
}
func (service *Service) RegisterAccount(phone types.Phone) (*types.Account, error) {
	for _, account := range service.accounts {
		if account.Phone == phone {
			return nil, Error("phone already registered")
		}
	}

	service.nextAccountID++
	account := &types.Account{
		ID:      service.nextAccountID,
		Phone:   phone,
		Balance: 0,
	}
	service.accounts = append(service.accounts, account)

	return account, nil
}

/*
func (receiver *Service) Deposit(accountID int64, amount types.Money) error {
	if amount <= 0 {
		return Error("amount must be greater than 0")
	}
	var account *types.Account
	for _, acc := range receiver.accounts {
		if acc.ID == accountID {
			ac++c+ount = acc
			break
		}
	}

	if account == nil {
		return Error("account not found")
	}

	account.Balance += amount
	return nil
}
*/
