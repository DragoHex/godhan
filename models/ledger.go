package models

type BoLedgerEntry struct {
	DhanClientID string `json:"dhanClientId"`
	Narration    string `json:"narration"`
	VoucherDate  string `json:"voucherdate"`
	Exchange     string `json:"exchange"`
	VoucherDesc  string `json:"voucherdesc"`
	VoucherNum   string `json:"vouchernumber"`
	Debit        string `json:"debit"`
	Credit       string `json:"credit"`
	RunBal       string `json:"runbal"`
}

