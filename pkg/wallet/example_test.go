package wallet

import (
	"reflect"
	"testing"

	"github.com/Tursunkhuja/wallet/pkg/types"
	"github.com/google/uuid"
)

//"github.com/Tursunkhuja/wallet/pkg/wallet"

func Test_FindAccountById_Exist(t *testing.T) {
	accountID := uuid.New().String()
	svc := Service{}
	acc1 := types.Account{ID: accountID, Phone: "992928333783", Balance: 100}
	svc.accounts = append(svc.accounts, &acc1)

	account, error := svc.FindAccountByID(accountID)

	if !reflect.DeepEqual(&acc1, account) {
		t.Errorf("There should be an account with ID = 1")
	}
	if !reflect.DeepEqual(nil, error) {
		t.Errorf("There should not be any error!")
	}
}

func Test_FindAccountById_NotExist(t *testing.T) {

	accountID := uuid.New().String()
	svc := Service{}
	acc2 := types.Account{ID: accountID, Phone: "992928303783", Balance: 200}
	svc.accounts = append(svc.accounts, &acc2)

	account, error := svc.FindAccountByID("1")

	if !reflect.DeepEqual(ErrAccountNotFound("account not found!"), error) {
		t.Errorf("There should not be any error!")
	}
	if reflect.DeepEqual(nil, account) {
		t.Errorf("There should not be an account with ID = 1")
	}
}
