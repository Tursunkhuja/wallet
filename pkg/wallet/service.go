package wallet

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/Tursunkhuja/wallet/pkg/types"
	"github.com/google/uuid"
)

var ErrPhoneRegistered = errors.New("phone already registered")
var ErrAmountMustBePositive = errors.New("amount must be greater than zero")
var ErrAccountNotFound = errors.New("account not found")
var ErrNotEnoughBalance = errors.New("not enough balance")
var ErrPaymentNotFound = errors.New("payment not found")
var ErrFavoriteNotFound = errors.New("favorite not found")

type Service struct {
	nextAccountID int64 // to generate a unique account number
	accounts      []*types.Account
	payments      []*types.Payment
	favorites     []*types.Favorite
}

type Error string

func (e Error) Error() string {
	return string(e)
}

func (s *Service) RegisterAccount(phone types.Phone) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.Phone == phone {
			return nil, ErrPhoneRegistered // if there is such a phone, then just leave
		}
	}
	s.nextAccountID++

	account := &types.Account{
		ID:      s.nextAccountID,
		Phone:   phone,
		Balance: 0,
	}

	s.accounts = append(s.accounts, account)
	return account, nil
}

func (s *Service) Deposit(accontID int64, amount types.Money) error {
	if amount <= 0 {
		return ErrAmountMustBePositive
	}

	var account *types.Account

	for _, acc := range s.accounts {
		if acc.ID == accontID {
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

func (s *Service) Pay(accontID int64, amount types.Money, category types.PaymentCategory) (*types.Payment, error) {
	if amount <= 0 {
		return nil, ErrAmountMustBePositive
	}

	var account *types.Account

	for _, acc := range s.accounts {
		if acc.ID == accontID {
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
		AccountID: accontID,
		Amount:    amount,
		Category:  category,
		Status:    types.PaymentStatusInProgress,
	}

	s.payments = append(s.payments, payment)

	return payment, nil
}

func (s *Service) FindAccountByID(accountID int64) (*types.Account, error) {
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

	return account, nil
}

func (s *Service) FindPaymentByID(paymentID string) (*types.Payment, error) {

	for _, payment := range s.payments {
		if paymentID == payment.ID {
			return payment, nil
		}
	}
	return nil, ErrPaymentNotFound
}

func (s *Service) Reject(paymentID string) error {
	targetPayment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return err
	}

	targetAccount, err := s.FindAccountByID(targetPayment.AccountID)
	if err != nil {
		return err
	}
	targetPayment.Status = types.PaymentStatusFail
	targetAccount.Balance += targetPayment.Amount

	return nil
}

func (s *Service) Repeat(paymentID string) (*types.Payment, error) {

	existingPayment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, err
	}

	// account, err := s.FindAccountByID(existingPayment.AccountID)
	// if err != nil{
	// 	return nil, err
	// }

	repeatedPayment, err := s.Pay(existingPayment.AccountID, existingPayment.Amount, existingPayment.Category)
	if err != nil {
		return nil, err
	}

	return repeatedPayment, nil
}

// creates favorites from a specific payment
func (s *Service) FavoritePayment(paymentID string, name string) (*types.Favorite, error) {
	//var favorite *types.Favorite

	payment, err := s.FindPaymentByID(paymentID)

	if err != nil {
		return nil, err
	}

	id := uuid.New().String()
	favorite := &types.Favorite{
		ID:        id,
		AccountID: payment.AccountID,
		Name:      name,
		Amount:    payment.Amount,
		Category:  payment.Category,
	}
	s.favorites = append(s.favorites, favorite)

	return favorite, nil
}

func (s *Service) FindFavoriteByID(favoriteID string) (*types.Favorite, error) {
	for _, favorite := range s.favorites {
		if favoriteID == favorite.ID {
			return favorite, nil
		}
	}
	return nil, ErrFavoriteNotFound
}

// makes a payment from a specific favorite
func (s *Service) PayFromFavorite(favoriteID string) (*types.Payment, error) {

	favorite, err := s.FindFavoriteByID(favoriteID)

	if err != nil {
		return nil, err
	}

	payment, err := s.Pay(favorite.AccountID, favorite.Amount, favorite.Category)

	if err != nil {
		return nil, err
	}

	return payment, nil
}

func (s *Service) ExportToFile(path string) error {

	file, err := os.Create(path)
	if err != nil {
		log.Println(err)
		return err
	}
	defer closeFile(file)

	content := make([]byte, 0)

	for _, accInfo := range s.accounts {
		accString := fmt.Sprintf("%v;%v;%v", accInfo.ID, accInfo.Phone, accInfo.Balance)
		if len(content) > 0 {
			content = append(content, []byte("|")...)
		}
		content = append(content, []byte(accString)...)
	}

	_, err = file.Write(content)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func closeFile(file *os.File) {
	err := file.Close()

	if err != nil {
		log.Print(err)
	}

}

func (s *Service) ImportFromFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		log.Println(err)
		return err
	}

	defer func() {
		if err = file.Close(); err != nil {
			log.Println(err)
		}
	}()

	content := make([]byte, 0)
	buf := make([]byte, 4)
	for {
		read, err := file.Read(buf)
		if err == io.EOF {
			content = append(content, buf[:read]...)
			break
		}

		if err != nil {
			log.Println(err)
			return err
		}
		content = append(content, buf[:read]...)
	}
	data := string(content)

	for _, v := range strings.Split(data, "|") {
		accS := strings.Split(v, ";")
		id, err := strconv.ParseInt(accS[0], 10, 64)
		if err != nil {
			log.Println(err)
			return err
		}
		phone := accS[1]
		balance, err := strconv.ParseInt(accS[2], 10, 64)
		if err != nil {
			log.Println(err)
			return err
		}
		account := types.Account{
			ID:      id,
			Phone:   types.Phone(phone),
			Balance: types.Money(balance),
		}
		s.accounts = append(s.accounts, &account)
	}
	return nil
}

func (s *Service) Export(dir string) error {

	err := s.ExportAccounts(dir)
	if err != nil {
		return err
	}

	err = s.ExportPayments(dir)
	if err != nil {
		return err
	}

	err = s.ExportFavorites(dir)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) ExportAccounts(dir string) error {

	if len(s.accounts) == 0 {
		return nil
	}

	content := make([]byte, 0)
	for _, v := range s.accounts {
		accString := fmt.Sprintf("%v;%v;%v", v.ID, v.Phone, v.Balance)
		if len(content) > 0 {
			content = append(content, []byte("\n")...)
		}
		content = append(content, []byte(accString)...)
	}
	err := os.WriteFile(dir+"/accounts.dump", content, 0666)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) ExportPayments(dir string) error {

	if len(s.payments) == 0 {
		return nil
	}

	content := make([]byte, 0)
	for _, v := range s.payments {
		payString := fmt.Sprintf("%v;%v;%v;%v;%v", v.ID, v.AccountID, v.Amount, v.Category, v.Status)
		if len(content) > 0 {
			content = append(content, []byte("\n")...)
		}
		content = append(content, []byte(payString)...)
	}
	err := os.WriteFile(dir+"/payments.dump", content, 0666)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) ExportFavorites(dir string) error {

	if len(s.favorites) == 0 {
		return nil
	}

	content := make([]byte, 0)
	for _, v := range s.favorites {
		favString := fmt.Sprintf("%v;%v;%v;%v;%v", v.ID, v.AccountID, v.Name, v.Amount, v.Category)
		if len(content) > 0 {
			content = append(content, []byte("\n")...)
		}
		content = append(content, []byte(favString)...)
	}
	err := os.WriteFile(dir+"/favorites.dump", content, 0666)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) Import(dir string) error {
	err := s.ImportAccounts(dir)
	if err != nil {
		return err
	}

	err = s.ImportPayments(dir)
	if err != nil {
		return err
	}

	err = s.ImportFavorites(dir)
	if err != nil {
		return err
	}
	return nil

}

func (s *Service) ImportAccounts(dir string) error {

	content, err := os.ReadFile(dir + "/accounts.dump")
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}

	for _, v := range strings.Split(string(content), "\n") {
		rec := strings.Split(v, ";")
		id, err := strconv.ParseInt(rec[0], 10, 64)
		if err != nil {
			return err
		}
		phone := rec[1]
		balance, err := strconv.ParseInt(rec[2], 10, 34)
		if err != nil {
			return err
		}
		acc, err := s.FindAccountByID(id)
		if err != nil {
			account := types.Account{
				ID:      id,
				Phone:   types.Phone(phone),
				Balance: types.Money(balance),
			}
			s.accounts = append(s.accounts, &account)
			s.nextAccountID++
			continue
		}
		acc.Phone = types.Phone(phone)
		acc.Balance = types.Money(balance)

	}

	return nil
}

