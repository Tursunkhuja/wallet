package wallet

import (
	"errors"

	"github.com/Tursunkhuja/wallet/pkg/types"
	"github.com/google/uuid"
	//"github.com/Tursunkhuja/wallet/pkg/wallet"
)

type Service struct {
	nextAccountID int64
	accounts      []*types.Account
	payments      []*types.Payment
}

var ErrAccountNotFound = errors.New("account not found")
var ErrNotEnoughBalance = errors.New("not enough balance")
var ErrAmountMustBePositive = errors.New("amount must be greater that zero")
var ErrPaymentNotFound = errors.New("payment not found")

type Error string

func (e Error) Error() string {
	return string(e)
}

func (service *Service) FindAccountByID(accountID int64) (*types.Account, error) {
	for _, account := range service.accounts {
		if account.ID == accountID {
			return account, nil
		}
	}
	return nil, ErrAccountNotFound
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

func (s *Service) Pay(accountID int64, amount types.Money, category types.PaymentCategory) (*types.Payment, error) {
	if amount <= 0 {
		return nil, ErrAmountMustBePositive
	}

	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}

	if account == nil {
		return nil, ErrAccountNotFound
	}

	if account.Balance < amount {
		return nil, ErrNotEnoughBalance
	}

	account.Balance -= amount
	paymentID := uuid.New().String()

	payment := &types.Payment{
		ID:        paymentID,
		AccountID: accountID,
		Amount:    amount,
		Category:  category,
		Status:    types.PaymentStatusInProgress,
	}

	s.payments = append(s.payments, payment)

	return payment, nil
}

func (s *Service) Reject(paymentID string) error {
	payment, err := s.FindPaymentByID(paymentID)
	if payment != nil {
		payment.Status = types.PaymentStatusFail
		account, _ := s.FindAccountByID(payment.AccountID)
		if account != nil {
			account.Balance += payment.Amount
		}
		return nil
	}

	return err
}

func (s *Service) FindPaymentByID(paymentID string) (*types.Payment, error) {
	for _, payment := range s.payments {
		if payment.ID == paymentID {
			return payment, nil
		}
	}
	return nil, ErrPaymentNotFound
}

func (receiver *Service) Deposit(accountID int64, amount types.Money) error {
	if amount <= 0 {
		return ErrAmountMustBePositive
	}
	var account *types.Account
	for _, acc := range receiver.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}

	if account == nil {
		return ErrAccountNotFound
	}

	account.Balance += amount
	return nil
}
