package models

type UserProfile struct {
	DhanClientID string `json:"dhanClientId"`
	IPFlag       bool   `json:"ipFlag"`
	IP           string `json:"ip"`
}

