package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const N26APIUrl = "https://api.tech26.de"

type N26TokenStruct struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	Expires      int    `json:"expires_in"`
	Scope        string `json:"scope"`
}

type N26Account struct {
	AvailableBalance float64 `json:"availableBalance"`
	UsableBalance    float64 `json:"usableBalance"`
	BankBalance      float64 `json:"bankBalance"`
	Iban             string  `json:"iban"`
	Bic              string  `json:"bic"`
	BankName         string  `json:"bankName"`
	Seized           bool    `json:"seized"`
	ID               string  `json:"id"`
}

type N26Transactions []N26Transaction

type N26Transaction struct {
	ID                 string  `json:"id"`
	UserID             string  `json:"userId"`
	Type               string  `json:"type"`
	Amount             float64 `json:"amount"`
	CurrencyCode       string  `json:"currencyCode"`
	VisibleTS          int64   `json:"visibleTS"`
	Recurring          bool    `json:"recurring"`
	PartnerBic         string  `json:"partnerBic"`
	PartnerName        string  `json:"partnerName"`
	AccountID          string  `json:"accountId"`
	PartnerIban        string  `json:"partnerIban"`
	Category           string  `json:"category"`
	ReferenceText      string  `json:"referenceText"`
	UserCertified      int64   `json:"userCertified"`
	Pending            bool    `json:"pending"`
	TransactionNature  string  `json:"transactionNature"`
	CreatedTS          int64   `json:"createdTS"`
	MandateID          string  `json:"mandateId"`
	CreditorIdentifier string  `json:"creditorIdentifier"`
	CreditorName       string  `json:"creditorName"`
	SmartLinkID        string  `json:"smartLinkId"`
	LinkID             string  `json:"linkId"`
	Confirmed          int64   `json:"confirmed"`
}

type N26Contacts []struct {
	UserID   string `json:"userId"`
	ID       string `json:"id"`
	Name     string `json:"name"`
	Subtitle string `json:"subtitle"`
	Account  struct {
		AccountType string `json:"accountType"`
		Iban        string `json:"iban"`
		Bic         string `json:"bic"`
	} `json:"account"`
}

type N26AccountLimit []struct {
	Limit  string  `json:"limit"`
	Amount float64 `json:"amount"`
}

type N26AccountInfo struct {
	ID                        string `json:"id"`
	Email                     string `json:"email"`
	FirstName                 string `json:"firstName"`
	LastName                  string `json:"lastName"`
	KycFirstName              string `json:"kycFirstName"`
	KycLastName               string `json:"kycLastName"`
	Title                     string `json:"title"`
	Gender                    string `json:"gender"`
	BirthDate                 int64  `json:"birthDate"`
	SignupCompleted           bool   `json:"signupCompleted"`
	Nationality               string `json:"nationality"`
	MobilePhoneNumber         string `json:"mobilePhoneNumber"`
	ShadowUserID              string `json:"shadowUserId"`
	TransferWiseTermsAccepted bool   `json:"transferWiseTermsAccepted"`
}

type N26BankStatements []struct {
	ID        string `json:"id"`
	URL       string `json:"url"`
	VisibleTS int64  `json:"visibleTS"`
	Month     int    `json:"month"`
	Year      int    `json:"year"`
}

// N26Interface includes all possible API Calls
type N26Interface interface {
	Transactions(amount string) *N26Transactions
	Balance() *N26Account
	Contacts() *N26Contacts
	Statements() *N26BankStatements
}

// N26API contains Email and Password to get a Token
type N26API struct {
	Email    string
	Password string
}

// Contacts returns all customer contacts
func (api *N26API) Contacts() *N26Contacts {
	contacts := &N26Contacts{}
	token := api.token()
	client := &http.Client{}
	req, _ := http.NewRequest("GET", N26APIUrl, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.URL.Path = "/api/smrt/contacts"
	resp, _ := client.Do(req)
	err := json.NewDecoder(resp.Body).Decode(contacts)
	if err != nil {
		fmt.Println("Dummy")
	}
	return contacts
}

// Transactions returns the latest transactions from customers bank account
func (api *N26API) Transactions(amount string) *N26Transactions {
	transactions := &N26Transactions{}
	token := api.token()
	client := &http.Client{}
	v := url.Values{}
	v.Add("limit", amount)
	req, _ := http.NewRequest("GET", N26APIUrl, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.URL.Path = "/api/smrt/transactions"
	req.URL.RawQuery = v.Encode()
	resp, _ := client.Do(req)
	err := json.NewDecoder(resp.Body).Decode(transactions)
	if err != nil {
		fmt.Println("Dummy")
	}
	return transactions
}

// Balance returns customers current balance
func (api *N26API) Balance() *N26Account {
	account := &N26Account{}
	token := api.token()
	client := &http.Client{}
	req, _ := http.NewRequest("GET", N26APIUrl, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.URL.Path = "/api/accounts"
	resp, _ := client.Do(req)
	err := json.NewDecoder(resp.Body).Decode(account)
	if err != nil {
		fmt.Println("Dummy")
	}
	return account
}

func (api *N26API) AccountLimit() *N26AccountLimit {
	accountLimit := &N26AccountLimit{}
	token := api.token()
	client := &http.Client{}
	req, _ := http.NewRequest("GET", N26APIUrl, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.URL.Path = "/api/settings/account/limits"
	resp, _ := client.Do(req)
	err := json.NewDecoder(resp.Body).Decode(accountLimit)
	if err != nil {
		fmt.Println("Dummy")
	}
	return accountLimit
}

func (api *N26API) AccountInfo() *N26AccountInfo {
	accountInfo := &N26AccountInfo{}
	token := api.token()
	client := &http.Client{}
	req, _ := http.NewRequest("GET", N26APIUrl, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.URL.Path = "/api/me"
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	err = json.NewDecoder(resp.Body).Decode(accountInfo)
	if err != nil {
		fmt.Println("Dummy")
	}
	return accountInfo
}

func (api *N26API) Statements() *N26BankStatements {
	bankStatements := &N26BankStatements{}
	token := api.token()
	client := &http.Client{}
	req, _ := http.NewRequest("GET", N26APIUrl, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.URL.Path = "/api/statements"
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	err = json.NewDecoder(resp.Body).Decode(bankStatements)
	if err != nil {
		fmt.Println("Dummy")
	}
	return bankStatements
}

func (api *N26API) Statement(statementID string) {
	token := api.token()
	client := &http.Client{}
	req, _ := http.NewRequest("GET", N26APIUrl, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.URL.Path = "/api/statements/" + statementID
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	byt, _ := ioutil.ReadAll(resp.Body)
	ioutil.WriteFile(
		fmt.Sprintf("%s.pdf", statementID),
		byt,
		0750,
	)
}

func (api *N26API) token() string {
	tk := &N26TokenStruct{}
	v := url.Values{}
	v.Set("grant_type", "password")
	v.Add("username", api.Email)
	v.Add("password", api.Password)
	req, err := http.NewRequest("GET", N26APIUrl, nil)
	if err != nil {
		fmt.Println("Dummy")
	}
	req.Header.Add("Authorization", "Basic YW5kcm9pZDpzZWNyZXQ=")
	req.URL.Path = "/oauth/token"
	req.URL.RawQuery = v.Encode()
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	err = json.NewDecoder(resp.Body).Decode(tk)
	if err != nil {
		fmt.Println(err)
	}
	return tk.AccessToken
}
