package wallet

import (
	"errors"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
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
	nextAccountID int64
	accounts      []*types.Account
	payments      []*types.Payment
	favorites     []*types.Favorite
}

func (s *Service) RegisterAccount(phone types.Phone) (account *types.Account, err error) {
	for _, account := range s.accounts {
		if account.Phone == phone {
			return nil, ErrPhoneRegistered // if there is such a phone, then just leave
		}
	}
	s.nextAccountID++

	account = &types.Account{
		ID:      s.nextAccountID,
		Phone:   phone,
		Balance: 0,
	}

	s.accounts = append(s.accounts, account)
	return account, nil
}

func (s *Service) FindAccountByID(accountID int64) (account *types.Account, err error) {
	for _, account := range s.accounts {
		if account.ID == accountID {
			return account, nil
		}
	}
	return nil, ErrAccountNotFound
}

func (s *Service) Deposit(accountID int64, amount types.Money) error {
	if amount <= 0 {
		return ErrAmountMustBePositive
	}

	account, err := s.FindAccountByID(accountID)
	if err != nil {
		return err
	}

	account.Balance += amount
	return nil
}

func (s *Service) Pay(accountID int64, amount types.Money, category types.PaymentCategory) (payment *types.Payment, err error) {
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
	payment = &types.Payment{
		ID:        paymentID,
		AccountID: accountID,
		Amount:    amount,
		Category:  category,
		Status:    types.PaymentStatusInProgress,
	}

	s.payments = append(s.payments, payment)

	return payment, nil
}

func (s *Service) FindPaymentByID(paymentID string) (payment *types.Payment, err error) {
	for _, payment := range s.payments {
		if payment.ID == paymentID {
			return payment, nil
		}
	}

	return nil, ErrPaymentNotFound
}

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

	return nil
}

func (s *Service) Repeat(paymentID string) (payment *types.Payment, err error) {
	payment, err = s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, err
	}
	return s.Pay(payment.AccountID, payment.Amount, payment.Category)
}

func (s *Service) FavoritePayment(paymentID string, name string) (favorite *types.Favorite, err error) {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, err
	}

	favorite = &types.Favorite{
		ID:        uuid.New().String(),
		AccountID: payment.AccountID,
		Amount:    payment.Amount,
		Name:      name,
		Category:  payment.Category,
	}

	s.favorites = append(s.favorites, favorite)
	return favorite, nil
}

func (s *Service) FindFavoriteByID(favoriteID string) (favorite *types.Favorite, err error) {
	for _, favorite := range s.favorites {
		if favorite.ID == favoriteID {
			return favorite, nil
		}
	}

	return nil, ErrFavoriteNotFound
}

func (s *Service) PayFromFavorite(favoriteID string) (payment *types.Payment, err error) {
	favorite, err := s.FindFavoriteByID(favoriteID)
	if err != nil {
		return nil, err
	}
	return s.Pay(favorite.AccountID, favorite.Amount, favorite.Category)
}

func (s *Service) ExportToFile(path string) error {

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	defer func() {
		err := file.Close()
		if err != nil {
			log.Print(err)
			return
		}
	}()

	str := ""
	for _, account := range s.accounts {
		str = str + strconv.FormatInt(account.ID, 10) + ";" + string(account.Phone) + ";" + strconv.FormatInt(int64(account.Balance), 10) + "|"
	}

	str = strings.TrimRight(str, "|")
	_, err = file.WriteString(str)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) ImportFromFile(path string) error {

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Print(err)
			return
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
			log.Print(err)
			return err
		}

		content = append(content, buf[:read]...)
	}

	rows := strings.Split(string(content), "|")

	for _, row := range rows {
		accountString := strings.Split(row, ";")
		if accountString != nil {
			id, _ := strconv.Atoi(accountString[0])
			balance, _ := strconv.Atoi(accountString[2])

			account := &types.Account{
				ID:      int64(id),
				Phone:   types.Phone(accountString[1]),
				Balance: types.Money(balance),
			}
			s.accounts = append(s.accounts, account)
		}
	}

	return nil
}

