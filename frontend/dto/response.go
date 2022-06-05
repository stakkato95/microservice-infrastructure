package dto

type ResponseDto struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

type MiddleResponse struct {
	Data string `json:"data"`
}

type ServiceResponseDto struct {
	Service       string      `json:"service,omitempty"`
	X_request_id  string      `json:"x-request-id,omitempty"`
	X_api_user_id string      `json:"x-api-user-id,omitempty"`
	Nested        interface{} `json:"nested,omitempty"`
}