func (s *Service) ImportPayments(dir string) error {

	content, err := os.ReadFile(dir + "/payments.dump")
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}

	for _, v := range strings.Split(string(content), "\n") {
		rec := strings.Split(v, ";")
		id := rec[0]

		accid, err := strconv.ParseInt(rec[1], 10, 64)
		if err != nil {
			return err
		}

		amount, err := strconv.ParseInt(rec[2], 10, 64)
		if err != nil {
			return err
		}

		category := rec[3]
		status := rec[4]

		pay, err := s.FindPaymentByID(id)
		if err != nil {
			payment := types.Payment{
				ID:        id,
				AccountID: accid,
				Amount:    types.Money(amount),
				Category:  types.PaymentCategory(category),
				Status:    types.PaymentStatus(status),
			}
			s.payments = append(s.payments, &payment)
			continue

		}
		pay.AccountID = accid
		pay.Amount = types.Money(amount)
		pay.Category = types.PaymentCategory(category)
		pay.Status = types.PaymentStatus(status)

	}

	return nil
}

func (s *Service) ImportFavorites(dir string) error {

	content, err := os.ReadFile(dir + "/favorites.dump")
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}

	for _, v := range strings.Split(string(content), "\n") {
		rec := strings.Split(v, ";")
		id := rec[0]
		accid, err := strconv.ParseInt(rec[1], 10, 64)
		if err != nil {
			return err
		}

		name := rec[2]

		amount, err := strconv.ParseInt(rec[3], 10, 64)
		if err != nil {
			return err
		}

		category := rec[4]

		fav, err := s.FindFavoriteByID(id)
		if err != nil {
			favorite := types.Favorite{
				ID:        id,
				AccountID: accid,
				Name:      name,
				Amount:    types.Money(amount),
				Category:  types.PaymentCategory(category),
			}
			s.favorites = append(s.favorites, &favorite)
			continue
		}
		fav.AccountID = accid
		fav.Amount = types.Money(amount)
		fav.Name = name
		fav.Category = types.PaymentCategory(category)

	}

	return nil
}

