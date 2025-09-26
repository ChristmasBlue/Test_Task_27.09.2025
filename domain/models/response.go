package models

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    []byte `json:"data,omitempty"`
	Details string `json:"details,omitempty"`
}
