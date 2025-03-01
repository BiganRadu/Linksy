package models

type RedirectLinkRequest struct {
	Ip      uint32 `json:"ip"`
	Country string `json:"country"`
	Device  string `json:"device"`
	Os      string `json:"os"`
}

type RedirectLinkResponse struct {
	RedirectLink string `json:"redirect_link"`
}

type LinkDetails struct {
	ID             string `json:"id"`
	Title          string `json:"title"`
	Icon           string `json:"icon"`
	ReferencedLink string `json:"referenced_link"`
	CreatedAt      int64  `json:"created_at"`
}
