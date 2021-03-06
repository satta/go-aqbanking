package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	. "github.com/umsatz/go-aqbanking"
)

func loadPins(filename string) []Pin {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal("%v", err)
		return nil
	}

	var _pins []pin
	if err = json.NewDecoder(f).Decode(&_pins); err != nil {
		log.Fatal("%v", err)
		return nil
	}

	var pins = make([]Pin, len(_pins))
	for i, pin := range _pins {
		pins[i] = Pin(&pin)
	}

	return pins
}

type pin struct {
	Blz string `json:"blz"`
	UID string `json:"uid"`
	PIN string `json:"pin"`
}

func (p *pin) BankCode() string {
	return p.Blz
}

func (p *pin) UserID() string {
	return p.UID
}

func (p *pin) Pin() string {
	return p.PIN
}

func listAccounts(ab *AQBanking) {
	accountCollection, err := ab.Accounts()
	if err != nil {
		log.Fatal("unable to list accounts: %v", err)
	}

	fmt.Println("%%\nAccounts")
	for _, account := range accountCollection.Accounts {
		fmt.Printf(`
## %v
Owner: %v
Currency: %v
Country: %v
AccountNumber: %v
BankCode: %v
Bank: %v
IBAN: %v
BIC: %v
`,
			account.Name,
			account.Owner,
			account.Currency,
			account.Country,
			account.AccountNumber,
			account.Bank.BankCode,
			account.Bank.Name,
			account.IBAN,
			account.BIC,
		)
	}
}

func listUsers(ab *AQBanking) {
	userCollection, err := ab.Users()
	if err != nil {
		log.Fatal("unable to list users: %v", err)
	}
	defer userCollection.Free()

	fmt.Println("%%\nUsers")
	for _, user := range userCollection.Users {
		fmt.Printf(`
## %v
Name: %v
UserId: %v
CustomerId: %v
`,
			user.ID,
			user.Name,
			user.UserID,
			user.CustomerID,
		)
	}
}

func listTransactionsFor(ab *AQBanking, account *Account) {
	fromDate := time.Date(2014, 05, 14, 0, 0, 0, 0, time.UTC)
	toDate := time.Date(2014, 05, 16, 0, 0, 0, 0, time.UTC)
	transactions, err := ab.Transactions(account, &fromDate, &toDate)
	// or
	// transactions, err := ab.AllTransactions(account)
	if err != nil {
		log.Fatalf("unable to get transactions!: %v", err)
	}

	for _, t := range transactions {
		fmt.Printf(`
## %v
'%v'
Status: %v
MandateReference: %v
CustomerReference: %v
LocalBankCode: %v
LocalAccountNumber: %v
LocalIBAN: %v
LocalBIC: %v
LocalName: %v
RemoteBankCode: %v
RemoteAccountNumber: %v
RemoteIBAN: %v
RemoteBIC: %v
RemoteName: %v
Date: %v
ValutaDate: %v
Total: %2.2f %v
Fee: %2.2f %v
`, t.Purpose,
			t.Text,
			t.Status,
			t.CustomerReference,
			t.LocalBankCode,
			t.LocalAccountNumber,
			t.LocalIBAN,
			t.LocalBIC,
			t.LocalName,
			t.RemoteBankCode,
			t.RemoteAccountNumber,
			t.RemoteIBAN,
			t.RemoteBIC,
			t.RemoteName,
			t.Date,
			t.ValutaDate,
			t.Total, t.TotalCurrency,
			t.Fee, t.FeeCurrency,
		)
	}
}

func listTransactions(ab *AQBanking) {
	accountList, err := ab.Accounts()
	if err != nil {
		log.Fatal("unable to list accounts: %v", err)
	}

	for _, account := range accountList.Accounts {
		listTransactionsFor(ab, &account)
	}
}

func main() {
	aq, err := NewAQBanking("custom", "./tmp")
	if err != nil {
		log.Fatal("unable to init aqbanking: %v", err)
	}
	defer aq.Free()

	fmt.Printf("using aqbanking %d.%d.%d\n",
		aq.Version.Major,
		aq.Version.Minor,
		aq.Version.Patchlevel,
	)

	for _, pin := range loadPins("pins.json") {
		aq.RegisterPin(pin)
	}

	listUsers(aq)
	listAccounts(aq)
	listTransactions(aq)
}
