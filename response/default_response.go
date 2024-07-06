package response

// AddressResponse represents address data without timestamps
type DefaultResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
}
