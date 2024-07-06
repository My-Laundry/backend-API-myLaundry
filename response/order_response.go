package response

type OrderResponse struct {
	ID        uint            `json:"id"`
	Status    string          `json:"status"`
	CreatedAt string          `json:"created_at"`
	UpdatedAt string          `json:"updated_at"`
	Customer  UserResponse    `json:"customer"`
	Courier   UserResponse    `json:"courier"`
	Admin     UserResponse    `json:"admin"`
	Service   ServiceResponse `json:"service"`
	Address   AddressResponse `json:"address"`
}
