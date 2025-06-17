package models

type PrivilegeUpdateRequest struct {
	Email   string `json:"email" bson:"name"`
	Product string `json:"product"`
	Role    string `json:"role"`
}
