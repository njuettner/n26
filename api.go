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

type N26Transactions []struct {
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

type N26Cards []struct {
	ID                                  string      `json:"id"`
	PublicToken                         interface{} `json:"publicToken"`
	Pan                                 interface{} `json:"pan"`
	MaskedPan                           string      `json:"maskedPan"`
	ExpirationDate                      int64       `json:"expirationDate"`
	CardType                            string      `json:"cardType"`
	Status                              string      `json:"status"`
	CardProduct                         interface{} `json:"cardProduct"`
	CardProductType                     string      `json:"cardProductType"`
	PinDefined                          int64       `json:"pinDefined"`
	CardActivated                       int64       `json:"cardActivated"`
	UsernameOnCard                      string      `json:"usernameOnCard"`
	ExceetExpressCardDelivery           interface{} `json:"exceetExpressCardDelivery"`
	Membership                          interface{} `json:"membership"`
	ExceetActualDeliveryDate            interface{} `json:"exceetActualDeliveryDate"`
	ExceetExpressCardDeliveryEmailSent  interface{} `json:"exceetExpressCardDeliveryEmailSent"`
	ExceetCardStatus                    interface{} `json:"exceetCardStatus"`
	ExceetExpectedDeliveryDate          interface{} `json:"exceetExpectedDeliveryDate"`
	ExceetExpressCardDeliveryTrackingID interface{} `json:"exceetExpressCardDeliveryTrackingId"`
	CardSettingsID                      interface{} `json:"cardSettingsId"`
	MptsCard                            bool        `json:"mptsCard"`
}

type N26AccountStatus struct {
	ID                           string `json:"id"`
	Created                      int64  `json:"created"`
	Updated                      int64  `json:"updated"`
	SingleStepSignup             int64  `json:"singleStepSignup"`
	EmailValidationInitiated     int64  `json:"emailValidationInitiated"`
	EmailValidationCompleted     int64  `json:"emailValidationCompleted"`
	ProductSelectionCompleted    int64  `json:"productSelectionCompleted"`
	PhonePairingInitiated        int64  `json:"phonePairingInitiated"`
	PhonePairingCompleted        int64  `json:"phonePairingCompleted"`
	KycInitiated                 int64  `json:"kycInitiated"`
	KycCompleted                 int64  `json:"kycCompleted"`
	KycPostIdentInitiated        int64  `json:"kycPostIdentInitiated"`
	KycWebIDInitiated            int64  `json:"kycWebIDInitiated"`
	KycWebIDCompleted            int64  `json:"kycWebIDCompleted"`
	CardActivationCompleted      int64  `json:"cardActivationCompleted"`
	PinDefinitionCompleted       int64  `json:"pinDefinitionCompleted"`
	BankAccountCreationInitiated int64  `json:"bankAccountCreationInitiated"`
	BankAccountCreationSucceded  int64  `json:"bankAccountCreationSucceded"`
	CoreDataUpdated              int64  `json:"coreDataUpdated"`
	FirstIncomingTransaction     int64  `json:"firstIncomingTransaction"`
	FlexAccount                  bool   `json:"flexAccount"`
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
	resp, _ := api.call("/api/smrt/contacts", nil)
	err := json.NewDecoder(resp.Body).Decode(contacts)
	if err != nil {
		fmt.Println("Dummy")
	}
	return contacts
}

// Transactions returns the latest transactions from customers bank account
func (api *N26API) Transactions(amount string) *N26Transactions {
	transactions := &N26Transactions{}
	v := &url.Values{}
	v.Add("limit", amount)
	resp, _ := api.call("/api/smrt/transactions", v)
	err := json.NewDecoder(resp.Body).Decode(transactions)
	if err != nil {
		fmt.Println("Dummy")
	}
	return transactions
}

// Balance returns customers current balance
func (api *N26API) Balance() *N26Account {
	account := &N26Account{}
	resp, _ := api.call("/api/accounts", nil)
	err := json.NewDecoder(resp.Body).Decode(account)
	if err != nil {
		fmt.Println("Dummy")
	}
	return account
}

func (api *N26API) AccountLimit() *N26AccountLimit {
	accountLimit := &N26AccountLimit{}
	resp, _ := api.call("/api/settings/account/limits", nil)
	err := json.NewDecoder(resp.Body).Decode(accountLimit)
	if err != nil {
		fmt.Println("Dummy")
	}
	return accountLimit
}

func (api *N26API) AccountInfo() *N26AccountInfo {
	accountInfo := &N26AccountInfo{}
	resp, err := api.call("/api/me", nil)
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
	resp, err := api.call("/api/statements", nil)
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
	resp, err := api.call("/api/statements/"+statementID, nil)
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

func (api *N26API) Stats() {
	v := &url.Values{}
	v.Set("type", "acct")
	v.Add("from", "1451602800")
	v.Add("to", "1451732400")
	v.Add("numSlices", "25")
	resp, err := api.call("/api/accounts/stats/", v)
	if err != nil {
		fmt.Println(err.Error())
	}
	byt, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(byt))
}

func (api *N26API) Cards() *N26Cards {
	cards := &N26Cards{}
	resp, err := api.call("/api/v2/cards", nil)
	if err != nil {
		fmt.Println(err.Error())
	}
	err = json.NewDecoder(resp.Body).Decode(cards)
	if err != nil {
		fmt.Println(err)
	}
	return cards
}

func (api *N26API) Status() *N26AccountStatus {
	accountStatus := &N26AccountStatus{}
	resp, err := api.call("/api/me/statuses", nil)
	if err != nil {
		fmt.Println(err.Error())
	}
	err = json.NewDecoder(resp.Body).Decode(accountStatus)
	if err != nil {
		fmt.Println(err)
	}
	return accountStatus
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

func (api *N26API) call(path string, v *url.Values) (*http.Response, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", N26APIUrl, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", api.token()))
	req.URL.Path = path
	if v != nil {
		req.URL.RawQuery = v.Encode()
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
