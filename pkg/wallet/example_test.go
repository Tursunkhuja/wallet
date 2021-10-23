package wallet

import (
	"reflect"
	"testing"

	"github.com/Tursunkhuja/wallet/pkg/types"
	"github.com/google/uuid"
)

//"github.com/Tursunkhuja/wallet/pkg/wallet"

func Test_FindAccountById_Exist(t *testing.T) {
	//accountID := uuid.New().String()
	svc := Service{}
	acc1 := types.Account{ID: 1, Phone: "992928333783", Balance: 100}
	svc.accounts = append(svc.accounts, &acc1)

	account, error := svc.FindAccountByID(1)

	if !reflect.DeepEqual(&acc1, account) {
		t.Errorf("There should be an account with ID = 1")
	}
	if !reflect.DeepEqual(nil, error) {
		t.Errorf("There should not be any error!")
	}
}

func Test_FindAccountById_NotExist(t *testing.T) {

	//accountID := uuid.New().String()
	svc := Service{}
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

func Test_Reject_PaymentExist(t *testing.T) {

	paymentID := uuid.New().String()
	svc := Service{}
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
	svc := Service{}
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
