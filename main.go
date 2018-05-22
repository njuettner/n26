package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/howeyc/gopass"
	"github.com/mitchellh/go-homedir"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v2"
)

var (
	version            = "0.3"
	n26                = kingpin.New("n26", "A command-line to interact with your N26 bank account")
	initialize         = n26.Command("init", "Setting up the configuration to use N26 CLI")
	transactions       = n26.Command("transactions", "N26 latest transactions (Number by Default: 5)")
	transactionsNumber = transactions.Arg("amount", "Number of transactions").Default("5").String()
	balance            = n26.Command("balance", "N26 balance")
	contacts           = n26.Command("contacts", "N26 contacts")
	account            = n26.Command("account", "N26 account")
	statements         = n26.Command("statements", "N26 statements, will be saved as PDF files")
	savings            = n26.Command("savings", "N26 savings and investments")
	statementID        = statements.Arg("statementID", "statement-YEAR-MONTH, e.g. statement-2017-05").String()
	info               = account.Command("info", "N26 account information")
	limit              = account.Command("limit", "N26 account limit")
	stats              = account.Command("stats", "N26 account statistics")
	status             = account.Command("status", "N26 account status")
	cards              = n26.Command("cards", "N26 cards")
	config             = Config()
	table              = tablewriter.NewWriter(os.Stdout)
	configFilePath     = "~/.config/n26.yaml"
)

func main() {

	n26.Version(version).Author("Nick Jüttner")

	switch kingpin.MustParse(n26.Parse(os.Args[1:])) {
	case initialize.FullCommand():
		var email string
		fmt.Print("N26 Email: ")
		fmt.Scanln(&email)
		fmt.Print("N26 Password: ")
		pass, err := gopass.GetPasswdMasked()
		if err != nil {
			renderErrorTable(err)
		}
		cfg := NewConfig(email, string(pass))
		data, err := yaml.Marshal(cfg)
		if err != nil {
			renderErrorTable(err)
		}
		filePath, err := homedir.Expand(configFilePath)
		if err != nil {
			renderErrorTable(err)
		}
		ioutil.WriteFile(filePath, data, 0700)
		if err != nil {
			renderErrorTable(err)
		}
		return

	case transactions.FullCommand():
		transactions, err := config.Transactions(*transactionsNumber)
		if err != nil {
			renderErrorTable(err)
			return
		}
		data := [][]string{}
		for _, transaction := range *transactions {
			amount := strconv.FormatFloat(transaction.Amount, 'f', -1, 64)
			data = append(data,
				[]string{
					transaction.PartnerName,
					fmt.Sprintf("%s %s", amount, transaction.CurrencyCode),
					strings.Replace(transaction.Category, "micro-v2-", "", -1)})
		}
		table.SetHeader([]string{"Partner Name", "Amount", "Category"})
		table.SetBorder(false)
		table.AppendBulk(data)
		table.Render()

	case balance.FullCommand():
		balance, err := config.Balance()
		if err != nil {
			renderErrorTable(err)
			return
		}
		available := strconv.FormatFloat(balance.AvailableBalance, 'f', -1, 64)
		usable := strconv.FormatFloat(balance.UsableBalance, 'f', -1, 64)
		data := [][]string{[]string{available, usable}}
		table.SetHeader([]string{"Available Balance", "Usable Balance"})
		table.SetBorder(false)
		table.AppendBulk(data)
		table.Render()

	case contacts.FullCommand():
		contacts, err := config.Contacts()
		if err != nil {
			renderErrorTable(err)
			return
		}
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
		table.SetBorder(false)
		table.AppendBulk(data)
		table.Render()

	case limit.FullCommand():
		limits, err := config.AccountLimit()
		if err != nil {
			renderErrorTable(err)
			return
		}
		data := [][]string{}
		for _, limit := range *limits {
			amount := strconv.FormatFloat(limit.Amount, 'f', -1, 64)
			data = append(data,
				[]string{
					limit.Limit,
					amount})
		}
		table.SetHeader([]string{"Limit", "Amount"})
		table.SetBorder(false)
		table.AppendBulk(data)
		table.Render()

	case info.FullCommand():
		accountInfo, err := config.AccountInfo()
		if err != nil {
			renderErrorTable(err)
			return
		}
		data := [][]string{[]string{accountInfo.ID,
			accountInfo.FirstName,
			accountInfo.LastName,
			accountInfo.Email,
			accountInfo.MobilePhoneNumber,
			accountInfo.Gender,
			accountInfo.Nationality,
		}}
		table.SetHeader([]string{"ID", "First Name", "Last Name", "Email", "Mobile", "Gender", "Nationality"})
		table.SetBorder(false)
		table.AppendBulk(data)
		table.Render()

	case savings.FullCommand():
		savings, err := config.Savings()
		if err != nil {
			renderErrorTable(err)
			return
		}
		data := [][]string{}
		for _, account := range savings.Accounts {
			data = append(data,
				[]string{account.Name,
					fmt.Sprintf("%.2f", account.Balance),
					fmt.Sprintf("%.2f", account.TotalDeposit),
					fmt.Sprintf("%.2f", account.Performance*100),
					fmt.Sprintf("%.2f", account.Profit),
					fmt.Sprintf("%.2f", account.MonthlyAmount),
					account.OptionID,
					account.Status})
		}
		table.SetHeader([]string{"Account Name", "Balance", "Total Deposit", "Performance (%)", "Profit", "Monthly Amount", "Option", "Status"})
		table.SetBorder(false)
		table.AppendBulk(data)
		table.Render()

	case statements.FullCommand():
		if len(*statementID) == 0 {
			bankStatements, err := config.Statements()
			if err != nil {
				renderErrorTable(err)
				return
			}
			data := [][]string{}
			for _, bankStatement := range *bankStatements {
				data = append(data,
					[]string{
						bankStatement.ID,
					})
			}
			table.SetHeader([]string{"ID"})
			table.SetBorder(false)
			table.AppendBulk(data)
			table.Render()
		} else {
			config.Statement(*statementID)
		}

	case stats.FullCommand():
		config.Stats()

	case status.FullCommand():
		accountStatus, err := config.Status()
		if err != nil {
			renderErrorTable(err)
			return
		}
		fmt.Println(*accountStatus)

	case cards.FullCommand():
		cards, err := config.Cards()
		if err != nil {
			renderErrorTable(err)
			return
		}
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
		table.SetBorder(false)
		table.AppendBulk(data)
		table.Render()
	}
}

func renderErrorTable(err error) {
	errorData := []string{err.Error()}
	table.SetHeader([]string{"Error"})
	table.SetBorder(false)
	table.Append(errorData)
	table.Render()
}
