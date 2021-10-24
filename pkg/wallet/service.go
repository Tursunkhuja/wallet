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
var ErrPhoneAlreadyRegistered = errors.New("phone already registered")

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
			return nil, ErrPhoneAlreadyRegistered
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

	account, err := s.FindAccountByID(accountID)
	if err != nil {
		return nil, err
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

/*
func (s *Service) Reject(paymentID string) error {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return err
	}
	account, err := s.FindAccountByID(payment.AccountID)
	if err != nil {
		return err
	}
	payment.Status = types.PaymentStatusFail
	account.Balance += payment.Amount
	return err
}
*/
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
	account, err := receiver.FindAccountByID(accountID)
	if err != nil {
		return err
	}
	account.Balance += amount
	return nil
}

func (s *Service) Repeat(paymentID string) (*types.Payment, error) {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, err
	}

	newpayment, err := s.Pay(payment.AccountID, payment.Amount, payment.Category)
	if err != nil {
		return nil, err
	}
	return newpayment, nil
}
