package models

// User represents a user in the system.
type User struct {
	ID           string            `json:"id,omitempty" bson:"_id,omitempty"`
	Name         string            `json:"name" bson:"name"`
	Email        string            `json:"email" bson:"email"`
	Password     string            `json:"password,omitempty" bson:"password"`
	ProductRoles map[string]string `json:"productRoles" bson:"productRoles"`
	SuperAdmin   bool              `json:"superadmin" bson:"superadmin"`
}
