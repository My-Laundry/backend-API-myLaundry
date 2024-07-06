package response

// UserResponse represents user data without timestamps
//
//	type UserResponse struct {
//		ID        uint              `json:"id"`
//		Username  string            `json:"username"`
//		Email     string            `json:"email"`
//		Role      string            `json:"role"`
//		Addresses []AddressResponse `json:"addresses"`
//	}
type UserResponse struct {
	ID        uint               `json:"id"`
	Username  string             `json:"username"`
	Email     string             `json:"email"`
	Role      *string            `json:"role,omitempty"`
	Addresses *[]AddressResponse `json:"addresses,omitempty"`
}
