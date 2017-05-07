package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	n26                = kingpin.New("n26", "A command-line to interact with N26")
	transactions       = n26.Command("transactions", "N26 Transactions (Number by Default: 5)")
	transactionsNumber = transactions.Arg("amount", "Number of Transactions").Default("5").String()
	balance            = n26.Command("balance", "N26 Balance")
	contacts           = n26.Command("contacts", "N26 Contacts")
	account            = n26.Command("account", "N26 Account")
	statements         = n26.Command("statements", "N26 Bank Statements")
	statementID        = statements.Arg("statementID", "statement-YEAR-MONTH e.g. statement-2017-05").String()
	info               = account.Command("info", "Info")
	limit              = account.Command("limit", "Limit")
	stats              = account.Command("stats", "Statistics")
	status             = account.Command("status", "Status")
	cards              = n26.Command("cards", "Cards")
	config             = NewConfig()
)

func main() {
	table := tablewriter.NewWriter(os.Stdout)
	switch kingpin.MustParse(n26.Parse(os.Args[1:])) {
	case transactions.FullCommand():
		transactions := config.Transactions(*transactionsNumber)
		data := [][]string{}
		for _, transaction := range *transactions {
			amount := strconv.FormatFloat(transaction.Amount, 'f', -1, 64)
			data = append(data,
				[]string{
					transaction.PartnerName,
					fmt.Sprintf("%s %s", amount, transaction.CurrencyCode),
					transaction.Category})
		}
		table.SetHeader([]string{"Partner Name", "Amount", "Category"})
		table.SetBorder(false) // Set Border to false
		table.AppendBulk(data) // Add Bulk Data
		table.Render()

	case balance.FullCommand():
		balance := config.Balance()
		available := strconv.FormatFloat(balance.AvailableBalance, 'f', -1, 64)
		usable := strconv.FormatFloat(balance.UsableBalance, 'f', -1, 64)
		data := [][]string{[]string{available, usable}}
		table.SetHeader([]string{"Available Balance", "Usable Balance"})
		table.SetBorder(false) // Set Border to false
		table.AppendBulk(data) // Add Bulk Data
		table.Render()

	case contacts.FullCommand():
		contacts := config.Contacts()
		data := [][]string{}
		for _, contact := range *contacts {
			data = append(data,
				[]string{
					contact.Name,
					contact.Account.Iban,
					contact.Account.Bic,
					contact.Account.AccountType})
		}
		table.SetHeader([]string{"Contact Name", "IBAN", "BIC", "Account Type"})
		table.SetBorder(false) // Set Border to false
		table.AppendBulk(data) // Add Bulk Data
		table.Render()
	case limit.FullCommand():
		limits := config.AccountLimit()
		data := [][]string{}
		for _, limit := range *limits {
			amount := strconv.FormatFloat(limit.Amount, 'f', -1, 64)
			data = append(data,
				[]string{
					limit.Limit,
					amount})
		}
		table.SetHeader([]string{"Limit", "Amount"})
		table.SetBorder(false) // Set Border to false
		table.AppendBulk(data) // Add Bulk Data
		table.Render()
	case info.FullCommand():
		accountInfo := config.AccountInfo()
		data := [][]string{[]string{accountInfo.ID,
			accountInfo.FirstName,
			accountInfo.LastName,
			accountInfo.Email,
			accountInfo.MobilePhoneNumber,
			accountInfo.Gender,
			accountInfo.Nationality,
		}}
		table.SetHeader([]string{"ID", "First Name", "Last Name", "Email", "Mobile", "Gender", "Nationality"})
		table.SetBorder(false) // Set Border to false
		table.AppendBulk(data) // Add Bulk Data
		table.Render()
	case statements.FullCommand():
		if len(*statementID) == 0 {
			bankStatements := config.Statements()
			data := [][]string{}
			for _, bankStatement := range *bankStatements {
				data = append(data,
					[]string{
						bankStatement.ID,
					})
			}
			table.SetHeader([]string{"ID"})
			table.SetBorder(false) // Set Border to false
			table.AppendBulk(data) // Add Bulk Data
			table.Render()
		} else {
			config.Statement(*statementID)
		}
	case stats.FullCommand():
		config.Stats()
	case status.FullCommand():
		accountStatus := config.Status()
		fmt.Println(*accountStatus)
	case cards.FullCommand():
		cards := config.Cards()
		data := [][]string{}
		for _, card := range *cards {
			data = append(data,
				[]string{
					card.ID,
					card.CardType,
					card.CardProductType,
					card.Status,
					card.UsernameOnCard,
				})
		}
		table.SetHeader([]string{"ID", "Card Type", "Card Product Type", "Status", "Username on card"})
		table.SetBorder(false) // Set Border to false
		table.AppendBulk(data) // Add Bulk Data
		table.Render()
	}
}
