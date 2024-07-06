package response

// AddressResponse represents address data without timestamps
type AddressResponse struct {
	ID            uint   `json:"id"`
	CustomerID    uint   `json:"customer_id"`
	ReceiverName  string `json:"receiver_name"`
	PhoneNumber   string `json:"phone_number"`
	HouseNumber   string `json:"house_number"`
	ResidenceName string `json:"residence_name"`
	AddressNotes  string `json:"address_notes"`
	StreetName    string `json:"street_name"`
	District      string `json:"district"`
	SubDistrict   string `json:"sub_district"`
	City          string `json:"city"`
	Area          string `json:"area"`
}
