package wallet

import (
	"reflect"
	"testing"

	"github.com/Tursunkhuja/wallet/pkg/types"
)

type testService struct {
	*Service
}

func newTestService() *testService {
	return &testService{Service: &Service{}}
}

func Test_FindAccountById_Exist(t *testing.T) {
	svc := newTestService()
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
	svc := newTestService()
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

func Test_Repeat(t *testing.T) {
	s := newTestService()
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

	newpayment, err := s.Repeat("5")
	if err == nil {
		t.Errorf("Regect(): can't copy payment, error = %v", err)
		return
	}
	if newpayment != nil {
		t.Errorf("dsfds")
		return
	}
	newpayment, err = s.Repeat(payment.ID)
	if err != nil {
		t.Errorf("Regect(): can't copy payment, error = %v", err)
		return
	}

	if !reflect.DeepEqual(types.Money(100), newpayment.Amount) {
		t.Errorf("New payment amount should be 100, because it is copied payment")
	}
}

func Test_FavoritePayment(t *testing.T) {
	s := newTestService()

	phone := types.Phone("992928303783")
	account, err := s.RegisterAccount(phone)
	if err != nil {
		t.Errorf("Regect(): can't register account, error = %v", err)
		return
	}
	//repeating register with the same phone number
	phone = types.Phone("992928303783")
	_, err2 := s.RegisterAccount(phone)
	if err2 == nil {
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
	//Pay with more than account amount
	payment, err2 := s.Pay(account.ID, 2000, "auto")
	if err2 == nil {
		t.Errorf("Regect(): can't create payment, error = %v", err)
		return
	}
	if payment != nil {
		t.Errorf("Regect(): can't create payment, error = %v", err)
		return
	}
	payment, err = s.Pay(account.ID, 100, "auto")
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

func generateTestData(num int) (*testService, error) {
	s := newTestService()

	acc, err := s.RegisterAccount("+992928333783")
	if err != nil {
		return nil, err
	}

	err = s.Deposit(acc.ID, 1_000_000)
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
