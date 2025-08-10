package model

import "net/http"

type Response struct {
	Status  string      `json:"status"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

type PaginationResponse struct {
	Page       int `json:"page"`
	Size       int `json:"size"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type ServiceResponse struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Data       any    `json:"data,omitempty"`
}

func (r *ServiceResponse) IsEmptyData() bool {
	return r.Data == nil || (r.Data != nil && r.Data == "")
}

func BadRequest(message string, data any) *ServiceResponse {
	return &ServiceResponse{
		StatusCode: http.StatusBadRequest,
		Message:    message,
		Data:       data,
	}
}

func DefaultError(message string, data any) *ServiceResponse {
	return &ServiceResponse{
		StatusCode: http.StatusInternalServerError,
		Message:    message,
		Data:       data,
	}
}

func Success() *ServiceResponse {
	return &ServiceResponse{
		StatusCode: http.StatusOK,
		Message:    "Success",
	}
}

func NotFound(message string) *ServiceResponse {
	return &ServiceResponse{
		StatusCode: http.StatusNotFound,
		Message:    message,
	}
}
