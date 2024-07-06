package response

// ServiceResponse represents address data without timestamps
type ServiceResponse struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
	Price uint   `json:"price"`
}