func (s *Service) Export(dir string) error {
	filepath1 := filepath.Join(dir, "accounts.dump")
	err := s.ExportToFileByServiceType(filepath1, "Account")
	if err != nil {
		return err
	}
	filepath2 := filepath.Join(dir, "payments.dump")
	err = s.ExportToFileByServiceType(filepath2, "Payment")
	if err != nil {
		return err
	}

	filepath3 := filepath.Join(dir, "favorites.dump")
	err = s.ExportToFileByServiceType(filepath3, "Favorite")
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) ExportAccountHistory(accountID int64) ([]types.Payment, error) {
	accFound := false
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			accFound = true
		}
	}
	if !accFound {
		return nil, ErrAccountNotFound
	}

	payments := []types.Payment{}
	paymFound := false
	for _, payment := range s.payments {
		if payment.AccountID == accountID {
			payments = append(payments, *payment)
			paymFound = true
		}
	}
	if paymFound {
		return payments, nil
	}
	return nil, nil
}

func (s *Service) HistoryToFiles(payments []types.Payment, dir string, records int) error {

	paymentCount := len(payments)
	c := paymentCount / records
	ostatka := paymentCount - c*records
	iterCount := c
	if ostatka > 0 {
		iterCount = c + 1
	}

	str := ""
	if iterCount == 0 {
		return nil
	}
	fileName := ""
	if iterCount == 1 {

		fileName = "payments.dump"
		for _, payment := range payments {
			str = str +
				string(payment.ID) + ";" +
				strconv.FormatInt(payment.AccountID, 10) + ";" +
				strconv.FormatInt(int64(payment.Amount), 10) + ";" +
				string(payment.Category) + ";" +
				string(payment.Status) + "\n"
		}
		filepathPayment := filepath.Join(dir, fileName)
		file, err := os.Create(filepathPayment)
		if err != nil {
			return err
		}

		str = strings.TrimRight(str, "\n")
		_, err = file.WriteString(str)
		if err != nil {
			return err
		}
		err = file.Close()
		if err != nil {
			return err
		}
	}

	if iterCount > 1 {
		currentPaymentCount := 0
		currentFileIndex := 0
		str = ""
		for _, payment := range payments {
			currentPaymentCount++
			str = str +
				string(payment.ID) + ";" +
				strconv.FormatInt(payment.AccountID, 10) + ";" +
				strconv.FormatInt(int64(payment.Amount), 10) + ";" +
				string(payment.Category) + ";" +
				string(payment.Status) + "\n"

			if currentPaymentCount%records == 0 || paymentCount == currentPaymentCount {
				currentFileIndex++
				fileName = string("payments" + strconv.FormatInt(int64(currentFileIndex), 10) + ".dump")
				filepathPayment := filepath.Join(dir, fileName)
				file, err := os.Create(filepathPayment)
				if err != nil {
					return err
				}

				str = strings.TrimRight(str, "\n")
				_, err = file.WriteString(str)
				if err != nil {
					return err
				}
				err = file.Close()
				if err != nil {
					return err
				}
				str = ""

			}
		}
	}
	return nil
}

