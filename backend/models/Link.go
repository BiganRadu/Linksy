package models

type Link struct {
	ID             string        `bson:"_id"`
	Title          string        `bson:"title"`
	MemberEmail    string        `bson:"member_email"`
	AccessMode     int           `bson:"access_mode"`
	AllowedIps     []int32       `bson:"allowed_ips"`
	BlackListedIps []int32       `bson:"black_listed_ips"`
	AccessEntries  []AccessEntry `bson:"access_entries"`
	ReferencedLink string        `bson:"referenced_link"`
	CreatedAt      int64         `bson:"created_at"`
}

type AccessEntry struct {
	HourStart  int64          `bson:"hour_start"`
	CountryMap map[string]int `bson:"country"`
	DeviceMap  map[string]int `bson:"device"`
	OsMap      map[string]int `bson:"os"`
	Accesses   int            `bson:"accesses"`
}
