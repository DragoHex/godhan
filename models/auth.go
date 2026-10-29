package models

import "encoding/json"

// RenewTokenResponse contains extended credentials and account metadata returned upon token renewal.
// DhanHQ returns {"token": "...", "createTime": "...", "expiryTime": "..."} in the live gateway,
// while documentation also references "accessToken". Both fields are kept in sync.
type RenewTokenResponse struct {
	DhanClientID         string     `json:"dhanClientId"`
	DhanClientName       string     `json:"dhanClientName"`
	DhanClientUCC        string     `json:"dhanClientUcc"`
	GivenPowerOfAttorney bool       `json:"givenPowerOfAttorney"`
	AccessToken          string     `json:"accessToken"`
	Token                string     `json:"token"`
	CreateTime           BrokerTime `json:"createTime"`
	ExpiryTime           BrokerTime `json:"expiryTime"`
}

// UnmarshalJSON normalizes DhanHQ renewal responses, ensuring both Token and AccessToken are populated.
func (r *RenewTokenResponse) UnmarshalJSON(data []byte) error {
	type Alias RenewTokenResponse
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(r),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	if r.AccessToken == "" && r.Token != "" {
		r.AccessToken = r.Token
	} else if r.Token == "" && r.AccessToken != "" {
		r.Token = r.AccessToken
	}
	return nil
}
