package models

type MemberInfoResponse struct {
	Email     string `json:"email"`
	Username  string `json:"username"`
	CreatedAt int64  `json:"createdAt"`
}

type MemberLinksResponse struct {
	Links []*LinkDetails `json:"links"`
}