func (s *Service) ExportAccountHistory(accountID int64) ([]types.Payment, error) {

	_, err := s.FindAccountByID(accountID)
	if err != nil {
		return nil, err
	}

	payments := []types.Payment{}

	for _, pay := range s.payments {
		if accountID == pay.AccountID {

			payments = append(payments, *pay)
		}
	}

	return payments, nil
}

func (s *Service) HistoryToFiles(payments []types.Payment, dir string, records int) error {

	if len(payments) == 0 {
		return nil
	}
	count := 0
	for i := 0; i < len(payments); {

		content := make([]byte, 0)
		for j := 0; j < records; j, i = j+1, i+1 {
			if i == len(payments) {
				break
			}
			pay := payments[i]
			str := fmt.Sprintf("%v;%v;%v;%v;%v", pay.ID, pay.AccountID, pay.Amount, pay.Category, pay.Status)
			if len(content) > 0 {
				content = append(content, []byte("\n")...)
			}
			content = append(content, []byte(str)...)

		}
		fileName := "payments.dump"
		if len(payments) > records {
			count++
			fileName = fmt.Sprintf("payments%v.dump", count)
		}

		err := os.WriteFile(dir+"/"+fileName, content, 0666)
		if err != nil {
			log.Println(err)
			return err
		}

	}

	return nil
}

func (s *Service) SumPayments(goroutines int) types.Money {
	if goroutines <= 1 {
		return s.SumPaymentsRegular()
	}

	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	sum := types.Money(0)

	numElem := int(
		math.Ceil(
			float64(len(s.payments)) / float64(goroutines)))

	indexStart := 0
	for i := 0; i < goroutines; i++ {
		wg.Add(1)

		go func(index, num int) {
			defer wg.Done()
			tmpSum := types.Money(0)
			for _, v := range s.payments[index:] {
				if num == 0 {
					break
				}
				num--
				tmpSum += v.Amount
			}
			mu.Lock()
			sum += tmpSum
			mu.Unlock()

		}(indexStart, numElem)

		indexStart += numElem
	}

	wg.Wait()
	return sum
}

