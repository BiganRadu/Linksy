package models

// RedirectLinkRequest represents the request structure for getting a redirect link.
type RedirectLinkRequest struct {
	Ip      uint32 `json:"ip"`
	Country string `json:"country"`
	Device  string `json:"device"`
	Os      string `json:"os"`
}

// RedirectLinkResponse represents the response structure for a redirect link.
type RedirectLinkResponse struct {
	RedirectLink string `json:"redirect_link"`
}

// LinkDetails represents the details of a link.
type LinkDetails struct {
	// ID is the unique identifier for the link.
	ID string `json:"id"`
	// Title is the title of the link, typically extracted from the webpage.
	Title string `json:"title"`
	// Icon is the URL of the link's icon or favicon.
	Icon string `json:"icon"`
	// ReferencedLink is the original URL that the link points to.
	ReferencedLink string `json:"referenced_link"`
	// CreatedAt is the timestamp when the link was created in Unix format.
	CreatedAt int64 `json:"created_at"`
	// QRPicture is the URL of the QR code image for the link.
	QRPicture string `json:"qr_picture"`
}
