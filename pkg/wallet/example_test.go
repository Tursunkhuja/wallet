package wallet

import (
	"reflect"
	"testing"

	"github.com/Tursunkhuja/wallet/pkg/types"
)

//"github.com/Tursunkhuja/wallet/pkg/wallet"

func Test_FindAccountById_Exist(t *testing.T) {
	//accountID := uuid.New().String()
	svc := &Service{}
	account, err := svc.RegisterAccount(types.Phone("992928333783"))
	if err != nil {
		t.Errorf("Error on registering account, error = %v", err)
		return
	}

	svc.accounts = append(svc.accounts, account)

	account2, error := svc.FindAccountByID(account.ID)

	if !reflect.DeepEqual(account, account2) {
		t.Errorf("There should be an account with this ID")
	}
	if !reflect.DeepEqual(nil, error) {
		t.Errorf("There should not be any error!")
	}
}

func Test_FindAccountById_NotExist(t *testing.T) {

	//accountID := uuid.New().String()
	svc := &Service{}
	acc2 := types.Account{ID: 2, Phone: "992928303783", Balance: 200}
	svc.accounts = append(svc.accounts, &acc2)

	account, error1 := svc.FindAccountByID(1)

	if reflect.DeepEqual(ErrAccountNotFound.Error(), error1) {
		t.Errorf("There should not be any error!")
	}
	if reflect.DeepEqual(nil, account) {
		t.Errorf("There should not be an account with ID = 1")
	}
}

/*
func Test_Reject_PaymentExist(t *testing.T) {

	paymentID := uuid.New().String()
	svc := &Service{}
	account := types.Account{ID: 1, Phone: "992928303783", Balance: 200}
	payment := types.Payment{
		ID:        paymentID,
		AccountID: 1,
		Amount:    200,
		Status:    types.PaymentStatusOK,
	}

	svc.accounts = append(svc.accounts, &account)
	svc.payments = append(svc.payments, &payment)

	error := svc.Reject(paymentID)

	if !reflect.DeepEqual(nil, error) {
		t.Errorf("There should not be an error, because payment id exists")
	}
}

func Test_Reject_PaymentNotExist(t *testing.T) {

	paymentID := uuid.New().String()
	svc := &Service{}
	account := types.Account{ID: 1, Phone: "992928303783", Balance: 200}
	payment := types.Payment{
		ID:        paymentID,
		AccountID: 1,
		Amount:    200,
		Status:    types.PaymentStatusOK,
	}

	svc.accounts = append(svc.accounts, &account)
	svc.payments = append(svc.payments, &payment)

	error := svc.Reject("wronPaymentID")

	if reflect.DeepEqual(ErrPaymentNotFound.Error(), error) {
		t.Errorf("payment id not found")
	}
}
*/

func Test_Repeat(t *testing.T) {

	//paymentID := uuid.New().String()

	s := &Service{}

	phone := types.Phone("992928303783")
	account, err := s.RegisterAccount(phone)
	if err != nil {
		t.Errorf("Regect(): can't register account, error = %v", err)
		return
	}

	//deposit with negative amont
	err = s.Deposit(account.ID, -1000)
	if err == nil {
		t.Errorf("Regect(): can't deposit with negative amount, error = %v", err)
		return
	}
	//deposit with wrong accout ID
	err = s.Deposit(123, 1000)
	if err == nil {
		t.Errorf("Regect(): can't deposit with wrong account ID, error = %v", err)
		return
	}
	//deposit with correct accout ID
	err = s.Deposit(account.ID, 1000)
	if err != nil {
		t.Errorf("Regect(): can't deposit account, error = %v", err)
		return
	}

	payment, err := s.Pay(account.ID, 100, "auto")
	if err != nil {
		t.Errorf("Regect(): can't create payment, error = %v", err)
		return
	}

	newpayment, err := s.Repeat(payment.ID)
	if err != nil {
		t.Errorf("Regect(): can't copy payment, error = %v", err)
		return
	}

	if !reflect.DeepEqual(types.Money(100), newpayment.Amount) {
		t.Errorf("New payment amount should be 100, because it is copied payment")
	}
}

func Test_FavoritePayment(t *testing.T) {

	//paymentID := uuid.New().String()
	s := &Service{}

	phone := types.Phone("992928303783")
	account, err := s.RegisterAccount(phone)
	if err != nil {
		t.Errorf("Regect(): can't register account, error = %v", err)
		return
	}

	//deposit with negative amont
	err = s.Deposit(account.ID, -1000)
	if err == nil {
		t.Errorf("Regect(): can't deposit with negative amount, error = %v", err)
		return
	}
	//deposit with wrong accout ID
	err = s.Deposit(123, 1000)
	if err == nil {
		t.Errorf("Regect(): can't deposit with wrong account ID, error = %v", err)
		return
	}
	//deposit with correct accout ID
	err = s.Deposit(account.ID, 1000)
	if err != nil {
		t.Errorf("Regect(): can't deposit account, error = %v", err)
		return
	}

	payment, err := s.Pay(account.ID, 100, "auto")
	if err != nil {
		t.Errorf("Regect(): can't create payment, error = %v", err)
		return
	}

	favorite, err := s.FavoritePayment(payment.ID, "New payment")
	if err != nil {
		t.Errorf("Regect(): can't add payment to favorite, error = %v", err)
		return
	}

	if !reflect.DeepEqual(types.Money(100), favorite.Amount) {
		t.Errorf("Favorite amount should be 100, because it is copied from payment")
	}

	//favorite with wrong ID
	_, hasError := s.PayFromFavorite("1")
	if hasError == nil {
		t.Errorf("Regect(): can't pay with wrong favorite ID, error = %v", err)
		return
	}
	//with correct favorite ID
	payFromFavorite, err := s.PayFromFavorite(favorite.ID)
	if err != nil {
		t.Errorf("Regect(): can't pay from favorite, error = %v", err)
		return
	}

	if !reflect.DeepEqual(types.Money(100), payFromFavorite.Amount) {
		t.Errorf("Payment amount should be 100, because it is copied from favorite")
	}
}
