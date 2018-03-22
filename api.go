package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const N26APIUrl = "https://api.tech26.de"

type N26Error struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type N26Token struct {
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

type N26Savings struct {
	TotalBalance float64 `json:"totalBalance"`
	CanOpenMore  bool    `json:"canOpenMore"`
	Accounts     []struct {
		ID            string  `json:"id"`
		Name          string  `json:"name"`
		MonthlyAmount float64 `json:"monthlyAmount"`
		NextDate      string  `json:"nextDate"`
		History       []struct {
			Name             string  `json:"name"`
			Date             string  `json:"date"`
			Value            float64 `json:"value"`
			Profit           float64 `json:"profit"`
			ProfitPercentage float64 `json:"profitPercentage"`
		} `json:"history"`
		Forecasts []struct {
			Name             string  `json:"name"`
			Date             string  `json:"date"`
			Value            float64 `json:"value"`
			PessimisticValue float64 `json:"pessimisticValue"`
			OptimisticValue  float64 `json:"optimisticValue"`
			Profit           float64 `json:"profit"`
			ProfitPercentage float64 `json:"profitPercentage"`
		} `json:"forecasts"`
		RiskDisclaimerURL     string  `json:"riskDisclaimerUrl"`
		ForecastDisclaimerURL string  `json:"forecastDisclaimerUrl"`
		OptionID              string  `json:"optionId"`
		StartingDate          string  `json:"startingDate"`
		Balance               float64 `json:"balance"`
		TotalDeposit          float64 `json:"totalDeposit"`
		Performance           float64 `json:"performance"`
		Profit                float64 `json:"profit"`
		Status                string  `json:"status"`
	} `json:"accounts"`
	PendingAccounts []interface{} `json:"pendingAccounts"`
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

// N26Credentials contains Email and Password to get a Token
type N26Credentials struct {
	Email    string
	Password string
}

// Contacts returns all customer contacts
func (n26 *N26Credentials) Contacts() (*N26Contacts, error) {
	contacts := &N26Contacts{}
	resp, err := n26.callAPI("/api/smrt/contacts", nil)
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(resp.Body).Decode(contacts)
	if err != nil {
		return nil, err
	}
	return contacts, nil
}

// Transactions returns the latest transactions from customers bank account
func (n26 *N26Credentials) Transactions(amount string) (*N26Transactions, error) {
	transactions := &N26Transactions{}
	v := &url.Values{}
	v.Add("limit", amount)
	resp, err := n26.callAPI("/api/smrt/transactions", v)
	if err != nil {
		return nil, err
	}
	err = checkHttpStatus(resp)
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(resp.Body).Decode(transactions)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

// Balance returns customers current balance
func (n26 *N26Credentials) Balance() (*N26Account, error) {
	account := &N26Account{}
	resp, err := n26.callAPI("/api/accounts", nil)
	if err != nil {
		return nil, err
	}
	err = checkHttpStatus(resp)
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(resp.Body).Decode(account)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (n26 *N26Credentials) AccountLimit() (*N26AccountLimit, error) {
	accountLimit := &N26AccountLimit{}
	resp, err := n26.callAPI("/api/settings/account/limits", nil)
	if err != nil {
		return nil, err
	}
	err = checkHttpStatus(resp)
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(resp.Body).Decode(accountLimit)
	if err != nil {
		return nil, err
	}
	return accountLimit, nil
}

func (n26 *N26Credentials) AccountInfo() (*N26AccountInfo, error) {
	accountInfo := &N26AccountInfo{}
	resp, err := n26.callAPI("/api/me", nil)
	if err != nil {
		return nil, err
	}
	err = checkHttpStatus(resp)
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(resp.Body).Decode(accountInfo)
	if err != nil {
		return nil, err
	}
	return accountInfo, nil
}

func (n26 *N26Credentials) Statements() (*N26BankStatements, error) {
	bankStatements := &N26BankStatements{}
	resp, err := n26.callAPI("/api/statements", nil)
	if err != nil {
		return nil, err
	}
	err = checkHttpStatus(resp)
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(resp.Body).Decode(bankStatements)
	if err != nil {
		return nil, err
	}
	return bankStatements, nil
}

func (n26 *N26Credentials) Statement(statementID string) {
	resp, err := n26.callAPI("/api/statements/"+statementID, nil)
	if err != nil {
		fmt.Println(err)
	}
	err = checkHttpStatus(resp)
	if err != nil {
		fmt.Println(err)
	}
	byt, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	ioutil.WriteFile(
		fmt.Sprintf("%s.pdf", statementID),
		byt,
		0750,
	)
}

func (n26 *N26Credentials) Stats() {
	v := &url.Values{}
	v.Set("type", "acct")
	v.Add("from", "1451602800")
	v.Add("to", "1451732400")
	v.Add("numSlices", "25")
	resp, err := n26.callAPI("/api/accounts/stats/", v)
	if err != nil {
		fmt.Println(err)
	}
	err = checkHttpStatus(resp)
	if err != nil {
		fmt.Println(err)
	}
	byt, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(byt))
}

func (n26 *N26Credentials) Cards() (*N26Cards, error) {
	cards := &N26Cards{}
	resp, err := n26.callAPI("/api/v2/cards", nil)
	if err != nil {
		return nil, err
	}
	err = checkHttpStatus(resp)
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(resp.Body).Decode(cards)
	if err != nil {
		return nil, err
	}
	return cards, nil
}

func (n26 *N26Credentials) Status() (*N26AccountStatus, error) {
	accountStatus := &N26AccountStatus{}
	resp, err := n26.callAPI("/api/me/statuses", nil)
	if err != nil {
		return nil, err
	}
	err = checkHttpStatus(resp)
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(resp.Body).Decode(accountStatus)
	if err != nil {
		return nil, err
	}
	return accountStatus, nil
}

func (n26 *N26Credentials) Savings() (*N26Savings, error) {
	savings := &N26Savings{}
	resp, err := n26.callAPI("/api/hub/savings/accounts", nil)
	if err != nil {
		return nil, err
	}
	err = checkHttpStatus(resp)
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(resp.Body).Decode(savings)
	if err != nil {
		return nil, err
	}
	return savings, nil
}

func (n26 *N26Credentials) getToken() (*N26Token, error) {
	tk := &N26Token{}
	v := url.Values{}
	v.Set("grant_type", "password")
	v.Add("username", n26.Email)
	v.Add("password", n26.Password)
	req, err := http.NewRequest("GET", N26APIUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Basic YW5kcm9pZDpzZWNyZXQ=")
	req.URL.Path = "/oauth/token"
	req.URL.RawQuery = v.Encode()
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(resp.Body).Decode(tk)
	if err != nil {
		return nil, err
	}
	return tk, nil
}

func (n26 *N26Credentials) callAPI(path string, v *url.Values) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", N26APIUrl, nil)
	if err != nil {
		return nil, err
	}
	tk, err := n26.getToken()
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tk.AccessToken))
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

func checkHttpStatus(resp *http.Response) error {
	if resp.StatusCode >= http.StatusBadRequest {
		n26Error := &N26Error{}
		err := json.NewDecoder(resp.Body).Decode(n26Error)
		if err != nil {
			return err
		}
		err = fmt.Errorf("%s", strings.Replace(n26Error.ErrorDescription, ":", "", -1))
		return err
	}
	return nil
}
