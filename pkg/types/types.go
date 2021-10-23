package types

type Money int64

type PaymentCategory string

type PaymentStatus string

const (
	PaymentStatusOK         PaymentStatus = "OK"
	PaymentStatusFail       PaymentStatus = "FAIL"
	PaymentStatusInProgress PaymentStatus = "INPROGRESS"
)

type Payment struct {
	ID        string
	AccountID int64
	Amount    Money
	Status    PaymentStatus
	Category  PaymentCategory
}

type Phone string

type Account struct {
	ID      string
	Phone   string
	Balance Money
}