func (s *Service) Import(dir string) error {
	filepath1 := filepath.Join(dir, "accounts.dump")
	if _, errOfExistingFile := os.Stat(filepath1); errOfExistingFile == nil {
		err := s.ImportFromFileByServiceType(filepath1, "Account")
		if err != nil {
			return err
		}
	}

	filepath2 := filepath.Join(dir, "payments.dump")
	if _, errOfExistingFile := os.Stat(filepath2); errOfExistingFile == nil {
		err := s.ImportFromFileByServiceType(filepath2, "Payment")
		if err != nil {
			return err
		}
	}

	filepath3 := filepath.Join(dir, "favorites.dump")
	if _, errOfExistingFile := os.Stat(filepath3); errOfExistingFile == nil {
		err := s.ImportFromFileByServiceType(filepath3, "Favorite")
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) ExportToFileByServiceType(path string, serviceType string) error {
	str := ""
	//Define here service type like accounts, payments or favorites
	hasData := false
	if serviceType == "Account" {
		for _, account := range s.accounts {
			hasData = true
			str = str +
				strconv.FormatInt(account.ID, 10) + ";" +
				string(account.Phone) + ";" +
				strconv.FormatInt(int64(account.Balance), 10) + "\n"
		}
	}
	if serviceType == "Payment" {
		for _, payment := range s.payments {
			hasData = true
			str = str +
				string(payment.ID) + ";" +
				strconv.FormatInt(payment.AccountID, 10) + ";" +
				strconv.FormatInt(int64(payment.Amount), 10) + ";" +
				string(payment.Category) + ";" +
				string(payment.Status) + "\n"
		}
	}
	if serviceType == "Favorite" {
		for _, favorite := range s.favorites {
			hasData = true
			str = str +
				string(favorite.ID) + ";" +
				strconv.FormatInt(favorite.AccountID, 10) + ";" +
				strconv.FormatInt(int64(favorite.Amount), 10) + ";" +
				string(favorite.Name) + ";" +
				string(favorite.Category) + "\n"
		}
	}

	if hasData {
		file, err := os.Create(path)
		if err != nil {
			return err
		}

		defer func() {
			err := file.Close()
			if err != nil {
				log.Print(err)
				return
			}
		}()

		str = strings.TrimRight(str, "\n")
		_, err = file.WriteString(str)
		if err != nil {
			return err
		}
	}

	return nil
}
func (s *Service) ImportFromFileByServiceType(path string, serviceType string) error {

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Print(err)
			return
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
			log.Print(err)
			return err
		}

		content = append(content, buf[:read]...)
	}

	rows := strings.Split(string(content), "\n")
	for _, row := range rows {
		stringDatas := strings.Split(row, ";")
		if stringDatas != nil {
			existID := false
			if serviceType == "Account" {
				id, _ := strconv.Atoi(stringDatas[0])
				phone := types.Phone(stringDatas[1])
				balance, _ := strconv.Atoi(stringDatas[2])

				for _, a := range s.accounts {
					if a.ID == int64(id) {
						a.Phone = phone
						a.Balance = types.Money(balance)
						s.nextAccountID = a.ID
						existID = true
					}
				}
				if !existID {
					account := &types.Account{
						ID:      int64(id),
						Phone:   phone,
						Balance: types.Money(balance),
					}
					s.accounts = append(s.accounts, account)
				}
			}
			if serviceType == "Payment" {
				accountID, _ := strconv.Atoi(stringDatas[1])
				amount, _ := strconv.Atoi(stringDatas[2])
				category := types.PaymentCategory(stringDatas[3])
				status := types.PaymentStatus(stringDatas[4])
				for _, p := range s.payments {
					if p.ID == stringDatas[0] {
						p.AccountID = int64(accountID)
						p.Amount = types.Money(amount)
						p.Category = category
						p.Status = status
						existID = true
					}
				}
				if !existID {
					payment := &types.Payment{
						ID:        stringDatas[0],
						AccountID: int64(accountID),
						Amount:    types.Money(amount),
						Category:  category,
						Status:    status,
					}
					s.payments = append(s.payments, payment)
				}
			}
			if serviceType == "Favorite" {
				accountID, _ := strconv.Atoi(stringDatas[1])
				amount, _ := strconv.Atoi(stringDatas[2])
				category := types.PaymentCategory(stringDatas[4])
				for _, f := range s.favorites {
					if f.ID == stringDatas[0] {
						f.AccountID = int64(accountID)
						f.Amount = types.Money(amount)
						f.Name = stringDatas[3]
						f.Category = category
						existID = true
					}
				}
				if !existID {
					favorite := &types.Favorite{
						ID:        stringDatas[0],
						AccountID: int64(accountID),
						Amount:    types.Money(amount),
						Name:      stringDatas[3],
						Category:  category,
					}
					s.favorites = append(s.favorites, favorite)
				}
			}
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
