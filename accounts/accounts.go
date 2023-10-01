package accounts

import (
	"errors"
	"fmt"
	"log"
)

var errorNoMoney = errors.New("인출할 수 없음. 돈 부족")

type account struct {
	owner   string
	balance int
}

func NewAccount(owner string) *account {
	account := account{owner: owner, balance: 0}
	return &account
}

func (theAccount *account) Deposit(amount int) {
	theAccount.balance += amount
}

func (theAccount *account) Withdraw(amount int) error {
	if theAccount.balance < amount {
		return errorNoMoney
	}

	theAccount.balance -= amount
	return nil
}

func (theAccount *account) ChangeOwner(newOwner string) {
	theAccount.owner = newOwner
}

func Example() {
	account := NewAccount("hoon")
	fmt.Println(account)

	account.ChangeOwner("hoon2")
	fmt.Println(account)

	account.Deposit(1000)
	fmt.Println(account)

	error := account.Withdraw(2000)
	if error != nil {
		log.Fatalln(error)
	} else {
		fmt.Println(account)
	}
}
