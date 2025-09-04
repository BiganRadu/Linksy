package models

// Link represents a link model in the system.
// It contains fields for the link's ID, title, icon, member email, access mode,
// allowed and blacklisted IPs, access entries, referenced link, creation timestamp,
// whether it has a QR code, and the QR code link.
// The access mode can be one of the predefined constants, such as Default, IpWhiteList,
// or IpBlackList, which determine how access to the link is controlled.
type Link struct {
	ID             string        `bson:"_id" json:"id"`
	Title          string        `bson:"title" json:"title"`
	Icon           string        `bson:"icon" json:"icon"`
	MemberEmail    string        `bson:"member_email" json:"member_email"`
	AccessMode     int           `bson:"access_mode" json:"access_mode"`
	AllowedIps     []uint32      `bson:"allowed_ips" json:"allowed_ips"`
	BlackListedIps []uint32      `bson:"black_listed_ips" json:"black_listed_ips"`
	AccessEntries  []AccessEntry `bson:"access_entries" json:"access_entries"`
	ReferencedLink string        `bson:"referenced_link" json:"referenced_link"`
	CreatedAt      int64         `bson:"created_at" json:"created_at"`
	HasQR          bool          `bson:"has_qr" json:"has_qr"`
	QRLink         string        `bson:"qr_link" json:"qr_link"`
}

// AccessEntry represents an entry in the access log for a link.
// It contains fields for the hour of access, maps for country, device, and OS statistics,
// and the total number of accesses during that hour.
type AccessEntry struct {
	HourStart  int64          `bson:"hour_start"`
	CountryMap map[string]int `bson:"country"`
	OsMap      map[string]int `bson:"os"`
	Accesses   int            `bson:"accesses"`
}
