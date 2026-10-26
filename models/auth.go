package models

// RenewTokenResponse contains extended credentials and account metadata returned upon token renewal.
type RenewTokenResponse struct {
	DhanClientID         string     `json:"dhanClientId"`
	DhanClientName       string     `json:"dhanClientName"`
	DhanClientUCC        string     `json:"dhanClientUcc"`
	GivenPowerOfAttorney bool       `json:"givenPowerOfAttorney"`
	AccessToken          string     `json:"accessToken"`
	ExpiryTime           BrokerTime `json:"expiryTime"`
}
