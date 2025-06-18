package models

/*
   PrivilegeUpdateRequest represents a request to update a user's role for a product.
   Used by superadmins and product admins to promote, demote, or remove privileges.
*/
type PrivilegeUpdateRequest struct {
	Email   string `json:"email" bson:"name"`
	Product string `json:"product"`
	Role    string `json:"role"`
}
