package pkgutil

import (
	"encoding/json"
	"net/http"
)

type HTTPResponse struct {
	Code    int    `json:"code" example:"200"`
	Status  string `json:"status" example:"OK"`
	Message string `json:"message,omitempty" example:"Success"`
	Data    any    `json:"data,omitempty" `
	Errors  []any  `json:"errors,omitempty" `
}

func (h HTTPResponse) MarshalJSON() ([]byte, error) {
	type Alias HTTPResponse
	alias := &struct {
		*Alias
	}{
		Alias: (*Alias)(&h),
	}

	alias.Status = http.StatusText(h.Code)
	if alias.Status == alias.Message {
		alias.Message = ""
	}

	return json.Marshal(alias)
}

type PaginationResponse struct {
	TotalData int `json:"total_data" example:"1"`
	TotalPage int `json:"total_page" example:"1"`
	Page      int `json:"page" example:"1"`
	Limit     int `json:"limit" example:"10"`
	Data      any `json:"data,omitempty" `
}

type ErrValidationResponse struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}
