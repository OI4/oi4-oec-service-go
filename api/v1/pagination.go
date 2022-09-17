package v1

type PaginationRequest struct {
	PerPage uint32 `json:"PerPage"`
	Page    uint32 `json:"Page"`
}

type PaginationResponse struct {
	TotalCount   uint32 `json:"TotalCount"`
	PerPage      uint32 `json:"PerPage"`
	Page         uint32 `json:"Page"`
	HasNext      bool   `json:"HasNext"`
	PaginationId string `json:"PaginationId"`
}
