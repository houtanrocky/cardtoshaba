package zarinhub

type ConvertRequest struct {
	Card string `json:"card"`
}

type Response struct {
	Data Data `json:"data"`
	Meta Meta `json:"meta"`
}

type Data struct {
	BankName           string `json:"bankName,omitempty"`
	Card               string `json:"card,omitempty"`
	Deposit            string `json:"deposit,omitempty"`
	DepositDescription string `json:"depositDescription,omitempty"`
	DepositOwners      string `json:"depositOwners,omitempty"`
	DepositStatus      string `json:"depositStatus,omitempty"`
	EnglishBankName    string `json:"englishBankName,omitempty"`
	Iban               string `json:"iban,omitempty"`
}

type Meta struct {
	Code         int         `json:"code"`
	ErrorMessage *string     `json:"errorMessage"`
	ErrorType    *string     `json:"errorType"`
	Errors       []ErrorItem `json:"errors"`
	IsSuccess    bool        `json:"isSuccess"`
	Message      string      `json:"message"`
	Status       string      `json:"status"`
	TrackID      *string     `json:"trackId"`
}

type ErrorItem struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message,omitempty"`
}
