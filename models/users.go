package models

// User represents a user in the system.
type User struct {
	ID              string            `json:"id,omitempty" bson:"_id,omitempty"`
	Name            string            `json:"name"`
	Email           string            `json:"email"`
	Password        string            `json:"password"`
	ConfirmPassword string            `json:"confirmPassword" bson:"-"` // <-- This will NOT be saved to MongoDB
	ProductRoles    map[string]string `json:"productRoles,omitempty"`
	SuperAdmin      bool              `json:"superadmin,omitempty"`
}