func (s *Service) SumPaymentsRegular() types.Money {
	sum := types.Money(0)

	for _, v := range s.payments {
		sum += v.Amount
	}

	return sum
}

func (s *Service) FilterPayments(accountID int64, goroutines int) ([]types.Payment, error) {
	if goroutines <= 1 {
		return s.FilterPaymentsRegular(accountID)
	}

	acc, err := s.FindAccountByID(accountID)
	if err != nil {
		return nil, err
	}

	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	payments := []types.Payment{}
	numElem := int(math.Ceil(float64(len(s.payments)) / float64(goroutines)))

	indexStart := 0
	for i := 0; i < goroutines; i++ {
		wg.Add(1)

		go func(index, num int) {
			defer wg.Done()
			tmpPays := []types.Payment{}
			for _, v := range s.payments[index:] {
				if num == 0 {
					break
				}
				num--

				if acc.ID == v.AccountID {
					tmpPays = append(tmpPays, *v)
				}
			}

			mu.Lock()
			payments = append(payments, tmpPays...)
			mu.Unlock()
		}(indexStart, numElem)

		indexStart += numElem
	}

	wg.Wait()

	return payments, nil
}

func (s *Service) FilterPaymentsRegular(accountID int64) ([]types.Payment, error) {
	acc, err := s.FindAccountByID(accountID)
	if err != nil {
		return nil, err
	}

	payments := []types.Payment{}

	for _, v := range s.payments {
		if acc.ID == v.AccountID {
			payment := types.Payment{
				ID:        v.ID,
				AccountID: v.AccountID,
				Category:  v.Category,
				Amount:    v.Amount,
				Status:    v.Status,
			}
			payments = append(payments, payment)
		}
	}
	return payments, nil
}

func (s *Service) FilterPaymentsByFn(
	filter func(payment types.Payment) bool,
	goroutines int,
) ([]types.Payment, error) {
	if goroutines <= 1 {
		return s.FilterPaymentsByFnRegular(filter)
	}

	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	payments := []types.Payment{}
	numElem := int(math.Ceil(float64(len(s.payments)) / float64(goroutines)))

	indexStart := 0
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(index, num int) {
			defer wg.Done()
			tmpPays := []types.Payment{}
			for _, v := range s.payments[index:] {
				if num == 0 {
					break
				}
				num--

				if filter(*v) {
					tmpPays = append(tmpPays, *v)
				}
			}
			mu.Lock()
			payments = append(payments, tmpPays...)
			mu.Unlock()
		}(indexStart, numElem)

		indexStart += numElem
	}
	wg.Wait()

	return payments, nil
}

func (s *Service) FilterPaymentsByFnRegular(
	filter func(payment types.Payment) bool) ([]types.Payment, error) {

	payments := []types.Payment{}
	for _, v := range s.payments {
		if filter(*v) {
			payments = append(payments, *v)
		}
	}
	return payments, nil
}

type Progress struct {
	Part   int
	Result types.Money
}

func (s *Service) SumPaymentsWithProgress() <-chan Progress {
	ch := make(chan Progress)
	size := 100_000
	parts := int(
		math.Ceil(
			float64(len(s.payments)) / float64(size)))

	wg := sync.WaitGroup{}
	wg.Add(parts)
	mu := sync.Mutex{}
	go func(ch chan Progress, wg *sync.WaitGroup) {
		wg.Wait()
		close(ch)
	}(ch, &wg)

	data := s.payments
	for i := 0; i < parts; i++ {
		go func(data []*types.Payment, part, size int) {
			defer wg.Done()
			tmpSum := types.Money(0)
			for _, v := range data {
				if size == 0 {
					break
				}
				size--
				tmpSum += v.Amount
			}
			mu.Lock()
			ch <- Progress{part, tmpSum}
			mu.Unlock()
		}(data[i*size:], parts, size)

	}
	return ch
}
