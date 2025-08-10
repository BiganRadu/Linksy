package models

// MemberInfoResponse represents the response structure for member information.
type MemberInfoResponse struct {
	Email     string `json:"email"`
	Username  string `json:"username"`
	CreatedAt int64  `json:"createdAt"`
}

// MemberLinksResponse represents the response structure for member links.
type MemberLinksResponse struct {
	Links []*LinkDetails `json:"links"`
}
