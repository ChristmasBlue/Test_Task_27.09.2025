package models

type Response struct {
	Status  string `json:"status"`
	ID      string `json:"id"`
	Message string `json:"message"`
	Data    []byte `json:"data,omitempty"`
	Details string `json:"details,omitempty"`
}
