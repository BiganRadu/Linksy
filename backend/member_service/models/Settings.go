package models

// ChangePasswordRequest represents the request structure for changing a member's password.
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// ChangeNameRequest represents the request structure for changing a member's name.
type ChangeNameRequest struct {
	NewName string `json:"new_name"`
}
