package models

type CreateAddressRequest struct {
	Label string  `json:"label"`
	Lat   float64 `json:"lat"`
	Lng   float64 `json:"lng"`
}
