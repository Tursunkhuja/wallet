package main

import (
	"fmt"

	"github.com/Tursunkhuja/wallet/pkg/types"
	"github.com/Tursunkhuja/wallet/pkg/wallet"
)

func main() {
	/*
		a := 4
		b := 3
		c := a / b
		ostatka := a - c*b
		iterCount := c
		if ostatka > 0 {
			iterCount = c + 1
		}
		println(iterCount)
	*/
	s := wallet.Service{}
	phone := types.Phone("992928303783")
	account, err := s.RegisterAccount(phone)
	if err != nil {
		fmt.Printf("Regect(): can't register account, error =")
	}
	account.Balance += 200
	payment, err := s.Pay(account.ID, 10, "auto")
	if err != nil {
		fmt.Printf("Regect(): can't create payment, error =")
	}

	_, err = s.Repeat("5")
	if err == nil {
		fmt.Printf("Regect(): can't copy payment, error = %v", err)
	}
	_, err = s.Repeat(payment.ID)
	if err != nil {
		fmt.Printf("Regect(): can't copy payment, error = ")
	}

	_, err = s.Repeat(payment.ID)
	if err != nil {
		fmt.Printf("Regect(): can't copy payment, error = ")
	}

	s.SumPayments(1)
	s.SumPayments(2)
}
