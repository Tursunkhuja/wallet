package wallet

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/Tursunkhuja/wallet/pkg/types"
	"github.com/google/uuid"
)

func TestService_RegisterAccount(t *testing.T) {
	s, err := generateTestData(10)
	if err != nil {
		t.Errorf("Error generate TEST data: %v", err)
		return
	}

	_, err = s.RegisterAccount("+7-777-455-82-01")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestService_RegisterAccount_fail(t *testing.T) {
	s, err := generateTestData(10)
	if err != nil {
		t.Errorf("Error generate TEST data: %v", err)
		return
	}

	_, err = s.RegisterAccount("+992-92-833-37-83")
	if err == nil {
		t.Errorf("err expected:%v, actual:%v", ErrPhoneRegistered, nil)
		return
	}
}

func TestService_Reject_success(t *testing.T) {
	s := newTestService()

	// check for registration
	_, payments, err := s.addAccount(defultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	// check for deposite
	payment := payments[0]
	err = s.Reject(payment.ID)
	if err != nil {
		t.Errorf("Reject(): error = %v", err)
		return
	}

	savedPayment, err := s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("Reject(): can't find payment by ID, error = %v", err)
		return
	}
	if savedPayment.Status != types.PaymentStatusFail {
		t.Errorf("Reject(): status didn't change, payment = %v", savedPayment)
		return
	}
	// check for account

	savedAccount, err := s.FindAccountByID(payment.AccountID)
	if err != nil {
		t.Errorf("Reject(): can't find account by ID, error = %v", err)
		return
	}

	if savedAccount.Balance != defultTestAccount.balance {
		t.Errorf("Reject(): balance didn't change, payment = %v", savedPayment)
		return
	}

}

func TestSevice_FindPaymentByID_seccess(t *testing.T) {

	s := newTestService()

	// check for registration
	_, payments, err := s.addAccount(defultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	// check for FindPaymentByID
	payment := payments[0]
	got, err := s.FindPaymentByID(payment.ID)

	if err != nil {
		t.Errorf("FindPaymentByID(): error = %v", err)
		return
	}

	if !reflect.DeepEqual(got, payment) {
		t.Errorf("FindPaymentByID(): payment is not found, got %v want %v ERROR = %v", err, got, payment)
		return
	}
}

func TestService_FindPaymentByID_fail(t *testing.T) {
	// creat a service
	s := newTestService()

	_, _, err := s.addAccount(defultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = s.FindPaymentByID(uuid.New().String())
	if err == nil {
		t.Errorf("FindPaymentByID(): must return error, returned nil")
	}

	if err != ErrPaymentNotFound {
		t.Errorf("FindPaymentByID(): must return ErrPaymentNotFound, returned = %v", err)
	}

}

func TestService_Repeated_success(t *testing.T) {
	s := newTestService()

	// check for registration
	_, payments, err := s.addAccount(defultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	// check for deposite
	payment := payments[0]
	repeatedPayment, err := s.Repeat(payment.ID)
	if err != nil {
		t.Errorf("Repeated(): error = %v", err)
		return
	}
	if payment.Amount != repeatedPayment.Amount {
		t.Errorf("Repeated(): payment amount has changed = %v", payment)
		return
	}

	if payment.AccountID != repeatedPayment.AccountID {
		t.Errorf("Repeated(): payment AccountID has changed = %v", payment)
		return
	}

	if payment.Category != repeatedPayment.Category {
		t.Errorf("Repeated(): payment category has changed = %v", payment)
		return
	}

	if payment.Status != repeatedPayment.Status {
		t.Errorf("Repeated(): payment status has changed = %v", payment)
		return
	}

}

func TestService_FindFavoriteByID_success(t *testing.T) {
	s := newTestService()

	_, payments, err := s.addAccount(defultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	payment := payments[0]

	favorite, err := s.FavoritePayment(payment.ID, "fast-food")

	if err != nil {
		t.Errorf("Favorite(): error = %v", err)
		return
	}

	got, err := s.FindFavoriteByID(favorite.ID)
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(got, favorite) {
		t.Errorf("FindFavoriteByID(): favorite is not found, got %v want %v ERROR = %v", err, got, favorite)
		return
	}
}

func TestService_Deposit(t *testing.T) {
	s, err := generateTestData(10)
	if err != nil {
		t.Errorf("Error generate TEST data: %v", err)
		return
	}

	accTest := s.accounts[0]
	err = s.Deposit(accTest.ID, 50_000)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestService_Pay(t *testing.T) {
	s, err := generateTestData(10)
	if err != nil {
		t.Errorf("Error generate TEST data: %v", err)
		return
	}

	accTest := s.accounts[0]
	_, err = s.Pay(accTest.ID, 5_000, "Sport")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestService_FindFavoriteByID_fail(t *testing.T) {
	s := newTestService()

	_, payments, err := s.addAccount(defultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	payment := payments[0]

	_, err = s.FavoritePayment(payment.ID, "fast-food")

	if err != nil {
		t.Errorf("Favorite(): error = %v", err)
		return
	}

	_, err = s.FindFavoriteByID(uuid.New().String())
	if err == nil {
		t.Errorf("FindFavoriteByID(): must return error, returned nil")
	}

	if err != ErrFavoriteNotFound {
		t.Errorf("FindFavoriteByID(): must return ErrFavoriteNotFound, returned = %v", err)
	}
}

func TestService_FavoritePayment_success(t *testing.T) {
	s := newTestService()

	_, payments, err := s.addAccount(defultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	payment := payments[0]

	favorite, err := s.FavoritePayment(payment.ID, "fast-food")

	if err != nil {
		t.Errorf("Favorite(): error = %v", err)
		return
	}

	if favorite.AccountID != payment.AccountID {
		t.Errorf("Favorite(): account ID changed = %v", favorite)
	}

}

func TestService_FavoritePayment_fail(t *testing.T) {

	s := newTestService()

	_, payments, err := s.addAccount(defultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	payment := payments[0]

	favorite, err := s.FavoritePayment(payment.ID, "fast-food")

	if err != nil {
		t.Errorf("Favorite(): error = %v", err)
		return
	}

	if favorite.AccountID != payment.AccountID {
		t.Errorf("Favorite(): account ID changed = %v", favorite)
	}

}

func TestService_PayFromFavorite_success(t *testing.T) {
	s := newTestService()

	_, payments, err := s.addAccount(defultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	payment := payments[0]

	favorite, err := s.FavoritePayment(payment.ID, "fast-food")

	if err != nil {
		t.Errorf("Favorite(): error = %v", err)
		return
	}

	if favorite.AccountID != payment.AccountID {
		t.Errorf("Favorite(): account ID changed = %v", favorite)
		return
	}

	favoritePayment, err := s.PayFromFavorite(favorite.ID)

	if err != nil {
		t.Errorf("PayFromFavorite(): error = %v", err)
		return
	}

	if favoritePayment.Amount != payment.Amount {
		t.Errorf("PayFromFavorite(): payment amount changed = %v", favorite)
		return
	}
}

type testService struct {
	*Service
}

func newTestService() *testService {
	return &testService{Service: &Service{}}
}

type testAccount struct {
	phone    types.Phone
	balance  types.Money
	payments []struct {
		amount   types.Money
		category types.PaymentCategory
	}
}

var defultTestAccount = testAccount{
	phone:   "+998970009113",
	balance: 10_000_00,
	payments: []struct {
		amount   types.Money
		category types.PaymentCategory
	}{
		{amount: 1_000_00, category: "auto"},
	},
}

func (s *testService) addAccount(data testAccount) (*types.Account, []*types.Payment, error) {
	// register the user
	account, err := s.RegisterAccount(data.phone)
	if err != nil {
		return nil, nil, fmt.Errorf("can't register account, error = %v", err)
	}

	err = s.Deposit(account.ID, data.balance)
	if err != nil {
		return nil, nil, fmt.Errorf("can't deposite account, error = %v", err)
	}

	payments := make([]*types.Payment, len(data.payments))
	for i, payment := range data.payments {
		payments[i], err = s.Pay(account.ID, payment.amount, payment.category)
		if err != nil {
			return nil, nil, fmt.Errorf("can't make payment, error = %v", err)
		}
	}
	return account, payments, nil
}

func TestService_SumPayments(t *testing.T) {
	s, err := generateTestData(10)
	if err != nil {
		t.Errorf("Error generate TEST data: %v", err)
		return
	}

	sum := s.SumPayments(5)
	if sum != types.Money(len(s.payments)) {
		t.Errorf("sum expected:%v, actual:%v", len(s.payments), sum)
	}

}
func TestService_SumPaymentsRegular(t *testing.T) {
	s, err := generateTestData(10)
	if err != nil {
		t.Errorf("Error generate TEST data: %v", err)
		return
	}

	sum := s.SumPaymentsRegular()
	if sum != types.Money(len(s.payments)) {
		t.Errorf("sum expected:%v, actual:%v", len(s.payments), sum)
	}
}

func BenchmarkService_SumPayments(b *testing.B) {
	s, err := generateTestData(10)
	if err != nil {
		b.Errorf("Error generate TEST data: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sum := s.SumPayments(5)
		b.StopTimer()
		if sum != types.Money(len(s.payments)) {
			b.Errorf("sum expected:%v, actual:%v", len(s.payments), sum)
		}
		b.StartTimer()
	}

}

func BenchmarkService_SumPaymentsRegular(b *testing.B) {
	s, err := generateTestData(10)

	if err != nil {
		b.Errorf("Error generate TEST data: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sum := s.SumPaymentsRegular()
		b.StopTimer()
		if sum != types.Money(len(s.payments)) {
			b.Errorf("sum expected:%v, actual:%v", len(s.payments), sum)
		}
		b.StartTimer()
	}

}

func generateTestData(num int) (*Service, error) {
	s := &Service{}

	acc, err := s.RegisterAccount("+992-92-833-37-83")
	if err != nil {
		return nil, err
	}

	err = s.Deposit(acc.ID, 1_000_000_000)
	if err != nil {
		return nil, err
	}

	var pay *types.Payment

	for i := 0; i < num; i++ {
		pay, err = s.Pay(acc.ID, 1, "Auto")
		if err != nil {
			return nil, err
		}
	}

	_, err = s.FavoritePayment(pay.ID, "testPayment")
	if err != nil {
		return nil, err
	}

	return s, nil
}

func TestService_Export(t *testing.T) {
	s, err := generateTestData(10)
	if err != nil {
		t.Errorf("Error generate TEST data: %v", err)
		return
	}

	err = s.Export(".")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestService_Import(t *testing.T) {
	s1, err := generateTestData(10)
	if err != nil {
		t.Errorf("Error generate TEST data: %v", err)
		return
	}

	err = s1.Export(".")
	if err != nil {
		t.Error(err)
		return
	}

	s2 := &Service{}
	err = s2.Import(".")
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(s1.accounts, s2.accounts) {
		t.Error("s1.accounts and s2.accounts must equals")
	}

	if !reflect.DeepEqual(s1.payments, s2.payments) {
		t.Error("s1.payments and s2.payments must equals")
	}

	if !reflect.DeepEqual(s1.favorites, s2.favorites) {
		t.Error("s1.favorites and s2.favorites must equals")
	}
}

func TestService_HistoryToFiles(t *testing.T) {
	s, err := generateTestData(10)
	if err != nil {
		t.Errorf("Error generate TEST data: %v", err)
		return
	}

	pays, err := s.ExportAccountHistory(1)
	if err != nil {
		t.Error(err)
		return
	}

	err = s.HistoryToFiles(pays, ".", 10)
	if err != nil {
		t.Error(err)
		return
	}

}
