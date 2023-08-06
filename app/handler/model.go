package handler

type Response struct {
	ErrorCode   int32       `json:"errorCode"`
	Message     string      `json:"message"`
	Description string      `json:"description,omitempty"`
	Data        interface{} `json:"data,omitempty"`
}