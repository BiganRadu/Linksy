package models

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
}

type AccessEntry struct {
	HourStart  int64          `bson:"hour_start"`
	CountryMap map[string]int `bson:"country"`
	DeviceMap  map[string]int `bson:"device"`
	OsMap      map[string]int `bson:"os"`
	Accesses   int            `bson:"accesses"`
}
