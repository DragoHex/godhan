package models

type TPINResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}


type EdisFormRequest struct {
	ISIN     string `json:"isin"`
	Qty      int32  `json:"qty"`
	Exchange string `json:"exchange"`
	Segment  string `json:"segment"`
	Bulk     bool   `json:"bulk"`
}


type EdisBulkFormRequest struct {
	ISIN     []string `json:"isin"`
	Exchange string   `json:"exchange"`
	Segment  string   `json:"segment"`
}


type EdisFormResponse struct {
	DhanClientID string `json:"dhanClientId"`
	EdisFormHTML string `json:"edisFormHtml"`
}


type EdisQtyStatusResponse struct {
	ClientID string `json:"clientId"`
	ISIN     string `json:"isin"`
	TotalQty string `json:"totalQty"`
	AprvdQty string `json:"aprvdQty"`
	Status   string `json:"status"`
	Remarks  string `json:"remarks"`
}

